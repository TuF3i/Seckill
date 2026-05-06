# Seckill 秒杀系统部署文档

## 环境依赖

| 组件 | 版本要求 | 用途 |
|------|----------|------|
| Go | >= 1.21 | 编译运行 |
| PostgreSQL | >= 12 | 持久化存储 |
| Redis Sentinel | >= 6 | 高性能缓存 |
| Apache Kafka | >= 2.8 | 消息队列 |
| Nacos | >= 2.0 | 配置中心和服务发现 |

---

## 快速启动

### 1. 启动基础设施

**PostgreSQL：**
```bash
docker run -d --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=adminer \
  -p 5432:5432 postgres:15
```

**Redis Sentinel（参考 `deployments/test/redis-sentinel/docker-compose.yaml`）：**
```bash
cd deployments/test/redis-sentinel
docker-compose up -d
```

**Kafka：**
```bash
docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  confluentinc/cp-kafka:latest
```

**Nacos：**
```bash
docker run -d --name nacos \
  -p 8848:8848 \
  -e MODE=standalone \
  nacos/nacos-server:v2.2.3
```

### 2. 在 Nacos 中创建配置

访问 Nacos 控制台 http://localhost:8848/nacos，创建配置：

- **Data ID**: `seckill`
- **Group**: `REDROCK`
- **格式**: JSON
- **内容**:

```json
{
  "postgresql": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "postgres",
    "defaultDB": "adminer",
    "sslMode": "disable"
  },
  "redis": {
    "masterName": "mymaster",
    "sentinelAddrs": ["localhost:26379"],
    "password": "root",
    "sentinelPassword": "root"
  },
  "kafka": {
    "brokers": ["localhost:9092"]
  },
  "gateway": {
    "listenAddr": "0.0.0.0",
    "listenPort": "8888",
    "monitoringPort": "8889"
  }
}
```

### 3. 初始化数据库

```bash
go run ./cmd/seckill initdb
```

将自动创建 `user_table`、`item_table`、`order_table`。

### 4. 启动服务

```bash
# 启动所有服务
go run ./cmd/seckill all

# 或分别启动
go run ./cmd/seckill gateway    # API 网关 (HTTP :8888)
go run ./cmd/seckill user       # 用户服务
go run ./cmd/seckill item       # 商品服务
go run ./cmd/seckill order      # 订单服务
go run ./cmd/seckill consumer   # 订单消费者
go run ./cmd/seckill payment    # 支付服务
```

---

## 环境变量配置

所有服务通过环境变量配置 Nacos 连接，默认值如下：

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `NACOS_ADDR` | `localhost` | Nacos 地址 |
| `NACOS_PORT` | `8848` | Nacos 端口 |
| `NACOS_USER` | `admin` | Nacos 用户名 |
| `NACOS_PASSWORD` | `admin` | Nacos 密码 |
| `CONFIG_ID` | `seckill` | 配置 Data ID |
| `CONFIG_GROUP` | `REDROCK` | 配置 Group |
| `CONTAINER_NAME` | 随机 7 位字符串 | 容器名（用于 Snowflake ID） |

```bash
# 自定义连接示例
export NACOS_ADDR=192.168.1.100
export NACOS_PORT=8848
export NACOS_USER=myuser
export NACOS_PASSWORD=mypass
export CONFIG_ID=seckill-prod
export CONFIG_GROUP=PROD
go run ./cmd/seckill all
```

---

## 服务注册与发现

所有微服务使用 **Nacos** 进行服务注册与发现：

- 服务端：通过 `server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient))` 注册
- 客户端（Gateway）：通过 `rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient))` 发现

Kitex RPC 服务端口由 Nacos 动态分配。

---

## Redis 分库

| 数据库 | 用途 | 服务 |
|--------|------|------|
| DB 0 | 用户 Token | userSvr |
| DB 1 | 商品库存 | itemSvr |
| DB 2 | 订单信息 | orderSvr, orderConsumer, paymentSvr |

---

## API 接口

完整 API 文档见 `api.md`。

| 服务 | 方法 | 路径 | 鉴权 |
|------|------|------|------|
| 用户 | POST | `/user/register` | 无 |
| 用户 | POST | `/user/login` | 无 |
| 用户 | GET | `/user/logout` | JWT |
| 用户 | GET | `/user/refresh` | JWT Refresh |
| 商品 | POST | `/item/add` | JWT |
| 商品 | POST | `/item/delete` | JWT |
| 商品 | GET | `/item/list` | JWT |
| 商品 | POST | `/item/flash/start` | JWT |
| 商品 | POST | `/item/flash/stop` | JWT |
| 订单 | POST | `/order/create` | JWT |
| 订单 | GET | `/order/paid` | JWT |
| 订单 | GET | `/order/unpaid` | JWT |
| 订单 | GET | `/order/cancelled` | JWT |
| 支付 | POST | `/payment/process` | JWT |

---

## 测试流程

```bash
# 1. 注册
curl -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'

# 2. 登录
curl -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'

# 3. 添加商品（使用返回的 accessToken）
curl -X POST http://localhost:8888/item/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"name":"秒杀商品","stock":100,"price":99.99,"description":"限时秒杀"}'

# 4. 启动秒杀
curl -X POST http://localhost:8888/item/flash/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"itemId":"{itemId}"}'

# 5. 创建订单
curl -X POST http://localhost:8888/order/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"itemId":"{itemId}"}'

# 6. 支付
curl -X POST http://localhost:8888/payment/process \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"orderId":"{orderId}"}'
```

---

## 测试

```bash
go test ./...
```

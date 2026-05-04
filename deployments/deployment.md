# Seckill 秒杀系统部署文档

## 环境依赖

| 组件 | 版本要求 | 用途 |
|------|----------|------|
| Go | >= 1.21 | 编译运行 |
| PostgreSQL | >= 12 | 持久化存储 |
| Redis Sentinel | >= 6 | 高性能缓存 |
| Apache Kafka | >= 2.8 | 消息队列 |
| Nacos | >= 2.0 | 配置中心 |

---

## 快速启动（开发环境）

### 1. 启动基础设施

使用 Docker Compose 一键启动所有依赖服务：

```bash
cd deployments/docker-compose
docker-compose up -d
```

> 如果 docker-compose.yaml 尚未配置，请手动启动以下服务：

**PostgreSQL（默认端口 5432）：**
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=adminer \
  -p 5432:5432 \
  postgres:15
```

**Redis Sentinel（默认端口 26379）：**
```bash
# 启动 Redis 主节点
docker run -d --name redis-master -p 6379:6379 redis:7

# 启动 Redis Sentinel
docker run -d --name redis-sentinel \
  -p 26379:26379 \
  -e REDIS_MASTER_NAME=mymaster \
  -e REDIS_MASTER_HOST=host.docker.internal \
  -e REDIS_MASTER_PORT=6379 \
  -e REDIS_SENTINEL_PASSWORD=root \
  bitnami/redis-sentinel:7
```

**Kafka（默认端口 9092）：**
```bash
docker run -d \
  --name kafka \
  -p 9092:9092 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  confluentinc/cp-kafka:latest
```

**Nacos（默认端口 8848）：**
```bash
docker run -d \
  --name nacos \
  -p 8848:8848 \
  -e MODE=standalone \
  nacos/nacos-server:v2.2.3
```

### 2. 在 Nacos 中创建配置

通过 Nacos 控制台（http://localhost:8848/nacos）或 API 创建配置：

- **Data ID**: `app-config`
- **Group**: `DEFAULT_GROUP`
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

> 配置文件模板见 [configs/config.json](configs/config.json)。

### 3. 初始化数据库

```bash
go run ./cmd/seckill initdb
```

成功输出：
```
Initializing database...
[ConfigLoader] config loaded successfully from nacos (dataId=app-config, group=DEFAULT_GROUP)
[InitDB] database initialized successfully
Database initialized successfully
```

将自动创建以下三张表：

| 表名 | 说明 |
|------|------|
| `user_table` | 用户表 |
| `item_table` | 商品表 |
| `order_table` | 订单表 |

> 如需自定义 Nacos 地址：
> ```bash
> go run ./cmd/seckill initdb \
>   --nacos-host 192.168.1.100 \
>   --nacos-port 8848 \
>   --nacos-user nacos \
>   --nacos-password nacos
> ```

### 4. 启动服务

#### 方式一：一键启动所有服务

```bash
go run ./cmd/seckill all
```

同时启动：UserSvr、ItemSvr、OrderSvr、OrderConsumer、PaymentSvr、API Gateway

#### 方式二：分别启动各服务

打开多个终端，分别执行：

```bash
# 终端 1 - API 网关
go run ./cmd/seckill gateway

# 终端 2 - 用户服务
go run ./cmd/seckill user

# 终端 3 - 商品服务
go run ./cmd/seckill item

# 终端 4 - 订单服务
go run ./cmd/seckill order

# 终端 5 - 订单消费者
go run ./cmd/seckill consumer

# 终端 6 - 支付服务
go run ./cmd/seckill payment
```

#### 方式三：独立运行（各微服务 runc.go）

```bash
go run ./internal/gateway       # API 网关
go run ./internal/userSvr       # 用户服务
go run ./internal/itemSvr       # 商品服务
go run ./internal/orderSvr      # 订单服务
go run ./internal/orderConsumer # 订单消费者
go run ./internal/paymentSvr    # 支付服务
```

---

## 服务端口说明

| 服务 | 协议 | 默认端口 | 说明 |
|------|------|----------|------|
| API Gateway | HTTP | 8888 | 对外 API 接口 |
| API Gateway Monitor | HTTP | 8889 | Prometheus 监控 |
| UserSvr | Kitex RPC | 随机 | 用户服务 |
| ItemSvr | Kitex RPC | 随机 | 商品服务 |
| OrderSvr | Kitex RPC | 随机 | 订单服务 |
| PaymentSvr | Kitex RPC | 随机 | 支付服务 |

> Kitex RPC 服务端口由 Nacos 服务发现动态分配，无需手动指定。

---

## 默认连接信息

### PostgreSQL

| 配置项 | 默认值 |
|--------|--------|
| Host | localhost |
| Port | 5432 |
| User | postgres |
| Password | postgres |
| Database | adminer |

### Redis Sentinel

| 配置项 | 默认值 |
|--------|--------|
| MasterName | mymaster |
| Sentinel Addrs | localhost:26379 |
| Password | root |
| Sentinel Password | root |

| 数据库 | 用途 |
|--------|------|
| DB 0 | 用户 Token 信息 |
| DB 1 | 商品预热信息 |
| DB 2 | 订单信息 |

### Kafka

| 配置项 | 默认值 |
|--------|--------|
| Brokers | localhost:9092 |
| Order Topic | order_topic |

### Nacos

| 配置项 | 默认值 |
|--------|--------|
| Host | localhost |
| Port | 8848 |
| User | nacos |
| Password | nacos |

---

## API 接口

启动后访问 `http://localhost:8888`。完整 API 文档见 [api.md](api.md)。

### 用户服务

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | `/user/register` | 用户注册 | 无 |
| POST | `/user/login` | 用户登录 | 无 |
| GET | `/user/logout` | 退出登录 | JWT |
| GET | `/user/refresh` | 刷新 Token | JWT Refresh |

### 商品服务

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | `/item/add` | 添加商品 | JWT |
| POST | `/item/delete` | 删除商品 | JWT |
| GET | `/item/list` | 获取商品列表 | JWT |
| POST | `/item/flash/start` | 启动秒杀 | JWT |
| POST | `/item/flash/stop` | 停止秒杀 | JWT |

### 订单服务

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | `/order/create` | 创建订单 | JWT |
| GET | `/order/paid` | 已支付订单 | JWT |
| GET | `/order/unpaid` | 未支付订单 | JWT |
| GET | `/order/cancelled` | 已取消订单 | JWT |

### 支付服务

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| POST | `/payment/process` | 处理支付 | JWT |

---

## 测试流程示例

### 1. 注册用户

```bash
curl -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'
```

### 2. 登录获取 Token

```bash
curl -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin123"}'
```

响应示例：
```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 3. 添加商品

```bash
curl -X POST http://localhost:8888/item/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"name":"秒杀商品","stock":100,"price":99.99,"description":"限时秒杀"}'
```

### 4. 启动秒杀

```bash
curl -X POST http://localhost:8888/item/flash/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"itemId":"{itemId}"}'
```

### 5. 创建订单

```bash
curl -X POST http://localhost:8888/order/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"itemId":"{itemId}"}'
```

### 6. 处理支付

```bash
curl -X POST http://localhost:8888/payment/process \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {accessToken}" \
  -d '{"orderId":"{orderId}"}'
```

---

## 测试

```bash
# 运行所有测试
go test ./...

# 运行指定服务测试
go test -v ./internal/itemSvr/...
go test -v ./internal/orderSvr/...
go test -v ./internal/paymentSvr/...

# 运行配置加载器测试
go test -v ./pkg/config/...
```

---

## 配置覆盖

### 通过 Nacos 控制台修改配置

1. 访问 http://localhost:8848/nacos
2. 进入"配置管理" → "配置列表"
3. 找到 `app-config`（DEFAULT_GROUP）
4. 编辑配置内容，点击"发布"

> **注意**：当前配置加载器仅在启动时加载一次，修改配置后需要重启服务。

### PostgreSQL 配置示例

```json
{
  "postgresql": {
    "host": "192.168.1.100",
    "port": 5432,
    "user": "myuser",
    "password": "mypassword",
    "defaultDB": "seckill",
    "sslMode": "require"
  }
}
```

---

## 常见问题

### 1. Nacos 连接失败

```
panic: initdb: create nacos client: ...
```

**解决**：确认 Nacos 已启动，检查 `--nacos-host` 和 `--nacos-port` 参数。

### 2. 数据库连接失败

```
panic: initdb: connect to postgresql: ...
```

**解决**：确认 PostgreSQL 已启动，检查 Nacos 配置中的数据库连接信息。

### 3. Redis 连接失败

```
panic: ... redis: connect: connection refused
```

**解决**：确认 Redis Sentinel 已启动，检查 Sentinel 地址和密码配置。

### 4. Kafka 连接失败

```
Error reading message: kafka: client has run out of available brokers
```

**解决**：确认 Kafka 已启动，检查 `configs/config.json` 中的 `brokers` 配置。

### 5. 端口冲突

```
panic: listen tcp :8888: bind: address already in use
```

**解决**：通过 Nacos 配置修改 `gateway.listenPort`，或关闭占用端口的进程。

---

## 项目结构

```
seckill/
├── cmd/seckill/              # 统一启动入口
├── api/kitex/                # Thrift IDL 定义
├── configs/                  # 配置模板
├── internal/
│   ├── gateway/              # API 网关
│   ├── userSvr/              # 用户服务
│   ├── itemSvr/              # 商品服务
│   ├── orderSvr/             # 订单服务
│   ├── orderConsumer/        # 订单消费者
│   ├── paymentSvr/           # 支付服务
│   └── initdb/               # 数据库初始化
├── infrastructures/          # 基础设施层
│   ├── postgres/             # PostgreSQL 客户端
│   ├── redis/                # Redis Sentinel 客户端
│   ├── kafka/                # Kafka 客户端
│   └── nacos/                # Nacos 客户端
├── pkg/                      # 公共工具包
│   ├── config/               # 配置加载器
│   ├── jwt/                  # JWT 工具
│   └── traceid/              # 链路追踪
└── deployments/              # 部署配置
```

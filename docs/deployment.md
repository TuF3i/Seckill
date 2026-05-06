# Seckill 秒杀系统部署文档

## 环境依赖

| 组件 | 版本要求 | 用途 |
|------|----------|------|
| Go | >= 1.21 | 编译运行 |
| Docker | >= 24.0 | 容器化部署 |
| PostgreSQL 17 + Pgpool 4.6 | - | 持久化存储（主从集群） |
| Redis Sentinel 7.2 | - | 高性能缓存 |
| Redpanda (Kafka 兼容) | latest | 消息队列 |
| Nacos (rnacos) | stable | 配置中心和服务发现 |

---

## 一键部署（推荐）

```bash
cd deployments/docker-compose/seckill
docker compose up -d
```

将自动启动以下 15 个容器：

| 容器名 | 服务 | 端口 |
|--------|------|------|
| `seckill-pg-0` | PostgreSQL 主节点 | 5433 |
| `seckill-pg-1` | PostgreSQL 从节点 | 5434 |
| `seckill-pgpool` | Pgpool 连接池 | 5432 |
| `seckill-redis-master` | Redis 主节点 | 6379 |
| `seckill-redis-slave1/2` | Redis 从节点 | 6380/6381 |
| `seckill-sentinel1/2/3` | Redis Sentinel | 26379-26381 |
| `seckill-redpanda` | Redpanda (Kafka) | 19092 |
| `seckill-nacos` | Nacos 配置中心 | 8848 |
| `seckill-init` | 初始化容器（一次性） | - |
| `seckill-user/item/order/payment-svr` | RPC 微服务 | Kitex RPC |
| `seckill-consumer` | 订单消费者 | - |
| `seckill-app` | API 网关 | **8888** |

> **初始化流程**：`seckill-init` 自动完成：
> 1. 等待 Nacos 就绪
> 2. 通过 REST API 创建配置（dataId=seckill, group=REDROCK）
> 3. 执行 `seckill initdb`（自动建表 + 创建 admin 用户）
>
> 各微服务容器等待 init 成功后自动启动。

### 镜像仓库

本项目使用以下镜像仓库（国内可访问）：

| 镜像 | 来源 |
|------|------|
| `golang:1.25-alpine` | `docker.1ms.run/library/golang` |
| `alpine:3.20` | `docker.1ms.run/library/alpine` |
| `redis:7.2-alpine` | `docker.1ms.run/library/redis` |
| `redpanda:latest` | `docker.1ms.run/redpandadata/redpanda` |
| `rnacos:stable` | `docker.1ms.run/qingpan/rnacos` |
| `postgresql-repmgr:17` | `swr.cn-north-4.myhuaweicloud.com/ddn-k8s/bitnamilegacy` |
| `pgpool:4.6` | `swr.cn-north-4.myhuaweicloud.com/ddn-k8s/bitnami` |

### 数据持久化

所有数据存储使用 Docker 命名卷：

| 卷名 | 用途 |
|------|------|
| `pg-0-data`, `pg-1-data` | PostgreSQL 数据 |
| `redis-master/slave1/slave2-data` | Redis 数据 |
| `sentinel1/2/3-data` | Sentinel 数据 |
| `nacos-data` | Nacos 数据 |

---

## 手动启动

### 1. 启动基础设施

参考 `deployments/docker-compose/` 下各子目录的配置：

```bash
# Nacos
cd deployments/docker-compose/nacos
docker compose up -d

# PostgreSQL + Pgpool
cd deployments/docker-compose/postgres
docker compose up -d

# Redis Sentinel
cd deployments/docker-compose/redis-sentinel
docker compose up -d

# Redpanda
cd deployments/docker-compose/redpanda
docker compose up -d
```

### 2. 在 Nacos 中创建配置

访问 http://localhost:8848 ，创建配置：

- **Data ID**: `seckill`
- **Group**: `REDROCK`
- **格式**: JSON

```json
{
  "postgresql": {
    "host": "localhost",
    "port": 5432,
    "user": "root",
    "password": "root",
    "defaultDB": "adminer",
    "sslMode": "disable"
  },
  "redis": {
    "masterName": "cluster-master",
    "sentinelAddrs": ["localhost:26379"],
    "password": "root",
    "sentinelPassword": "root"
  },
  "kafka": {
    "brokers": ["localhost:19092"]
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

创建三张表：`user_table`、`item_table`、`order_table`，并自动创建 admin 用户（`admin@seckill.com / admin123`）。

### 4. 启动服务

```bash
# 启动所有微服务
go run ./cmd/seckill all

# 或分别启动
go run ./cmd/seckill gateway   # API 网关 (:8888)
go run ./cmd/seckill user      # 用户服务
go run ./cmd/seckill item      # 商品服务
go run ./cmd/seckill order     # 订单服务
go run ./cmd/seckill consumer  # 订单消费者
go run ./cmd/seckill payment   # 支付服务
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
| `CONTAINER_NAME` | 随机 7 位字符串 | 容器名（Snowflake Node ID） |
| `BENCHMARK` | - | 设为 `true` 跳过 5 分钟限购 |

```bash
# 自定义连接示例
export NACOS_ADDR=192.168.1.100
export NACOS_PORT=8848
export NACOS_USER=myuser
export NACOS_PASSWORD=mypass
go run ./cmd/seckill all
```

---

## 服务注册与发现

使用 Nacos 进行服务注册与发现：

- **服务端**：`server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient))`
- **客户端（Gateway）**：`rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient))`

Kitex RPC 服务端口由 Nacos 动态分配。

---

## Redis 分库

| 数据库 | 用途 | 服务 |
|--------|------|------|
| DB 0 | 用户 Token | UserSvr |
| DB 1 | 商品库存 | ItemSvr |
| DB 2 | 订单信息 | OrderSvr, OrderConsumer, PaymentSvr |

---

## 快速验证

```bash
# 管理员登录（initdb 已自动创建）
curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}'
```

详细 API 接口见 README.md。

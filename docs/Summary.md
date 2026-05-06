# Seckill 秒杀系统 Demo

基于 Go 微服务架构的秒杀系统演示项目。

## 技术栈

| 组件 | 技术 | 用途 |
|------|------|------|
| RPC 框架 | Kitex (Thrift) | 微服务间 RPC 通信 |
| API 网关 | Hertz (HTTP) | 对外 HTTP 接口 |
| 数据库 | PostgreSQL + Pgpool + GORM | 持久化存储 |
| 缓存 | Redis Sentinel | 高性能缓存 |
| 消息队列 | Redpanda (Kafka 兼容) | 异步订单处理 |
| 配置/服务发现 | Nacos (rnacos) | 配置中心与注册 |
| JWT | golang-jwt v5 | 身份认证 |
| 统一入口 | Cobra | 命令行统一启动 |

## 项目结构

```
seckill/
├── api/kitex/                    # Thrift IDL 定义
│   ├── userSvr.thrift
│   ├── itemSvr.thrift
│   ├── orderSvr.thrift
│   └── paymentSvr.thrift
├── cmd/seckill/                  # Cobra 统一入口
│   ├── root.go                   # 根命令 + main()
│   ├── all.go                    # 启动所有微服务
│   ├── initdb.go                 # 初始化数据库
│   ├── user/item/order/...       # 各服务启动命令
├── internal/
│   ├── gateway/                  # API 网关 (Hertz)
│   ├── userSvr/                  # 用户服务
│   ├── itemSvr/                  # 商品服务
│   ├── orderSvr/                 # 订单服务
│   ├── orderConsumer/            # 订单消费者 (Kafka)
│   ├── paymentSvr/               # 支付服务
│   └── initdb/                   # 数据库初始化
├── infrastructures/              # 基础设施层
│   ├── postgres/                 # PostgreSQL 客户端
│   ├── redis/                    # Redis Sentinel 客户端
│   ├── kafka/                    # Kafka 客户端
│   └── nacos/                    # Nacos 客户端
├── pkg/                          # 公共工具包
│   ├── config/                   # Nacos 配置加载器
│   ├── jwt/                      # JWT 工具
│   ├── env/                      # 环境变量加载
│   └── stringToNodeID/           # Snowflake ID 生成
├── deployments/
│   └── docker-compose/           # Docker Compose 部署
│       ├── seckill/              # 全量编排
│       ├── postgres/             # PostgreSQL 集群
│       ├── redis-sentinel/       # Redis Sentinel 集群
│       ├── nacos/                # Nacos 配置中心
│       └── redpanda/             # Redpanda 消息队列
├── scripts/
│   ├── benchmark/main.go         # 压测程序
│   └── stress-test-guide.md      # 压测指南
└── configs/                      # 配置模板
```

---

## 快速部署

### 一键启动（推荐）

```bash
cd deployments/docker-compose/seckill
docker compose up -d
```

将自动启动所有服务：

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
| `seckill-init` | 初始化容器 | 一次性，退出 |
| `seckill-user-svr` | 用户服务 | Kitex RPC |
| `seckill-item-svr` | 商品服务 | Kitex RPC |
| `seckill-order-svr` | 订单服务 | Kitex RPC |
| `seckill-payment-svr` | 支付服务 | Kitex RPC |
| `seckill-consumer` | 订单消费者 | - |
| `seckill-app` | API 网关 | **8888** |

> `seckill-init` 会自动在 Nacos 中创建配置并执行数据库初始化，无需手动操作。

### 手动启动

也可以不依赖 Docker，单独启动各组件，然后手动运行服务：

```bash
# 1. 初始化数据库（需要 Nacos 已就绪且有配置）
go run ./cmd/seckill initdb

# 2. 启动所有微服务
go run ./cmd/seckill all

# 或分别启动
go run ./cmd/seckill gateway
go run ./cmd/seckill user
go run ./cmd/seckill item
go run ./cmd/seckill order
go run ./cmd/seckill consumer
go run ./cmd/seckill payment
```

---

## 快速体验

### 1. 接口验证

```bash
# 管理员登录（initdb 已自动创建）
curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}'
```

### 2. 添加商品并启动秒杀

```bash
# 登录获取 Token
TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")

# 添加商品
ITEM_ID=$(curl -s -X POST http://localhost:8888/item/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"秒杀商品","stock":1000,"price":9.99,"description":"限时秒杀"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['itemId'])")
echo "Item ID: $ITEM_ID"

# 启动秒杀
curl -s -X POST http://localhost:8888/item/flash/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"itemId\":\"$ITEM_ID\"}"
```

### 3. 用户秒杀下单

```bash
# 注册用户
curl -s -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"user123"}'

# 登录
USER_TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"user123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")

# 创建订单
ORDER_ID=$(curl -s -X POST http://localhost:8888/order/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"itemId\":\"$ITEM_ID\"}" | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['orderId'])")
echo "Order ID: $ORDER_ID"

# 支付
curl -s -X POST http://localhost:8888/payment/process \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"orderId\":\"$ORDER_ID\"}"
```

---

## 压测指南

### 准备压测数据

```bash
# 登录管理员
TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")

# 添加高库存商品
ITEM_ID=$(curl -s -X POST http://localhost:8888/item/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"压测商品","stock":100000,"price":0.01,"description":"benchmark"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['itemId'])")

# 启动秒杀
curl -s -X POST http://localhost:8888/item/flash/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"itemId\":\"$ITEM_ID\"}"

# 注册压测用户并获取 Token
curl -s -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"email":"bench@test.com","password":"bench123"}'
USER_TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bench@test.com","password":"bench123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")
```

### 运行压测

```bash
cd scripts/benchmark

# 一键压测（自动准备数据 + 压测）
bash bench.sh

# 自定义参数
C=100 N=10000 bash bench.sh

# 自定义 Gateway
GATEWAY=http://192.168.1.100:8888 bash bench.sh

# 或手动指定参数
go run main.go --item $ITEM_ID --token $USER_TOKEN -c 100 -n 10000
```

> **注意**：默认有 5 分钟限购策略，同一用户只能秒杀一次。
> 如需无限购压测，启动 ItemSvr 时设置环境变量 `BENCHMARK=true`。

### 压测输出

```
=== Throughput ===
QPS (Queries Per Second):       694.44
TPS (Transactions Per Second):  687.50
Elapsed:                        7.2s

=== Latency (ms) ===
Min:    1.23
Avg:    5.67
Max:    45.89
P50:    4.21
P90:    9.87
P99:    18.45
P99.9:  32.10

=== Error Breakdown ===
  status=41010 info=Insufficient Stock: 50
```

---

## 核心业务逻辑

### 秒杀下单流程

```
POST /order/create { "itemId" }
  → Gateway: 从 JWT 解析 userId
  → ItemSvr.PrepareOrder(userId, itemId)
      ├── Lua 脚本原子执行：
      │   1. 检查秒杀状态
      │   2. 检查 5 分钟限购
      │   3. DECR 库存
      │   4. 设置限购标记 (TTL=5min)
      ├── 返回 price
  → OrderSvr.CreateOrder(userId, itemId, price)
      └── Kafka 异步持久化
```

### Redis 分库

| 数据库 | 用途 | 服务 |
|--------|------|------|
| DB 0 | 用户 Token | UserSvr |
| DB 1 | 商品库存/秒杀状态 | ItemSvr |
| DB 2 | 订单信息 | OrderSvr, OrderConsumer, PaymentSvr |

### 业务规则

- **5 分钟限购**：同一用户 5 分钟内禁止重复购买同一商品
- **10 分钟自动取消**：未支付订单超时自动取消

---

## 微服务总览

### Gateway (API 网关 · HTTP :8888)

| 路径 | 方法 | 说明 | 鉴权 |
|------|------|------|------|
| `/user/register` | POST | 用户注册 | 无 |
| `/user/login` | POST | 用户登录 | 无 |
| `/user/logout` | GET | 退出登录 | JWT |
| `/user/refresh` | GET | 刷新 Token | Refresh |
| `/item/add` | POST | 添加商品 | JWT |
| `/item/delete` | POST | 删除商品 | JWT |
| `/item/list` | GET | 商品列表 | JWT |
| `/item/flash/start` | POST | 启动秒杀 | JWT |
| `/item/flash/stop` | POST | 停止秒杀 | JWT |
| `/order/create` | POST | 创建订单 | JWT |
| `/order/paid` | GET | 已支付订单 | JWT |
| `/order/unpaid` | GET | 未支付订单 | JWT |
| `/order/cancelled` | GET | 已取消订单 | JWT |
| `/payment/process` | POST | 处理支付 | JWT |

### UserSvr (用户服务)

| RPC 方法 | 说明 |
|----------|------|
| `RegisterUser` | 注册 |
| `Login` | 登录，返回 JWT |
| `Logout` | 退出 |
| `RefreshAccessToken` | 刷新 Token |
| `VerifyAccessToken` | 验证 Token |
| `VerifyRefreshToken` | 验证 Refresh Token |

### ItemSvr (商品服务)

| RPC 方法 | 说明 |
|----------|------|
| `AddItem` | 添加商品（秒杀进行中禁止） |
| `DeleteItem` | 删除商品（秒杀进行中禁止） |
| `StartFlashSale` | 启动秒杀（预热 Redis） |
| `StopFlashSale` | 停止秒杀 |
| `ListItems` | 商品列表（秒杀时查 Redis） |
| `GetItem` | 商品详情 |
| `PrepareOrder` | 预下单（Lua 原子检查+扣库存） |

### OrderSvr (订单服务)

| RPC 方法 | 说明 |
|----------|------|
| `CreateOrder` | 创建订单（Kafka 异步） |
| `QueryPaidOrders` | 已支付订单 |
| `QueryUnpaidOrders` | 未支付订单 |
| `QueryCancelledOrders` | 已取消订单 |

### PaymentSvr (支付服务)

| RPC 方法 | 说明 |
|----------|------|
| `ProcessPayment` | 处理支付 |

### OrderConsumer (订单消费者)

- 消费 `order_topic` 消息 → 持久化订单
- 每分钟扫描未支付订单 → 超时 10 分钟自动取消

---

## 测试

```bash
# 运行所有测试
go test ./...

# 运行指定服务测试
go test -v ./internal/itemSvr/...
go test -v ./internal/orderSvr/...
go test -v ./internal/paymentSvr/...
go test -v ./pkg/...
```

---

## 详细文档

| 文档 | 路径 | 内容 |
|------|------|------|
| API 文档 | `api.md` | 所有 HTTP 接口定义 |
| 部署文档 | `deployments/deployment.md` | Docker Compose 详细说明 |
| 压测指南 | `scripts/stress-test-guide.md` | 压测原理与多方案 |
| 配置模板 | `configs/config.json` | Nacos 配置模板 |
| 秒杀 Lua | `internal/itemSvr/core/cache/script/` | Redis Lua 原子脚本 |

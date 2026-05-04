# Seckill 秒杀系统 Demo

基于 Go 微服务架构的秒杀系统演示项目。

## 技术栈

| 组件 | 技术 | 用途 |
|------|------|------|
| RPC 框架 | Kitex (Thrift) | 微服务间 RPC 通信 |
| API 网关 | Hertz (HTTP) | 对外 HTTP 接口 |
| 数据库 | PostgreSQL + GORM | 持久化存储 |
| 缓存 | Redis Sentinel | 高性能缓存 |
| 消息队列 | Kafka (segmentio/kafka-go) | 异步订单处理 |
| 服务发现 | Nacos | 服务注册与配置 |
| JWT | golang-jwt v5 | 身份认证 |
| 统一入口 | Cobra | 命令行统一启动 |

## 项目结构

```
seckill/
├── api/kitex/                    # Thrift IDL 定义
│   ├── userSvr.thrift            # 用户服务接口
│   ├── itemSvr.thrift            # 商品服务接口
│   ├── orderSvr.thrift           # 订单服务接口
│   └── paymentSvr.thrift         # 支付服务接口
├── cmd/seckill/                  # Cobra 统一入口
│   ├── root.go                   # 根命令 + main()
│   ├── all.go                    # 启动所有微服务
│   ├── user.go                   # 启动 UserSvr
│   ├── item.go                   # 启动 ItemSvr
│   ├── order.go                  # 启动 OrderSvr
│   ├── consumer.go               # 启动 OrderConsumer
│   ├── payment.go                # 启动 PaymentSvr
│   └── gateway.go                # 启动 API Gateway
├── internal/
│   ├── userSvr/                  # 用户服务 (已完成)
│   ├── itemSvr/                  # 商品服务
│   ├── orderSvr/                 # 订单服务
│   ├── orderConsumer/            # 订单消费者
│   ├── paymentSvr/               # 支付服务
│   └── gateway/                  # API 网关
├── infrastructures/              # 基础设施层
│   ├── postgres/client.go        # PostgreSQL 客户端
│   ├── redis/client.go           # Redis Sentinel 客户端
│   ├── kafka/                    # Kafka 客户端
│   └── nacos/client.go           # Nacos 客户端
├── pkg/                          # 公共工具包
│   ├── jwt/jwt.go                # JWT 工具
│   ├── traceid/traceid.go        # 链路追踪
│   └── enumTransfer/             # 枚举转换
├── deployments/docker-compose/   # 部署配置
└── go.mod
```

## 微服务总览

### 1. UserSvr（用户服务）

**端口**: Kitex RPC 默认端口

**RPC 接口**:

| 方法 | 说明 |
|------|------|
| `RegisterUser` | 用户注册 |
| `Login` | 用户登录，返回 JWT |
| `Logout` | 用户退出登录 |
| `RefreshAccessToken` | 刷新访问令牌 |
| `VerifyAccessToken` | 验证访问令牌 |
| `VerifyRefreshToken` | 验证刷新令牌 |

### 2. ItemSvr（商品服务）

**Redis DB**: 1 (商品信息)

**RPC 接口**:

| 方法 | 说明 | 请求参数 | 响应 |
|------|------|----------|------|
| `AddItem` | 添加商品 | name, stock, price, description | itemId |
| `DeleteItem` | 删除商品 | id | - |
| `StartFlashSale` | 启动秒杀(预热) | itemId | - |
| `StopFlashSale` | 停止秒杀 | itemId | - |

**秒杀启动流程**:
1. 校验商品存在且秒杀未开始
2. DB 标记秒杀状态为进行中 (flash_status=1)
3. Redis 预热商品库存 (WarmUpItemStock)
4. Redis 标记秒杀状态 (SetFlashStatus)

**业务错误码**:

| 状态码 | 说明 |
|--------|------|
| 41001 | 无效商品名称 |
| 41002 | 无效商品库存 |
| 41003 | 无效商品价格 |
| 41004 | 无效商品 ID |
| 41005 | 商品不存在 |
| 41006 | 秒杀已开始 |
| 41007 | 秒杀未开始 |

### 3. OrderSvr（订单服务）

**Redis DB**: 2 (订单信息)

**MQ Topic**: `order_topic`

**RPC 接口**:

| 方法 | 说明 | 请求参数 | 响应 |
|------|------|----------|------|
| `CreateOrder` | 创建订单(MQ异步) | userId, itemId, price | orderId |
| `QueryPaidOrders` | 查询已支付订单 | userId | OrderInfo[] |
| `QueryUnpaidOrders` | 查询未支付订单 | userId | OrderInfo[] |
| `QueryCancelledOrders` | 查询已取消订单 | userId | OrderInfo[] |

**订单创建流程**:
1. 校验请求参数
2. 生成 UUID 订单 ID
3. 构建 JSON 消息 (OrderId, UserId, ItemId, Price)
4. 发送到 Kafka `order_topic`
5. 返回 orderId

**业务错误码**:

| 状态码 | 说明 |
|--------|------|
| 42001 | 无效订单 ID |
| 42002 | 无效用户 ID |
| 42003 | 无效商品 ID |
| 42004 | 无效价格 |
| 42005 | 订单不存在 |
| 42007 | MQ 发送失败 |

### 4. OrderConsumer（订单消费者）

**Redis DB**: 2 (订单信息)

**功能**:
1. 消费 Kafka `order_topic` 消息
2. 将订单持久化到 PostgreSQL `order_table`
3. 向 Redis 写入订单状态 (供 OrderSvr 查询)
4. **订单超时自动取消**: 定期扫描未支付订单，超时 10 分钟自动取消

**超时检查机制**:
- 启动后台 goroutine `timeoutChecker`
- 每分钟扫描一次未支付订单
- 订单创建超过 10 分钟后自动更新状态为已取消 (3)
- 同步更新 Redis 缓存

### 5. PaymentSvr（支付服务）

**Redis DB**: 2 (订单信息)

**RPC 接口**:

| 方法 | 说明 | 请求参数 | 响应 |
|------|------|----------|------|
| `ProcessPayment` | 处理支付 | orderId, userId | bool |

**支付流程**:
1. 校验订单 ID
2. 查询订单是否存在且未支付
3. DB 更新订单状态为已支付 (2)
4. Redis 删除订单缓存 (DelOrderCache)
5. 返回 true

**业务错误码**:

| 状态码 | 说明 |
|--------|------|
| 43001 | 无效订单 ID |
| 43002 | 订单不存在 |
| 43003 | 订单已支付 |
| 43004 | 支付处理失败 |

### 6. Gateway（API 网关）

**端口**: Hertz HTTP 默认端口 8888

**路由**:

| 路径 | 方法 | 说明 |
|------|------|------|
| `/user/register` | POST | 用户注册 |
| `/user/login` | POST | 用户登录 |
| `/user/logout` | GET | 退出登录(JWT 鉴权) |
| `/user/refresh` | GET | 刷新 Token(JWT 刷新) |

## Redis 分库策略

| 数据库 | 用途 | 服务 |
|--------|------|------|
| DB 0 | 用户 Token 信息 | UserSvr |
| DB 1 | 商品预热信息 | ItemSvr |
| DB 2 | 订单信息 | OrderSvr, OrderConsumer, PaymentSvr |

## 业务规则

### 5 分钟限购

秒杀开始后，用户成功抢购同一件商品后，5 分钟内禁止再次购买该商品。

- `Cache.ExistsPurchaseLimit(userId, itemId)` 检查是否在限购期内
- `Cache.SetPurchaseLimit(userId, itemId)` 设置限购标记 (TTL=5min)
- 限购 Key 格式: `item:limit:{itemId}:{userId}`

### 10 分钟自动取消

订单创建后 10 分钟内未支付，系统自动取消。

- OrderConsumer 后台 `timeoutChecker` 每分钟扫描
- `dao.GetUnpaidOrders()` 查询所有未支付订单
- 超时后更新状态为已取消 (3)，同步更新 Redis

## 生命周期管理

每个微服务实现以下生命周期函数：

```go
func OnCreate()   // 资源创建：初始化 DB、Redis、Kafka 连接
func OnDestory()  // 资源销毁：关闭连接、清理资源
```

通过 `core/app/app.go` 导出，供 `runc.go`（独立运行）和 `cmd/seckill/`（统一入口）调用。

## 启动方式

### 前置条件

确保以下服务已启动:
- PostgreSQL (localhost:5432)
- Redis Sentinel (localhost:26379)
- Kafka (localhost:9092)
- Nacos (localhost:8848)

### 方式一：统一入口启动

```bash
# 启动所有服务
go run ./cmd/seckill all

# 启动单个服务
go run ./cmd/seckill user      # 用户服务
go run ./cmd/seckill item      # 商品服务
go run ./cmd/seckill order     # 订单服务
go run ./cmd/seckill consumer  # 订单消费者
go run ./cmd/seckill payment   # 支付服务
go run ./cmd/seckill gateway   # API 网关
```

### 方式二：独立启动（各服务 runc.go）

```bash
go run ./internal/userSvr
go run ./internal/itemSvr
go run ./internal/orderSvr
go run ./internal/orderConsumer
go run ./internal/paymentSvr
go run ./internal/gateway
```

## 测试

```bash
# 运行所有单元测试
go test -v ./internal/...

# 运行指定服务测试
go test -v ./internal/itemSvr/...
go test -v ./internal/orderSvr/...
go test -v ./internal/paymentSvr/...
```

## 部署

参考 `deployments/docker-compose/` 目录下的 Docker Compose 配置。

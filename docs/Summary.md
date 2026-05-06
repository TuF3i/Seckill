# Seckill 秒杀系统项目总结

## 架构总览

```
┌─────────────┐     HTTP :8888
│   Gateway   │◄──────────────── 客户端
│  (Hertz)    │
└──────┬──────┘
       │ Kitex RPC (Nacos 服务发现)
       │
  ┌────┼────┬────┬────┬────┐
  │    │    │    │    │    │
  ▼    ▼    ▼    ▼    ▼    ▼
User Item Order Paym. Cons.
Svr   Svr   Svr   Svr   mer
(DB0) (DB1) (DB2) (DB2) (DB2)
       │            │
       ▼            ▼
     Redis      Redpanda
    Sentinel     (Kafka)
       │            │
       └─────┬──────┘
             ▼
        PostgreSQL
        (Pgpool + 主从)
             │
             ▼
           Nacos
         (rnacos)
```

## 核心业务流程

### 秒杀下单

```
POST /order/create { itemId }
  ├── JWT 鉴权 → 获取 userId
  ├── ItemSvr.PrepareOrder(userId, itemId)
  │     └── Lua 脚本原子操作：
  │           1. GET flash status     → 未开始则拒绝
  │           2. EXISTS purchase limit → 已限购则拒绝
  │           3. DECR stock           → 不足则回滚
  │           4. SET purchase limit   → TTL=5min
  │         → 返回 price
  └── OrderSvr.CreateOrder(userId, itemId, price)
        └── Kafka 异步 → OrderConsumer 持久化入库
```

### 秒杀状态管理

```
StartFlashSale:
  DB 更新 flash_status=1
  Redis 预热库存 (WarmUpItemStock)
  Redis 设置秒杀状态 (SetFlashStatus)

StopFlashSale:
  DB 更新 flash_status=0
  Redis 删除缓存 (DelItemFlashCache)
```

### 限购与超时

| 规则 | 实现 | TTL |
|------|------|-----|
| 5 分钟限购 | Redis SET + Lua 原子检查 | 5 min |
| 10 分钟自动取消 | OrderConsumer 定时扫描 | 10 min |

## Redis 分库

| DB | 用途 | Key 示例 |
|----|------|----------|
| 0 | 用户 Token | `token:access:{uid}` |
| 1 | 商品库存/秒杀 | `item:stock:{id}`, `item:flash:{id}`, `item:limit:{id}:{uid}` |
| 2 | 订单状态 | `order:status:{orderId}` |

## 微服务列表

| 服务 | 框架 | 功能 |
|------|------|------|
| Gateway | Hertz | HTTP API 网关，JWT 鉴权，RPC 聚合 |
| UserSvr | Kitex | 用户注册/登录/JWT 验证 |
| ItemSvr | Kitex | 商品 CRUD、秒杀控制、Lua 预下单 |
| OrderSvr | Kitex | 订单创建(Kafka)、订单查询 |
| PaymentSvr | Kitex | 支付处理 |
| OrderConsumer | - | Kafka 消费、超时取消 |

## 基础设施

| 组件 | 版本 | 端口 |
|------|------|------|
| PostgreSQL 17 (主从) | bitnami/repmgr | 5433/5434 |
| Pgpool 4.6 | bitnami | 5432 |
| Redis Sentinel 7.2 | redis:7.2-alpine | 6379-6381, 26379-26381 |
| Redpanda | redpandadata/redpanda | 19092 |
| rnacos | qingpan/rnacos | 8848 |

## 错误码

| 范围 | 服务 | 错误码 |
|------|------|--------|
| 41001-41011 | ItemSvr | 商品/秒杀相关 |
| 42001-42007 | OrderSvr | 订单相关 |
| 43001-43004 | PaymentSvr | 支付相关 |
| 44001-44004 | UserSvr | 用户/JWT 相关 |
| 20000 | - | 成功 |
| 50000 | - | 内部错误 |

## 功能实现进度

| 模块 | 状态 |
|------|------|
| 用户注册/登录 | ✅ |
| JWT 鉴权 | ✅ |
| 商品 CRUD | ✅ |
| 秒杀控制（启动/停止） | ✅ |
| 秒杀下单（Lua 原子脚本） | ✅ |
| 5 分钟限购 | ✅ |
| Kafka 异步订单 | ✅ |
| 订单查询（支付/未支付/取消） | ✅ |
| 支付处理 | ✅ |
| 超时自动取消 | ✅ |
| BENCHMARK 模式（无限购） | ✅ |
| Nacos 配置中心 | ✅ |
| Nacos 服务发现 | ✅ |
| Docker Compose 全量编排 | ✅ |
| 一键压测脚本 | ✅ |
| 压测结果可视化 | ✅ |

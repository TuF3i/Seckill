# Seckill API 接口文档

## 通用说明

### 基础地址

```
http://{gateway_host}:{gateway_port}
```

### 统一响应格式

所有接口返回统一 JSON 结构：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": { ... }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | int32 | 业务状态码 |
| `info` | string | 状态描述 |
| `data` | object/array | 业务数据 |

### 鉴权方式

大部分接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer {accessToken}
```

### 通用状态码

| 状态码 | 说明 |
|--------|------|
| 20000 | 操作成功 |
| 10001 | 空的 JWT 字符串 |
| 50000 | 内部错误 |

---

## 一、用户服务

### 1.1 用户注册

```
POST /user/register
```

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `email` | string | 是 | 邮箱，长度 5-99 |
| `password` | string | 是 | 密码，长度 5-99 |

**请求示例**：

```json
{
  "email": "user@example.com",
  "password": "mypassword"
}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

**业务错误码**：

| 状态码 | 说明 |
|--------|------|
| 40001 | 无效邮箱格式 |
| 40002 | 无效密码格式 |

---

### 1.2 用户登录

```
POST /user/login
```

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `email` | string | 是 | 用户邮箱 |
| `password` | string | 是 | 用户密码 |

**请求示例**：

```json
{
  "email": "user@example.com",
  "password": "mypassword"
}
```

**成功响应**：

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

**业务错误码**：

| 状态码 | 说明 |
|--------|------|
| 40001 | 无效邮箱格式 |
| 40002 | 无效密码格式 |
| 40003 | 密码错误 |

---

### 1.3 退出登录

```
GET /user/logout
```

**鉴权**：JWT Access Token

**请求头**：

```
Authorization: Bearer {accessToken}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

---

### 1.4 刷新访问令牌

```
GET /user/refresh
```

**鉴权**：JWT Refresh Token

**请求头**：

```
Authorization: Bearer {refreshToken}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

## 二、商品服务

### 2.1 添加商品

```
POST /item/add
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 商品名称，1-128 字符 |
| `stock` | int64 | 是 | 商品库存，大于 0 |
| `price` | float64 | 是 | 商品价格，大于 0 |
| `description` | string | 否 | 商品描述 |

**请求示例**：

```json
{
  "name": "限量版运动鞋",
  "stock": 100,
  "price": 599.00,
  "description": "限时秒杀商品"
}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "itemId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }
}
```

---

### 2.2 删除商品

```
POST /item/delete
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `itemId` | string | 是 | 商品 ID |

**请求示例**：

```json
{
  "itemId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

---

### 2.3 获取商品列表

```
GET /item/list
```

**鉴权**：JWT Access Token（仅 ADMIN 角色可访问）

**请求头**：

```
Authorization: Bearer {accessToken}
```

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "限量版运动鞋",
      "stock": 100,
      "price": 599.00,
      "description": "限时秒杀商品"
    }
  ]
}
```

**说明**：
- 仅 ADMIN 角色用户可以调用，非 ADMIN 返回 `Permission Denied`
- 秒杀进行中的商品，`stock` 字段返回 Redis 中的实时库存
- 秒杀未开始的商品，`stock` 字段返回数据库中的原始库存

---

### 2.4 启动秒杀

```
POST /item/flash/start
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `itemId` | string | 是 | 商品 ID |

**请求示例**：

```json
{
  "itemId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**处理流程**：
1. 校验商品存在且秒杀未开始
2. 数据库标记秒杀状态为进行中
3. Redis 预热商品库存
4. Redis 标记秒杀状态

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

---

### 2.5 停止秒杀

```
POST /item/flash/stop
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `itemId` | string | 是 | 商品 ID |

**请求示例**：

```json
{
  "itemId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**处理流程**：
1. 校验秒杀正在进行
2. 数据库标记秒杀为已停止
3. Redis 清除秒杀缓存和库存缓存

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": null
}
```

### 商品服务业务错误码

| 状态码 | 说明 |
|--------|------|
| 41001 | 无效商品名称 |
| 41002 | 无效商品库存 |
| 41003 | 无效商品价格 |
| 41004 | 无效商品 ID |
| 41005 | 商品不存在 |
| 41006 | 秒杀已开始 |
| 41007 | 秒杀未开始 |
| 41008 | 权限不足（仅管理员可操作） |

---

## 三、订单服务

### 3.1 创建订单

```
POST /order/create
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `userId` | string | 是 | 用户 ID |
| `itemId` | string | 是 | 商品 ID |
| `price` | float64 | 是 | 订单金额，大于 0 |

**请求示例**：

```json
{
  "userId": "user-uuid-123",
  "itemId": "item-uuid-456",
  "price": 599.00
}
```

**处理流程**：
1. 校验请求参数
2. 生成 UUID 作为订单 ID
3. 将订单消息发送到 Kafka `order_topic`
4. 异步：OrderConsumer 消费消息并持久化到数据库

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "orderId": "order-uuid-789"
  }
}
```

**业务错误码**：

| 状态码 | 说明 |
|--------|------|
| 42002 | 无效用户 ID |
| 42003 | 无效商品 ID |
| 42004 | 无效价格 |
| 42007 | MQ 发送失败 |

---

### 3.2 查询已支付订单

```
GET /order/paid
```

**鉴权**：JWT Access Token

**请求参数** (URL Query)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `userId` | string | 是 | 用户 ID |

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": [
    {
      "orderId": "order-uuid-789",
      "userId": "user-uuid-123",
      "itemId": "item-uuid-456",
      "price": 599.00,
      "status": 2,
      "createTime": "2026-05-04 12:00:00"
    }
  ]
}
```

**订单状态说明**：

| 状态值 | 说明 |
|--------|------|
| 1 | 未支付 |
| 2 | 已支付 |
| 3 | 已取消 |

---

### 3.3 查询未支付订单

```
GET /order/unpaid
```

**鉴权**：JWT Access Token

**请求参数** (URL Query)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `userId` | string | 是 | 用户 ID |

**响应格式**同 3.2，仅 `status` 为 1。

---

### 3.4 查询已取消订单

```
GET /order/cancelled
```

**鉴权**：JWT Access Token

**请求参数** (URL Query)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `userId` | string | 是 | 用户 ID |

**响应格式**同 3.2，仅 `status` 为 3。

---

## 四、支付服务

### 4.1 处理支付

```
POST /payment/process
```

**鉴权**：JWT Access Token

**请求参数** (JSON Body)：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `orderId` | string | 是 | 订单 ID |
| `userId` | string | 是 | 用户 ID |

**请求示例**：

```json
{
  "orderId": "order-uuid-789",
  "userId": "user-uuid-123"
}
```

**处理流程**：
1. 校验订单存在且处于未支付状态
2. 数据库更新订单状态为已支付 (2)
3. Redis 删除订单缓存

**成功响应**：

```json
{
  "status": 20000,
  "info": "Operation Success",
  "data": {
    "success": true
  }
}
```

**业务错误码**：

| 状态码 | 说明 |
|--------|------|
| 43001 | 无效订单 ID |
| 43002 | 订单不存在 |
| 43003 | 订单已支付 |
| 43004 | 支付处理失败 |

---

## 五、业务规则说明

### 5.1 Redis 分库

| 数据库 | 用途 | 涉及服务 |
|--------|------|----------|
| DB 0 | 用户 Token 信息 | UserSvr |
| DB 1 | 商品预热信息 | ItemSvr |
| DB 2 | 订单信息 | OrderSvr, OrderConsumer, PaymentSvr |

### 5.2 秒杀业务规则

**5 分钟限购**：
- 用户成功抢购同一件商品后，5 分钟内禁止再次购买
- 限购标记存储在 Redis，TTL 为 5 分钟
- 判断接口：`Cache.ExistsPurchaseLimit(userId, itemId)`

**10 分钟自动取消**：
- 订单创建后 10 分钟内未支付，系统自动取消
- OrderConsumer 后台每分钟扫描未支付订单
- 超时后自动将订单状态更新为已取消 (3)
- 同步更新 Redis 缓存

### 5.3 商品信息查询策略

`ListItems` 接口查询商品时：
- 秒杀进行中的商品：从 Redis 获取实时库存
- 秒杀未开始的商品：从数据库获取原始库存
- 仅 ADMIN 角色用户可调用

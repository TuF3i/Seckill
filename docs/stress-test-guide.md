# 秒杀创建订单接口压测指南

## 请求链路

```
POST /order/create { "itemId" }
  → Gateway: 从 JWT 解析 userId
  → ItemSvr.PrepareOrder(userId, itemId)
      ├── Lua 脚本原子执行
      │   ├── 检查秒杀状态
      │   ├── 检查 5 分钟限购（BENCHMARK 模式跳过）
      │   ├── DECR 库存
      │   └── 设置限购标记 (TTL=5min)
      └── 返回 price
  → OrderSvr.CreateOrder(userId, itemId, price)
      └── Kafka → OrderConsumer 持久化入库
```

---

## 一键压测

```bash
cd scripts/benchmark

# 默认参数（50 并发，5000 请求）
bash bench.sh

# 自定义并发和请求数
C=100 N=10000 bash bench.sh

# 自定义 Gateway 地址
GATEWAY=http://192.168.1.100:8888 bash bench.sh

# 自定义库存和价格
STOCK=50000 PRICE=0.01 bash bench.sh
```

脚本自动执行：
1. 管理员登录（`admin@seckill.com / admin123`）
2. 添加商品（库存 10 万）
3. 启动秒杀
4. 注册压测用户 → 获取 Token
5. 运行 Go benchmark 程序

---

## 压测模式

### BENCHMARK 模式（压测推荐）

ItemSvr 启用 `BENCHMARK=true` 环境变量时，跳过 5 分钟限购检查，单用户可连续下单。

在 docker-compose 中设置：

```yaml
seckill-item:
  environment:
    BENCHMARK: "true"
```

或手动启动：

```bash
BENCHMARK=true go run ./cmd/seckill item
```

### 默认模式（有限购）

每用户 5 分钟内只能购买同一商品一次。单用户压测只有第一次请求成功，后续全部返回 41011。

---

## Benchmark 程序

`scripts/benchmark/main.go`

### 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-gateway` | `http://localhost:8888` | Gateway 地址 |
| `-item` | 必填 | 商品 ID |
| `-token` | 必填 | JWT access token |
| `-c` | `50` | 并发数 |
| `-n` | `2000` | 总请求数 |

### 输出指标

| 指标 | 说明 |
|------|------|
| **QPS** | 总请求数 / 总耗时，系统吞吐 |
| **TPS** | 成功请求数 / 总耗时，有效事务吞吐 |
| **Min / Avg / Max** | 延迟极值 (ms) |
| **P50 / P90 / P99 / P99.9** | 百分位延迟 (ms) |
| **Progress** | 实时进度和当前 RPS |
| **Error Breakdown** | 按错误码分类统计 |

### 实际压测结果（50 并发 5000 请求）

```
=== Throughput ===
QPS (Queries Per Second):       1845.02
TPS (Transactions Per Second):  1845.02
Elapsed:                        2.71s

=== Latency (ms) ===
Min:    1.21
Avg:    26.50
Max:    710.52
P50:    8.10
P90:    21.32
P99:    287.02
P99.9:  686.74
```

---

## 错误码说明

| 状态码 | 说明 | 原因 |
|--------|------|------|
| 20000 | 成功 | - |
| 41007 | 秒杀未开始 | 检查 `/item/flash/start` 是否已调用 |
| 41010 | 库存不足 | 增加库存或检查 Redis `GET item:stock:{itemId}` |
| 41011 | 超过限购 | 启用 BENCHMARK 模式或使用多用户 |
| 50000 | 内部错误 | 检查服务日志 |

---

## 多用户压测

不使用 BENCHMARK 模式时，批量注册用户模拟真实场景：

```bash
# 注册 100 个用户
for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8888/user/register \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"user${i}@bench.com\",\"password\":\"pass123\"}" &
done
wait

# 批量登录保存 token
for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8888/user/login \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"user${i}@bench.com\",\"password\":\"pass123\"}" | \
    python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])" >> tokens.txt &
done
wait

# wrk 压测（需自行编写 wrk-order.lua）
wrk -t10 -c100 -d30s -s wrk-order.lua http://localhost:8888/order/create
```

---

## 推荐流程

```
1. docker compose up -d
2. 设置 BENCHMARK=true 重启 ItemSvr
3. bash bench.sh
4. 观察 QPS/TPS/P99 延迟
5. 调整并发参数，重复测试
```

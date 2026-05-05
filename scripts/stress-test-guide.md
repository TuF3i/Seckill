# 秒杀创建订单接口压测指南

## 请求链路

```
POST /order/create { "itemId" }
  → Gateway: 从 JWT 解析 userId
  → ItemSvr.PrepareOrder(userId, itemId)   // Lua 脚本：检查秒杀状态 + 限购 + 扣库存
  → OrderSvr.CreateOrder(userId, itemId, price)  // 发送 Kafka 消息
  → 返回 orderId
```

---

## 总体流程

```
docker-compose 启动基础设施和所有服务
  └── initdb 自动创建 admin 用户 (admin@seckill.com / admin123)

手动操作：
  1. 用 admin 登录 → 添加商品 → 启动秒杀
  2. 注册普通用户 → 登录取 token
  3. 用 benchmark 程序压测

压测时 ItemSvr 支持 BENCHMARK 模式（跳过 5 分钟限购）
```

---

## 前置准备

### 1. 启动项目

```bash
cd deployments/docker-compose/seckill
docker compose up -d
```

容器就绪后（包含 initdb 自动执行），服务已全部启动。

### 2. 准备压测数据

```bash
# 1. 管理员登录（initdb 已创建 admin@seckill.com / admin123）
TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")

# 2. 添加商品（库存 10 万）
ITEM_ID=$(curl -s -X POST http://localhost:8888/item/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"压测商品","stock":100000,"price":0.01,"description":"benchmark"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['itemId'])")
echo "Item ID: $ITEM_ID"

# 3. 启动秒杀
curl -s -X POST http://localhost:8888/item/flash/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"itemId\":\"$ITEM_ID\"}"

# 4. 注册压测用户
curl -s -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{"email":"bench@test.com","password":"bench123"}'

# 5. 获取压测用户的 token
USER_TOKEN=$(curl -s -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bench@test.com","password":"bench123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])")
echo "Token: $USER_TOKEN"
```

---

## 一键压测

在 `scripts/benchmark/` 目录下提供了一键压测脚本，自动完成所有准备工作和压测：

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

## 压测模式选择

本压测脚本只做压测，**不负责启动服务和准备数据**，数据和服务需提前就绪。

### 模式一：默认（有限购）

每用户 5 分钟内只能购买同一商品一次。单用户压测只有在第一次请求成功，后续请求全部返回 41011。

```bash
# 需要先用 docker-compose 正常启动服务
cd scripts/benchmark

go run main.go \
  --item $ITEM_ID \
  --token $USER_TOKEN \
  -c 1 -n 1
```

### 模式二：BENCHMARK 模式（压测推荐，无限购）

ItemSvr 启用 `BENCHMARK=true` 环境变量时，跳过 5 分钟限购检查，单用户可连续下单。

在 docker-compose 中设置：

```yaml
seckill-item:
  environment:
    BENCHMARK: "true"
```

或手动启动 ItemSvr：

```bash
BENCHMARK=true go run ./cmd/seckill item
```

压测命令：

```bash
cd scripts/benchmark

go run main.go \
  --item $ITEM_ID \
  --token $USER_TOKEN \
  -c 50 -n 5000
```

---

## Benchmark 程序

`scripts/benchmark/main.go`：

### 参数

```
-gateway  string  Gateway 地址 (默认: http://localhost:8888)
-item     string  商品 ID（必填）
-token    string  JWT access token（必填）
-c        int     并发数 (默认: 50)
-n        int     总请求数 (默认: 2000)
```

### 使用示例

```bash
# 最小用法
go run main.go --item xxx --token yyy

# 自定义并发和请求数
go run main.go --item xxx --token yyy -c 100 -n 10000

# 自定义 Gateway 地址
go run main.go --gateway http://192.168.1.100:8888 --item xxx --token yyy
```

### 输出示例

```
Seckill Benchmark
  Gateway:     http://localhost:8888
  Item ID:     abc-123
  Token:       eyJhbGciOiJIUzI1...
  Concurrency: 50
  Requests:    5000

=== Benchmark Running ===
  Progress: 1200/5000 (24.0%) | RPS: 720

=== Summary ===
Total:          5000
Completed:      5000 (100.0%)
Success:        4950 (99.0%)
Failed:         50 (1.0%)

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

### 指标说明

| 指标 | 说明 |
|------|------|
| **QPS** | 总请求数 / 总耗时，衡量系统吞吐能力 |
| **TPS** | 成功请求数 / 总耗时，衡量有效事务处理能力 |
| **Min / Avg / Max** | 请求延迟的最小/平均/最大值 |
| **P50** | 中位数延迟，50% 请求在此值以下完成 |
| **P90** | 90% 请求在此值以下完成 |
| **P99** | 99% 请求在此值以下完成，关注长尾延迟 |
| **P99.9** | 99.9% 请求在此值以下完成，极端情况 |
| **Progress** | 运行中的实时进度和当前 RPS |

---

## 错误码说明

| 状态码 | 说明 | 原因 |
|--------|------|------|
| 20000 | 成功 | - |
| 41007 | 秒杀未开始 | 检查是否已调用 `/item/flash/start` |
| 41010 | 库存不足 | 增加商品库存，或检查 Redis `GET item:stock:{itemId}` |
| 41011 | 超过限购 | 启用 BENCHMARK 模式，或使用多用户 |
| 50000 | 内部错误 | 检查服务日志 |

---

## 多用户压测（不使用 BENCHMARK 模式）

如果需要模拟真实多用户场景（带限购），需要批量注册用户并获取 token：

```bash
# 注册 100 个用户
for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8888/user/register \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"user${i}@bench.com\",\"password\":\"pass123\"}" &
done
wait

# 批量登录并保存 token
rm -f tokens.txt
for i in $(seq 1 100); do
  curl -s -X POST http://localhost:8888/user/login \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"user${i}@bench.com\",\"password\":\"pass123\"}" | \
    python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])" >> tokens.txt &
done
wait

# 使用 wrk 压测
wrk -t10 -c100 -d30s -s wrk-order.lua http://localhost:8888/order/create
```

`wrk-order.lua` 示例（需自行适配 itemId）：

```lua
tokens = {}
function setup(thread)
  local file = io.open("tokens.txt", "r")
  local i = 1
  for line in file:lines() do
    tokens[i] = line
    i = i + 1
  end
  file:close()
  thread:set("tokens", tokens)
end

request = function()
  local idx = math.random(#tokens)
  wrk.headers["Authorization"] = "Bearer " .. tokens[idx]
  wrk.headers["Content-Type"] = "application/json"
  wrk.body = '{"itemId":"xxx"}'
  return wrk.format("POST")
end
```

---

## 推荐压测流程

```
1. docker compose up -d
2. 管理员登录 → 添加商品 → 启动秒杀
3. 注册普通用户 → 登录取 token
4. BENCHMARK=true 重启 ItemSvr
5. go run main.go --item xxx --token yyy -c 100 -n 10000
6. 观察 RPS 和错误率，分析瓶颈
```

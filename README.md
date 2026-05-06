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

## 快速部署

```bash
cd deployments/docker-compose/seckill
docker compose up -d
```

一键启动 15 个容器，`seckill-init` 自动在 Nacos 中创建配置并初始化数据库。

详细部署说明请参考 [docs/deployment.md](docs/deployment.md)。

---

## Bench Mark Result

```
========================================
  Seckill One-Click Benchmark
========================================
  Gateway:      http://localhost:8888
  Concurrency:  50
  Requests:     5000
========================================


=== Summary ===
Total:          5000
Completed:      5000 (100.0%)
Success:        5000 (100.0%)
Failed:         0 (0.0%)

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

> 50 并发 5000 请求，100% 成功率，QPS 1845，P50 延迟 8ms。
> 详细信息请参考 [docs/benchmark-result.md](docs/benchmark-result.md)。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [docs/deployment.md](docs/deployment.md) | 部署说明（Docker Compose / 手动） |
| [docs/stress-test-guide.md](docs/stress-test-guide.md) | 压测指南（一键脚本 / 多方案） |
| [docs/benchmark-result.md](docs/benchmark-result.md) | 压测结果详细数据 |
| [docs/NacosUsage.md](docs/NacosUsage.md) | Nacos 接入教程 |
| [docs/Summary.md](docs/Summary.md) | 项目总结与架构说明 |
| [configs/config.json](configs/config.json) | Nacos 配置模板 |

---

## 项目结构

```
seckill/
├── api/kitex/                    # Thrift IDL 定义
├── cmd/seckill/                  # Cobra 统一入口
├── internal/
│   ├── gateway/                  # API 网关 (Hertz)
│   ├── userSvr/                  # 用户服务
│   ├── itemSvr/                  # 商品服务
│   ├── orderSvr/                 # 订单服务
│   ├── orderConsumer/            # 订单消费者 (Kafka)
│   ├── paymentSvr/               # 支付服务
│   └── initdb/                   # 数据库初始化
├── infrastructures/              # 基础设施层
├── pkg/                          # 公共工具包
├── scripts/                      # 压测脚本
└── deployments/                  # Docker Compose 部署
```

## 测试

```bash
go test ./...
```

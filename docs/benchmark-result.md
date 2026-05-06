## Bench Mark Result
```
========================================
  Seckill One-Click Benchmark
========================================
  Gateway:      http://localhost:8888
  Concurrency:  50
  Requests:     5000
  Stock:        100000
========================================

[1/5] Admin login (admin@seckill.com)
  ✓ Token acquired
[2/5] Adding benchmark item
  ✓ Item created: 09c272a5-0729-41ae-92e9-1e9818478331
[3/5] Starting flash sale
  ✓ Flash sale started
[4/5] Registering benchmark user
  ✓ User ready (email: bench-1778044117@test.com)
[5/5] Running benchmark

Seckill Benchmark
  Gateway:     http://localhost:8888
  Item ID:     09c272a5-0729-41ae-92e9-1e9818478331
  Token:       eyJhbGciOiJIUzI1...
  Concurrency: 50
  Requests:    5000

=== Benchmark Running ===
  Progress: 4478/5000 (89.6%) | RPS: 2239

=== Summary ===
Total:          5000
Completed:      5000 (100.0%)
Success:        5000 (100.0%)
Failed:         0 (0.0%)

=== Throughput ===
QPS (Queries Per Second):       2206.35
TPS (Transactions Per Second):  2206.35
Elapsed:                        2.266189487s

=== Latency (ms) ===
Min:    1.11
Avg:    22.15
Max:    600.41
P50:    6.40
P90:    18.65
P99:    235.22
P99.9:  521.89
```

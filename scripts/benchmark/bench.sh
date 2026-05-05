#!/bin/bash
set -e

GATEWAY="${GATEWAY:-http://localhost:8888}"
CONCURRENCY="${C:-50}"
REQUESTS="${N:-5000}"
STOCK="${STOCK:-100000}"
PRICE="${PRICE:-0.01}"

echo "========================================"
echo "  Seckill One-Click Benchmark"
echo "========================================"
echo "  Gateway:      $GATEWAY"
echo "  Concurrency:  $CONCURRENCY"
echo "  Requests:     $REQUESTS"
echo "  Stock:        $STOCK"
echo "========================================"
echo ""

# ── Step 1: Admin Login ──────────────────────────────────────
echo "[1/5] Admin login (admin@seckill.com)"
ADMIN_TOKEN=$(curl -s -X POST "$GATEWAY/user/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@seckill.com","password":"admin123"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])" 2>/dev/null)

if [ -z "$ADMIN_TOKEN" ]; then
  echo "  ✗ Login failed. Is the service running?"
  exit 1
fi
echo "  ✓ Token acquired"

# ── Step 2: Add Item ─────────────────────────────────────────
echo "[2/5] Adding benchmark item"
ITEM_ID=$(curl -s -X POST "$GATEWAY/item/add" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"name\":\"Benchmark Item\",\"stock\":$STOCK,\"price\":$PRICE,\"description\":\"auto-benchmark\"}" | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['itemId'])" 2>/dev/null)

if [ -z "$ITEM_ID" ]; then
  echo "  ✗ Add item failed"
  exit 1
fi
echo "  ✓ Item created: $ITEM_ID"

# ── Step 3: Start Flash Sale ─────────────────────────────────
echo "[3/5] Starting flash sale"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$GATEWAY/item/flash/start" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"itemId\":\"$ITEM_ID\"}")

if [ "$HTTP_CODE" != "200" ]; then
  echo "  ✗ Start flash sale failed (HTTP $HTTP_CODE)"
  exit 1
fi
echo "  ✓ Flash sale started"

# ── Step 4: Register & Login Benchmark User ──────────────────
echo "[4/5] Registering benchmark user"
BENCH_USER="bench-$(date +%s)@test.com"

curl -s -X POST "$GATEWAY/user/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$BENCH_USER\",\"password\":\"bench123\"}" > /dev/null

USER_TOKEN=$(curl -s -X POST "$GATEWAY/user/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$BENCH_USER\",\"password\":\"bench123\"}" | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['data']['accessToken'])" 2>/dev/null)

if [ -z "$USER_TOKEN" ]; then
  echo "  ✗ User login failed"
  exit 1
fi
echo "  ✓ User ready (email: $BENCH_USER)"

# ── Step 5: Run Benchmark ────────────────────────────────────
echo "[5/5] Running benchmark"
echo ""

cd "$(dirname "$0")"
go run main.go \
  --gateway "$GATEWAY" \
  --item "$ITEM_ID" \
  --token "$USER_TOKEN" \
  -c "$CONCURRENCY" \
  -n "$REQUESTS"

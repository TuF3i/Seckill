#!/bin/sh
set -e

echo "Waiting for Nacos..."
until wget -q --spider http://nacos:8848/nacos/ 2>/dev/null; do
  sleep 2
done
echo "Nacos ready. Creating config..."

CONFIG_JSON='{
  "postgresql": {
    "host": "pgpool",
    "port": 5432,
    "user": "'"${POSTGRESQL_USERNAME}"'",
    "password": "'"${POSTGRESQL_PASSWORD}"'",
    "defaultDB": "'"${POSTGRESQL_DATABASE}"'",
    "sslMode": "disable"
  },
  "redis": {
    "masterName": "cluster-master",
    "sentinelAddrs": ["sentinel1:26379","sentinel2:26379","sentinel3:26379"],
    "password": "root",
    "sentinelPassword": "root"
  },
  "kafka": {
    "brokers": ["redpanda:19092"]
  },
  "gateway": {
    "listenAddr": "0.0.0.0",
    "listenPort": "8888",
    "monitoringPort": "8889"
  },
  "opentelemetry": {
    "exportEndpoint": "otel:4317"
  }
}'

curl -s -X POST http://nacos:8848/nacos/v1/cs/configs \
  -d "dataId=${CONFIG_ID}&group=${CONFIG_GROUP}&content=$(echo "${CONFIG_JSON}")&type=json"

echo "Config created. Running initdb..."
exec /app/seckill initdb

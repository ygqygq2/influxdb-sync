#!/bin/bash
# 查询 InfluxDB 3.x (v2 兼容模式) 数据

# 3.x v2 兼容模式配置
url="http://localhost:18088/api/v2/query"
token="test3xtoken"
org="testorg"

echo "查询 InfluxDB 3.x (v2 兼容模式) 数据..."

# Flux 查询：统计 measurements
echo "=== 查询所有 measurements ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> group(columns: ["_measurement"])
  |> distinct(column: "_measurement")
  |> keep(columns: ["_value"])
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

echo ""
echo "=== CPU 数据统计 ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "cpu")
  |> count()
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

echo ""
echo "=== Memory 数据统计 ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "memory")
  |> count()
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

echo ""
echo "=== Disk 数据统计 ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "disk")
  |> count()
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

echo ""
echo "=== 总数据点统计 ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> count()
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

echo ""
echo "=== 最新10条 CPU 数据 ==="
flux_query='
from(bucket: "testbucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "cpu")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 10)
'

curl -s -H "Authorization: Token $token" \
  -H "Content-Type: application/vnd.flux" \
  -H "Accept: application/csv" \
  -XPOST "$url?org=$org" \
  --data-raw "$flux_query"

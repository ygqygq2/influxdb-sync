#!/bin/bash
# 查询 influxdb2 源库和目标库的 cpu/mem 数据条数并对比
url_src="http://localhost:18086/api/v2/query?org=testorg"
url_tgt="http://localhost:18087/api/v2/query?org=testorg"
token="testtoken"
bucket="testbucket"

for m in cpu mem; do
  echo "统计 $m 行数..."
  src_count=$(curl -s -XPOST "$url_src" \
    -H "Authorization: Token $token" \
    -H "Content-Type: application/json" \
    -d '{"query":"from(bucket: \"'$bucket'\") |> range(start: 0) |> filter(fn: (r) => r._measurement == \"'$m'\") |> count()"}' \
    | grep -o 'result.*_value,[0-9]*' | grep -o '[0-9]*$' | head -n1)
  tgt_count=$(curl -s -XPOST "$url_tgt" \
    -H "Authorization: Token $token" \
    -H "Content-Type: application/json" \
    -d '{"query":"from(bucket: \"'$bucket'\") |> range(start: 0) |> filter(fn: (r) => r._measurement == \"'$m'\") |> count()"}' \
    | grep -o 'result.*_value,[0-9]*' | grep -o '[0-9]*$' | head -n1)
  src_count=${src_count:-0}
  tgt_count=${tgt_count:-0}
  echo "源库: $src_count"
  echo "目标库: $tgt_count"
  diff=$((src_count-tgt_count))
  if [ "$diff" -eq 0 ]; then
    echo "同步完成，数据一致"
  else
    echo "差异: $diff 条"
  fi
  echo
done

#!/bin/bash
# 向 influxdb2-src 批量写入大数据量测试数据
url="http://localhost:18086/api/v2/write?org=testorg&bucket=testbucket&precision=ns"
token="testtoken"


# 生成 1 万条 cpu 数据
echo "批量写入 cpu..."
for i in {1..10}; do
  batch=""
  for j in {1..1000}; do
    ts=$(( (i-1)*1000 + j ))
    val=$(awk "BEGIN {print 0.5 + ($ts % 100) * 0.01}")
    batch+="cpu,host=host$((ts%10)),region=uswest value=$val $(( 1690000000000000000 + ts * 1000000 ))\n"
  done
  echo -e "$batch" | curl -s -XPOST "$url" --header "Authorization: Token $token" --data-binary @-
done

# 生成 1 万条 mem 数据
echo "批量写入 mem..."
for i in {1..10}; do
  batch=""
  for j in {1..1000}; do
    ts=$(( (i-1)*1000 + j ))
    val=$(awk "BEGIN {print 20 + ($ts % 100) * 0.1}")
    batch+="mem,host=host$((ts%10)),region=uswest value=$val $(( 1690000000000000000 + ts * 1000000 ))\n"
  done
  echo -e "$batch" | curl -s -XPOST "$url" --header "Authorization: Token $token" --data-binary @-
done

echo "写入完成"

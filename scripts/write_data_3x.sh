#!/bin/bash
# 向 InfluxDB 3.x (v2 兼容模式) 批量写入测试数据

# 3.x v2 兼容模式配置
url="http://localhost:18088/api/v2/write"
token="test3xtoken"
org="testorg"
bucket="testbucket"

echo "向 InfluxDB 3.x (v2 兼容模式) 写入测试数据..."

# 生成 5000 条 cpu 数据
echo "批量写入 cpu 数据..."
for i in {1..5}; do
  batch=""
  for j in {1..1000}; do
    ts=$(( (i-1)*1000 + j ))
    val=$(awk "BEGIN {print 0.3 + ($ts % 80) * 0.01}")
    batch+="cpu,host=3x-host$((ts%8)),region=cloud,version=3x value=$val $(( 1690000000000000000 + ts * 1000000 ))\n"
  done
  echo -e "$batch" | curl -s -H "Authorization: Token $token" \
    -H "Content-Type: text/plain" \
    -XPOST "$url?org=$org&bucket=$bucket&precision=ns" \
    --data-binary @-
  echo "Batch $i/5 completed"
done

# 生成 5000 条 memory 数据
echo "批量写入 memory 数据..."
for i in {1..5}; do
  batch=""
  for j in {1..1000}; do
    ts=$(( (i-1)*1000 + j ))
    val=$(awk "BEGIN {print 15 + ($ts % 60) * 0.2}")
    batch+="memory,host=3x-host$((ts%8)),region=cloud,version=3x used=$val,available=$(awk "BEGIN {print 100-$val}") $(( 1690000000000000000 + ts * 1000000 ))\n"
  done
  echo -e "$batch" | curl -s -H "Authorization: Token $token" \
    -H "Content-Type: text/plain" \
    -XPOST "$url?org=$org&bucket=$bucket&precision=ns" \
    --data-binary @-
  echo "Batch $i/5 completed"
done

# 生成 3000 条 disk 数据
echo "批量写入 disk 数据..."
for i in {1..3}; do
  batch=""
  for j in {1..1000}; do
    ts=$(( (i-1)*1000 + j ))
    val=$(awk "BEGIN {print 40 + ($ts % 50) * 0.5}")
    batch+="disk,host=3x-host$((ts%8)),device=/dev/sda1,region=cloud,version=3x usage_percent=$val $(( 1690000000000000000 + ts * 1000000 ))\n"
  done
  echo -e "$batch" | curl -s -H "Authorization: Token $token" \
    -H "Content-Type: text/plain" \
    -XPOST "$url?org=$org&bucket=$bucket&precision=ns" \
    --data-binary @-
  echo "Batch $i/3 completed"
done

echo "InfluxDB 3.x 测试数据写入完成！"
echo "总计: 5000 cpu + 5000 memory + 3000 disk = 13000 条数据点"

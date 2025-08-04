#!/bin/bash
# 查询 influxdb1 源库和目标库的 cpu/mem 数据条数并对比

db="testdb"
user="admin"
pass="admin123"
src_container="influxdb1-src"
tgt_container="influxdb1-dst"

for m in cpu mem; do
  echo "统计 $m 行数..."
  src_count=$(docker exec $src_container influx -username $user -password $pass -database $db -execute "SELECT COUNT(value) FROM $m" -format csv | awk -F',' 'NR==2{print $3}')
  tgt_count=$(docker exec $tgt_container influx -username $user -password $pass -database $db -execute "SELECT COUNT(value) FROM $m" -format csv | awk -F',' 'NR==2{print $3}')
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

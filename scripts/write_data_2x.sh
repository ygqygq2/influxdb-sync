#!/bin/bash
# 向 influxdb2-src 写入测试数据
curl -i -XPOST 'http://localhost:18086/api/v2/write?org=testorg&bucket=testbucket&precision=s' \
  --header 'Authorization: Token testtoken' \
  --data-raw 'cpu,host=host1,region=uswest value=0.64'
curl -i -XPOST 'http://localhost:18086/api/v2/write?org=testorg&bucket=testbucket&precision=s' \
  --header 'Authorization: Token testtoken' \
  --data-raw 'cpu,host=host2,region=uswest value=0.72'
curl -i -XPOST 'http://localhost:18086/api/v2/write?org=testorg&bucket=testbucket&precision=s' \
  --header 'Authorization: Token testtoken' \
  --data-raw 'mem,host=host1,region=uswest value=33.1'

#!/bin/bash
# 向 influxdb1-src 写入测试数据
curl -i -u admin:admin123 -XPOST 'http://localhost:18086/write?db=testdb' --data-binary 'cpu,host=host1,region=uswest value=0.64'
curl -i -u admin:admin123 -XPOST 'http://localhost:18086/write?db=testdb' --data-binary 'cpu,host=host2,region=uswest value=0.72'
curl -i -u admin:admin123 -XPOST 'http://localhost:18086/write?db=testdb' --data-binary 'mem,host=host1,region=uswest value=33.1'


# 查询写入结果（带转码）
urlencode() {
    # 用于简单转码SQL语句
    local LANG=C
    local length="${#1}"
    for (( i = 0; i < length; i++ )); do
        local c="${1:i:1}"
        case $c in
            [a-zA-Z0-9.~_-]) printf "$c" ;;
            ' ') printf '%%20' ;;
            *) printf '%%%02X' "'${c}'" ;;
        esac
    done
}

echo -e "\n--- 查询 cpu ---"
curl -s -u admin:admin123 "http://localhost:18086/query?db=testdb&q=$(urlencode 'SELECT * FROM cpu')&pretty=true"
echo -e "\n--- 查询 mem ---"
curl -s -u admin:admin123 "http://localhost:18086/query?db=testdb&q=$(urlencode 'SELECT * FROM mem')&pretty=true"

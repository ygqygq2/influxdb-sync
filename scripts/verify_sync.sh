#!/bin/bash

# 通用的同步验证脚本
# 支持 1x-1x, 1x-2x, 2x-2x 等多种同步模式

set -e

# 默认值
SYNC_MODE=""
MEASUREMENTS="cpu mem"
VERBOSE=false

# 显示帮助信息
show_help() {
    cat << EOF
Usage: $0 -m <sync_mode> [options]

同步验证脚本，支持多种同步模式

Options:
  -m, --mode MODE       同步模式: 1x-1x, 1x-2x, 2x-2x
  -t, --measurements    要检查的measurements，默认: cpu mem
  -v, --verbose         详细输出
  -h, --help           显示帮助信息

Examples:
  $0 -m 1x-1x           验证1x到1x的同步
  $0 -m 1x-2x           验证1x到2x的同步
  $0 -m 2x-2x           验证2x到2x的同步
  $0 -m 1x-2x -t "cpu"  只验证cpu measurement
EOF
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--mode)
                SYNC_MODE="$2"
                shift 2
                ;;
            -t|--measurements)
                MEASUREMENTS="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                echo "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# 验证参数
validate_args() {
    if [[ -z "$SYNC_MODE" ]]; then
        echo "错误: 必须指定同步模式 (-m)"
        show_help
        exit 1
    fi

    case "$SYNC_MODE" in
        1x-1x|1x-2x|2x-2x)
            ;;
        *)
            echo "错误: 不支持的同步模式: $SYNC_MODE"
            echo "支持的模式: 1x-1x, 1x-2x, 2x-2x"
            exit 1
            ;;
    esac
}

# 查询InfluxDB 1.x数据量
query_influxdb1() {
    local container=$1
    local database=$2
    local measurement=$3
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo "  查询: docker exec $container influx -execute \"SELECT count(*) FROM $measurement\" -database $database"
    fi
    
    local result=$(docker exec "$container" influx -execute "SELECT count(*) FROM $measurement" -database "$database" 2>/dev/null)
    # 解析结果，跳过表头，取count_value列的值
    local count=$(echo "$result" | grep -E "^[0-9]" | awk '{print $2}')
    echo "${count:-0}"
}

# 查询InfluxDB 2.x数据量
query_influxdb2() {
    local container=$1
    local bucket=$2
    local measurement=$3
    local token=${4:-testtoken}
    local org=${5:-testorg}
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo "  查询: docker exec $container influx query 'from(bucket:\"$bucket\") |> range(start:2023-07-21T00:00:00Z, stop:2023-07-23T00:00:00Z) |> filter(fn:(r) => r._measurement == \"$measurement\") |> count()'"
    fi
    
    local query="from(bucket:\"$bucket\") |> range(start:2023-07-21T00:00:00Z, stop:2023-07-23T00:00:00Z) |> filter(fn:(r) => r._measurement == \"$measurement\") |> count()"
    local result=$(docker exec "$container" influx query "$query" --token "$token" --org "$org" 2>/dev/null)
    
    # 计算所有_value的总和
    local total=0
    local values=$(echo "$result" | grep -o '[0-9]\+$')
    for value in $values; do
        total=$((total + value))
    done
    echo "$total"
}

# 验证1x-1x同步
verify_1x_1x() {
    echo "验证 InfluxDB 1.x 到 1.x 同步..."
    echo "==============================================="
    
    for measurement in $MEASUREMENTS; do
        echo "检查 $measurement measurement:"
        
        src_count=$(query_influxdb1 "influxdb1-src" "testdb" "$measurement")
        dst_count=$(query_influxdb1 "influxdb1-dst" "testdb" "$measurement")
        
        echo "  源库: $src_count 条"
        echo "  目标库: $dst_count 条"
        
        # 确保数值有效
        src_count=${src_count:-0}
        dst_count=${dst_count:-0}
        diff=$((src_count - dst_count))
        if [[ $diff -eq 0 ]]; then
            echo "  ✅ 同步完成，数据一致"
        else
            echo "  ❌ 数据不一致，差异: $diff 条"
        fi
        echo
    done
}

# 验证1x-2x同步
verify_1x_2x() {
    echo "验证 InfluxDB 1.x 到 2.x 同步..."
    echo "==============================================="
    
    for measurement in $MEASUREMENTS; do
        echo "检查 $measurement measurement:"
        
        src_count=$(query_influxdb1 "influxdb1-src" "testdb" "$measurement")
        dst_count=$(query_influxdb2 "influxdb2-dst" "testdb" "$measurement")
        
        echo "  源库(1.x): $src_count 条"
        echo "  目标库(2.x): $dst_count 条"
        
        # 确保数值有效
        src_count=${src_count:-0}
        dst_count=${dst_count:-0}
        diff=$((src_count - dst_count))
        if [[ $diff -eq 0 ]]; then
            echo "  ✅ 同步完成，数据一致"
        else
            echo "  ❌ 数据不一致，差异: $diff 条"
        fi
        echo
    done
}

# 验证2x-2x同步
verify_2x_2x() {
    echo "验证 InfluxDB 2.x 到 2.x 同步..."
    echo "==============================================="
    
    for measurement in $MEASUREMENTS; do
        echo "检查 $measurement measurement:"
        
        src_count=$(query_influxdb2 "influxdb2-src" "testbucket" "$measurement")
        dst_count=$(query_influxdb2 "influxdb2-dst" "testbucket" "$measurement")
        
        echo "  源库(2.x): $src_count 条"
        echo "  目标库(2.x): $dst_count 条"
        
        # 确保数值有效
        src_count=${src_count:-0}
        dst_count=${dst_count:-0}
        diff=$((src_count - dst_count))
        if [[ $diff -eq 0 ]]; then
            echo "  ✅ 同步完成，数据一致"
        else
            echo "  ❌ 数据不一致，差异: $diff 条"
        fi
        echo
    done
}

# 检查容器是否运行
check_containers() {
    local containers=()
    
    case "$SYNC_MODE" in
        1x-1x)
            containers=("influxdb1-src" "influxdb1-dst")
            ;;
        1x-2x)
            containers=("influxdb1-src" "influxdb2-dst")
            ;;
        2x-2x)
            containers=("influxdb2-src" "influxdb2-dst")
            ;;
    esac
    
    for container in "${containers[@]}"; do
        if ! docker ps --format "table {{.Names}}" | grep -q "^$container$"; then
            echo "错误: 容器 $container 未运行"
            echo "请先启动相应的 docker-compose 环境"
            exit 1
        fi
    done
}

# 主函数
main() {
    parse_args "$@"
    validate_args
    
    echo "同步验证脚本"
    echo "模式: $SYNC_MODE"
    echo "Measurements: $MEASUREMENTS"
    echo
    
    check_containers
    
    case "$SYNC_MODE" in
        1x-1x)
            verify_1x_1x
            ;;
        1x-2x)
            verify_1x_2x
            ;;
        2x-2x)
            verify_2x_2x
            ;;
    esac
}

# 运行主函数
main "$@"

# influxdb-sync

Go 语言编写的 InfluxDB 数据迁移/同步工具

## 功能

- 支持配置源和目标 InfluxDB 地址、认证信息
- InfluxDB 1.x/2.x 支持
- 从源同步数据到目标
- 支持断点续传
- 不依赖物理目录

## 使用方法

1. 编辑 `config.yaml` 配置源和目标信息
2. 运行 `go run main.go`

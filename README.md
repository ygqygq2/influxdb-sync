# influxdb-sync

[![Test](https://github.com/ygqygq2/influxdb-sync/actions/workflows/test.yml/badge.svg)](https://github.com/ygqygq2/influxdb-sync/actions/workflows/test.yml)
[![Code Quality](https://github.com/ygqygq2/influxdb-sync/actions/workflows/quality.yml/badge.svg)](https://github.com/ygqygq2/influxdb-sync/actions/workflows/quality.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygqygq2/influxdb-sync)](https://goreportcard.com/report/github.com/ygqygq2/influxdb-sync)
[![License](https://img.shields.io/github/license/ygqygq2/influxdb-sync)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ygqygq2/influxdb-sync)](https://github.com/ygqygq2/influxdb-sync/releases)

Go 语言编写的 InfluxDB 数据迁移/同步工具，支持断点续传和配置化

## 功能

- 支持配置源和目标 InfluxDB 地址、认证信息
- InfluxDB 1.x/2.x 支持
- 从源同步数据到目标
- 支持断点续传
- 不依赖物理目录

## 使用方法

1. 编辑 `config.yaml` 配置源和目标信息
2. 运行 `go run main.go`

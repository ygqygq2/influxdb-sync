# InfluxDB 同步工具

[![Test](https://github.com/ygqygq2/influxdb-sync/actions/workflows/test.yml/badge.svg)](https://github.com/ygqygq2/influxdb-sync/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygqygq2/influxdb-sync)](https://goreportcard.com/report/github.com/ygqygq2/influxdb-sync)
[![License](https://img.shields.io/github/license/ygqygq2/influxdb-sync)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ygqygq2/influxdb-sync)](https://github.com/ygqygq2/influxdb-sync/releases)

一个高性能、功能完整的 InfluxDB 数据迁移和同步工具，用 Go 语言编写，支持多种同步模式和断点续传。

## ✨ 核心特性

### 🔄 多版本同步支持

- **1x → 1x**: InfluxDB 1.x 到 1.x 的数据同步
- **1x → 2x**: InfluxDB 1.x 到 2.x 的跨版本迁移
- **2x → 2x**: InfluxDB 2.x 到 2.x 的数据同步

### 🚀 高性能设计

- **并发处理**: 支持多 measurement 并行同步
- **批量传输**: 可配置的批次大小优化网络效率
- **断点续传**: 基于时间戳的增量同步，支持中断恢复
- **内存优化**: 流式处理，避免大数据集内存溢出

### ⚙️ 灵活配置

- **YAML 配置**: 人性化的配置文件格式
- **认证支持**: 支持用户名/密码和 Token 认证
- **过滤功能**: 支持数据库包含/排除规则
- **自定义命名**: 支持目标数据库前缀/后缀

### 🔒 可靠性保证

- **数据校验**: 传输前后的数据完整性检查
- **错误重试**: 可配置的重试次数和间隔
- **详细日志**: 分级日志输出，便于问题排查
- **进度跟踪**: 实时显示同步进度和性能指标

## 📦 快速开始

### 下载安装

从 [Releases](https://github.com/ygqygq2/influxdb-sync/releases) 页面下载适合你系统的预编译二进制文件：

```bash
# Linux
wget https://github.com/ygqygq2/influxdb-sync/releases/latest/download/influxdb-sync_Linux_x86_64.zip
unzip influxdb-sync_Linux_x86_64.zip

# Windows
# 下载 influxdb-sync_Windows_x86_64.zip 并解压
````

### 配置文件

项目提供了多种同步模式的示例配置文件：

- `config.yaml` - InfluxDB 1.x → 1.x 同步
- `config_1x2x.yaml` - InfluxDB 1.x → 2.x 同步
- `config_2x2x.yaml` - InfluxDB 2.x → 2.x 同步

根据你的需求复制对应的配置文件并修改相关参数。

**同步模式自动判断**：程序会根据配置文件中的 `source.type` 和 `target.type` 字段自动选择同步模式。

### 运行同步

```bash
# 运行同步 (模式由配置文件中的 type 字段自动判断)
./influxdb-sync config.yaml

# 不同模式使用对应的配置文件
./influxdb-sync config_1x2x.yaml  # 1x → 2x 同步
./influxdb-sync config_2x2x.yaml  # 2x → 2x 同步
```

## 🛠️ 开发和构建

### 本地开发

```bash
# 克隆代码
git clone https://github.com/ygqygq2/influxdb-sync.git
cd influxdb-sync

# 安装依赖 (需要先安装 Task)
task deps

# 运行测试
task test

# 本地构建
task build
```

### 自动化任务

项目使用 [Task](https://taskfile.dev/) 进行自动化：

```bash
task test          # 运行测试
task test-coverage # 测试覆盖率
task build         # 构建二进制
task release       # 发布构建 (多平台)
task clean         # 清理构建产物

# Git hooks 管理
task install-hooks # 安装 pre-commit hooks
task test-hooks    # 测试 hooks 状态
```

### Git Hooks 自动化

项目配置了 pre-commit hooks，每次提交时自动：

- 🔧 格式化 Go 代码
- 🔍 运行静态分析
- 📦 检查依赖状态

无需手动记住运行 `task fmt`！详见 [Git Hooks 说明](docs/GIT_HOOKS.md)。

## 📋 配置参考

支持的同步模式和对应的配置参数：

| 同步模式 | 说明               | 主要配置项                                                       |
| -------- | ------------------ | ---------------------------------------------------------------- |
| `1x1x`   | InfluxDB 1.x → 1.x | `source.type: 1`, `target.type: 1`, 用户名密码认证               |
| `1x2x`   | InfluxDB 1.x → 2.x | `source.type: 1`, `target.type: 2`, 源端用户名密码，目标端 Token |
| `2x2x`   | InfluxDB 2.x → 2.x | `source.type: 2`, `target.type: 2`, Token 认证，组织和桶配置     |

详细配置说明请参考项目中的示例配置文件。

## 📈 性能特点

- **压缩优化**: 使用 UPX 压缩，二进制文件减少 57% 大小
- **并发处理**: 支持配置并发数量，充分利用系统资源
- **批量传输**: 智能批次大小，平衡内存使用和网络效率
- **断点续传**: 意外中断后可从上次位置继续，避免重复传输

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！请确保：

1. 代码通过所有测试：`task test`
2. 代码格式符合规范：`task fmt-check`
3. 测试覆盖率不低于 50%：`task test-coverage`

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

# InfluxDB 同步工具架构文档

## 项目概述

本项目是一个用 Go 语言编写的 InfluxDB 数据迁移/同步工具，支持以下同步模式：

- **1x1x**: InfluxDB 1.x 到 1.x
- **1x2x**: InfluxDB 1.x 到 2.x
- **2x2x**: InfluxDB 2.x 到 2.x
- **1x3x**: InfluxDB 1.x 到 3.x 🆕
- **2x3x**: InfluxDB 2.x 到 3.x 🆕
- **3x3x**: InfluxDB 3.x 到 3.x 🆕

## 架构设计

### 核心模块结构

```
influxdb-sync/
├── main.go                     # 程序入口点，处理命令行参数和模式分发
├── cmd/                        # 命令层，处理不同同步模式的调度
│   ├── sync.go                 # 同步模式分发和配置转换
│   └── sync_test.go           # 同步功能测试
├── internal/                   # 内部包，核心业务逻辑
│   ├── common/                 # 通用组件和接口定义
│   │   ├── syncer.go          # 核心同步引擎
│   │   ├── types.go           # 通用数据类型和接口
│   │   └── *_test.go         # 完整的功能测试
│   ├── config/                 # 配置管理
│   │   ├── config.go          # YAML配置文件解析
│   │   └── config_test.go     # 配置测试（100%覆盖率）
│   ├── utils/                  # 通用工具函数
│   │   ├── progress.go        # 进度条显示工具
│   │   ├── strings.go         # 字符串处理工具
│   │   └── *_test.go         # 工具函数测试（100%覆盖率）
│   ├── influxdb1/             # InfluxDB 1.x 特定实现
│   │   ├── adapter.go         # 1.x 数据适配器
│   │   ├── client.go          # 1.x 客户端封装
│   │   ├── sync.go            # 1.x 通用同步逻辑
│   │   ├── sync_1x2x.go       # 1.x 到 2.x 专用同步
│   │   ├── filter.go          # 1.x 特定的过滤和转义
│   │   ├── types.go           # 1.x 类型定义和别名
│   │   └── *_test.go         # 完整测试套件
│   ├── influxdb2/             # InfluxDB 2.x 特定实现
│   │   ├── adapter.go         # 2.x 数据适配器
│   │   ├── sync_2x2x.go       # 2.x 到 2.x 同步逻辑
│   │   └── *_test.go         # 完整测试套件
│   ├── influxdb3/             # InfluxDB 3.x 特定实现 🆕
│   │   ├── types.go           # 3.x 配置类型定义
│   │   ├── client.go          # 3.x 多模式客户端
│   │   ├── adapter.go         # 3.x 数据适配器
│   │   ├── sync_1x3x.go       # 1.x 到 3.x 同步逻辑
│   │   ├── sync_2x3x.go       # 2.x 到 3.x 同步逻辑
│   │   ├── sync_3x3x.go       # 3.x 到 3.x 同步逻辑
│   │   └── *_test.go         # 完整测试套件
│   └── logx/                   # 日志组件
│       ├── logx.go            # 轻量级日志实现
│       └── logx_test.go       # 日志测试（85%覆盖率）
└── docs/                       # 文档目录
    └── ARCHITECTURE.md        # 本架构文档
```

### 设计原则

#### 1. 清晰的关注点分离

- **common/**: 定义通用接口和核心同步逻辑，版本无关
- **influxdb1/**: InfluxDB 1.x 特定的实现细节
- **influxdb2/**: InfluxDB 2.x 特定的实现细节
- **utils/**: 可复用的工具函数，消除代码重复

#### 2. 接口驱动设计

```go
// 核心接口定义在 common/types.go
type DataSource interface {
    GetDatabases() ([]string, error)
    GetMeasurements(database string) ([]string, error)
    QueryData(database, measurement string, startTime int64, batchSize int) ([]DataPoint, error)
    Connect() error
    Close() error
}

type DataTarget interface {
    WritePoints(database string, points []DataPoint) error
    Connect() error
    Close() error
}
```

#### 3. 配置系统设计

- **统一配置结构**: `common.SyncConfig` 作为所有同步模式的基础
- **类型别名兼容**: `influxdb1.SyncConfig = common.SyncConfig` 保持向后兼容
- **灵活的配置转换**: 支持不同版本间的配置自动转换

#### 4. 并发和性能优化

- **并发处理**: 支持多个 measurement 并行同步
- **批量操作**: 可配置的批次大小优化网络传输
- **断点续传**: 基于时间戳的增量同步
- **进度显示**: 实时同步进度反馈

### 核心流程

#### 1. 程序启动流程

```
main.go → 解析命令行参数 → cmd/sync.go → 模式分发
                                      ↓
            ┌─────────────────────────────────────────────┐
            │              模式分发                          │
            └─────────────────────────────────────────────┘
                    │        │        │        │        │        │
                  1x1x     1x2x     2x2x     1x3x     2x3x     3x3x
                    │        │        │        │        │        │
         influxdb1/sync.go sync_1x2x.go sync_2x2x.go sync_1x3x.go sync_2x3x.go sync_3x3x.go
```

#### 2. 同步执行流程

```
配置加载 → 客户端连接 → 数据库发现 → Measurement发现
                                            ↓
    ┌─── 并发处理 ───┐    ┌─── 并发处理 ───┐    ┌─── 并发处理 ───┐
    │  Measurement1  │    │  Measurement2  │    │  MeasurementN  │
    │      ↓         │    │      ↓         │    │      ↓         │
    │   数据查询      │    │   数据查询      │    │   数据查询      │
    │      ↓         │    │      ↓         │    │      ↓         │
    │   数据写入      │    │   数据写入      │    │   数据写入      │
    │      ↓         │    │      ↓         │    │      ↓         │
    │   进度更新      │    │   进度更新      │    │   进度更新      │
    └───────────────┘    └───────────────┘    └───────────────┘
                                ↓
                           同步完成统计
```

#### 3. 错误处理和重试机制

- **连接失败**: 自动重试机制，支持连接超时配置
- **数据传输失败**: 批量操作失败时的部分重试
- **断点续传**: 基于 `resume.state` 文件的状态恢复

### 关键特性

#### 1. 多版本支持

- **InfluxDB 1.x**: 使用 SQL-like 查询语法和 REST API
- **InfluxDB 2.x**: 使用 Flux 查询语言和 v2 API
- **跨版本同步**: 1.x 到 2.x 的数据格式转换

#### 2. 性能优化

- **并发控制**: 可配置的并发 goroutine 数量
- **批量处理**: 可调节的批次大小
- **内存管理**: 流式处理避免大数据集内存溢出
- **网络优化**: 连接复用和超时控制

#### 3. 可观测性

- **详细日志**: 分级日志输出，支持 Debug、Info、Warn、Error 级别
- **进度跟踪**: 实时显示同步进度和预估完成时间
- **性能指标**: 查询耗时、写入耗时、数据点数量统计

#### 4. 可靠性保证

- **数据校验**: 写入前的数据格式验证
- **事务安全**: 批量操作的原子性保证
- **状态持久化**: 同步状态的文件持久化
- **错误恢复**: 异常情况下的自动恢复机制

### 测试策略

#### 1. 测试覆盖率

- **整体覆盖率**: 53.9%（超过 50% 要求）
- **关键模块**:
  - `config`: 100% 覆盖率
  - `utils`: 100% 覆盖率
  - `common`: 80.7% 覆盖率
  - `cmd`: 72.5% 覆盖率

#### 2. 测试类型

- **单元测试**: 针对每个函数和方法的独立测试
- **集成测试**: 跨模块的功能测试
- **边界测试**: 异常情况和边界条件测试
- **性能测试**: 大数据量和高并发场景测试

#### 3. 测试数据

- **Mock 数据**: 使用内存中的模拟数据源进行测试
- **真实连接**: 针对真实 InfluxDB 实例的集成测试
- **边界情况**: 空数据、大数据、异常数据的处理测试

### 扩展性设计

#### 1. 新版本支持

添加新 InfluxDB 版本支持的步骤：

1. 在 `internal/` 下创建新的版本目录（如 `influxdb3/`）
2. 实现 `DataSource` 和 `DataTarget` 接口
3. 在 `cmd/sync.go` 中添加新的同步模式
4. 添加相应的测试用例

#### 2. 新功能添加

- **过滤器扩展**: 在各版本模块中扩展特定的过滤逻辑
- **转换器添加**: 新增数据格式转换器
- **监控集成**: 添加 Prometheus 监控指标
- **配置扩展**: 在 `common.SyncConfig` 中添加新的配置项

### 最佳实践

#### 1. 代码组织

- **包边界清晰**: 每个包有明确的职责边界
- **接口优先**: 使用接口定义模块间的交互
- **错误处理**: 统一的错误处理和日志记录模式
- **资源管理**: 正确的连接和资源释放

#### 2. 性能调优

- **批次大小**: 根据网络和服务器性能调整 `BatchSize`
- **并发数量**: 根据系统资源调整 `ConcurrentQueries`
- **超时配置**: 合理设置连接和查询超时时间
- **内存使用**: 避免大量数据的内存缓存

#### 3. 运维部署

- **配置管理**: 使用 YAML 配置文件管理不同环境
- **日志级别**: 生产环境使用 Info 级别，调试时使用 Debug 级别
- **监控告警**: 监控同步进度和错误率
- **备份策略**: 重要数据同步前的备份确认

## InfluxDB 3.x 支持特性 🆕

### 兼容模式设计

InfluxDB 3.x 支持三种兼容模式，本工具充分利用这些特性：

1. **v1 兼容模式**: 支持 InfluxQL 查询和 Line Protocol 写入
2. **v2 兼容模式**: 支持 Flux 查询和 v2 API 访问
3. **原生模式**: 支持 SQL 查询和原生 v3 API

### 多模式适配器

```go
// influxdb3/client.go - 多模式客户端
type Client3x struct {
    v1Client    client.Client      // v1 兼容客户端
    v2Client    influxdb2.Client   // v2 兼容客户端
    compatMode  string             // "v1", "v2", "native"
}

// influxdb3/adapter.go - 智能适配器
type DataSource3x struct {
    client *Client3x
    config interface{} // V1CompatConfig, V2CompatConfig, NativeConfig
}
```

### 配置兼容性

- **灵活配置**: 支持混合兼容模式（如 v1 源 → v2 目标）
- **自动适配**: 根据配置自动选择合适的 API 接口
- **向前兼容**: 现有 1.x/2.x 配置可直接用于 3.x 目标

## 总结

这个架构设计实现了：

1. **高内聚低耦合**: 清晰的模块边界和职责分离
2. **高度可扩展**: 基于接口的设计便于添加新功能
3. **性能优化**: 并发处理和批量操作提升同步效率
4. **可靠性保证**: 完善的错误处理和状态恢复机制
5. **良好的可测试性**: 高测试覆盖率和多层次测试策略

通过这个架构，项目能够稳定、高效地处理各种 InfluxDB 数据同步场景，同时保持良好的可维护性和扩展性。

package influxdb3

import "github.com/ygqygq2/influxdb-sync/internal/common"

// 为了保持一致性，为 common.SyncConfig 创建一个别名
type SyncConfig = common.SyncConfig

// InfluxDB 3.x 特定配置
type Config3x struct {
	Database      string `yaml:"database"`        // 3.x database name
	Namespace     string `yaml:"namespace"`       // 3.x namespace
	UseV1Compat   bool   `yaml:"use_v1_compat"`   // 使用 v1 兼容模式
	UseV2Compat   bool   `yaml:"use_v2_compat"`   // 使用 v2 兼容模式
	UseSQL        bool   `yaml:"use_sql"`         // 使用 SQL 查询
	HTTPTimeout   int    `yaml:"http_timeout"`    // HTTP 超时时间(秒)
	FlightSQLPort int    `yaml:"flight_sql_port"` // Flight SQL 端口
}

// V1 兼容模式配置
type V1CompatConfig struct {
	Addr     string
	User     string
	Pass     string
	Database string
}

// V2 兼容模式配置
type V2CompatConfig struct {
	URL      string
	Token    string
	Org      string
	Bucket   string
	Database string
}

// 原生 3.x 配置
type NativeConfig struct {
	URL       string
	Token     string
	Database  string
	Namespace string
	UseSQL    bool
}

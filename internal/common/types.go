package common

import "time"

// 通用同步配置
type SyncConfig struct {
	SourceAddr      string
	SourceUser      string
	SourcePass      string
	SourceDB        string
	SourceDBExclude []string
	SourceToken     string
	SourceOrg       string
	SourceBucket    string
	TargetAddr      string
	TargetUser      string
	TargetPass      string
	TargetDB        string
	TargetDBPrefix  string
	TargetDBSuffix  string
	TargetToken     string
	TargetOrg       string
	TargetBucket    string
	BatchSize       int
	Start           string
	End             string
	ResumeFile      string
	Parallel        int
	RetryCount      int
	RetryInterval   int
	RateLimit       int
	LogLevel        string
}

// 数据点结构
type DataPoint struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	Time        time.Time
}

// 同步结果
type SyncResult struct {
	Measurement string
	Error       error
}

// 数据源接口
type DataSource interface {
	Connect() error
	Close() error
	GetDatabases() ([]string, error)
	GetMeasurements(db string) ([]string, error)
	GetTagKeys(db, measurement string) (map[string]bool, error)
	QueryData(db, measurement string, startTime int64, batchSize int) ([]DataPoint, int64, error)
}

// 数据目标接口
type DataTarget interface {
	Connect() error
	Close() error
	WritePoints(db string, points []DataPoint) error
}

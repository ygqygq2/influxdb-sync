package influxdb1

type SyncConfig struct {
	SourceAddr      string
	SourceUser      string
	SourcePass      string
	SourceDB        string
	SourceDBExclude []string
	TargetAddr      string
	TargetUser      string
	TargetPass      string
	TargetDB        string
	TargetDBPrefix  string
	TargetDBSuffix  string
	BatchSize       int
	Start           string // 起始时间
	End             string // 结束时间
	ResumeFile      string
	Parallel        int    // 并发同步表数，默认4
	RetryCount      int    // 写入失败重试次数，默认3次
	RetryInterval   int    // 重试间隔毫秒数，默认500ms
	RateLimit       int    // 每批写入后限流毫秒数，默认50ms，0表示不限流
	LogLevel        string // 日志级别: debug, info, warn, error
}

type ExtraConfig struct {
	SourceDBExclude []string
	TargetDBPrefix  string
	TargetDBSuffix  string
}

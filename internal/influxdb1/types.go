package influxdb1

type SyncConfig struct {
   SourceAddr string
   SourceUser string
   SourcePass string
   SourceDB   string
   SourceDBExclude []string
   TargetAddr string
   TargetUser string
   TargetPass string
   TargetDB   string
   TargetDBPrefix string
   TargetDBSuffix string
   BatchSize  int
   Start      string // 起始时间
   End        string // 结束时间
   ResumeFile string
   Parallel   int    // 并发同步表数，默认4
}

type ExtraConfig struct {
	SourceDBExclude []string
	TargetDBPrefix  string
	TargetDBSuffix  string
}

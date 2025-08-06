package influxdb1

import "github.com/ygqygq2/influxdb-sync/internal/common"

// 为了向后兼容，为 common.SyncConfig 创建一个别名
type SyncConfig = common.SyncConfig

type ExtraConfig struct {
	SourceDBExclude []string
	TargetDBPrefix  string
	TargetDBSuffix  string
}

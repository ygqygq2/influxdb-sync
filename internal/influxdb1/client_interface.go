package influxdb1

import client "github.com/influxdata/influxdb1-client/v2"

// InfluxClient 接口，用于 mock 测试
type InfluxClient interface {
	Query(q client.Query) (*client.Response, error)
	Write(bp client.BatchPoints) error
	Close() error
}

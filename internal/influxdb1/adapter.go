package influxdb1

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

// InfluxDB 1.x 数据源实现
type DataSource struct {
	cli    client.Client
	config DataSourceConfig
}

type DataSourceConfig struct {
	Addr string
	User string
	Pass string
}

// InfluxDB 1.x 数据目标实现
type DataTarget struct {
	cli    client.Client
	config DataTargetConfig
}

type DataTargetConfig struct {
	Addr string
	User string
	Pass string
}

// 创建数据源
func NewDataSource(config DataSourceConfig) *DataSource {
	return &DataSource{config: config}
}

// 创建数据目标
func NewDataTarget(config DataTargetConfig) *DataTarget {
	return &DataTarget{config: config}
}

// 数据源接口实现
func (ds *DataSource) Connect() error {
	// 设置30秒超时，避免长时间阻塞
	cli, err := NewClient(ds.config.Addr, ds.config.User, ds.config.Pass, 30*time.Second)
	if err != nil {
		return err
	}
	ds.cli = cli.cli
	return nil
}

func (ds *DataSource) Close() error {
	if ds.cli != nil {
		return ds.cli.Close()
	}
	return nil
}

func (ds *DataSource) GetDatabases() ([]string, error) {
	dbRes, err := ds.cli.Query(client.NewQuery("SHOW DATABASES", "", ""))
	if err != nil {
		return nil, err
	}
	if dbRes.Error() != nil {
		return nil, dbRes.Error()
	}

	var dbs []string
	for _, result := range dbRes.Results {
		for _, series := range result.Series {
			for _, v := range series.Values {
				if len(v) > 0 {
					if name, ok := v[0].(string); ok {
						dbs = append(dbs, name)
					}
				}
			}
		}
	}
	return dbs, nil
}

func (ds *DataSource) GetMeasurements(db string) ([]string, error) {
	showRes, err := ds.cli.Query(client.NewQuery("SHOW MEASUREMENTS", db, ""))
	if err != nil {
		return nil, err
	}
	if showRes.Error() != nil {
		return nil, showRes.Error()
	}

	var measurements []string
	for _, result := range showRes.Results {
		for _, series := range result.Series {
			for _, v := range series.Values {
				if len(v) > 0 {
					if name, ok := v[0].(string); ok {
						measurements = append(measurements, name)
					}
				}
			}
		}
	}
	return measurements, nil
}

func (ds *DataSource) GetTagKeys(db, measurement string) (map[string]bool, error) {
	tagKeys := make(map[string]bool)
	q := fmt.Sprintf("SHOW TAG KEYS FROM %s", escapeMeasurement(measurement))
	res, err := ds.cli.Query(client.NewQuery(q, db, ""))
	if err != nil {
		return tagKeys, err
	}
	if res.Error() != nil {
		return tagKeys, res.Error()
	}

	for _, result := range res.Results {
		for _, series := range result.Series {
			for _, row := range series.Values {
				if len(row) > 0 {
					if tagKey, ok := row[0].(string); ok {
						tagKeys[tagKey] = true
					}
				}
			}
		}
	}
	return tagKeys, nil
}

func (ds *DataSource) QueryData(db, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	em := escapeMeasurement(measurement)
	var q string
	if startTime == 0 {
		q = fmt.Sprintf("SELECT * FROM %s ORDER BY time ASC LIMIT %d", em, batchSize)
	} else {
		q = fmt.Sprintf("SELECT * FROM %s WHERE time > %d ORDER BY time ASC LIMIT %d", em, startTime, batchSize)
	}

	logx.Debug(fmt.Sprintf("执行查询: %s", q))
	queryStart := time.Now()
	res, err := ds.cli.Query(client.NewQuery(q, db, "ns"))
	logx.Debug(fmt.Sprintf("查询耗时: %v", time.Since(queryStart)))
	if err != nil {
		return nil, 0, err
	}
	if res.Error() != nil {
		return nil, 0, res.Error()
	}

	var points []common.DataPoint
	var maxTime int64 = startTime

	// 获取标签字段
	tagKeys, err := ds.GetTagKeys(db, measurement)
	if err != nil {
		logx.Warn("获取标签字段失败，使用默认:", err)
		tagKeys = map[string]bool{"host": true, "region": true}
	}

	for _, result := range res.Results {
		for _, series := range result.Series {
			colIdx := map[string]int{}
			for i, col := range series.Columns {
				colIdx[col] = i
			}

			for _, row := range series.Values {
				tags := map[string]string{}
				fields := map[string]interface{}{}
				var t time.Time
				var tUnix int64 = 0
				skip := false

				for col, idx := range colIdx {
					switch col {
					case "time":
						switch v := row[idx].(type) {
						case string:
							t, _ = time.Parse(time.RFC3339Nano, v)
							tUnix = t.UnixNano()
						case time.Time:
							t = v
							tUnix = t.UnixNano()
						case int64:
							t = time.Unix(0, v)
							tUnix = v
						case float64:
							t = time.Unix(0, int64(v))
							tUnix = int64(v)
						case json.Number:
							if ns, err := v.Int64(); err == nil {
								t = time.Unix(0, ns)
								tUnix = ns
							}
						default:
							logx.Warn(fmt.Sprintf("未知time类型: %T, value: %v, measurement: %s, db: %s, 跳过该点", v, v, series.Name, db))
							skip = true
						}
					default:
						// 动态判断是标签还是字段
						if tagKeys[col] {
							// 这是标签字段
							if s, ok := row[idx].(string); ok {
								tags[col] = s
							}
						} else {
							// 这是普通字段
							if val := row[idx]; val != nil {
								if sv, ok := val.(string); !ok || sv != "" {
									fields[col] = val
								}
							}
						}
					}
				}

				if skip {
					continue
				}

				if tUnix > maxTime {
					maxTime = tUnix
				}

				points = append(points, common.DataPoint{
					Measurement: series.Name,
					Tags:        tags,
					Fields:      fields,
					Time:        t,
				})
			}
		}
	}

	return points, maxTime, nil
}

// 数据目标接口实现
func (dt *DataTarget) Connect() error {
	// 设置30秒超时，避免长时间阻塞
	cli, err := NewClient(dt.config.Addr, dt.config.User, dt.config.Pass, 30*time.Second)
	if err != nil {
		return err
	}
	dt.cli = cli.cli
	return nil
}

func (dt *DataTarget) Close() error {
	if dt.cli != nil {
		return dt.cli.Close()
	}
	return nil
}

func (dt *DataTarget) WritePoints(db string, points []common.DataPoint) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: db, Precision: "ns"})
	if err != nil {
		return err
	}

	for _, point := range points {
		pt, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
		if err != nil {
			continue // 跳过有问题的点
		}
		bp.AddPoint(pt)
	}

	return dt.cli.Write(bp)
}

// measurement 名称转义，双引号包裹并转义内部双引号
func escapeMeasurement(m string) string {
	return "\"" + strings.ReplaceAll(m, "\"", "\\\"") + "\""
}

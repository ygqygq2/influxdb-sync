package influxdb2

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// 通用适配器，可同时作为数据源和数据目标
type Adapter struct {
	URL    string
	Token  string
	Org    string
	Bucket string
	client influxdb2.Client
}

// 数据源接口实现
func (a *Adapter) Connect() error {
	client := influxdb2.NewClient(a.URL, a.Token)
	a.client = client
	return nil
}

func (a *Adapter) Close() error {
	if a.client != nil {
		a.client.Close()
	}
	return nil
}

func (a *Adapter) GetDatabases() ([]string, error) {
	// InfluxDB 2.x 使用 bucket，这里返回指定的 bucket
	if a.Bucket != "" {
		return []string{a.Bucket}, nil
	}

	// 如果没有指定 bucket，查询所有 buckets
	bucketsAPI := a.client.BucketsAPI()
	buckets, err := bucketsAPI.GetBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range *buckets {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}

func (a *Adapter) GetMeasurements(bucket string) ([]string, error) {
	queryAPI := a.client.QueryAPI(a.Org)

	// 查询所有 measurement
	query := fmt.Sprintf(`
		import "influxdata/influxdb/schema"
		schema.measurements(bucket: "%s")
	`, bucket)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var measurements []string
	for result.Next() {
		if result.Record().ValueByKey("_value") != nil {
			if measurement, ok := result.Record().ValueByKey("_value").(string); ok {
				measurements = append(measurements, measurement)
			}
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return measurements, nil
}

func (a *Adapter) GetTagKeys(bucket, measurement string) (map[string]bool, error) {
	queryAPI := a.client.QueryAPI(a.Org)

	// 查询指定 measurement 的所有 tag keys
	query := fmt.Sprintf(`
		import "influxdata/influxdb/schema"
		schema.tagKeys(bucket: "%s", predicate: (r) => r._measurement == "%s")
	`, bucket, measurement)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	tagKeys := make(map[string]bool)
	for result.Next() {
		if result.Record().ValueByKey("_value") != nil {
			if tagKey, ok := result.Record().ValueByKey("_value").(string); ok {
				tagKeys[tagKey] = true
			}
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return tagKeys, nil
}

func (a *Adapter) QueryData(bucket, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	queryAPI := a.client.QueryAPI(a.Org)

	// 构建查询
	var startFilter string
	if startTime > 0 {
		startTimeRFC := time.Unix(0, startTime).UTC().Format(time.RFC3339Nano)
		startFilter = fmt.Sprintf(`|> filter(fn: (r) => r._time > time(v: "%s"))`, startTimeRFC)
	} else {
		startFilter = ""
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -100y)
		|> filter(fn: (r) => r._measurement == "%s")
		%s
		|> sort(columns: ["_time"])
		|> limit(n: %d)
	`, bucket, measurement, startFilter, batchSize)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, 0, err
	}

	var points []common.DataPoint
	var maxTime int64 = startTime
	pointsMap := make(map[string]*common.DataPoint) // 按时间戳分组

	for result.Next() {
		record := result.Record()
		timeKey := record.Time().Format(time.RFC3339Nano)

		// 获取或创建数据点
		if _, exists := pointsMap[timeKey]; !exists {
			pointsMap[timeKey] = &common.DataPoint{
				Measurement: measurement,
				Tags:        make(map[string]string),
				Fields:      make(map[string]interface{}),
				Time:        record.Time(),
			}
		}

		point := pointsMap[timeKey]

		// 复制所有 tags
		for key, value := range record.Values() {
			if key[0] != '_' && key != "result" && key != "table" {
				// 这是 tag
				if strVal, ok := value.(string); ok {
					point.Tags[key] = strVal
				}
			}
		}

		// 设置 field
		if record.Field() != "" && record.Value() != nil {
			point.Fields[record.Field()] = record.Value()
		}

		// 更新最大时间
		if record.Time().UnixNano() > maxTime {
			maxTime = record.Time().UnixNano()
		}
	}

	if result.Err() != nil {
		return nil, 0, result.Err()
	}

	// 转换为切片
	for _, point := range pointsMap {
		points = append(points, *point)
	}

	return points, maxTime, nil
}

// 数据目标接口实现
func (a *Adapter) WritePoints(bucket string, points []common.DataPoint) error {
	writeAPI := a.client.WriteAPIBlocking(a.Org, bucket)

	for _, point := range points {
		p := influxdb2.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			return err
		}
	}

	return nil
}

package influxdb3

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

// DataSource3x InfluxDB 3.x 数据源实现
type DataSource3x struct {
	client *Client3x
	config interface{} // V1CompatConfig, V2CompatConfig, 或 NativeConfig
}

// DataTarget3x InfluxDB 3.x 数据目标实现
type DataTarget3x struct {
	client *Client3x
	config interface{} // V1CompatConfig, V2CompatConfig, 或 NativeConfig
}

// NewDataSource3x 创建 3.x 数据源（通用构造函数）
func NewDataSource3x(config interface{}) (*DataSource3x, error) {
	switch cfg := config.(type) {
	case V1CompatConfig:
		return NewV1CompatDataSource(cfg), nil
	case V2CompatConfig:
		return NewV2CompatDataSource(cfg), nil
	case NativeConfig:
		return NewNativeDataSource(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported config type: %T", config)
	}
}

// NewDataTarget3x 创建 3.x 数据目标（通用构造函数）
func NewDataTarget3x(config interface{}) (*DataTarget3x, error) {
	switch cfg := config.(type) {
	case V1CompatConfig:
		return NewV1CompatDataTarget(cfg), nil
	case V2CompatConfig:
		return NewV2CompatDataTarget(cfg), nil
	case NativeConfig:
		return NewNativeDataTarget(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported config type: %T", config)
	}
}

// NewV1CompatDataSource 创建 v1 兼容模式数据源
func NewV1CompatDataSource(config V1CompatConfig) *DataSource3x {
	return &DataSource3x{
		config: config,
	}
}

// NewV2CompatDataSource 创建 v2 兼容模式数据源
func NewV2CompatDataSource(config V2CompatConfig) *DataSource3x {
	return &DataSource3x{
		config: config,
	}
}

// NewNativeDataSource 创建原生 3.x 数据源
func NewNativeDataSource(config NativeConfig) *DataSource3x {
	return &DataSource3x{
		config: config,
	}
}

// NewV1CompatDataTarget 创建 v1 兼容模式数据目标
func NewV1CompatDataTarget(config V1CompatConfig) *DataTarget3x {
	return &DataTarget3x{
		config: config,
	}
}

// NewV2CompatDataTarget 创建 v2 兼容模式数据目标
func NewV2CompatDataTarget(config V2CompatConfig) *DataTarget3x {
	return &DataTarget3x{
		config: config,
	}
}

// NewNativeDataTarget 创建原生 3.x 数据目标
func NewNativeDataTarget(config NativeConfig) *DataTarget3x {
	return &DataTarget3x{
		config: config,
	}
}

// DataSource 接口实现
func (ds *DataSource3x) Connect() error {
	var err error

	switch config := ds.config.(type) {
	case V1CompatConfig:
		ds.client, err = NewV1CompatClient(config)
	case V2CompatConfig:
		ds.client, err = NewV2CompatClient(config)
	case NativeConfig:
		ds.client, err = NewClient3x(config)
	default:
		return fmt.Errorf("unsupported config type: %T", config)
	}

	if err != nil {
		return err
	}

	// 测试连接
	return ds.client.Ping(30 * time.Second)
}

func (ds *DataSource3x) Close() error {
	if ds.client != nil {
		return ds.client.Close()
	}
	return nil
}

func (ds *DataSource3x) GetDatabases() ([]string, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("client not connected")
	}
	return ds.client.GetDatabases()
}

func (ds *DataSource3x) GetMeasurements(database string) ([]string, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("client not connected")
	}

	switch ds.client.compatMode {
	case "v1":
		q := client.NewQuery("SHOW MEASUREMENTS", database, "")
		resp, err := ds.client.QueryInfluxQL(q.Command, database)
		if err != nil {
			return nil, err
		}
		if resp.Error() != nil {
			return nil, resp.Error()
		}

		var measurements []string
		if len(resp.Results) > 0 && len(resp.Results[0].Series) > 0 {
			for _, row := range resp.Results[0].Series[0].Values {
				if len(row) > 0 {
					if measurement, ok := row[0].(string); ok {
						measurements = append(measurements, measurement)
					}
				}
			}
		}
		return measurements, nil

	case "v2":
		query := fmt.Sprintf(`
			import "influxdata/influxdb/schema"
			schema.measurements(bucket: "%s")
		`, database)

		// 从 config 获取 org
		org := ""
		if v2cfg, ok := ds.config.(V2CompatConfig); ok {
			org = v2cfg.Org
		}
		
		result, err := ds.client.QueryFlux(query, org)
		if err != nil {
			return nil, err
		}

		var measurements []string
		for result.Next() {
			record := result.Record()
			if measurementValue := record.ValueByKey("_value"); measurementValue != nil {
				if measurement, ok := measurementValue.(string); ok {
					measurements = append(measurements, measurement)
				}
			}
		}
		return measurements, nil

	case "native":
		// 使用 SQL 查询 measurements
		query := fmt.Sprintf("SHOW MEASUREMENTS FROM \"%s\"", database)
		data, err := ds.client.QuerySQL(query)
		if err != nil {
			return nil, err
		}

		// 解析 JSON 响应
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		// 提取 measurements (这个实现可能需要根据实际的 3.x API 响应格式调整)
		var measurements []string
		if data, ok := result["data"].([]interface{}); ok {
			for _, row := range data {
				if rowMap, ok := row.(map[string]interface{}); ok {
					if measurement, ok := rowMap["measurement"].(string); ok {
						measurements = append(measurements, measurement)
					}
				}
			}
		}
		return measurements, nil
	}

	return []string{}, nil
}

func (ds *DataSource3x) GetTagKeys(database, measurement string) (map[string]bool, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("client not connected")
	}

	tagKeys := make(map[string]bool)

	switch ds.client.compatMode {
	case "v1":
		query := fmt.Sprintf("SHOW TAG KEYS FROM \"%s\"", escapeMeasurement(measurement))
		resp, err := ds.client.QueryInfluxQL(query, database)
		if err != nil {
			return nil, err
		}
		if resp.Error() != nil {
			return nil, resp.Error()
		}

		if len(resp.Results) > 0 && len(resp.Results[0].Series) > 0 {
			for _, row := range resp.Results[0].Series[0].Values {
				if len(row) > 0 {
					if tagKey, ok := row[0].(string); ok {
						tagKeys[tagKey] = true
					}
				}
			}
		}

	case "v2":
		query := fmt.Sprintf(`
			import "influxdata/influxdb/schema"
			schema.tagKeys(
				bucket: "%s",
				predicate: (r) => r._measurement == "%s"
			)
		`, database, measurement)

		// 从 config 获取 org
		org := ""
		if v2cfg, ok := ds.config.(V2CompatConfig); ok {
			org = v2cfg.Org
		}
		
		result, err := ds.client.QueryFlux(query, org)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			record := result.Record()
			if tagKeyValue := record.ValueByKey("_value"); tagKeyValue != nil {
				if tagKey, ok := tagKeyValue.(string); ok {
					tagKeys[tagKey] = true
				}
			}
		}

	case "native":
		// 对于原生 3.x，我们可以查询一小批数据来确定 tag keys
		query := fmt.Sprintf("SELECT * FROM \"%s\" LIMIT 1", measurement)
		data, err := ds.client.QuerySQL(query)
		if err != nil {
			logx.Debug(fmt.Sprintf("获取标签字段失败，使用默认字段: %v", err))
		} else {
			// 解析响应来确定哪些是 tag fields
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				// 这里需要根据实际的 3.x API 响应格式来实现
				// 暂时返回一些常见的系统字段
			}
		}
	}

	// 添加通用的系统字段
	systemFields := []string{"_field", "_measurement", "_start", "_stop", "time"}
	for _, field := range systemFields {
		tagKeys[field] = true
	}

	return tagKeys, nil
}

func (ds *DataSource3x) QueryData(database, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	if ds.client == nil {
		return nil, 0, fmt.Errorf("client not connected")
	}

	switch ds.client.compatMode {
	case "v1":
		return ds.queryDataV1(database, measurement, startTime, batchSize)
	case "v2":
		return ds.queryDataV2(database, measurement, startTime, batchSize)
	case "native":
		return ds.queryDataNative(database, measurement, startTime, batchSize)
	}

	return nil, 0, fmt.Errorf("unsupported compatibility mode: %s", ds.client.compatMode)
}

// v1 兼容模式查询数据
func (ds *DataSource3x) queryDataV1(database, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	em := escapeMeasurement(measurement)
	var query string
	if startTime == 0 {
		query = fmt.Sprintf("SELECT * FROM %s ORDER BY time ASC LIMIT %d", em, batchSize)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE time > %d ORDER BY time ASC LIMIT %d", em, startTime, batchSize)
	}

	logx.Debug(fmt.Sprintf("执行查询: %s", query))
	resp, err := ds.client.QueryInfluxQL(query, database)
	if err != nil {
		return nil, 0, err
	}
	if resp.Error() != nil {
		return nil, 0, resp.Error()
	}

	return ds.parseInfluxQLResponse(resp, startTime)
}

// v2 兼容模式查询数据
func (ds *DataSource3x) queryDataV2(database, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	var startFilter string
	if startTime > 0 {
		startTimeRFC := time.Unix(0, startTime).UTC().Format(time.RFC3339Nano)
		startFilter = fmt.Sprintf(`|> filter(fn: (r) => r._time > time(v: "%s"))`, startTimeRFC)
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -100y)
		|> filter(fn: (r) => r._measurement == "%s")
		%s
		|> sort(columns: ["_time"])
		|> limit(n: %d)
	`, database, measurement, startFilter, batchSize)

	logx.Debug(fmt.Sprintf("执行 Flux 查询: %s", query))
	
	// 从 config 获取 org
	org := ""
	if v2cfg, ok := ds.config.(V2CompatConfig); ok {
		org = v2cfg.Org
	}
	
	result, err := ds.client.QueryFlux(query, org)
	if err != nil {
		return nil, 0, err
	}

	return ds.parseFluxResponse(result, startTime)
}

// 原生 3.x 模式查询数据
func (ds *DataSource3x) queryDataNative(database, measurement string, startTime int64, batchSize int) ([]common.DataPoint, int64, error) {
	var whereClause string
	if startTime > 0 {
		whereClause = fmt.Sprintf(" WHERE time > %d", startTime)
	}

	query := fmt.Sprintf("SELECT * FROM \"%s\"%s ORDER BY time ASC LIMIT %d", measurement, whereClause, batchSize)

	logx.Debug(fmt.Sprintf("执行 SQL 查询: %s", query))
	data, err := ds.client.QuerySQL(query)
	if err != nil {
		return nil, 0, err
	}

	return ds.parseSQLResponse(data, startTime)
}

// DataTarget 接口实现
func (dt *DataTarget3x) Connect() error {
	var err error

	switch config := dt.config.(type) {
	case V1CompatConfig:
		dt.client, err = NewV1CompatClient(config)
	case V2CompatConfig:
		dt.client, err = NewV2CompatClient(config)
	case NativeConfig:
		dt.client, err = NewClient3x(config)
	default:
		return fmt.Errorf("unsupported config type: %T", config)
	}

	if err != nil {
		return err
	}

	// 测试连接
	return dt.client.Ping(30 * time.Second)
}

func (dt *DataTarget3x) Close() error {
	if dt.client != nil {
		return dt.client.Close()
	}
	return nil
}

func (dt *DataTarget3x) WritePoints(database string, points []common.DataPoint) error {
	if dt.client == nil {
		return fmt.Errorf("client not connected")
	}

	// 转换为 Line Protocol 格式
	var lines []string
	for _, point := range points {
		line := formatLineProtocol(point)
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	data := strings.Join(lines, "\n")
	return dt.client.WriteLineProtocol(data)
}

// 工具函数
func escapeMeasurement(m string) string {
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(m, "\"", "\\\""))
}

func formatLineProtocol(point common.DataPoint) string {
	if len(point.Fields) == 0 {
		return ""
	}

	// 构建 measurement 和 tags
	var tagParts []string
	for key, value := range point.Tags {
		if key != "_field" && key != "_measurement" && key != "_start" && key != "_stop" && key != "time" {
			tagParts = append(tagParts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	measurement := point.Measurement
	if len(tagParts) > 0 {
		measurement += "," + strings.Join(tagParts, ",")
	}

	// 构建 fields
	var fieldParts []string
	for key, value := range point.Fields {
		switch v := value.(type) {
		case string:
			fieldParts = append(fieldParts, fmt.Sprintf("%s=\"%s\"", key, strings.ReplaceAll(v, "\"", "\\\"")))
		case float64:
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%g", key, v))
		case int64:
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%di", key, v))
		case bool:
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%t", key, v))
		default:
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", key, v))
		}
	}

	if len(fieldParts) == 0 {
		return ""
	}

	// 构建完整的 line protocol
	timestamp := point.Time.UnixNano()
	return fmt.Sprintf("%s %s %d", measurement, strings.Join(fieldParts, ","), timestamp)
}

// 响应解析函数 - 解析 InfluxQL 响应
func (ds *DataSource3x) parseInfluxQLResponse(resp *client.Response, startTime int64) ([]common.DataPoint, int64, error) {
	var points []common.DataPoint
	var maxTime int64 = startTime

	if len(resp.Results) == 0 || len(resp.Results[0].Series) == 0 {
		return points, maxTime, nil
	}

	for _, series := range resp.Results[0].Series {
		for _, row := range series.Values {
			tags := make(map[string]string)
			fields := make(map[string]interface{})
			var t time.Time
			var tUnix int64
			skip := false

			// 复制 series tags
			for k, v := range series.Tags {
				tags[k] = v
			}

			// 解析列值
			for idx, col := range series.Columns {
				switch col {
				case "time":
					v := row[idx]
					switch val := v.(type) {
					case string:
						if parsedTime, err := time.Parse(time.RFC3339, val); err == nil {
							t = parsedTime
							tUnix = t.UnixNano()
						}
					case json.Number:
						if ns, err := val.Int64(); err == nil {
							t = time.Unix(0, ns)
							tUnix = ns
						}
					default:
						logx.Warn(fmt.Sprintf("未知time类型: %T, measurement: %s", val, series.Name))
						skip = true
					}
				default:
					if val := row[idx]; val != nil {
						if sv, ok := val.(string); !ok || sv != "" {
							fields[col] = val
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

	return points, maxTime, nil
}

func (ds *DataSource3x) parseFluxResponse(result *api.QueryTableResult, startTime int64) ([]common.DataPoint, int64, error) {
	var points []common.DataPoint
	var maxTime int64 = startTime
	pointsMap := make(map[string]*common.DataPoint) // 按时间戳分组

	for result.Next() {
		record := result.Record()
		timeKey := record.Time().Format(time.RFC3339Nano)
		measurement := record.Measurement()

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

func (ds *DataSource3x) parseSQLResponse(data []byte, startTime int64) ([]common.DataPoint, int64, error) {
	// 解析 SQL 查询的 JSON 响应
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, 0, err
	}

	// 这里需要根据实际的 3.x SQL API 响应格式来实现
	// 目前返回空结果
	return []common.DataPoint{}, startTime, nil
}

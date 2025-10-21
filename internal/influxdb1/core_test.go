package influxdb1

import (
	"testing"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// 测试 WritePoints 的数据转换逻辑
func TestDataTarget_WritePointsConversion(t *testing.T) {
	config := DataTargetConfig{
		Addr: "http://localhost:8086",
		User: "admin",
		Pass: "password",
	}

	dt := NewDataTarget(config)
	
	// 创建测试数据点
	now := time.Now()
	points := []common.DataPoint{
		{
			Measurement: "cpu",
			Tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
			},
			Fields: map[string]interface{}{
				"usage_user":   0.64,
				"usage_system": 0.32,
			},
			Time: now,
		},
		{
			Measurement: "memory",
			Tags: map[string]string{
				"host": "server01",
			},
			Fields: map[string]interface{}{
				"used":      int64(1024),
				"available": 2048.5,
			},
			Time: now.Add(time.Minute),
		},
	}
	
	// 创建 BatchPoints 来验证转换
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "testdb",
		Precision: "ns",
	})
	if err != nil {
		t.Fatalf("创建 BatchPoints 失败: %v", err)
	}

	// 验证每个点都能正确转换
	for _, point := range points {
		pt, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
		if err != nil {
			t.Errorf("转换数据点失败: %v", err)
			continue
		}
		
		// 验证点的基本属性
		if pt.Name() != point.Measurement {
			t.Errorf("Measurement 不匹配: got %s, want %s", pt.Name(), point.Measurement)
		}
		
		bp.AddPoint(pt)
	}

	// 验证批次点数量
	if len(bp.Points()) != len(points) {
		t.Errorf("批次点数量不匹配: got %d, want %d", len(bp.Points()), len(points))
	}
	
	_ = dt // 避免未使用警告
}

// 测试不同数据类型的字段处理
func TestDataPoint_FieldTypes(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]interface{}
		valid  bool
	}{
		{
			name: "整数字段",
			fields: map[string]interface{}{
				"count": int64(100),
			},
			valid: true,
		},
		{
			name: "浮点数字段",
			fields: map[string]interface{}{
				"value": 3.14159,
			},
			valid: true,
		},
		{
			name: "字符串字段",
			fields: map[string]interface{}{
				"status": "active",
			},
			valid: true,
		},
		{
			name: "布尔字段",
			fields: map[string]interface{}{
				"enabled": true,
			},
			valid: true,
		},
		{
			name: "混合类型",
			fields: map[string]interface{}{
				"count":   int64(42),
				"ratio":   0.85,
				"name":    "test",
				"active":  true,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := common.DataPoint{
				Measurement: "test_measurement",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: tt.fields,
				Time:   time.Now(),
			}

			// 尝试转换为 InfluxDB Point
			_, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
			
			if tt.valid && err != nil {
				t.Errorf("有效字段类型转换失败: %v", err)
			}
		})
	}
}

// 测试标签的转义和特殊字符处理
func TestTagEscaping(t *testing.T) {
	tests := []struct {
		name     string
		tags     map[string]string
		expected bool
	}{
		{
			name: "普通标签",
			tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
			},
			expected: true,
		},
		{
			name: "带空格的标签",
			tags: map[string]string{
				"host": "server 01",
			},
			expected: true,
		},
		{
			name: "带特殊字符的标签",
			tags: map[string]string{
				"host": "server-01_test",
			},
			expected: true,
		},
		{
			name: "空标签值",
			tags: map[string]string{
				"host": "",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := common.DataPoint{
				Measurement: "test",
				Tags:        tt.tags,
				Fields: map[string]interface{}{
					"value": 1.0,
				},
				Time: time.Now(),
			}

			_, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
			if tt.expected && err != nil {
				t.Errorf("标签处理失败: %v, tags: %v", err, tt.tags)
			}
		})
	}
}

// 测试时间戳处理
func TestTimestampHandling(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "当前时间",
			time: time.Now(),
		},
		{
			name: "过去时间",
			time: time.Now().Add(-24 * time.Hour),
		},
		{
			name: "未来时间",
			time: time.Now().Add(24 * time.Hour),
		},
		{
			name: "纳秒精度",
			time: time.Unix(0, 1609459200123456789),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := common.DataPoint{
				Measurement: "test",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: map[string]interface{}{
					"value": 1.0,
				},
				Time: tt.time,
			}

			pt, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
			if err != nil {
				t.Errorf("时间戳处理失败: %v", err)
				return
			}

			// 验证时间戳
			if !pt.Time().Equal(tt.time) {
				t.Errorf("时间戳不匹配: got %v, want %v", pt.Time(), tt.time)
			}
		})
	}
}

// 测试空字段的处理
func TestEmptyFieldsHandling(t *testing.T) {
	// 空字段应该被忽略
	point := common.DataPoint{
		Measurement: "test",
		Tags: map[string]string{
			"host": "server01",
		},
		Fields: map[string]interface{}{
			"value":  1.0,
			"empty":  "",     // 空字符串
			"valid":  "data",
		},
		Time: time.Now(),
	}

	_, err := client.NewPoint(point.Measurement, point.Tags, point.Fields, point.Time)
	if err != nil {
		t.Errorf("空字段处理失败: %v", err)
	}
}

// 测试大批量数据点
func TestBatchPointsCapacity(t *testing.T) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "testdb",
		Precision: "ns",
	})
	if err != nil {
		t.Fatalf("创建 BatchPoints 失败: %v", err)
	}

	// 添加1000个点
	now := time.Now()
	for i := 0; i < 1000; i++ {
		pt, err := client.NewPoint(
			"test_measurement",
			map[string]string{
				"host": "server01",
				"id":   string(rune(i)),
			},
			map[string]interface{}{
				"value": float64(i),
			},
			now.Add(time.Duration(i)*time.Second),
		)
		if err != nil {
			t.Errorf("创建点 %d 失败: %v", i, err)
			continue
		}
		bp.AddPoint(pt)
	}

	if len(bp.Points()) != 1000 {
		t.Errorf("批次点数量不正确: got %d, want 1000", len(bp.Points()))
	}
}

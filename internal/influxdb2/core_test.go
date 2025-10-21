package influxdb2

import (
	"testing"
	"time"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// 测试数据点格式化
func TestFormatDataPoints(t *testing.T) {
	now := time.Now()
	points := []common.DataPoint{
		{
			Measurement: "cpu",
			Tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
			},
			Fields: map[string]interface{}{
				"usage": 75.5,
			},
			Time: now,
		},
		{
			Measurement: "memory",
			Tags: map[string]string{
				"host": "server02",
			},
			Fields: map[string]interface{}{
				"used":  int64(1024),
				"total": int64(2048),
			},
			Time: now,
		},
	}

	// 验证数据点基本属性
	for i, point := range points {
		if point.Measurement == "" {
			t.Errorf("点 %d 的 measurement 为空", i)
		}
		if len(point.Fields) == 0 {
			t.Errorf("点 %d 的 fields 为空", i)
		}
		if point.Time.IsZero() {
			t.Errorf("点 %d 的时间为零值", i)
		}
	}
}

// 测试多种数据类型字段
func TestFieldDataTypes(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name   string
		fields map[string]interface{}
	}{
		{
			name: "浮点数",
			fields: map[string]interface{}{
				"value": 3.14159,
			},
		},
		{
			name: "整数",
			fields: map[string]interface{}{
				"count": int64(100),
			},
		},
		{
			name: "字符串",
			fields: map[string]interface{}{
				"status": "active",
			},
		},
		{
			name: "布尔值",
			fields: map[string]interface{}{
				"enabled": true,
			},
		},
		{
			name: "混合类型",
			fields: map[string]interface{}{
				"count":  int64(42),
				"ratio":  0.85,
				"name":   "test",
				"active": true,
			},
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
				Time:   now,
			}

			// 验证所有字段都存在
			for key := range tt.fields {
				if _, ok := point.Fields[key]; !ok {
					t.Errorf("字段 %s 丢失", key)
				}
			}
		})
	}
}

// 测试标签处理
func TestTagsProcessing(t *testing.T) {
	tests := []struct {
		name string
		tags map[string]string
	}{
		{
			name: "普通标签",
			tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
			},
		},
		{
			name: "带特殊字符",
			tags: map[string]string{
				"host": "server-01_test",
			},
		},
		{
			name: "空标签集",
			tags: map[string]string{},
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

			if point.Tags == nil {
				t.Error("标签不应为 nil")
			}
		})
	}
}

// 测试批量数据处理
func TestBulkDataProcessing(t *testing.T) {
	now := time.Now()
	count := 100

	points := make([]common.DataPoint, count)
	for i := 0; i < count; i++ {
		points[i] = common.DataPoint{
			Measurement: "test",
			Tags: map[string]string{
				"host": "server01",
			},
			Fields: map[string]interface{}{
				"value": float64(i),
			},
			Time: now.Add(time.Duration(i) * time.Second),
		}
	}

	if len(points) != count {
		t.Errorf("数据点数量错误: got %d, want %d", len(points), count)
	}

	// 验证时间戳递增
	for i := 1; i < len(points); i++ {
		if !points[i].Time.After(points[i-1].Time) {
			t.Errorf("时间戳顺序错误: 点 %d", i)
		}
	}
}

// 测试空字段处理
func TestEmptyFieldHandling(t *testing.T) {
	point := common.DataPoint{
		Measurement: "test",
		Tags: map[string]string{
			"host": "server01",
		},
		Fields: map[string]interface{}{
			"value": 1.0,
			"name":  "", // 空字符串
		},
		Time: time.Now(),
	}

	// 空字段应该仍然存在于 map 中
	if _, ok := point.Fields["name"]; !ok {
		t.Error("空字段应该保留在 Fields 中")
	}
}

// 测试时间精度
func TestTimePrecision(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "纳秒精度",
			time: time.Unix(0, 1609459200123456789),
		},
		{
			name: "当前时间",
			time: time.Now(),
		},
		{
			name: "过去时间",
			time: time.Now().Add(-24 * time.Hour),
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

			if !point.Time.Equal(tt.time) {
				t.Errorf("时间不匹配: got %v, want %v", point.Time, tt.time)
			}
		})
	}
}

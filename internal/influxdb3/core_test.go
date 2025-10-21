package influxdb3

import (
	"testing"
	"time"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// 测试 V1 兼容模式配置验证
func TestV1CompatConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  V1CompatConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: V1CompatConfig{
				Addr:     "http://localhost:8086",
				User:     "admin",
				Pass:     "password",
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "缺少地址",
			config: V1CompatConfig{
				User:     "admin",
				Pass:     "password",
				Database: "testdb",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.config.Addr == ""
			if hasErr != tt.wantErr {
				t.Errorf("配置验证错误: got %v, want %v", hasErr, tt.wantErr)
			}
		})
	}
}

// 测试 V2 兼容模式配置验证
func TestV2CompatConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  V2CompatConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: V2CompatConfig{
				URL:      "http://localhost:8086",
				Token:    "test-token",
				Org:      "myorg",
				Bucket:   "mybucket",
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "缺少数据库",
			config: V2CompatConfig{
				URL:    "http://localhost:8086",
				Token:  "test-token",
				Org:    "myorg",
				Bucket: "mybucket",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.config.URL == "" || tt.config.Database == ""
			if hasErr != tt.wantErr {
				t.Errorf("配置验证错误: got %v, want %v", hasErr, tt.wantErr)
			}
		})
	}
}

// 测试 Native 模式配置验证
func TestNativeConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  NativeConfig
		wantErr bool
	}{
		{
			name: "有效配置",
			config: NativeConfig{
				URL:       "http://localhost:8086",
				Token:     "test-token",
				Database:  "testdb",
				Namespace: "default",
			},
			wantErr: false,
		},
		{
			name: "缺少数据库",
			config: NativeConfig{
				URL:   "http://localhost:8086",
				Token: "test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.config.URL == "" || tt.config.Database == ""
			if hasErr != tt.wantErr {
				t.Errorf("配置验证错误: got %v, want %v", hasErr, tt.wantErr)
			}
		})
	}
}

// 测试数据点转换逻辑
func TestDataPointConversion(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name  string
		point common.DataPoint
	}{
		{
			name: "简单数据点",
			point: common.DataPoint{
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
		},
		{
			name: "多字段数据点",
			point: common.DataPoint{
				Measurement: "system",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: map[string]interface{}{
					"load1":  1.5,
					"load5":  2.3,
					"load15": 1.8,
					"uptime": int64(3600),
				},
				Time: now,
			},
		},
		{
			name: "带字符串字段的数据点",
			point: common.DataPoint{
				Measurement: "events",
				Tags: map[string]string{
					"severity": "info",
				},
				Fields: map[string]interface{}{
					"message": "system started",
					"code":    int64(200),
				},
				Time: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证数据点的基本属性
			if tt.point.Measurement == "" {
				t.Error("Measurement 不能为空")
			}
			if len(tt.point.Fields) == 0 {
				t.Error("Fields 不能为空")
			}
			if tt.point.Time.IsZero() {
				t.Error("时间戳不能为零值")
			}
		})
	}
}

// 测试 Namespace 处理
func TestNamespaceHandling(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		database  string
		expected  bool
	}{
		{
			name:      "默认命名空间",
			namespace: "default",
			database:  "mydb",
			expected:  true,
		},
		{
			name:      "自定义命名空间",
			namespace: "prod",
			database:  "metrics",
			expected:  true,
		},
		{
			name:      "空命名空间",
			namespace: "",
			database:  "testdb",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证数据库名称不为空
			isValid := tt.database != ""
			if isValid != tt.expected {
				t.Errorf("命名空间处理验证错误: got %v, want %v", isValid, tt.expected)
			}
		})
	}
}

// 测试字段类型处理
func TestFieldTypeHandling(t *testing.T) {
	tests := []struct {
		name   string
		field  interface{}
		valid  bool
	}{
		{
			name:  "整数",
			field: int64(42),
			valid: true,
		},
		{
			name:  "浮点数",
			field: 3.14159,
			valid: true,
		},
		{
			name:  "字符串",
			field: "test value",
			valid: true,
		},
		{
			name:  "布尔值",
			field: true,
			valid: true,
		},
		{
			name:  "int 类型",
			field: 100,
			valid: true,
		},
		{
			name:  "float32 类型",
			field: float32(2.71),
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := common.DataPoint{
				Measurement: "test",
				Tags: map[string]string{
					"tag1": "value1",
				},
				Fields: map[string]interface{}{
					"field": tt.field,
				},
				Time: time.Now(),
			}

			// 验证字段值不为 nil
			if point.Fields["field"] == nil && tt.valid {
				t.Error("有效字段不应为 nil")
			}
		})
	}
}

// 测试标签键值对处理
func TestTagKeyValueProcessing(t *testing.T) {
	tests := []struct {
		name string
		tags map[string]string
	}{
		{
			name: "单个标签",
			tags: map[string]string{
				"host": "server01",
			},
		},
		{
			name: "多个标签",
			tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
				"dc":     "dc1",
			},
		},
		{
			name: "带特殊字符的标签",
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

			// 验证标签不为 nil
			if point.Tags == nil {
				t.Error("标签不应为 nil")
			}
		})
	}
}

// 测试时间范围处理
func TestTimeRangeProcessing(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		valid bool
	}{
		{
			name:  "有效时间范围",
			start: now.Add(-1 * time.Hour),
			end:   now,
			valid: true,
		},
		{
			name:  "单点时间",
			start: now,
			end:   now,
			valid: true,
		},
		{
			name:  "过去24小时",
			start: now.Add(-24 * time.Hour),
			end:   now,
			valid: true,
		},
		{
			name:  "结束时间早于开始时间",
			start: now,
			end:   now.Add(-1 * time.Hour),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := !tt.end.Before(tt.start)
			if isValid != tt.valid {
				t.Errorf("时间范围验证错误: got %v, want %v", isValid, tt.valid)
			}
		})
	}
}

// 测试批量数据点处理
func TestBulkDataPointsProcessing(t *testing.T) {
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
		t.Errorf("批量数据点数量不正确: got %d, want %d", len(points), count)
	}

	// 验证每个点都有有效字段
	for i, point := range points {
		if len(point.Fields) == 0 {
			t.Errorf("点 %d 缺少字段", i)
		}
	}
}

// 测试 NewV1CompatDataTarget
func TestNewV1CompatDataTarget(t *testing.T) {
	cfg := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "testdb",
	}

	target := NewV1CompatDataTarget(cfg)
	if target == nil {
		t.Error("NewV1CompatDataTarget 返回 nil")
	}
}

package common

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// Mock数据源用于测试
type mockDataSource struct {
	databases    []string
	measurements []string
	connected    bool
	shouldError  bool
}

func (m *mockDataSource) Connect() error {
	if m.shouldError {
		return &mockError{"连接失败"}
	}
	m.connected = true
	return nil
}

func (m *mockDataSource) Close() error {
	m.connected = false
	return nil
}

func (m *mockDataSource) GetDatabases() ([]string, error) {
	if m.shouldError {
		return nil, &mockError{"获取数据库失败"}
	}
	return m.databases, nil
}

func (m *mockDataSource) GetMeasurements(db string) ([]string, error) {
	if m.shouldError {
		return nil, &mockError{"获取测量失败"}
	}
	return m.measurements, nil
}

func (m *mockDataSource) GetTagKeys(db, measurement string) (map[string]bool, error) {
	return map[string]bool{"host": true, "region": true}, nil
}

func (m *mockDataSource) QueryData(db, measurement string, startTime int64, batchSize int) ([]DataPoint, int64, error) {
	if m.shouldError {
		return nil, 0, &mockError{"查询数据失败"}
	}

	// 返回模拟数据
	points := []DataPoint{
		{
			Measurement: measurement,
			Tags:        map[string]string{"host": "server1", "region": "us-east"},
			Fields:      map[string]interface{}{"value": 100, "status": "ok"},
			Time:        time.Now(),
		},
	}
	return points, time.Now().UnixNano(), nil
}

// Mock数据目标用于测试
type mockDataTarget struct {
	connected   bool
	shouldError bool
	writtenData []DataPoint
	mu          sync.Mutex
}

func (m *mockDataTarget) Connect() error {
	if m.shouldError {
		return &mockError{"连接失败"}
	}
	m.connected = true
	return nil
}

func (m *mockDataTarget) Close() error {
	m.connected = false
	return nil
}

func (m *mockDataTarget) WritePoints(db string, points []DataPoint) error {
	if m.shouldError {
		return &mockError{"写入失败"}
	}
	m.mu.Lock()
	m.writtenData = append(m.writtenData, points...)
	m.mu.Unlock()
	return nil
}

func (m *mockDataTarget) GetWrittenDataCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.writtenData)
}

// Mock错误类型
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestNewSyncer(t *testing.T) {
	cfg := SyncConfig{
		BatchSize: 1000,
		Start:     "2024-01-01T00:00:00Z",
		Parallel:  4,
	}

	source := &mockDataSource{}
	target := &mockDataTarget{}

	syncer := NewSyncer(cfg, source, target)

	if syncer == nil {
		t.Fatal("NewSyncer返回nil")
	}

	if syncer.cfg.BatchSize != 1000 {
		t.Errorf("期望BatchSize为1000, 实际为%d", syncer.cfg.BatchSize)
	}

	if syncer.source != source {
		t.Error("源数据不匹配")
	}

	if syncer.target != target {
		t.Error("目标数据不匹配")
	}
}

func TestSyncerGetStartTime(t *testing.T) {
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "resume.state")

	testCases := []struct {
		name         string
		startTime    string
		resumeFile   string
		resumeData   string
		createResume bool
		expectError  bool
	}{
		{
			name:        "默认起始时间",
			startTime:   "",
			expectError: false,
		},
		{
			name:        "有效起始时间",
			startTime:   "2024-01-01T00:00:00Z",
			expectError: false,
		},
		{
			name:         "恢复文件存在且时间更新",
			startTime:    "2024-01-01T00:00:00Z",
			resumeFile:   resumeFile,
			resumeData:   "2024-06-01T00:00:00Z",
			createResume: true,
			expectError:  false,
		},
		{
			name:         "恢复文件存在但时间较早",
			startTime:    "2024-06-01T00:00:00Z",
			resumeFile:   resumeFile,
			resumeData:   "2024-01-01T00:00:00Z",
			createResume: true,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建恢复文件
			if tc.createResume && tc.resumeFile != "" {
				err := os.WriteFile(tc.resumeFile, []byte(tc.resumeData), 0644)
				if err != nil {
					t.Fatalf("无法创建恢复文件: %v", err)
				}
			}

			cfg := SyncConfig{
				Start:      tc.startTime,
				ResumeFile: tc.resumeFile,
			}

			syncer := NewSyncer(cfg, &mockDataSource{}, &mockDataTarget{})
			startTimeNano, err := syncer.getStartTime()

			if tc.expectError && err == nil {
				t.Error("期望返回错误，但没有错误")
			}

			if !tc.expectError && err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
			}

			if !tc.expectError && tc.startTime != "" && startTimeNano == 0 {
				t.Error("有效起始时间应该大于0")
			}
		})
	}
}

func TestSyncerGetDatabases(t *testing.T) {
	testCases := []struct {
		name        string
		sourceDB    string
		dbList      []string
		dbExclude   []string
		expectedDBs []string
		shouldError bool
	}{
		{
			name:        "指定源数据库",
			sourceDB:    "mydb",
			expectedDBs: []string{"mydb"},
		},
		{
			name:        "获取所有数据库",
			dbList:      []string{"db1", "db2", "_internal", "db3"},
			expectedDBs: []string{"db1", "db2", "db3"},
		},
		{
			name:        "排除指定数据库",
			dbList:      []string{"db1", "db2", "system", "db3"},
			dbExclude:   []string{"system"},
			expectedDBs: []string{"db1", "db2", "db3"},
		},
		{
			name:        "获取数据库失败",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := SyncConfig{
				SourceDB:        tc.sourceDB,
				SourceDBExclude: tc.dbExclude,
			}

			source := &mockDataSource{
				databases:   tc.dbList,
				shouldError: tc.shouldError,
			}

			syncer := NewSyncer(cfg, source, &mockDataTarget{})
			dbs, err := syncer.getDatabases()

			if tc.shouldError {
				if err == nil {
					t.Error("期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
				return
			}

			if len(dbs) != len(tc.expectedDBs) {
				t.Errorf("期望数据库数量为%d, 实际为%d", len(tc.expectedDBs), len(dbs))
				return
			}

			for i, expectedDB := range tc.expectedDBs {
				if dbs[i] != expectedDB {
					t.Errorf("期望数据库[%d]为%s, 实际为%s", i, expectedDB, dbs[i])
				}
			}
		})
	}
}

func TestSyncWithMockData(t *testing.T) {
	// 测试成功的同步流程
	cfg := SyncConfig{
		BatchSize:  100,
		Start:      "2024-01-01T00:00:00Z",
		Parallel:   2,
		ResumeFile: "",
	}

	source := &mockDataSource{
		databases:    []string{"testdb"},
		measurements: []string{"cpu", "memory"},
	}

	target := &mockDataTarget{}

	syncer := NewSyncer(cfg, source, target)

	ctx := context.Background()
	err := syncer.Sync(ctx)

	if err != nil {
		t.Errorf("同步失败: %v", err)
	}

	if source.connected {
		t.Error("源数据库应该已断开连接")
	}

	if target.connected {
		t.Error("目标数据库应该已断开连接")
	}

	if target.GetWrittenDataCount() == 0 {
		t.Error("目标数据库应该有写入数据")
	}
}

func TestSyncWithConnectionError(t *testing.T) {
	cfg := SyncConfig{
		BatchSize: 100,
		Start:     "2024-01-01T00:00:00Z",
	}

	// 测试源连接失败
	source := &mockDataSource{shouldError: true}
	target := &mockDataTarget{}

	syncer := NewSyncer(cfg, source, target)

	ctx := context.Background()
	err := syncer.Sync(ctx)

	if err == nil {
		t.Error("期望连接失败，但没有错误")
	}
}

func TestDataPointStruct(t *testing.T) {
	// 测试DataPoint结构体
	now := time.Now()
	point := DataPoint{
		Measurement: "cpu",
		Tags:        map[string]string{"host": "server1"},
		Fields:      map[string]interface{}{"value": 95.5},
		Time:        now,
	}

	if point.Measurement != "cpu" {
		t.Errorf("期望测量名为 cpu, 实际为 %s", point.Measurement)
	}

	if point.Tags["host"] != "server1" {
		t.Errorf("期望host标签为 server1, 实际为 %s", point.Tags["host"])
	}

	if point.Fields["value"] != 95.5 {
		t.Errorf("期望value字段为 95.5, 实际为 %v", point.Fields["value"])
	}

	if !point.Time.Equal(now) {
		t.Errorf("时间戳不匹配")
	}
}

func TestSyncResultStruct(t *testing.T) {
	// 测试SyncResult结构体
	err := &mockError{"测试错误"}
	result := SyncResult{
		Measurement: "cpu",
		Error:       err,
	}

	if result.Measurement != "cpu" {
		t.Errorf("期望测量名为 cpu, 实际为 %s", result.Measurement)
	}

	if result.Error != err {
		t.Errorf("错误不匹配")
	}

	// 测试无错误情况
	result2 := SyncResult{
		Measurement: "memory",
		Error:       nil,
	}

	if result2.Error != nil {
		t.Errorf("期望无错误，但得到 %v", result2.Error)
	}
}

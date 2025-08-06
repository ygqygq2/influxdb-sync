package common

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSyncerWithRealConfiguration(t *testing.T) {
	// 创建临时resume文件用于测试
	tmpDir := t.TempDir()
	resumeFile := filepath.Join(tmpDir, "test_resume.state")

	// 写入resume时间
	resumeTime := "2024-06-01T12:00:00Z"
	err := os.WriteFile(resumeFile, []byte(resumeTime), 0644)
	if err != nil {
		t.Fatalf("无法创建resume文件: %v", err)
	}

	cfg := SyncConfig{
		SourceDB:        "specific_db",
		SourceDBExclude: []string{"exclude1", "exclude2"},
		BatchSize:       500,
		Start:           "2024-01-01T00:00:00Z",
		ResumeFile:      resumeFile,
		Parallel:        2,
		RetryCount:      3,
		RetryInterval:   100,
		RateLimit:       10,
		LogLevel:        "debug",
	}

	source := &mockDataSource{
		databases:    []string{"db1", "db2", "exclude1", "specific_db", "_internal"},
		measurements: []string{"cpu", "memory", "disk"},
	}

	target := &mockDataTarget{}
	syncer := NewSyncer(cfg, source, target)

	// 测试各种同步器方法
	t.Run("测试起始时间获取", func(t *testing.T) {
		startTime, err := syncer.getStartTime()
		if err != nil {
			t.Errorf("获取起始时间失败: %v", err)
		}
		if startTime <= 0 {
			t.Error("起始时间应该大于0")
		}
	})

	t.Run("测试数据库列表获取", func(t *testing.T) {
		dbs, err := syncer.getDatabases()
		if err != nil {
			t.Errorf("获取数据库列表失败: %v", err)
		}
		if len(dbs) != 1 || dbs[0] != "specific_db" {
			t.Errorf("期望数据库列表为 [specific_db], 实际为 %v", dbs)
		}
	})

	t.Run("测试同步执行", func(t *testing.T) {
		ctx := context.Background()
		err := syncer.Sync(ctx)
		if err != nil {
			t.Errorf("同步执行失败: %v", err)
		}

		// 验证数据已写入
		if len(target.writtenData) == 0 {
			t.Error("应该有数据被写入目标")
		}
	})
}

func TestSyncerEdgeCases(t *testing.T) {
	t.Run("空数据库列表", func(t *testing.T) {
		cfg := SyncConfig{BatchSize: 100, Parallel: 1}
		source := &mockDataSource{databases: []string{}}
		target := &mockDataTarget{}
		syncer := NewSyncer(cfg, source, target)

		ctx := context.Background()
		err := syncer.Sync(ctx)
		if err != nil {
			t.Errorf("空数据库列表同步失败: %v", err)
		}
	})

	t.Run("空measurement列表", func(t *testing.T) {
		cfg := SyncConfig{BatchSize: 100, Parallel: 1}
		source := &mockDataSource{
			databases:    []string{"testdb"},
			measurements: []string{},
		}
		target := &mockDataTarget{}
		syncer := NewSyncer(cfg, source, target)

		ctx := context.Background()
		err := syncer.Sync(ctx)
		if err != nil {
			t.Errorf("空measurement列表同步失败: %v", err)
		}
	})

	t.Run("源连接失败", func(t *testing.T) {
		cfg := SyncConfig{BatchSize: 100}
		source := &mockDataSource{shouldError: true}
		target := &mockDataTarget{}
		syncer := NewSyncer(cfg, source, target)

		ctx := context.Background()
		err := syncer.Sync(ctx)
		if err == nil {
			t.Error("期望源连接失败返回错误")
		}
	})

	t.Run("目标连接失败", func(t *testing.T) {
		cfg := SyncConfig{BatchSize: 100}
		source := &mockDataSource{databases: []string{"testdb"}}
		target := &mockDataTarget{shouldError: true}
		syncer := NewSyncer(cfg, source, target)

		ctx := context.Background()
		err := syncer.Sync(ctx)
		if err == nil {
			t.Error("期望目标连接失败返回错误")
		}
	})
}

func TestSyncerConfigurationVariations(t *testing.T) {
	testCases := []struct {
		name   string
		config SyncConfig
	}{
		{
			name: "最小配置",
			config: SyncConfig{
				BatchSize: 100,
				Parallel:  1,
			},
		},
		{
			name: "大批次配置",
			config: SyncConfig{
				BatchSize: 10000,
				Parallel:  1,
			},
		},
		{
			name: "高并发配置",
			config: SyncConfig{
				BatchSize: 1000,
				Parallel:  16,
			},
		},
		{
			name: "低延迟配置",
			config: SyncConfig{
				BatchSize:     100,
				Parallel:      4,
				RetryCount:    1,
				RetryInterval: 10,
				RateLimit:     0,
			},
		},
		{
			name: "高可靠性配置",
			config: SyncConfig{
				BatchSize:     500,
				Parallel:      2,
				RetryCount:    10,
				RetryInterval: 2000,
				RateLimit:     1000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			source := &mockDataSource{
				databases:    []string{"testdb"},
				measurements: []string{"metric1"},
			}
			target := &mockDataTarget{}
			syncer := NewSyncer(tc.config, source, target)

			ctx := context.Background()
			err := syncer.Sync(ctx)
			if err != nil {
				t.Errorf("配置 %s 同步失败: %v", tc.name, err)
			}
		})
	}
}

func TestDataSourceInterface(t *testing.T) {
	// 测试DataSource接口的所有方法
	source := &mockDataSource{
		databases:    []string{"db1", "db2"},
		measurements: []string{"cpu", "memory"},
	}

	// 测试Connect方法
	err := source.Connect()
	if err != nil {
		t.Errorf("Connect失败: %v", err)
	}

	// 测试GetDatabases方法
	dbs, err := source.GetDatabases()
	if err != nil {
		t.Errorf("GetDatabases失败: %v", err)
	}
	if len(dbs) != 2 {
		t.Errorf("期望2个数据库，实际 %d 个", len(dbs))
	}

	// 测试GetMeasurements方法
	measurements, err := source.GetMeasurements("db1")
	if err != nil {
		t.Errorf("GetMeasurements失败: %v", err)
	}
	if len(measurements) != 2 {
		t.Errorf("期望2个measurement，实际 %d 个", len(measurements))
	}

	// 测试GetTagKeys方法
	tagKeys, err := source.GetTagKeys("db1", "cpu")
	if err != nil {
		t.Errorf("GetTagKeys失败: %v", err)
	}
	if len(tagKeys) == 0 {
		t.Error("期望有tag keys返回")
	}

	// 测试QueryData方法
	points, maxTime, err := source.QueryData("db1", "cpu", time.Now().UnixNano(), 100)
	if err != nil {
		t.Errorf("QueryData失败: %v", err)
	}
	if len(points) == 0 {
		t.Error("期望有数据点返回")
	}
	if maxTime <= 0 {
		t.Error("maxTime应该大于0")
	}

	// 测试Close方法
	err = source.Close()
	if err != nil {
		t.Errorf("Close失败: %v", err)
	}
}

func TestDataTargetInterface(t *testing.T) {
	// 测试DataTarget接口的所有方法
	target := &mockDataTarget{}

	// 测试Connect方法
	err := target.Connect()
	if err != nil {
		t.Errorf("Connect失败: %v", err)
	}

	// 测试WritePoints方法
	points := []DataPoint{
		{
			Measurement: "test",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]interface{}{"value": 100},
			Time:        time.Now(),
		},
	}
	err = target.WritePoints("testdb", points)
	if err != nil {
		t.Errorf("WritePoints失败: %v", err)
	}

	// 验证数据写入
	if len(target.writtenData) != 1 {
		t.Errorf("期望写入1个点，实际写入 %d 个", len(target.writtenData))
	}

	// 测试Close方法
	err = target.Close()
	if err != nil {
		t.Errorf("Close失败: %v", err)
	}
}

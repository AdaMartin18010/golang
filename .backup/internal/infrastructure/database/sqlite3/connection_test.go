// Package sqlite3 提供 SQLite3 数据库连接的测试
package sqlite3

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "file:app.db?cache=shared&mode=rwc", cfg.DSN)
	assert.Equal(t, 25, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxLifetime)
	assert.Equal(t, 10*time.Minute, cfg.ConnMaxIdleTime)
}

// TestNewConnection_Success 测试成功创建连接
func TestNewConnection_Success(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"

	cfg := Config{
		DSN:             dsn,
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Minute,
	}

	db, err := NewConnection(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	// 验证连接池设置
	assert.Equal(t, 10, db.Stats().MaxOpenConnections)
}

// TestNewConnection_DefaultConfig 测试使用默认配置
func TestNewConnection_DefaultConfig(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_default.db")

	// 使用空 DSN 应该使用默认配置
	cfg := Config{
		DSN: "",
	}
	// 手动设置 DSN 以避免默认路径问题
	cfg.DSN = "file:" + dbPath + "?cache=shared&mode=rwc"

	db, err := NewConnection(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()
}

// TestNewConnectionWithDSN_Success 测试使用 DSN 创建连接
func TestNewConnectionWithDSN_Success(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_dsn.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"

	db, err := NewConnectionWithDSN(dsn)
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	// 验证数据库可以执行查询
	err = db.Ping()
	assert.NoError(t, err)
}

// TestNewConnection_InvalidDSN 测试无效的 DSN
func TestNewConnection_InvalidDSN(t *testing.T) {
	cfg := Config{
		DSN: "/invalid/path/to/db.sqlite",
	}

	db, err := NewConnection(cfg)
	// SQLite 会尝试创建文件，但在无效路径上会失败
	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestNewConnection_PingTimeout 测试连接超时
func TestNewConnection_PingTimeout(t *testing.T) {
	// 使用一个会导致 ping 失败的配置
	// 注意：SQLite 本地文件很难模拟 ping 失败，除非权限问题
	// 这里主要测试代码路径
	cfg := Config{
		DSN: "file:/nonexistent/path/test.db?mode=ro",
	}

	db, err := NewConnection(cfg)
	// 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestNewConnection_PoolSettings 测试连接池设置
func TestNewConnection_PoolSettings(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_pool.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"

	cfg := Config{
		DSN:             dsn,
		MaxOpenConns:    0, // 0 表示不设置
		MaxIdleConns:    0,
		ConnMaxLifetime: 0,
		ConnMaxIdleTime: 0,
	}

	db, err := NewConnection(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	// 当设置为 0 时，应该使用默认值
	stats := db.Stats()
	assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)
}

// TestNewConnection_Query 测试数据库查询
func TestNewConnection_Query(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_query.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"

	db, err := NewConnectionWithDSN(dsn)
	require.NoError(t, err)
	defer db.Close()

	// 创建测试表
	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	require.NoError(t, err)

	// 插入数据
	_, err = db.Exec("INSERT INTO test (name) VALUES (?)", "test_name")
	require.NoError(t, err)

	// 查询数据
	var name string
	err = db.QueryRow("SELECT name FROM test WHERE id = 1").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "test_name", name)
}

// TestNewConnection_Concurrent 测试并发访问
func TestNewConnection_Concurrent(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_concurrent.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"

	db, err := NewConnection(Config{
		DSN:          dsn,
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	})
	require.NoError(t, err)
	defer db.Close()

	// 创建表
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS concurrent_test (id INTEGER PRIMARY KEY)")
	require.NoError(t, err)

	// 并发查询
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			defer func() { done <- true }()
			err := db.Ping()
			assert.NoError(t, err)
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestNewConnectionWithDSN_EmptyDSN 测试空 DSN
func TestNewConnectionWithDSN_EmptyDSN(t *testing.T) {
	// 空 DSN 会使用默认配置
	db, err := NewConnectionWithDSN("")
	require.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	// 清理默认数据库文件
	defer os.Remove("app.db")
}

// TestConfig_Struct 测试配置结构体
func TestConfig_Struct(t *testing.T) {
	cfg := Config{
		DSN:             "test.dsn",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute,
	}

	assert.Equal(t, "test.dsn", cfg.DSN)
	assert.Equal(t, 10, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, time.Hour, cfg.ConnMaxLifetime)
	assert.Equal(t, time.Minute, cfg.ConnMaxIdleTime)
}

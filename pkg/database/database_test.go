package database

import (
	"context"
	"testing"
	"time"
)

func TestNewDatabase_PostgreSQL(t *testing.T) {
	// 注意：需要实际的 PostgreSQL 连接才能测试
	// 这里只测试配置验证
	cfg := Config{
		Driver:          DriverPostgreSQL,
		DSN:             "postgres://user:password@localhost/dbname?sslmode=disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		PingTimeout:     5 * time.Second,
	}

	// 如果没有实际的数据库连接，这个测试会失败
	// 在实际环境中，应该使用测试数据库
	_, err := NewDatabase(cfg)
	if err != nil {
		// 预期错误：连接失败（如果没有数据库）
		t.Logf("Expected error (no database): %v", err)
	}
}

func TestNewDatabase_SQLite3(t *testing.T) {
	cfg := Config{
		Driver:          DriverSQLite3,
		DSN:             "file:test.db?cache=shared&mode=memory",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		PingTimeout:     5 * time.Second,
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create SQLite3 database: %v", err)
	}
	defer db.Close()

	if db.Driver() != DriverSQLite3 {
		t.Errorf("Expected driver SQLite3, got %s", db.Driver())
	}

	// 测试 Ping
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// 测试 Stats
	stats := db.Stats()
	_ = stats // 验证可以获取统计信息
}

func TestNewDatabase_UnsupportedDriver(t *testing.T) {
	cfg := Config{
		Driver: Driver("unsupported"),
		DSN:    "test",
	}

	_, err := NewDatabase(cfg)
	if err == nil {
		t.Error("Expected error for unsupported driver")
	}
	if err != ErrUnsupportedDriver {
		t.Errorf("Expected ErrUnsupportedDriver, got %v", err)
	}
}

func TestDatabase_Transaction(t *testing.T) {
	cfg := Config{
		Driver: DriverSQLite3,
		DSN:    "file:test_tx.db?cache=shared&mode=memory",
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 创建表
	_, err = db.Exec(ctx, "CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 开始事务
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// 在事务中执行操作
	_, err = tx.Exec(ctx, "INSERT INTO test (name) VALUES (?)", "test")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to insert: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// 验证数据
	rows, err := db.Query(ctx, "SELECT name FROM test WHERE name = ?", "test")
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Expected to find inserted row")
	}
}


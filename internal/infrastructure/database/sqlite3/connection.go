package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Config SQLite3 连接配置
type Config struct {
	DSN             string        // 数据源名称（例如: "file:example.db?cache=shared&mode=rwc"）
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生存时间
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		DSN:             "file:app.db?cache=shared&mode=rwc",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// NewConnection 创建 SQLite3 数据库连接
func NewConnection(config Config) (*sql.DB, error) {
	if config.DSN == "" {
		config = DefaultConfig()
	}

	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 database: %w", err)
	}

	// 设置连接池参数
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite3 database: %w", err)
	}

	return db, nil
}

// NewConnectionWithDSN 使用 DSN 创建连接（快捷方式）
func NewConnectionWithDSN(dsn string) (*sql.DB, error) {
	return NewConnection(Config{DSN: dsn})
}

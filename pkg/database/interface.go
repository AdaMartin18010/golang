package database

import (
	"context"
	"database/sql"
	"time"
)

// Driver 数据库驱动类型
type Driver string

const (
	DriverPostgreSQL Driver = "postgresql"
	DriverSQLite3    Driver = "sqlite3"
	DriverMySQL      Driver = "mysql"
)

// Config 数据库配置
type Config struct {
	Driver          Driver        // 数据库驱动类型
	DSN             string        // 数据源名称
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生存时间
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
	PingTimeout     time.Duration // Ping 超时时间
}

// Database 通用数据库接口
// 提供统一的数据库操作接口，支持多种数据库驱动
type Database interface {
	// Driver 返回数据库驱动类型
	Driver() Driver

	// DB 返回底层的 *sql.DB 连接
	DB() *sql.DB

	// Ping 检查数据库连接
	Ping(ctx context.Context) error

	// Close 关闭数据库连接
	Close() error

	// Stats 返回连接池统计信息
	Stats() sql.DBStats

	// Begin 开始一个事务
	Begin(ctx context.Context) (Transaction, error)

	// Exec 执行 SQL 语句（不返回结果）
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Query 执行查询 SQL 语句
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow 执行查询 SQL 语句，返回单行
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Prepare 准备 SQL 语句
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
}

// Transaction 事务接口
type Transaction interface {
	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// Exec 在事务中执行 SQL 语句
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Query 在事务中执行查询 SQL 语句
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow 在事务中执行查询 SQL 语句，返回单行
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Prepare 在事务中准备 SQL 语句
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
}

// NewDatabase 创建数据库连接
// 根据配置自动选择对应的数据库驱动
func NewDatabase(cfg Config) (Database, error) {
	switch cfg.Driver {
	case DriverPostgreSQL:
		return NewPostgreSQL(cfg)
	case DriverSQLite3:
		return NewSQLite3(cfg)
	case DriverMySQL:
		return NewMySQL(cfg)
	default:
		return nil, ErrUnsupportedDriver
	}
}

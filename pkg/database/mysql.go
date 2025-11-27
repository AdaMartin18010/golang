package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
)

// mysql 实现 MySQL 数据库
type mysql struct {
	driver Driver
	db     *sql.DB
	config Config
}

// NewMySQL 创建 MySQL 数据库连接
func NewMySQL(cfg Config) (Database, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("mysql DSN is required")
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// 设置连接池参数
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// 测试连接
	pingTimeout := cfg.PingTimeout
	if pingTimeout == 0 {
		pingTimeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %v", ErrPingFailed, err)
	}

	return &mysql{
		driver: DriverMySQL,
		db:     db,
		config: cfg,
	}, nil
}

func (m *mysql) Driver() Driver {
	return m.driver
}

func (m *mysql) DB() *sql.DB {
	return m.db
}

func (m *mysql) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

func (m *mysql) Close() error {
	return m.db.Close()
}

func (m *mysql) Stats() sql.DBStats {
	return m.db.Stats()
}

func (m *mysql) Begin(ctx context.Context) (Transaction, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

func (m *mysql) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *mysql) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *mysql) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m *mysql) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return m.db.PrepareContext(ctx, query)
}

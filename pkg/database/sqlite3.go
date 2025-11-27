package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite3 驱动
)

// sqlite3 实现 SQLite3 数据库
type sqlite3 struct {
	driver Driver
	db     *sql.DB
	config Config
}

// NewSQLite3 创建 SQLite3 数据库连接
func NewSQLite3(cfg Config) (Database, error) {
	if cfg.DSN == "" {
		cfg.DSN = "file:app.db?cache=shared&mode=rwc"
	}

	db, err := sql.Open("sqlite3", cfg.DSN)
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

	return &sqlite3{
		driver: DriverSQLite3,
		db:     db,
		config: cfg,
	}, nil
}

func (s *sqlite3) Driver() Driver {
	return s.driver
}

func (s *sqlite3) DB() *sql.DB {
	return s.db
}

func (s *sqlite3) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *sqlite3) Close() error {
	return s.db.Close()
}

func (s *sqlite3) Stats() sql.DBStats {
	return s.db.Stats()
}

func (s *sqlite3) Begin(ctx context.Context) (Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

func (s *sqlite3) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqlite3) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *sqlite3) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *sqlite3) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}

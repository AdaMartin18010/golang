package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL 驱动
)

// postgresql 实现 PostgreSQL 数据库
type postgresql struct {
	driver Driver
	db     *sql.DB
	config Config
}

// NewPostgreSQL 创建 PostgreSQL 数据库连接
func NewPostgreSQL(cfg Config) (Database, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("postgresql DSN is required")
	}

	db, err := sql.Open("postgres", cfg.DSN)
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

	return &postgresql{
		driver: DriverPostgreSQL,
		db:     db,
		config: cfg,
	}, nil
}

func (p *postgresql) Driver() Driver {
	return p.driver
}

func (p *postgresql) DB() *sql.DB {
	return p.db
}

func (p *postgresql) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *postgresql) Close() error {
	return p.db.Close()
}

func (p *postgresql) Stats() sql.DBStats {
	return p.db.Stats()
}

func (p *postgresql) Begin(ctx context.Context) (Transaction, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

func (p *postgresql) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

func (p *postgresql) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

func (p *postgresql) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

func (p *postgresql) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return p.db.PrepareContext(ctx, query)
}

package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config 数据库配置
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// Connection 数据库连接
type Connection struct {
	pool *pgxpool.Pool
}

// NewConnection 创建数据库连接
func NewConnection(ctx context.Context, cfg Config) (*Connection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 测试连接
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{pool: pool}, nil
}

// Close 关闭连接
func (c *Connection) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// Pool 获取连接池
func (c *Connection) Pool() *pgxpool.Pool {
	return c.pool
}

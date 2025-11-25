package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver

	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/internal/infrastructure/database/ent/schema"
)

// Client Ent 客户端
type Client struct {
	*ent.Client
}

// NewClient 创建 Ent 客户端
func NewClient(cfg *config.DatabaseConfig) (*Client, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	client := ent.NewClient(ent.Driver(drv))

	// 运行迁移
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Client{Client: client}, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	return c.Client.Close()
}

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/yourusername/golang/internal/infrastructure/database/ent/migrate"
)

// NewClientFromConfig 从配置创建 Ent 客户端
func NewClientFromConfig(
	ctx context.Context,
	host, port, user, password, dbname, sslmode string,
) (*Client, error) {
	// 构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// 使用 ent.Open 创建客户端
	client, err := Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return client, nil
}

// Migrate 运行数据库迁移
func (c *Client) Migrate(ctx context.Context) error {
	return c.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true))
}

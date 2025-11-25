package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver

	"github.com/yourusername/golang/internal/config"
)

// Client Ent 客户端接口
// 注意：需要先运行 go generate 生成 Ent 客户端代码
// 然后替换 *ent.Client 为实际生成的客户端类型
type Client struct {
	// TODO: 替换为生成的 Ent 客户端
	// 例如: *ent.Client (从生成的代码导入)
	db interface{} // 临时占位
}

// NewClient 创建 Ent 客户端
func NewClient(cfg *config.DatabaseConfig) (*Client, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// TODO: 使用生成的 Ent 客户端
	// client := ent.NewClient(ent.Driver(drv))
	//
	// // 运行迁移
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	return nil, fmt.Errorf("failed to create schema: %w", err)
	// }

	// 临时返回
	return &Client{db: drv}, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	// TODO: 实现关闭逻辑
	if closer, ok := c.db.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

package ent

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// Client Ent 客户端包装器
// 注意：需要先运行 make generate-ent 生成 Ent 代码
// 生成代码后，取消注释下面的导入和实现
type Client struct {
	// TODO: 生成 Ent 代码后，取消注释并导入 gen 包
	// import "github.com/yourusername/golang/internal/infrastructure/database/ent/gen"
	// client *gen.Client
	db *sql.DB
}

// NewClient 创建 Ent 客户端
// 注意：需要先运行 make generate-ent 生成 Ent 代码
func NewClient(dsn string) (*Client, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// TODO: 生成 Ent 代码后，取消注释以下代码
	// drv := entsql.OpenDB(dialect.Postgres, db)
	// genClient := gen.NewClient(gen.Driver(drv))
	// return &Client{client: genClient, db: db}, nil

	// 临时返回，等待生成 Ent 代码
	_ = entsql.OpenDB(dialect.Postgres, db) // 避免未使用导入
	return nil, fmt.Errorf("Ent code not generated. Please run 'make generate-ent' first")
}

// NewClientFromConfig 从配置创建 Ent 客户端
func NewClientFromConfig(ctx context.Context, host, port, user, password, dbname, sslmode string) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	return NewClient(dsn)
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c != nil && c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Client 返回底层 Ent 客户端（生成代码后使用）
// TODO: 生成 Ent 代码后，取消注释
// func (c *Client) Client() *gen.Client {
// 	return c.client
// }

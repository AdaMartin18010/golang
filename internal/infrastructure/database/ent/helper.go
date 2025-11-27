// Package ent 提供 Ent ORM 客户端和辅助函数
//
// 设计原理：
// 1. 提供便捷的客户端创建函数
// 2. 封装数据库迁移逻辑
// 3. 简化 Ent 客户端的使用
//
// 架构位置：
// - 位置：Infrastructure Layer (internal/infrastructure/database/ent/)
// - 职责：数据库访问、ORM 客户端管理
// - 依赖：Ent ORM 库
package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/yourusername/golang/internal/infrastructure/database/ent/migrate"
)

// NewClientFromConfig 从配置创建 Ent 客户端
//
// 设计原理：
// 1. 从配置参数构建 PostgreSQL 连接字符串
// 2. 使用 Ent Open 函数创建客户端
// 3. 返回配置好的 Ent 客户端
//
// 参数：
//   - ctx: 上下文
//   - host: 数据库主机地址
//   - port: 数据库端口
//   - user: 数据库用户名
//   - password: 数据库密码
//   - dbname: 数据库名称
//   - sslmode: SSL 模式（disable、require、verify-full 等）
//
// 返回：
//   - *Client: Ent 客户端实例
//   - error: 连接失败时返回错误
//
// 连接字符串格式：
//   host=%s port=%s user=%s password=%s dbname=%s sslmode=%s
//
// 使用示例：
//   client, err := ent.NewClientFromConfig(
//       ctx,
//       "localhost",
//       "5432",
//       "postgres",
//       "password",
//       "mydb",
//       "disable",
//   )
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer client.Close()
func NewClientFromConfig(
	ctx context.Context,
	host, port, user, password, dbname, sslmode string,
) (*Client, error) {
	// 构建 PostgreSQL 连接字符串（DSN）
	// 格式：host=... port=... user=... password=... dbname=... sslmode=...
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// 使用 ent.Open 创建客户端
	// dialect.Postgres 指定使用 PostgreSQL 驱动
	client, err := Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return client, nil
}

// Migrate 运行数据库迁移
//
// 设计原理：
// 1. 根据 Ent Schema 定义创建或更新数据库表结构
// 2. 支持删除旧索引和列（开发环境）
// 3. 生产环境应该谨慎使用 DropIndex 和 DropColumn
//
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - error: 迁移失败时返回错误
//
// 迁移选项：
// - WithDropIndex(true): 删除旧索引（开发环境）
// - WithDropColumn(true): 删除旧列（开发环境）
//
// 注意事项：
// - 生产环境应该禁用 DropIndex 和 DropColumn
// - 应该先备份数据库再运行迁移
// - 建议在部署前手动运行迁移
//
// 使用示例：
//   client, err := ent.NewClientFromConfig(...)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // 运行迁移
//   if err := client.Migrate(ctx); err != nil {
//       log.Fatalf("Failed to run migrations: %v", err)
//   }
func (c *Client) Migrate(ctx context.Context) error {
	// 运行数据库迁移
	// WithDropIndex(true): 删除旧索引（开发环境使用）
	// WithDropColumn(true): 删除旧列（开发环境使用）
	// 生产环境应该设置为 false，避免数据丢失
	return c.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true))
}

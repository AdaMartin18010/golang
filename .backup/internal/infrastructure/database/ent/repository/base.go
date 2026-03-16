// Package repository 提供基于 Ent 的通用仓储实现
//
// 设计原理：
// 1. 实现框架定义的 Repository 接口
// 2. 使用 Ent ORM 作为底层数据访问技术
// 3. 提供通用的 CRUD 操作实现
// 4. 支持事务和上下文传播
//
// 架构位置：
// - 位置：Infrastructure Layer (internal/infrastructure/database/ent/repository/)
// - 职责：Ent 仓储实现
// - 依赖：Ent ORM、Domain Layer 接口
//
// 使用方式：
// 1. 用户需要为每个实体创建具体的仓储实现
// 2. 继承 BaseRepository 并实现业务特定的查询方法
// 3. 使用 Ent 客户端进行数据访问
//
// 示例：
//   type UserRepository struct {
//       *BaseRepository[*User, *ent.User]
//       client *ent.Client
//   }
//
//   func NewUserRepository(client *ent.Client) *UserRepository {
//       return &UserRepository{
//           BaseRepository: NewBaseRepository(client),
//           client: client,
//       }
//   }
//
//   func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
//       entUser, err := r.client.User.Query().
//           Where(user.EmailEQ(email)).
//           Only(ctx)
//       if err != nil {
//           return nil, err
//       }
//       return r.toDomain(entUser), nil
//   }
package repository

import (
	"context"
	"fmt"
	"strings"

	entclient "github.com/yourusername/golang/internal/infrastructure/database/ent"
	apperrors "github.com/yourusername/golang/pkg/errors"
)

// BaseRepository 基于 Ent 的通用仓储实现
//
// 设计原理：
// 1. 使用泛型 T 表示领域实体类型
// 2. 使用泛型 E 表示 Ent 实体类型（实现 ent.Entity 接口）
// 3. 提供通用的 CRUD 操作实现
// 4. 需要用户实现实体转换方法（toDomain、toEnt）
//
// 类型参数：
//   - T: 领域实体类型（Domain Entity）
//   - E: Ent 实体类型（Ent Entity，需要实现 ent.Entity 接口）
//
// 注意：
// - 这是一个基础实现，用户需要继承并实现业务特定的方法
// - 用户需要实现实体转换逻辑（toDomain、toEnt）
// - 用户需要实现 ID 获取和设置方法
type BaseRepository[T any, E interface{}] struct {
	client *entclient.Client
	// toDomain 将 Ent 实体转换为领域实体
	toDomain func(E) (*T, error)
	// toEnt 将领域实体转换为 Ent 实体
	toEnt func(*T) (E, error)
	// getID 从领域实体获取 ID
	getID func(*T) (string, error)
	// setID 设置领域实体的 ID
	setID func(*T, string) error
}

// NewBaseRepository 创建基础仓储实例
//
// 参数：
//   - client: Ent 客户端
//   - toDomain: Ent 实体到领域实体的转换函数
//   - toEnt: 领域实体到 Ent 实体的转换函数
//   - getID: 从领域实体获取 ID 的函数
//   - setID: 设置领域实体 ID 的函数
//
// 返回：
//   - *BaseRepository: 基础仓储实例
func NewBaseRepository[T any, E interface{}](
	client *entclient.Client,
	toDomain func(E) (*T, error),
	toEnt func(*T) (E, error),
	getID func(*T) (string, error),
	setID func(*T, string) error,
) *BaseRepository[T, E] {
	return &BaseRepository[T, E]{
		client:   client,
		toDomain: toDomain,
		toEnt:    toEnt,
		getID:    getID,
		setID:    setID,
	}
}

// Create 创建实体
//
// 实现说明：
// 1. 将领域实体转换为 Ent 实体
// 2. 使用 Ent 客户端创建实体
// 3. 将创建的 Ent 实体转换回领域实体
// 4. 设置实体的 ID 和时间戳
//
// 注意：这是一个基础实现，用户需要在具体的仓储中重写此方法
func (r *BaseRepository[T, E]) Create(ctx context.Context, entity *T) error {
	// 注意：这里需要根据具体的 Ent 实体类型调用相应的创建方法
	// 用户需要在具体的仓储实现中重写此方法
	_ = entity // 避免未使用变量警告
	return fmt.Errorf("Create method must be implemented in concrete repository")
}

// FindByID 根据 ID 查找实体
//
// 实现说明：
// 1. 使用 Ent 客户端查询实体
// 2. 将 Ent 实体转换为领域实体
// 3. 如果实体不存在，返回 nil 和 nil error
func (r *BaseRepository[T, E]) FindByID(ctx context.Context, id string) (*T, error) {
	// 注意：这里需要根据具体的 Ent 实体类型调用相应的查询方法
	// 用户需要在具体的仓储实现中重写此方法
	return nil, fmt.Errorf("FindByID method must be implemented in concrete repository")
}

// Update 更新实体
//
// 实现说明：
// 1. 获取实体的 ID
// 2. 使用 Ent 客户端更新实体
// 3. 更新实体的 UpdatedAt 时间戳
func (r *BaseRepository[T, E]) Update(ctx context.Context, entity *T) error {
	id, err := r.getID(entity)
	if err != nil {
		return fmt.Errorf("failed to get entity ID: %w", err)
	}

	// 注意：这里需要根据具体的 Ent 实体类型调用相应的更新方法
	// 用户需要在具体的仓储实现中重写此方法
	_ = id // 避免未使用变量警告
	return fmt.Errorf("Update method must be implemented in concrete repository")
}

// Delete 删除实体
//
// 实现说明：
// 1. 使用 Ent 客户端删除实体
// 2. 支持硬删除和软删除（根据业务需求）
func (r *BaseRepository[T, E]) Delete(ctx context.Context, id string) error {
	// 注意：这里需要根据具体的 Ent 实体类型调用相应的删除方法
	// 用户需要在具体的仓储实现中重写此方法
	return fmt.Errorf("Delete method must be implemented in concrete repository")
}

// List 列出实体（支持分页）
//
// 实现说明：
// 1. 使用 Ent 客户端查询实体列表
// 2. 应用分页限制和偏移量
// 3. 将 Ent 实体列表转换为领域实体列表
func (r *BaseRepository[T, E]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	// 注意：这里需要根据具体的 Ent 实体类型调用相应的查询方法
	// 用户需要在具体的仓储实现中重写此方法
	return nil, fmt.Errorf("List method must be implemented in concrete repository")
}

// Client 获取 Ent 客户端
//
// 返回：
//   - *entclient.Client: Ent 客户端实例
//
// 使用场景：
// - 在具体的仓储实现中使用客户端进行复杂查询
// - 在事务中使用客户端
func (r *BaseRepository[T, E]) Client() *entclient.Client {
	return r.client
}

// WithTx 在事务中执行操作
//
// 参数：
//   - ctx: 上下文
//   - fn: 在事务中执行的函数
//
// 返回：
//   - error: 执行失败时返回错误
//
// 使用示例：
//   err := repo.WithTx(ctx, func(tx *entclient.Tx) error {
//       // 在事务中执行操作
//       return nil
//   })
func (r *BaseRepository[T, E]) WithTx(ctx context.Context, fn func(*entclient.Tx) error) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// handleEntError 处理 Ent 错误
//
// 功能：
// - 将 Ent 错误转换为框架错误
// - 处理常见的数据库错误（如唯一约束冲突、记录不存在等）
func handleEntError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否是记录不存在错误
	// 注意：需要使用具体的 ent 包中的 IsNotFound 函数
	// 这里使用字符串匹配作为临时方案
	errStr := err.Error()
	if contains(errStr, "not found") || contains(errStr, "no rows") {
		return apperrors.NewNotFoundError("entity", "")
	}

	// 检查是否是唯一约束冲突
	// 注意：需要使用具体的 sql 包中的 IsConstraintError 函数
	// 这里使用字符串匹配作为临时方案
	if contains(errStr, "unique constraint") || contains(errStr, "duplicate key") {
		return apperrors.NewConflictError("entity already exists")
	}

	// 其他错误
	return apperrors.NewInternalError(fmt.Sprintf("database error: %v", err), err)
}

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// 导入必要的包
// 注意：由于泛型限制，这里无法直接验证接口实现
// 用户需要在具体的仓储实现中确保实现了 Repository 接口

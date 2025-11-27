package user

import (
	"context"

	"github.com/yourusername/golang/internal/domain/interfaces"
)

// Repository 是用户仓储接口。
//
// 设计原理：
// 1. 仓储接口定义在领域层，实现放在基础设施层
// 2. 符合依赖倒置原则（Dependency Inversion Principle）
// 3. 领域层不依赖具体的数据访问技术
//
// 架构位置：
// - 接口定义：Domain Layer (internal/domain/user/)
// - 接口实现：Infrastructure Layer (internal/infrastructure/database/)
//
// 使用场景：
// 1. 应用层通过仓储接口访问用户实体
// 2. 基础设施层实现仓储接口，封装数据访问细节
// 3. 测试时可以使用 Mock 仓储实现
//
// 示例：
//
//	// 基础设施层实现
//	type EntUserRepository struct {
//	    client *ent.Client
//	}
//
//	func (r *EntUserRepository) Create(ctx context.Context, user *User) error {
//	    // 使用 Ent 实现
//	}
//
//	// 应用层使用
//	type UserService struct {
//	    repo Repository  // 依赖接口，不依赖具体实现
//	}
type Repository interface {
	// 继承通用仓储接口
	// 注意：使用 User 而不是 *User，因为接口中已经是指针类型
	interfaces.Repository[User]

	// FindByEmail 根据邮箱查找用户。
	//
	// 功能说明：
	// - 根据邮箱查找用户
	// - 邮箱应该是唯一的
	//
	// 参数：
	// - ctx: 上下文
	// - email: 用户邮箱
	//
	// 返回：
	// - *User: 找到的用户，如果不存在返回 nil
	// - error: 查找失败时返回错误
	//
	// 业务规则：
	// - 如果用户不存在，应该返回 nil 和 nil error（不是错误）
	// - 如果用户被软删除，应该根据业务规则决定是否返回
	//
	// 实现要求：
	// - 应该使用索引优化查询性能
	// - 应该处理查询超时
	// - 应该处理数据库连接错误
	FindByEmail(ctx context.Context, email string) (*User, error)
}

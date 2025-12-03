package interfaces

import "context"

// Specification 规约模式接口
//
// 设计原理：
// 1. 规约模式（Specification Pattern）用于封装业务规则和查询条件
// 2. 将复杂的查询逻辑从仓储中分离出来，提高可复用性和可测试性
// 3. 支持组合多个规约，构建复杂的查询条件
// 4. 符合开闭原则：对扩展开放，对修改关闭
//
// 架构位置：
// - 接口定义：Domain Layer (internal/domain/interfaces/)
// - 具体规约：Domain Layer (internal/domain/*/specifications/)
// - 规约实现：Infrastructure Layer (转换为数据库查询)
//
// 使用场景：
// 1. 复杂的业务查询条件
// 2. 需要组合多个条件的场景
// 3. 需要复用查询逻辑的场景
// 4. 需要在内存中过滤实体的场景
//
// 示例：
//   // 定义规约
//   type ActiveUserSpec struct{}
//   func (s ActiveUserSpec) IsSatisfiedBy(user *User) bool {
//       return user.Status == "active"
//   }
//
//   type EmailSpec struct {
//       Email string
//   }
//   func (s EmailSpec) IsSatisfiedBy(user *User) bool {
//       return user.Email == s.Email
//   }
//
//   // 组合规约
//   spec := And(ActiveUserSpec{}, EmailSpec{Email: "test@example.com"})
//
//   // 使用规约查询
//   users, err := repo.FindBySpecification(ctx, spec)
//
// 参考：
// - Martin Fowler: Specification Pattern
// - Eric Evans: Domain-Driven Design
type Specification[T any] interface {
	// IsSatisfiedBy 检查实体是否满足规约
	//
	// 参数：
	//   - entity: 要检查的实体
	//
	// 返回：
	//   - bool: 是否满足规约
	//
	// 用途：
	// 1. 在内存中过滤实体
	// 2. 验证业务规则
	// 3. 单元测试规约逻辑
	IsSatisfiedBy(entity *T) bool
}

// CompositeSpecification 组合规约接口
// 支持 And、Or、Not 逻辑组合
type CompositeSpecification[T any] interface {
	Specification[T]

	// And 与操作
	And(other Specification[T]) Specification[T]

	// Or 或操作
	Or(other Specification[T]) Specification[T]

	// Not 非操作
	Not() Specification[T]
}

// AndSpecification And 规约实现
type AndSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// NewAndSpecification 创建 And 规约
func NewAndSpecification[T any](left, right Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: left, right: right}
}

func (s *AndSpecification[T]) IsSatisfiedBy(entity *T) bool {
	return s.left.IsSatisfiedBy(entity) && s.right.IsSatisfiedBy(entity)
}

func (s *AndSpecification[T]) And(other Specification[T]) Specification[T] {
	return NewAndSpecification[T](s, other)
}

func (s *AndSpecification[T]) Or(other Specification[T]) Specification[T] {
	return NewOrSpecification[T](s, other)
}

func (s *AndSpecification[T]) Not() Specification[T] {
	return NewNotSpecification[T](s)
}

// OrSpecification Or 规约实现
type OrSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// NewOrSpecification 创建 Or 规约
func NewOrSpecification[T any](left, right Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: left, right: right}
}

func (s *OrSpecification[T]) IsSatisfiedBy(entity *T) bool {
	return s.left.IsSatisfiedBy(entity) || s.right.IsSatisfiedBy(entity)
}

func (s *OrSpecification[T]) And(other Specification[T]) Specification[T] {
	return NewAndSpecification[T](s, other)
}

func (s *OrSpecification[T]) Or(other Specification[T]) Specification[T] {
	return NewOrSpecification[T](s, other)
}

func (s *OrSpecification[T]) Not() Specification[T] {
	return NewNotSpecification[T](s)
}

// NotSpecification Not 规约实现
type NotSpecification[T any] struct {
	spec Specification[T]
}

// NewNotSpecification 创建 Not 规约
func NewNotSpecification[T any](spec Specification[T]) Specification[T] {
	return &NotSpecification[T]{spec: spec}
}

func (s *NotSpecification[T]) IsSatisfiedBy(entity *T) bool {
	return !s.spec.IsSatisfiedBy(entity)
}

func (s *NotSpecification[T]) And(other Specification[T]) Specification[T] {
	return NewAndSpecification[T](s, other)
}

func (s *NotSpecification[T]) Or(other Specification[T]) Specification[T] {
	return NewOrSpecification[T](s, other)
}

func (s *NotSpecification[T]) Not() Specification[T] {
	return s.spec // 双重否定
}

// RepositoryWithSpecification 支持规约的仓储接口
//
// 扩展了基础 Repository 接口，添加了基于规约的查询方法
type RepositoryWithSpecification[T any] interface {
	Repository[T]

	// FindBySpecification 根据规约查找实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - spec: 规约条件
	//
	// 返回：
	//   - []*T: 满足规约的实体列表
	//   - error: 查询失败时返回错误
	//
	// 实现要求：
	// 1. 基础设施层需要将规约转换为数据库查询
	// 2. 如果无法转换为数据库查询，可以先查询所有数据再在内存中过滤
	// 3. 应该考虑性能，避免全表扫描
	FindBySpecification(ctx context.Context, spec Specification[T]) ([]*T, error)

	// CountBySpecification 根据规约统计实体数量
	//
	// 参数：
	//   - ctx: 上下文
	//   - spec: 规约条件
	//
	// 返回：
	//   - int: 满足规约的实体数量
	//   - error: 统计失败时返回错误
	CountBySpecification(ctx context.Context, spec Specification[T]) (int, error)
}

// 辅助函数：简化规约组合

// And 创建 And 规约
func And[T any](left, right Specification[T]) Specification[T] {
	return NewAndSpecification[T](left, right)
}

// Or 创建 Or 规约
func Or[T any](left, right Specification[T]) Specification[T] {
	return NewOrSpecification[T](left, right)
}

// Not 创建 Not 规约
func Not[T any](spec Specification[T]) Specification[T] {
	return NewNotSpecification[T](spec)
}


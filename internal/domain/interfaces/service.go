package interfaces

import "context"

// DomainService 通用领域服务接口（框架抽象）
// 用户需要根据业务需求定义具体的领域服务接口
type DomainService interface {
	// Validate 验证实体
	Validate(ctx context.Context, entity interface{}) error

	// Process 处理业务逻辑
	Process(ctx context.Context, entity interface{}) error
}

// Validator 验证器接口
type Validator interface {
	// Validate 验证实体
	Validate(ctx context.Context, entity interface{}) error
}

// Processor 处理器接口
type Processor interface {
	// Process 处理业务逻辑
	Process(ctx context.Context, entity interface{}) error
}

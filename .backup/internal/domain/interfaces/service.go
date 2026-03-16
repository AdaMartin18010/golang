package interfaces

import "context"

// DomainService 通用领域服务接口（框架抽象）
//
// 设计原理：
// 1. 领域服务（Domain Service）用于处理不属于单个实体的业务逻辑
// 2. 领域服务是无状态的，不保存业务状态
// 3. 领域服务处理跨聚合的操作或复杂的业务规则
// 4. 领域服务接口定义在领域层，实现可以在领域层或基础设施层
//
// 架构位置：
// - 接口定义：Domain Layer (internal/domain/interfaces/)
// - 服务实现：Domain Layer 或 Infrastructure Layer
//
// 使用场景：
// 1. 跨聚合的业务逻辑（如转账涉及两个账户）
// 2. 复杂的业务规则计算（如价格计算、折扣计算）
// 3. 需要多个仓储协调的操作
// 4. 不适合放在实体中的业务逻辑
//
// 与实体方法的区别：
// - 实体方法：处理单个实体的业务逻辑
// - 领域服务：处理跨实体或复杂的业务逻辑
//
// 示例：
//   // 领域服务接口
//   type PricingService interface {
//       CalculatePrice(ctx context.Context, product *Product, quantity int) (*Money, error)
//   }
//
//   // 领域服务实现
//   type pricingService struct {
//       discountRepo DiscountRepository
//   }
//
//   func (s *pricingService) CalculatePrice(ctx context.Context, product *Product, quantity int) (*Money, error) {
//       // 计算价格逻辑，可能涉及折扣、促销等
//   }
//
// 注意事项：
// - 领域服务应该是无状态的
// - 应该优先考虑将业务逻辑放在实体中
// - 只有在业务逻辑不适合放在实体中时才使用领域服务
//
// 用户需要根据业务需求定义具体的领域服务接口，例如：
//   type EmailValidationService interface {
//       ValidateEmail(ctx context.Context, email string) error
//   }
type DomainService interface {
	// Validate 验证实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - entity: 要验证的实体
	//
	// 返回：
	//   - error: 验证失败时返回错误
	//
	// 业务规则：
	//   - 验证应该包括业务规则验证，不仅仅是格式验证
	//   - 验证应该返回明确的错误信息
	//   - 验证应该是幂等的
	//
	// 实现要求：
	//   - 应该验证实体的所有业务规则
	//   - 应该返回详细的验证错误
	//   - 应该考虑性能，避免重复验证
	Validate(ctx context.Context, entity interface{}) error

	// Process 处理业务逻辑
	//
	// 参数：
	//   - ctx: 上下文
	//   - entity: 要处理的实体
	//
	// 返回：
	//   - error: 处理失败时返回错误
	//
	// 业务规则：
	//   - 处理应该符合业务规则
	//   - 处理应该是幂等的（如果可能）
	//   - 处理应该记录操作日志（可选）
	//
	// 实现要求：
	//   - 应该处理业务异常
	//   - 应该考虑事务管理
	//   - 应该考虑并发安全
	Process(ctx context.Context, entity interface{}) error
}

// Validator 验证器接口
//
// 设计原理：
// 1. 验证器专门用于验证业务规则
// 2. 验证器可以是领域服务的一部分，也可以是独立的服务
// 3. 验证器应该专注于验证逻辑，不包含业务处理逻辑
//
// 使用场景：
// 1. 复杂的验证逻辑（如邮箱唯一性验证）
// 2. 需要访问外部资源的验证（如调用外部 API 验证）
// 3. 跨实体的验证（如订单金额验证需要访问产品信息）
//
// 示例：
//   type EmailValidator interface {
//       Validate(ctx context.Context, email string) error
//   }
//
//   type emailValidator struct {
//       userRepo UserRepository
//   }
//
//   func (v *emailValidator) Validate(ctx context.Context, email string) error {
//       // 验证邮箱格式
//       if !isValidEmailFormat(email) {
//           return ErrInvalidEmailFormat
//       }
//
//       // 验证邮箱唯一性
//       existing, err := v.userRepo.FindByEmail(ctx, email)
//       if err == nil && existing != nil {
//           return ErrEmailAlreadyExists
//       }
//
//       return nil
//   }
type Validator interface {
	// Validate 验证实体
	//
	// 参数：
	//   - ctx: 上下文
	//   - entity: 要验证的实体
	//
	// 返回：
	//   - error: 验证失败时返回错误
	//
	// 业务规则：
	//   - 验证应该包括所有相关的业务规则
	//   - 验证错误应该明确，便于用户理解
	//   - 验证应该是幂等的
	//
	// 实现要求：
	//   - 应该返回详细的验证错误
	//   - 应该考虑性能
	//   - 应该处理验证过程中的异常
	Validate(ctx context.Context, entity interface{}) error
}

// Processor 处理器接口
//
// 设计原理：
// 1. 处理器专门用于处理业务逻辑
// 2. 处理器可以是领域服务的一部分，也可以是独立的服务
// 3. 处理器应该专注于业务处理，不包含验证逻辑（验证应该提前完成）
//
// 使用场景：
// 1. 复杂的业务处理逻辑（如订单处理、支付处理）
// 2. 需要多个步骤的业务流程
// 3. 需要协调多个实体的操作
//
// 示例：
//   type OrderProcessor interface {
//       Process(ctx context.Context, order *Order) error
//   }
//
//   type orderProcessor struct {
//       orderRepo    OrderRepository
//       productRepo  ProductRepository
//       paymentService PaymentService
//   }
//
//   func (p *orderProcessor) Process(ctx context.Context, order *Order) error {
//       // 1. 验证订单
//       if err := order.Validate(); err != nil {
//           return err
//       }
//
//       // 2. 处理支付
//       if err := p.paymentService.ProcessPayment(ctx, order); err != nil {
//           return err
//       }
//
//       // 3. 更新库存
//       for _, item := range order.Items {
//           if err := p.productRepo.DecreaseStock(ctx, item.ProductID, item.Quantity); err != nil {
//               return err
//           }
//       }
//
//       // 4. 保存订单
//       return p.orderRepo.Update(ctx, order)
//   }
type Processor interface {
	// Process 处理业务逻辑
	//
	// 参数：
	//   - ctx: 上下文
	//   - entity: 要处理的实体
	//
	// 返回：
	//   - error: 处理失败时返回错误
	//
	// 业务规则：
	//   - 处理应该符合业务规则
	//   - 处理应该是事务性的（如果涉及多个操作）
	//   - 处理应该记录操作日志（可选）
	//
	// 实现要求：
	//   - 应该处理业务异常
	//   - 应该考虑事务管理
	//   - 应该考虑并发安全
	//   - 应该考虑幂等性
	Process(ctx context.Context, entity interface{}) error
}

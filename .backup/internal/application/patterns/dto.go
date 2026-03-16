package patterns

import "time"

// DTO 数据传输对象基类（框架抽象）
//
// 设计原理：
// 1. DTO（Data Transfer Object）用于在不同层之间传输数据
// 2. DTO 不包含业务逻辑，只包含数据
// 3. DTO 用于隔离领域模型和外部接口
// 4. DTO 可以包含验证规则和序列化标签
//
// 架构位置：
// - 定义：Application Layer (internal/application/patterns/)
// - 使用：Application Layer 和 Interfaces Layer
//
// DTO 的作用：
// 1. 隔离领域模型：领域模型不直接暴露给外部
// 2. 版本控制：可以独立演进 DTO，不影响领域模型
// 3. 性能优化：可以只传输必要的数据
// 4. 安全性：可以隐藏敏感信息
//
// 使用场景：
// 1. API 请求/响应
// 2. 跨服务通信
// 3. 数据序列化/反序列化
//
// 示例：
//   // 定义 DTO
//   type UserDTO struct {
//       patterns.DTO
//       Email string `json:"email" validate:"required,email"`
//       Name  string `json:"name" validate:"required,min=2,max=50"`
//   }
//
//   // 从领域对象转换为 DTO
//   func toUserDTO(user *domain.User) *UserDTO {
//       return &UserDTO{
//           DTO: patterns.DTO{
//               ID:        user.ID,
//               CreatedAt: user.CreatedAt,
//               UpdatedAt: user.UpdatedAt,
//           },
//           Email: user.Email,
//           Name:  user.Name,
//       }
//   }
//
//   // 从 DTO 转换为领域对象
//   func (dto *UserDTO) ToDomain() *domain.User {
//       return &domain.User{
//           ID:        dto.ID,
//           Email:     dto.Email,
//           Name:      dto.Name,
//           CreatedAt: dto.CreatedAt,
//           UpdatedAt: dto.UpdatedAt,
//       }
//   }
//
// 注意事项：
// - DTO 应该只包含必要的数据
// - DTO 应该包含验证规则
// - DTO 应该使用 JSON 标签
// - DTO 不应该包含业务逻辑
//
// 与领域对象的区别：
// - DTO：数据传输对象，不包含业务逻辑
// - 领域对象：包含业务逻辑和业务规则
type DTO struct {
	// ID 唯一标识
	ID string `json:"id"`

	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// ToDTO 转换函数接口（框架抽象）
//
// 设计原理：
// 1. 提供统一的转换接口
// 2. 使用泛型 T 表示目标 DTO 类型
// 3. 支持从领域对象转换为 DTO
//
// 使用场景：
// 1. 领域对象转 DTO
// 2. 统一转换逻辑
//
// 示例：
//   // 领域对象实现接口
//   func (u *User) ToDTO() *UserDTO {
//       return &UserDTO{
//           DTO: patterns.DTO{
//               ID:        u.ID,
//               CreatedAt: u.CreatedAt,
//               UpdatedAt: u.UpdatedAt,
//           },
//           Email: u.Email,
//           Name:  u.Name,
//       }
//   }
//
//   // 使用
//   user := domain.NewUser("test@example.com", "Test User")
//   dto := user.ToDTO()
//
// 注意事项：
// - 转换逻辑应该简单直接
// - 应该只转换必要的数据
// - 应该处理 nil 值
type ToDTO[T any] interface {
	// ToDTO 转换为 DTO
	//
	// 返回：
	//   - T: 转换后的 DTO，类型为 T
	//
	// 实现要求：
	//   - 应该返回非 nil 的 DTO
	//   - 应该处理所有必要的字段
	//   - 应该处理时间格式等细节
	ToDTO() T
}

// FromDTO 从 DTO 转换函数接口（框架抽象）
//
// 设计原理：
// 1. 提供统一的转换接口
// 2. 使用泛型 T 表示源 DTO 类型
// 3. 支持从 DTO 转换为领域对象
//
// 使用场景：
// 1. DTO 转领域对象
// 2. 统一转换逻辑
//
// 示例：
//   // 领域对象实现接口
//   func (u *User) FromDTO(dto *UserDTO) error {
//       u.ID = dto.ID
//       u.Email = dto.Email
//       u.Name = dto.Name
//       u.CreatedAt = dto.CreatedAt
//       u.UpdatedAt = dto.UpdatedAt
//
//       // 验证
//       return u.Validate()
//   }
//
//   // 使用
//   dto := &UserDTO{Email: "test@example.com", Name: "Test User"}
//   user := &domain.User{}
//   if err := user.FromDTO(dto); err != nil {
//       return err
//   }
//
// 注意事项：
// - 转换时应该验证数据
// - 应该处理必填字段
// - 应该处理时间格式等细节
type FromDTO[T any] interface {
	// FromDTO 从 DTO 转换
	//
	// 参数：
	//   - dto: 源 DTO，类型为 T
	//
	// 返回：
	//   - error: 转换失败时返回错误
	//
	// 实现要求：
	//   - 应该验证 DTO 数据
	//   - 应该处理所有必要的字段
	//   - 应该返回明确的错误信息
	FromDTO(dto T) error
}

// PaginatedDTO 分页 DTO
//
// 设计原理：
// 1. 提供标准化的分页响应格式
// 2. 使用泛型 T 表示数据项类型
// 3. 包含分页信息和数据列表
//
// 使用场景：
// 1. 列表查询响应
// 2. 搜索查询响应
// 3. 分页数据返回
//
// 示例：
//   // 定义分页 DTO
//   type UserListDTO struct {
//       patterns.PaginatedDTO[*UserDTO]
//   }
//
//   // 使用
//   users := []*domain.User{...}
//   dtos := make([]*UserDTO, len(users))
//   for i, user := range users {
//       dtos[i] = toUserDTO(user)
//   }
//
//   result := &UserListDTO{
//       PaginatedDTO: patterns.PaginatedDTO[*UserDTO]{
//           Data:  dtos,
//           Total: total,
//           Page:  page,
//           Size:  size,
//       },
//   }
//
// 注意事项：
// - Data 应该是切片，即使只有一条数据
// - Total 表示总记录数，不是当前页的数量
// - Page 和 Size 应该与请求参数一致
type PaginatedDTO[T any] struct {
	// Data 数据列表
	// 注意：即使只有一条数据，也应该是切片
	Data []T `json:"data"`

	// Total 总记录数（用于分页）
	// 注意：这是所有符合条件的记录数，不是当前页的数量
	Total int `json:"total"`

	// Page 当前页码（从 1 开始）
	Page int `json:"page"`

	// Size 每页数量
	Size int `json:"size"`
}

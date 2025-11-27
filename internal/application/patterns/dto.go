package patterns

import "time"

// DTO 数据传输对象基类（框架抽象）
type DTO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToDTO 转换函数接口（框架抽象）
type ToDTO[T any] interface {
	// ToDTO 转换为 DTO
	ToDTO() T
}

// FromDTO 从 DTO 转换函数接口（框架抽象）
type FromDTO[T any] interface {
	// FromDTO 从 DTO 转换
	FromDTO(dto T) error
}

// PaginatedDTO 分页 DTO
type PaginatedDTO[T any] struct {
	Data  []T   `json:"data"`
	Total int   `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

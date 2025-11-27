package patterns

import "context"

// Query 查询接口（框架抽象）
// 用于处理读操作（Read、List、Search）
type Query interface {
	// Execute 执行查询
	Execute(ctx context.Context) (interface{}, error)
}

// QueryHandler 查询处理器接口（框架抽象）
type QueryHandler[T Query, R any] interface {
	// Handle 处理查询
	Handle(ctx context.Context, query T) (R, error)
}

// QueryResult 查询结果
type QueryResult[T any] struct {
	Data  []T
	Total int
	Page  int
	Size  int
}

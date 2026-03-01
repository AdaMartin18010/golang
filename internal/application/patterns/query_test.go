package patterns

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockQuery 测试用的 Mock Query
type MockQuery struct {
	ID string
}

func (m MockQuery) Execute(ctx context.Context) (interface{}, error) {
	if m.ID == "" {
		return nil, errors.New("ID不能为空")
	}
	return map[string]string{"id": m.ID}, nil
}

// MockQueryResult 测试用的结果类型
type MockQueryResult struct {
	ID   string
	Name string
}

// MockQueryHandler 测试用的 Query Handler
type MockQueryHandler struct {
	returnResult *MockQueryResult
	returnError  error
}

func (m *MockQueryHandler) Handle(ctx context.Context, query MockQuery) (*MockQueryResult, error) {
	if query.ID == "" {
		return nil, errors.New("查询ID不能为空")
	}
	return m.returnResult, m.returnError
}

// ==================== QueryResult 测试 ====================

// TestQueryResult 测试 QueryResult 结构
func TestQueryResult(t *testing.T) {
	tests := []struct {
		name        string
		result      QueryResult[string]
		expectData  []string
		expectTotal int
		expectPage  int
		expectSize  int
	}{
		{
			name: "基本分页结果",
			result: QueryResult[string]{
				Data:  []string{"item1", "item2", "item3"},
				Total: 100,
				Page:  1,
				Size:  10,
			},
			expectData:  []string{"item1", "item2", "item3"},
			expectTotal: 100,
			expectPage:  1,
			expectSize:  10,
		},
		{
			name: "空结果",
			result: QueryResult[string]{
				Data:  []string{},
				Total: 0,
				Page:  1,
				Size:  10,
			},
			expectData:  []string{},
			expectTotal: 0,
			expectPage:  1,
			expectSize:  10,
		},
		{
			name: "第二页结果",
			result: QueryResult[string]{
				Data:  []string{"item11", "item12"},
				Total: 12,
				Page:  2,
				Size:  10,
			},
			expectData:  []string{"item11", "item12"},
			expectTotal: 12,
			expectPage:  2,
			expectSize:  10,
		},
		{
			name: "最后一页",
			result: QueryResult[string]{
				Data:  []string{"item91", "item92", "item93"},
				Total: 93,
				Page:  10,
				Size:  10,
			},
			expectData:  []string{"item91", "item92", "item93"},
			expectTotal: 93,
			expectPage:  10,
			expectSize:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectData, tt.result.Data)
			assert.Equal(t, tt.expectTotal, tt.result.Total)
			assert.Equal(t, tt.expectPage, tt.result.Page)
			assert.Equal(t, tt.expectSize, tt.result.Size)
		})
	}
}

// TestQueryResultWithStruct 测试包含结构体的 QueryResult
func TestQueryResultWithStruct(t *testing.T) {
	type UserDTO struct {
		ID    string
		Email string
		Name  string
	}

	users := []UserDTO{
		{ID: "1", Email: "user1@example.com", Name: "User 1"},
		{ID: "2", Email: "user2@example.com", Name: "User 2"},
		{ID: "3", Email: "user3@example.com", Name: "User 3"},
	}

	result := QueryResult[UserDTO]{
		Data:  users,
		Total: 50,
		Page:  1,
		Size:  10,
	}

	assert.Len(t, result.Data, 3)
	assert.Equal(t, 50, result.Total)
	assert.Equal(t, "user1@example.com", result.Data[0].Email)
}

// TestQueryResultWithInt 测试整数类型的 QueryResult
func TestQueryResultWithInt(t *testing.T) {
	result := QueryResult[int]{
		Data:  []int{1, 2, 3, 4, 5},
		Total: 100,
		Page:  1,
		Size:  5,
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, result.Data)
	assert.Equal(t, 100, result.Total)
}

// TestQueryResultWithPointer 测试指针类型的 QueryResult
func TestQueryResultWithPointer(t *testing.T) {
	type Item struct {
		Name  string
		Value int
	}

	item1 := &Item{Name: "item1", Value: 10}
	item2 := &Item{Name: "item2", Value: 20}

	result := QueryResult[*Item]{
		Data:  []*Item{item1, item2},
		Total: 2,
		Page:  1,
		Size:  10,
	}

	assert.Len(t, result.Data, 2)
	assert.Equal(t, "item1", result.Data[0].Name)
	assert.Equal(t, 20, result.Data[1].Value)
}

// TestQueryResultEmptyData 测试空数据的 QueryResult
func TestQueryResultEmptyData(t *testing.T) {
	result := QueryResult[string]{
		Data:  nil,
		Total: 0,
		Page:  1,
		Size:  10,
	}

	assert.Nil(t, result.Data)
	assert.Equal(t, 0, result.Total)
}

// ==================== MockQuery 测试 ====================

// TestMockQuery_Execute 测试 Query 执行
func TestMockQuery_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("有效ID", func(t *testing.T) {
		query := MockQuery{ID: "123"}
		result, err := query.Execute(ctx)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"id": "123"}, result)
	})

	t.Run("空ID", func(t *testing.T) {
		query := MockQuery{ID: ""}
		result, err := query.Execute(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ID不能为空")
	})

	t.Run("特殊字符ID", func(t *testing.T) {
		query := MockQuery{ID: "id-with-special-chars_123"}
		result, err := query.Execute(ctx)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"id": "id-with-special-chars_123"}, result)
	})
}

// TestMockQuery_ExecuteWithContext 测试带上下文的 Query 执行
func TestMockQuery_ExecuteWithContext(t *testing.T) {
	t.Run("正常上下文", func(t *testing.T) {
		ctx := context.Background()
		query := MockQuery{ID: "123"}
		result, err := query.Execute(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("取消的上下文", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消
		query := MockQuery{ID: "123"}
		result, err := query.Execute(ctx)
		// MockQuery 没有检查上下文，所以应该成功
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("带值的上下文", func(t *testing.T) {
		type contextKey string
		key := contextKey("user-id")
		ctx := context.WithValue(context.Background(), key, "user-123")
		query := MockQuery{ID: "123"}
		result, err := query.Execute(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// ==================== MockQueryHandler 测试 ====================

// TestMockQueryHandler_Handle 测试 Query Handler
func TestMockQueryHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("查询成功", func(t *testing.T) {
		expectedResult := &MockQueryResult{ID: "123", Name: "Test"}
		handler := &MockQueryHandler{returnResult: expectedResult, returnError: nil}
		query := MockQuery{ID: "123"}

		result, err := handler.Handle(ctx, query)

		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("查询失败", func(t *testing.T) {
		handler := &MockQueryHandler{returnResult: nil, returnError: errors.New("查询失败")}
		query := MockQuery{ID: "123"}

		result, err := handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "查询失败")
	})

	t.Run("空ID验证", func(t *testing.T) {
		handler := &MockQueryHandler{returnResult: nil, returnError: nil}
		query := MockQuery{ID: ""}

		result, err := handler.Handle(ctx, query)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "查询ID不能为空")
	})

	t.Run("返回nil结果但无错误", func(t *testing.T) {
		handler := &MockQueryHandler{returnResult: nil, returnError: nil}
		query := MockQuery{ID: "123"}

		result, err := handler.Handle(ctx, query)

		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

// TestQueryHandlerInterface 测试 QueryHandler 接口
func TestQueryHandlerInterface(t *testing.T) {
	t.Run("接口实现", func(t *testing.T) {
		// 验证 MockQueryHandler 实现了 QueryHandler 接口
		var _ QueryHandler[MockQuery, *MockQueryResult] = &MockQueryHandler{}
	})

	t.Run("泛型处理", func(t *testing.T) {
		handler := &MockQueryHandler{
			returnResult: &MockQueryResult{ID: "123", Name: "Test"},
			returnError:  nil,
		}
		ctx := context.Background()
		query := MockQuery{ID: "123"}

		var h QueryHandler[MockQuery, *MockQueryResult] = handler
		result, err := h.Handle(ctx, query)

		require.NoError(t, err)
		assert.Equal(t, "123", result.ID)
	})
}

// ==================== 分页逻辑测试 ====================

// TestQueryResultPagination 测试分页逻辑
func TestQueryResultPagination(t *testing.T) {
	tests := []struct {
		name          string
		total         int
		page          int
		size          int
		expectPages   int
		expectHasNext bool
		expectHasPrev bool
	}{
		{
			name:          "第一页",
			total:         100,
			page:          1,
			size:          10,
			expectPages:   10,
			expectHasNext: true,
			expectHasPrev: false,
		},
		{
			name:          "中间页",
			total:         100,
			page:          5,
			size:          10,
			expectPages:   10,
			expectHasNext: true,
			expectHasPrev: true,
		},
		{
			name:          "最后一页",
			total:         100,
			page:          10,
			size:          10,
			expectPages:   10,
			expectHasNext: false,
			expectHasPrev: true,
		},
		{
			name:          "少于每页数量",
			total:         5,
			page:          1,
			size:          10,
			expectPages:   1,
			expectHasNext: false,
			expectHasPrev: false,
		},
		{
			name:          "非整页总数",
			total:         95,
			page:          10,
			size:          10,
			expectPages:   10,
			expectHasNext: false,
			expectHasPrev: true,
		},
		{
			name:          "零条记录",
			total:         0,
			page:          1,
			size:          10,
			expectPages:   0,
			expectHasNext: false,
			expectHasPrev: false,
		},
		{
			name:          "大页面尺寸",
			total:         100,
			page:          1,
			size:          100,
			expectPages:   1,
			expectHasNext: false,
			expectHasPrev: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 计算总页数
			pages := 0
			if tt.size > 0 {
				pages = (tt.total + tt.size - 1) / tt.size
			}
			assert.Equal(t, tt.expectPages, pages)

			// 是否有下一页
			hasNext := tt.page < pages
			assert.Equal(t, tt.expectHasNext, hasNext)

			// 是否有上一页
			hasPrev := tt.page > 1
			assert.Equal(t, tt.expectHasPrev, hasPrev)
		})
	}
}

// TestQueryResultPaginationEdgeCases 测试分页边界条件
func TestQueryResultPaginationEdgeCases(t *testing.T) {
	t.Run("size为0", func(t *testing.T) {
		// 避免除零错误
		total := 100
		size := 0
		if size > 0 {
			_ = (total + size - 1) / size
		}
		// 不应该panic
		assert.True(t, true)
	})

	t.Run("负数参数", func(t *testing.T) {
		// 测试负数参数的处理
		total := -10
		size := 10
		pages := (total + size - 1) / size
		assert.Equal(t, 0, pages)
	})

	t.Run("超大页面", func(t *testing.T) {
		total := 100
		page := 1000
		size := 10
		pages := (total + size - 1) / size
		hasNext := page < pages
		assert.False(t, hasNext)
	})
}

// ==================== 复杂查询场景测试 ====================

// TestComplexQueryScenario 测试复杂查询场景
func TestComplexQueryScenario(t *testing.T) {
	type SearchUsersQuery struct {
		Keyword   string
		Status    string
		StartDate string
		EndDate   string
		Page      int
		Size      int
	}

	type UserDTO struct {
		ID     string
		Name   string
		Email  string
		Status string
	}

	query := SearchUsersQuery{
		Keyword:   "test",
		Status:    "active",
		StartDate: "2024-01-01",
		EndDate:   "2024-12-31",
		Page:      1,
		Size:      20,
	}

	result := QueryResult[UserDTO]{
		Data: []UserDTO{
			{ID: "1", Name: "Test User 1", Email: "test1@example.com", Status: "active"},
			{ID: "2", Name: "Test User 2", Email: "test2@example.com", Status: "active"},
		},
		Total: 2,
		Page:  query.Page,
		Size:  query.Size,
	}

	assert.Equal(t, query.Page, result.Page)
	assert.Equal(t, query.Size, result.Size)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "active", result.Data[0].Status)
}

// TestQueryResultWithNestedStruct 测试嵌套结构体的 QueryResult
func TestQueryResultWithNestedStruct(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type UserWithAddress struct {
		ID      string
		Name    string
		Email   string
		Address Address
	}

	result := QueryResult[UserWithAddress]{
		Data: []UserWithAddress{
			{
				ID:    "1",
				Name:  "User 1",
				Email: "user1@example.com",
				Address: Address{
					Street:  "123 Main St",
					City:    "New York",
					Country: "USA",
				},
			},
			{
				ID:    "2",
				Name:  "User 2",
				Email: "user2@example.com",
				Address: Address{
					Street:  "456 Oak Ave",
					City:    "Los Angeles",
					Country: "USA",
				},
			},
		},
		Total: 2,
		Page:  1,
		Size:  10,
	}

	assert.Len(t, result.Data, 2)
	assert.Equal(t, "New York", result.Data[0].Address.City)
	assert.Equal(t, "Los Angeles", result.Data[1].Address.City)
}

// ==================== 性能测试 ====================

// BenchmarkQueryResult_Create QueryResult 创建性能测试
func BenchmarkQueryResult_Create(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = "item" + string(rune(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = QueryResult[string]{
			Data:  data,
			Total: 1000,
			Page:  1,
			Size:  100,
		}
	}
}

// BenchmarkQueryHandler_Handle Query Handler 性能测试
func BenchmarkQueryHandler_Handle(b *testing.B) {
	ctx := context.Background()
	handler := &MockQueryHandler{
		returnResult: &MockQueryResult{ID: "123", Name: "Test"},
		returnError:  nil,
	}
	query := MockQuery{ID: "123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.Handle(ctx, query)
	}
}

// BenchmarkQuery_Execute Query 执行性能测试
func BenchmarkQuery_Execute(b *testing.B) {
	ctx := context.Background()
	query := MockQuery{ID: "123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = query.Execute(ctx)
	}
}

// BenchmarkQueryResultPagination 分页计算性能测试
func BenchmarkQueryResultPagination(b *testing.B) {
	total := 1000000
	size := 20

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = (total + size - 1) / size
	}
}

// ==================== 示例函数 ====================

// ExampleQueryResult QueryResult 使用示例
func ExampleQueryResult() {
	// 创建分页查询结果
	result := QueryResult[string]{
		Data:  []string{"user1", "user2", "user3"},
		Total: 100,
		Page:  1,
		Size:  10,
	}

	_ = result
	// Output:
}

// ExampleMockQuery MockQuery 使用示例
func ExampleMockQuery() {
	ctx := context.Background()
	query := MockQuery{ID: "123"}
	_, _ = query.Execute(ctx)
	// Output:
}

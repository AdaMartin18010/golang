package patterns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==================== DTO 测试 ====================

// TestDTO 测试 DTO 结构
func TestDTO(t *testing.T) {
	t.Run("创建DTO", func(t *testing.T) {
		now := time.Now()
		dto := DTO{
			ID:        "test-id-123",
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, "test-id-123", dto.ID)
		assert.Equal(t, now, dto.CreatedAt)
		assert.Equal(t, now, dto.UpdatedAt)
	})

	t.Run("空DTO", func(t *testing.T) {
		dto := DTO{}

		assert.Empty(t, dto.ID)
		assert.True(t, dto.CreatedAt.IsZero())
		assert.True(t, dto.UpdatedAt.IsZero())
	})

	t.Run("不同时间戳", func(t *testing.T) {
		createdAt := time.Now().Add(-24 * time.Hour)
		updatedAt := time.Now()
		dto := DTO{
			ID:        "test-id",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		assert.True(t, dto.CreatedAt.Before(dto.UpdatedAt))
	})

	t.Run("特殊字符ID", func(t *testing.T) {
		dto := DTO{
			ID:        "id-with-special-chars_123.456",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, "id-with-special-chars_123.456", dto.ID)
	})

	t.Run("UUID格式ID", func(t *testing.T) {
		uuid := "550e8400-e29b-41d4-a716-446655440000"
		dto := DTO{
			ID:        uuid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.Equal(t, uuid, dto.ID)
	})
}

// TestDTOWithTime 测试 DTO 的时间处理
func TestDTOWithTime(t *testing.T) {
	t.Run("UTC时间", func(t *testing.T) {
		now := time.Now().UTC()
		dto := DTO{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, now, dto.CreatedAt)
		assert.Equal(t, "UTC", dto.CreatedAt.Location().String())
	})

	t.Run("不同时区", func(t *testing.T) {
		shanghai, _ := time.LoadLocation("Asia/Shanghai")
		now := time.Now().In(shanghai)
		dto := DTO{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, "Asia/Shanghai", dto.CreatedAt.Location().String())
	})

	t.Run("零值时间", func(t *testing.T) {
		dto := DTO{
			ID:        "test-id",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		assert.True(t, dto.CreatedAt.IsZero())
		assert.True(t, dto.UpdatedAt.IsZero())
	})
}

// ==================== PaginatedDTO 测试 ====================

// TestPaginatedDTO 测试 PaginatedDTO 结构
func TestPaginatedDTO(t *testing.T) {
	t.Run("基本分页DTO", func(t *testing.T) {
		dto := PaginatedDTO[string]{
			Data:  []string{"item1", "item2", "item3"},
			Total: 100,
			Page:  1,
			Size:  10,
		}

		assert.Len(t, dto.Data, 3)
		assert.Equal(t, 100, dto.Total)
		assert.Equal(t, 1, dto.Page)
		assert.Equal(t, 10, dto.Size)
	})

	t.Run("空分页DTO", func(t *testing.T) {
		dto := PaginatedDTO[string]{
			Data:  []string{},
			Total: 0,
			Page:  1,
			Size:  10,
		}

		assert.Empty(t, dto.Data)
		assert.Equal(t, 0, dto.Total)
	})

	t.Run("nil数据分页DTO", func(t *testing.T) {
		dto := PaginatedDTO[string]{
			Data:  nil,
			Total: 0,
			Page:  1,
			Size:  10,
		}

		assert.Nil(t, dto.Data)
		assert.Equal(t, 0, dto.Total)
	})

	t.Run("结构体分页DTO", func(t *testing.T) {
		type UserDTO struct {
			ID   string
			Name string
		}

		dto := PaginatedDTO[UserDTO]{
			Data: []UserDTO{
				{ID: "1", Name: "User 1"},
				{ID: "2", Name: "User 2"},
				{ID: "3", Name: "User 3"},
			},
			Total: 50,
			Page:  2,
			Size:  20,
		}

		assert.Len(t, dto.Data, 3)
		assert.Equal(t, 50, dto.Total)
		assert.Equal(t, 2, dto.Page)
		assert.Equal(t, 20, dto.Size)
		assert.Equal(t, "User 1", dto.Data[0].Name)
	})

	t.Run("最后一页", func(t *testing.T) {
		dto := PaginatedDTO[int]{
			Data:  []int{91, 92, 93},
			Total: 93,
			Page:  10,
			Size:  10,
		}

		assert.Equal(t, 10, dto.Page)
		assert.Len(t, dto.Data, 3)
	})

	t.Run("大量数据", func(t *testing.T) {
		data := make([]int, 1000)
		for i := range data {
			data[i] = i
		}

		dto := PaginatedDTO[int]{
			Data:  data,
			Total: 10000,
			Page:  1,
			Size:  1000,
		}

		assert.Len(t, dto.Data, 1000)
		assert.Equal(t, 10000, dto.Total)
	})
}

// TestPaginatedDTOWithInt 测试整数类型的 PaginatedDTO
func TestPaginatedDTOWithInt(t *testing.T) {
	dto := PaginatedDTO[int]{
		Data:  []int{1, 2, 3, 4, 5},
		Total: 100,
		Page:  1,
		Size:  5,
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, dto.Data)
	assert.Equal(t, 100, dto.Total)
}

// TestPaginatedDTOWithPointer 测试指针类型的 PaginatedDTO
func TestPaginatedDTOWithPointer(t *testing.T) {
	type Item struct {
		Name  string
		Value int
	}

	item1 := &Item{Name: "item1", Value: 10}
	item2 := &Item{Name: "item2", Value: 20}

	dto := PaginatedDTO[*Item]{
		Data:  []*Item{item1, item2},
		Total: 2,
		Page:  1,
		Size:  10,
	}

	assert.Len(t, dto.Data, 2)
	assert.Equal(t, "item1", dto.Data[0].Name)
	assert.Equal(t, 20, dto.Data[1].Value)
}

// TestPaginatedDTOPaginationLogic 测试分页逻辑
func TestPaginatedDTOPaginationLogic(t *testing.T) {
	tests := []struct {
		name        string
		total       int
		page        int
		size        int
		expectPages int
	}{
		{
			name:        "正好整页",
			total:       100,
			page:        1,
			size:        10,
			expectPages: 10,
		},
		{
			name:        "非整页",
			total:       95,
			page:        1,
			size:        10,
			expectPages: 10, // 95/10 = 9.5 -> 10
		},
		{
			name:        "少于每页大小",
			total:       5,
			page:        1,
			size:        10,
			expectPages: 1,
		},
		{
			name:        "零条记录",
			total:       0,
			page:        1,
			size:        10,
			expectPages: 0,
		},
		{
			name:        "size为0",
			total:       100,
			page:        1,
			size:        0,
			expectPages: 0, // 避免除零
		},
		{
			name:        "大尺寸分页",
			total:       1000000,
			page:        1,
			size:        1000,
			expectPages: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := 0
			if tt.size > 0 {
				pages = (tt.total + tt.size - 1) / tt.size
			}
			assert.Equal(t, tt.expectPages, pages)
		})
	}
}

// TestPaginatedDTOPaginationNavigation 测试分页导航
func TestPaginatedDTOPaginationNavigation(t *testing.T) {
	tests := []struct {
		name          string
		total         int
		page          int
		size          int
		expectHasNext bool
		expectHasPrev bool
		expectIsFirst bool
		expectIsLast  bool
	}{
		{
			name:          "第一页",
			total:         100,
			page:          1,
			size:          10,
			expectHasNext: true,
			expectHasPrev: false,
			expectIsFirst: true,
			expectIsLast:  false,
		},
		{
			name:          "中间页",
			total:         100,
			page:          5,
			size:          10,
			expectHasNext: true,
			expectHasPrev: true,
			expectIsFirst: false,
			expectIsLast:  false,
		},
		{
			name:          "最后一页",
			total:         100,
			page:          10,
			size:          10,
			expectHasNext: false,
			expectHasPrev: true,
			expectIsFirst: false,
			expectIsLast:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := (tt.total + tt.size - 1) / tt.size
			hasNext := tt.page < pages
			hasPrev := tt.page > 1
			isFirst := tt.page == 1
			isLast := tt.page == pages

			assert.Equal(t, tt.expectHasNext, hasNext)
			assert.Equal(t, tt.expectHasPrev, hasPrev)
			assert.Equal(t, tt.expectIsFirst, isFirst)
			assert.Equal(t, tt.expectIsLast, isLast)
		})
	}
}

// ==================== ToDTO / FromDTO 接口测试 ====================

// MockEntity 模拟领域实体
type MockEntity struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MockEntityDTO 模拟 DTO
type MockEntityDTO struct {
	DTO
	Name  string
	Email string
}

// ToDTO 转换为 DTO（模拟 ToDTO 接口）
func (e *MockEntity) ToDTO() *MockEntityDTO {
	return &MockEntityDTO{
		DTO: DTO{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		Name:  e.Name,
		Email: e.Email,
	}
}

// TestToDTO 测试 ToDTO 转换
func TestToDTO(t *testing.T) {
	now := time.Now()
	entity := &MockEntity{
		ID:        "entity-123",
		Name:      "Test Entity",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := entity.ToDTO()

	assert.Equal(t, "entity-123", dto.ID)
	assert.Equal(t, "Test Entity", dto.Name)
	assert.Equal(t, "test@example.com", dto.Email)
	assert.Equal(t, now, dto.CreatedAt)
	assert.Equal(t, now, dto.UpdatedAt)
}

// TestToDTOWithEmptyEntity 测试空实体转换
func TestToDTOWithEmptyEntity(t *testing.T) {
	entity := &MockEntity{}

	dto := entity.ToDTO()

	assert.Empty(t, dto.ID)
	assert.Empty(t, dto.Name)
	assert.Empty(t, dto.Email)
	assert.True(t, dto.CreatedAt.IsZero())
}

// TestToDTOInterface 测试 ToDTO 接口
func TestToDTOInterface(t *testing.T) {
	// 验证 MockEntity 实现了 ToDTO 接口
	var _ ToDTO[*MockEntityDTO] = &MockEntity{}

	now := time.Now()
	entity := &MockEntity{
		ID:        "test-id",
		Name:      "Test",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	var converter ToDTO[*MockEntityDTO] = entity
	dto := converter.ToDTO()

	assert.NotNil(t, dto)
	assert.Equal(t, "test-id", dto.ID)
}

// TestFromDTOInterface 测试 FromDTO 接口
func TestFromDTOInterface(t *testing.T) {
	// 模拟 FromDTO 接口实现
	type MockEntityWithFromDTO struct {
		MockEntity
	}

	now := time.Now()
	dto := &MockEntityDTO{
		DTO: DTO{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:  "Test Name",
		Email: "test@example.com",
	}

	// 手动实现 FromDTO
	entity := &MockEntity{
		ID:        dto.ID,
		Name:      dto.Name,
		Email:     dto.Email,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	assert.Equal(t, dto.ID, entity.ID)
	assert.Equal(t, dto.Name, entity.Name)
	assert.Equal(t, dto.Email, entity.Email)
}

// ==================== 复杂 DTO 测试 ====================

// TestDTOWithComplexData 测试复杂数据的 DTO
func TestDTOWithComplexData(t *testing.T) {
	type AddressDTO struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type UserDetailDTO struct {
		DTO
		Name    string     `json:"name"`
		Email   string     `json:"email"`
		Address AddressDTO `json:"address"`
		Tags    []string   `json:"tags"`
	}

	now := time.Now()
	dto := UserDetailDTO{
		DTO: DTO{
			ID:        "user-123",
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:  "John Doe",
		Email: "john@example.com",
		Address: AddressDTO{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
		Tags: []string{"premium", "verified"},
	}

	assert.Equal(t, "user-123", dto.ID)
	assert.Equal(t, "John Doe", dto.Name)
	assert.Equal(t, "New York", dto.Address.City)
	assert.Len(t, dto.Tags, 2)
	assert.Contains(t, dto.Tags, "premium")
}

// TestDTOWithNestedPaginatedDTO 测试嵌套的分页 DTO
func TestDTOWithNestedPaginatedDTO(t *testing.T) {
	type OrderItemDTO struct {
		ProductID string
		Quantity  int
		Price     float64
	}

	type OrderDTO struct {
		DTO
		UserID string
		Items  PaginatedDTO[OrderItemDTO]
		Total  float64
	}

	now := time.Now()
	order := OrderDTO{
		DTO: DTO{
			ID:        "order-123",
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID: "user-456",
		Items: PaginatedDTO[OrderItemDTO]{
			Data: []OrderItemDTO{
				{ProductID: "prod-1", Quantity: 2, Price: 29.99},
				{ProductID: "prod-2", Quantity: 1, Price: 49.99},
			},
			Total: 2,
			Page:  1,
			Size:  10,
		},
		Total: 109.97,
	}

	assert.Equal(t, "order-123", order.ID)
	assert.Equal(t, "user-456", order.UserID)
	assert.Len(t, order.Items.Data, 2)
	assert.Equal(t, 109.97, order.Total)
}

// ==================== 性能测试 ====================

// BenchmarkDTO_Create DTO 创建性能测试
func BenchmarkDTO_Create(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DTO{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
}

// BenchmarkPaginatedDTO_Create PaginatedDTO 创建性能测试
func BenchmarkPaginatedDTO_Create(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = "item"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PaginatedDTO[string]{
			Data:  data,
			Total: 1000,
			Page:  1,
			Size:  100,
		}
	}
}

// BenchmarkToDTO_ToDTO 转换性能测试
func BenchmarkToDTO_ToDTO(b *testing.B) {
	now := time.Now()
	entity := &MockEntity{
		ID:        "entity-123",
		Name:      "Test Entity",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = entity.ToDTO()
	}
}

// ==================== 示例函数 ====================

// ExampleDTO DTO 使用示例
func ExampleDTO() {
	dto := DTO{
		ID:        "example-id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_ = dto
	// Output:
}

// ExamplePaginatedDTO PaginatedDTO 使用示例
func ExamplePaginatedDTO() {
	result := PaginatedDTO[string]{
		Data:  []string{"item1", "item2"},
		Total: 100,
		Page:  1,
		Size:  10,
	}

	_ = result
	// Output:
}

// ExampleToDTO ToDTO 接口使用示例
func ExampleToDTO() {
	entity := &MockEntity{
		ID:    "example-id",
		Name:  "Example",
		Email: "example@example.com",
	}

	dto := entity.ToDTO()
	_ = dto
	// Output:
}

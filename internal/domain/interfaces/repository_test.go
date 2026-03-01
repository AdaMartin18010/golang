package interfaces

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEntity 测试用的实体类型
type TestRepoEntity struct {
	ID    string
	Name  string
	Value int
}

// mockRepository 模拟仓储实现，用于测试接口
type mockRepository[T any] struct {
	data    map[string]*T
	listData []*T
	err     error
}

func newMockRepository[T any]() *mockRepository[T] {
	return &mockRepository[T]{
		data:     make(map[string]*T),
		listData: make([]*T, 0),
	}
}

func (m *mockRepository[T]) Create(ctx context.Context, entity *T) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	if m.err != nil {
		return nil, m.err
	}
	if entity, ok := m.data[id]; ok {
		return entity, nil
	}
	return nil, nil
}

func (m *mockRepository[T]) Update(ctx context.Context, entity *T) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockRepository[T]) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.data, id)
	return nil
}

func (m *mockRepository[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listData, nil
}

// mockRepositoryWithQuery 支持查询的模拟仓储
type TestQuery struct {
	Name string
	Min  int
	Max  int
}

type mockRepositoryWithQuery[T any, Q any] struct {
	mockRepository[T]
	queryResult []*T
	count       int
}

func (m *mockRepositoryWithQuery[T, Q]) FindByQuery(ctx context.Context, query Q) ([]*T, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.queryResult, nil
}

func (m *mockRepositoryWithQuery[T, Q]) Count(ctx context.Context) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.count, nil
}

func (m *mockRepositoryWithQuery[T, Q]) CountByQuery(ctx context.Context, query Q) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return len(m.queryResult), nil
}

// TestRepository_Create 测试 Create 方法
func TestRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()

	entity := &TestRepoEntity{
		ID:    "1",
		Name:  "Test",
		Value: 100,
	}

	err := repo.Create(ctx, entity)
	assert.NoError(t, err)
}

// TestRepository_Create_Error 测试 Create 错误处理
func TestRepository_Create_Error(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()
	repo.err = errors.New("database error")

	entity := &TestRepoEntity{ID: "1", Name: "Test"}
	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

// TestRepository_FindByID 测试 FindByID 方法
func TestRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()

	// 测试找不到的情况
	result, err := repo.FindByID(ctx, "non-existent")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// 测试找到的情况
	entity := &TestRepoEntity{ID: "1", Name: "Test", Value: 100}
	repo.data["1"] = entity

	result, err = repo.FindByID(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Name)
}

// TestRepository_FindByID_Error 测试 FindByID 错误处理
func TestRepository_FindByID_Error(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()
	repo.err = errors.New("connection failed")

	result, err := repo.FindByID(ctx, "1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestRepository_Update 测试 Update 方法
func TestRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()

	entity := &TestRepoEntity{ID: "1", Name: "Updated", Value: 200}
	err := repo.Update(ctx, entity)

	assert.NoError(t, err)
}

// TestRepository_Update_Error 测试 Update 错误处理
func TestRepository_Update_Error(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()
	repo.err = errors.New("update failed")

	entity := &TestRepoEntity{ID: "1", Name: "Test"}
	err := repo.Update(ctx, entity)

	assert.Error(t, err)
}

// TestRepository_Delete 测试 Delete 方法
func TestRepository_Delete(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()

	// 先添加数据
	repo.data["1"] = &TestRepoEntity{ID: "1", Name: "Test"}

	err := repo.Delete(ctx, "1")
	assert.NoError(t, err)
	assert.Nil(t, repo.data["1"])
}

// TestRepository_Delete_Error 测试 Delete 错误处理
func TestRepository_Delete_Error(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()
	repo.err = errors.New("delete failed")

	err := repo.Delete(ctx, "1")

	assert.Error(t, err)
}

// TestRepository_List 测试 List 方法
func TestRepository_List(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()

	// 空列表
	results, err := repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, results)
	assert.NotNil(t, results)

	// 有数据的列表
	repo.listData = []*TestRepoEntity{
		{ID: "1", Name: "First"},
		{ID: "2", Name: "Second"},
	}

	results, err = repo.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

// TestRepository_List_Error 测试 List 错误处理
func TestRepository_List_Error(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository[TestRepoEntity]()
	repo.err = errors.New("query failed")

	results, err := repo.List(ctx, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, results)
}

// TestRepositoryWithQuery_FindByQuery 测试 FindByQuery 方法
func TestRepositoryWithQuery_FindByQuery(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{
		queryResult: []*TestRepoEntity{
			{ID: "1", Name: "Test1", Value: 100},
			{ID: "2", Name: "Test2", Value: 200},
		},
	}

	query := TestQuery{Name: "Test", Min: 0, Max: 100}
	results, err := repo.FindByQuery(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

// TestRepositoryWithQuery_FindByQuery_Error 测试 FindByQuery 错误处理
func TestRepositoryWithQuery_FindByQuery_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{}
	repo.err = errors.New("query execution failed")

	query := TestQuery{Name: "Test"}
	results, err := repo.FindByQuery(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, results)
}

// TestRepositoryWithQuery_Count 测试 Count 方法
func TestRepositoryWithQuery_Count(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{}
	repo.count = 100

	count, err := repo.Count(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 100, count)
}

// TestRepositoryWithQuery_Count_Error 测试 Count 错误处理
func TestRepositoryWithQuery_Count_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{}
	repo.err = errors.New("count failed")

	count, err := repo.Count(ctx)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

// TestRepositoryWithQuery_CountByQuery 测试 CountByQuery 方法
func TestRepositoryWithQuery_CountByQuery(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{}
	repo.queryResult = make([]*TestRepoEntity, 5)

	query := TestQuery{Name: "Test"}
	count, err := repo.CountByQuery(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

// TestRepositoryWithQuery_CountByQuery_Error 测试 CountByQuery 错误处理
func TestRepositoryWithQuery_CountByQuery_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepositoryWithQuery[TestRepoEntity, TestQuery]{}
	repo.err = errors.New("count by query failed")

	query := TestQuery{Name: "Test"}
	count, err := repo.CountByQuery(ctx, query)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

// TestRepository_ContextCancellation 测试上下文取消
func TestRepository_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	repo := newMockRepository[TestRepoEntity]()

	// 测试各种方法在取消的上下文下的行为
	err := repo.Create(ctx, &TestRepoEntity{ID: "1"})
	assert.NoError(t, err) // mock 不检查上下文

	_, err = repo.FindByID(ctx, "1")
	assert.NoError(t, err)
}

// TestRepository_InterfaceCompliance 验证接口实现
func TestRepository_InterfaceCompliance(t *testing.T) {
	// 验证 mockRepository 实现了 Repository 接口
	var _ Repository[TestRepoEntity] = (*mockRepository[TestRepoEntity])(nil)

	// 验证 mockRepositoryWithQuery 实现了 RepositoryWithQuery 接口
	var _ RepositoryWithQuery[TestRepoEntity, TestQuery] = (*mockRepositoryWithQuery[TestRepoEntity, TestQuery])(nil)
}

// TestRepository_GenericType 测试泛型类型支持
func TestRepository_GenericType(t *testing.T) {
	ctx := context.Background()

	// 测试不同类型的实体
	type StringEntity struct {
		ID   string
		Data string
	}

	type IntEntity struct {
		ID   string
		Data int
	}

	stringRepo := newMockRepository[StringEntity]()
	intRepo := newMockRepository[IntEntity]()

	// StringEntity 操作
	stringEntity := &StringEntity{ID: "1", Data: "test"}
	err := stringRepo.Create(ctx, stringEntity)
	assert.NoError(t, err)

	// IntEntity 操作
	intEntity := &IntEntity{ID: "1", Data: 42}
	err = intRepo.Create(ctx, intEntity)
	assert.NoError(t, err)
}

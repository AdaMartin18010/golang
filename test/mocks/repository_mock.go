package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository 通用 Repository Mock
// 使用 testify/mock 实现
type MockRepository[T any] struct {
	mock.Mock
}

// Create Mock 实现
func (m *MockRepository[T]) Create(ctx context.Context, entity *T) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

// FindByID Mock 实现
func (m *MockRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*T), args.Error(1)
}

// Update Mock 实现
func (m *MockRepository[T]) Update(ctx context.Context, entity *T) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

// Delete Mock 实现
func (m *MockRepository[T]) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List Mock 实现
func (m *MockRepository[T]) List(ctx context.Context, limit, offset int) ([]*T, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*T), args.Error(1)
}

// ExpectCreate 设置 Create 期望
func (m *MockRepository[T]) ExpectCreate(ctx context.Context, entity *T, err error) *mock.Call {
	return m.On("Create", ctx, entity).Return(err)
}

// ExpectFindByID 设置 FindByID 期望
func (m *MockRepository[T]) ExpectFindByID(ctx context.Context, id string, result *T, err error) *mock.Call {
	return m.On("FindByID", ctx, id).Return(result, err)
}

// ExpectUpdate 设置 Update 期望
func (m *MockRepository[T]) ExpectUpdate(ctx context.Context, entity *T, err error) *mock.Call {
	return m.On("Update", ctx, entity).Return(err)
}

// ExpectDelete 设置 Delete 期望
func (m *MockRepository[T]) ExpectDelete(ctx context.Context, id string, err error) *mock.Call {
	return m.On("Delete", ctx, id).Return(err)
}

// ExpectList 设置 List 期望
func (m *MockRepository[T]) ExpectList(ctx context.Context, limit, offset int, result []*T, err error) *mock.Call {
	return m.On("List", ctx, limit, offset).Return(result, err)
}

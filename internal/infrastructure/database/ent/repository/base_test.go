// Package repository provides comprehensive tests for Ent base repository.
//
// 测试策略：
// 1. 使用 SQLite 内存数据库进行单元测试
// 2. 测试泛型仓储的基础功能
// 3. 覆盖事务管理和错误处理
//
// 运行测试：
//   - go test -v ./internal/infrastructure/database/ent/repository/...
package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/infrastructure/database/ent"
	"github.com/yourusername/golang/internal/infrastructure/database/ent/migrate"
	entuser "github.com/yourusername/golang/internal/infrastructure/database/ent/user"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestClient 创建测试用的 Ent 客户端
func setupTestClient(t *testing.T) *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to open ent client: %v", err)
	}

	// 执行 schema 迁移
	ctx := context.Background()
	if err := migrate.Create(ctx, client.Schema, migrate.Tables); err != nil {
		client.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	return client
}

// setupTestRepository 创建测试用的仓储和客户端
func setupTestRepository(t *testing.T) (*BaseRepository[user.User, *ent.User], *ent.Client) {
	client := setupTestClient(t)

	// 创建转换函数
	toDomain := func(e *ent.User) (*user.User, error) {
		return &user.User{
			ID:        e.ID,
			Email:     e.Email,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}, nil
	}

	toEnt := func(u *user.User) (*ent.User, error) {
		return &ent.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}, nil
	}

	getID := func(u *user.User) (string, error) {
		return u.ID, nil
	}

	setID := func(u *user.User, id string) error {
		u.ID = id
		return nil
	}

	repo := NewBaseRepository[user.User, *ent.User](client, toDomain, toEnt, getID, setID)
	return repo, client
}

// TestNewBaseRepository 测试创建基础仓储
func TestNewBaseRepository(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.Client())
	assert.Equal(t, client, repo.Client())
}

// TestBaseRepository_Client 测试获取客户端
func TestBaseRepository_Client(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	c := repo.Client()
	assert.NotNil(t, c)
	assert.Equal(t, client, c)
}

// TestBaseRepository_WithTx_Commit 测试事务提交
func TestBaseRepository_WithTx_Commit(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 在事务中创建用户
	err := repo.WithTx(ctx, func(tx *ent.Tx) error {
		_, err := tx.User.Create().
			SetID("tx-user-id").
			SetEmail("tx@example.com").
			SetName("Transaction User").
			Save(ctx)
		return err
	})
	require.NoError(t, err)

	// 验证用户已创建
	u, err := client.User.Query().
		Where(entuser.EmailEQ("tx@example.com")).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Transaction User", u.Name)
}

// TestBaseRepository_WithTx_Rollback 测试事务回滚
func TestBaseRepository_WithTx_Rollback(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 在事务中创建用户，然后返回错误触发回滚
	testErr := errors.New("test error")
	err := repo.WithTx(ctx, func(tx *ent.Tx) error {
		_, err := tx.User.Create().
			SetID("rollback-user-id").
			SetEmail("rollback@example.com").
			SetName("Rollback User").
			Save(ctx)
		if err != nil {
			return err
		}
		return testErr
	})
	assert.Error(t, err)
	// 错误可能是被包装后的错误，检查错误消息包含 test error
	assert.Contains(t, err.Error(), "test error")

	// 验证用户未创建
	exists, _ := client.User.Query().
		Where(entuser.EmailEQ("rollback@example.com")).
		Exist(ctx)
	assert.False(t, exists)
}

// TestBaseRepository_WithTx_Panic 测试事务 panic 回滚
func TestBaseRepository_WithTx_Panic(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 使用 defer/recover 捕获 panic
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "test panic", r)
	}()

	// 在事务中触发 panic
	_ = repo.WithTx(ctx, func(tx *ent.Tx) error {
		_, err := tx.User.Create().
			SetID("panic-user-id").
			SetEmail("panic@example.com").
			SetName("Panic User").
			Save(ctx)
		if err != nil {
			return err
		}
		panic("test panic")
	})
}

// TestBaseRepository_Create_NotImplemented 测试 Create 未实现
func TestBaseRepository_Create_NotImplemented(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()
	u := user.NewUser("test@example.com", "Test User")

	err := repo.Create(ctx, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be implemented")
}

// TestBaseRepository_FindByID_NotImplemented 测试 FindByID 未实现
func TestBaseRepository_FindByID_NotImplemented(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	_, err := repo.FindByID(ctx, "test-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be implemented")
}

// TestBaseRepository_Update_NotImplemented 测试 Update 未实现
func TestBaseRepository_Update_NotImplemented(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()
	u := user.NewUser("test@example.com", "Test User")
	u.ID = "test-id"

	err := repo.Update(ctx, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be implemented")
}

// TestBaseRepository_Update_GetIDError 测试 Update 获取 ID 错误
func TestBaseRepository_Update_GetIDError(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	getIDError := errors.New("get id error")

	repo := NewBaseRepository[user.User, *ent.User](
		client,
		func(e *ent.User) (*user.User, error) { return nil, nil },
		func(u *user.User) (*ent.User, error) { return nil, nil },
		func(u *user.User) (string, error) { return "", getIDError },
		func(u *user.User, id string) error { return nil },
	)

	ctx := context.Background()
	u := &user.User{}

	err := repo.Update(ctx, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get entity ID")
}

// TestBaseRepository_Delete_NotImplemented 测试 Delete 未实现
func TestBaseRepository_Delete_NotImplemented(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	err := repo.Delete(ctx, "test-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be implemented")
}

// TestBaseRepository_List_NotImplemented 测试 List 未实现
func TestBaseRepository_List_NotImplemented(t *testing.T) {
	repo, client := setupTestRepository(t)
	defer client.Close()

	ctx := context.Background()

	_, err := repo.List(ctx, 10, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be implemented")
}

// TestHandleEntError_NotFound 测试处理未找到错误
func TestHandleEntError_NotFound(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: "",
		},
		{
			name:     "not found error",
			input:    errors.New("record not found"),
			expected: "not found",
		},
		{
			name:     "no rows error",
			input:    errors.New("sql: no rows in result set"),
			expected: "not found",
		},
		{
			name:     "unique constraint error",
			input:    errors.New("unique constraint violation"),
			expected: "already exists",
		},
		{
			name:     "duplicate key error",
			input:    errors.New("duplicate key value violates unique constraint"),
			expected: "already exists",
		},
		{
			name:     "other error",
			input:    errors.New("some other error"),
			expected: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleEntError(tt.input)
			if tt.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expected)
			}
		})
	}
}

// TestContains 测试字符串包含检查
func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "WORLD", true}, // 不区分大小写
		{"hello world", "foo", false},
		{"", "foo", false},
		{"foo", "", true},
		{"Hello World", "hello", true},
	}

	for _, tt := range tests {
		result := contains(tt.s, tt.substr)
		assert.Equal(t, tt.expected, result, "contains(%q, %q)", tt.s, tt.substr)
	}
}

// TestBaseRepository_WithTx_BeginError 测试事务开始错误
func TestBaseRepository_WithTx_BeginError(t *testing.T) {
	// 这个测试需要模拟客户端，我们使用一个已经关闭的客户端来触发错误
	client := setupTestClient(t)

	repo := NewBaseRepository[user.User, *ent.User](
		client,
		func(e *ent.User) (*user.User, error) { return nil, nil },
		func(u *user.User) (*ent.User, error) { return nil, nil },
		func(u *user.User) (string, error) { return "", nil },
		func(u *user.User, id string) error { return nil },
	)

	// 关闭客户端
	client.Close()

	ctx := context.Background()
	err := repo.WithTx(ctx, func(tx *ent.Tx) error {
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")
}

// BenchmarkWithTx 基准测试：事务性能
func BenchmarkWithTx(b *testing.B) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	repo := NewBaseRepository[user.User, *ent.User](
		client,
		func(e *ent.User) (*user.User, error) { return nil, nil },
		func(u *user.User) (*ent.User, error) { return nil, nil },
		func(u *user.User) (string, error) { return "", nil },
		func(u *user.User, id string) error { return nil },
	)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.WithTx(ctx, func(tx *ent.Tx) error {
			return nil
		})
	}
}

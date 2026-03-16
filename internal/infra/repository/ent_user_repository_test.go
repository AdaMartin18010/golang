// Package repository provides comprehensive tests for Ent user repository.
//
// 测试策略：
// 1. 使用 SQLite 内存数据库进行单元测试
// 2. 测试完整的 CRUD 操作
// 3. 测试错误处理和边界情况
//
// 运行测试：
//   - go test -v ./internal/infrastructure/repository/...
package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
	"github.com/yourusername/golang/internal/infra/database/ent"
	"github.com/yourusername/golang/internal/infra/database/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

// setupEntUserRepository 创建测试用的 EntUserRepository
func setupEntUserRepository(t *testing.T) (*EntUserRepository, *ent.Client) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	repo := NewEntUserRepository(client)
	return repo, client
}

// TestNewEntUserRepository 测试创建 EntUserRepository
func TestNewEntUserRepository(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntUserRepository(client)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.BaseRepository)
	assert.NotNil(t, repo.client)
}

// TestEntUserRepository_Create 测试创建用户
func TestEntUserRepository_Create(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()
	u := user.NewUser("test@example.com", "Test User")

	err := repo.Create(ctx, u)
	require.NoError(t, err)

	// 验证时间戳被设置
	assert.False(t, u.CreatedAt.IsZero())
	assert.False(t, u.UpdatedAt.IsZero())

	// 验证用户已保存到数据库
	found, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.Email, found.Email)
	assert.Equal(t, u.Name, found.Name)
}

// TestEntUserRepository_Create_DuplicateEmail 测试创建重复邮箱用户
func TestEntUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建第一个用户
	u1 := user.NewUser("duplicate@example.com", "First User")
	err := repo.Create(ctx, u1)
	require.NoError(t, err)

	// 尝试创建相同邮箱的用户
	u2 := user.NewUser("duplicate@example.com", "Second User")
	err = repo.Create(ctx, u2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

// TestEntUserRepository_FindByID 测试根据 ID 查找用户
func TestEntUserRepository_FindByID(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	u := user.NewUser("find@example.com", "Find User")
	err := repo.Create(ctx, u)
	require.NoError(t, err)

	// 查找用户
	found, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
	assert.Equal(t, u.Email, found.Email)
	assert.Equal(t, u.Name, found.Name)
}

// TestEntUserRepository_FindByID_NotFound 测试查找不存在的用户
func TestEntUserRepository_FindByID_NotFound(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 查找不存在的用户
	found, err := repo.FindByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, user.ErrUserNotFound, err)
	assert.Nil(t, found)
}

// TestEntUserRepository_FindByEmail 测试根据邮箱查找用户
func TestEntUserRepository_FindByEmail(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	u := user.NewUser("email@example.com", "Email User")
	err := repo.Create(ctx, u)
	require.NoError(t, err)

	// 通过邮箱查找
	found, err := repo.FindByEmail(ctx, "email@example.com")
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
	assert.Equal(t, u.Email, found.Email)
}

// TestEntUserRepository_FindByEmail_NotFound 测试查找不存在的邮箱
func TestEntUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 查找不存在的邮箱
	found, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, user.ErrUserNotFound, err)
	assert.Nil(t, found)
}

// TestEntUserRepository_Update 测试更新用户
func TestEntUserRepository_Update(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	u := user.NewUser("update@example.com", "Original Name")
	err := repo.Create(ctx, u)
	require.NoError(t, err)

	// 更新用户
	u.Name = "Updated Name"
	err = repo.Update(ctx, u)
	require.NoError(t, err)

	// 验证更新
	found, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", found.Name)
}

// TestEntUserRepository_Update_NotFound 测试更新不存在的用户
func TestEntUserRepository_Update_NotFound(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 更新不存在的用户
	u := &user.User{
		ID:    "non-existent-id",
		Email: "test@example.com",
		Name:  "Test User",
	}
	err := repo.Update(ctx, u)
	assert.Error(t, err)
	assert.Equal(t, user.ErrUserNotFound, err)
}

// TestEntUserRepository_Delete 测试删除用户
func TestEntUserRepository_Delete(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	u := user.NewUser("delete@example.com", "Delete User")
	err := repo.Create(ctx, u)
	require.NoError(t, err)

	// 删除用户
	err = repo.Delete(ctx, u.ID)
	require.NoError(t, err)

	// 验证用户已删除
	_, err = client.User.Get(ctx, u.ID)
	assert.Error(t, err)
}

// TestEntUserRepository_Delete_NotFound 测试删除不存在的用户
func TestEntUserRepository_Delete_NotFound(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 删除不存在的用户
	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, user.ErrUserNotFound, err)
}

// TestEntUserRepository_List 测试列出用户
func TestEntUserRepository_List(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 创建多个测试用户
	for i := 0; i < 5; i++ {
		u := user.NewUser("list"+string(rune('0'+i))+"@example.com", "List User "+string(rune('0'+i)))
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	// 列出所有用户
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, users, 5)

	// 分页：限制数量
	users, err = repo.List(ctx, 3, 0)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// 分页：偏移量
	users, err = repo.List(ctx, 3, 3)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

// TestEntUserRepository_List_Empty 测试空列表
func TestEntUserRepository_List_Empty(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 空数据库
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, users)
}

// TestEntUserRepository_WithTx 测试事务支持
func TestEntUserRepository_WithTx(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 在事务中创建用户
	err := repo.WithTx(ctx, func(tx *ent.Tx) error {
		// 使用事务客户端创建临时仓储
		txRepo := NewEntUserRepository(tx.Client())
		u := user.NewUser("tx@example.com", "Transaction User")
		return txRepo.Create(ctx, u)
	})
	require.NoError(t, err)

	// 验证用户已创建
	found, err := repo.FindByEmail(ctx, "tx@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Transaction User", found.Name)
}

// TestEntUserRepository_CRUDSequence 测试完整的 CRUD 流程
func TestEntUserRepository_CRUDSequence(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 1. Create
	u := user.NewUser("crud@example.com", "CRUD User")
	err := repo.Create(ctx, u)
	require.NoError(t, err)
	assert.NotEmpty(t, u.ID)

	// 2. Read
	found, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.Email, found.Email)
	assert.Equal(t, u.Name, found.Name)

	// 通过邮箱查找
	foundByEmail, err := repo.FindByEmail(ctx, u.Email)
	require.NoError(t, err)
	assert.Equal(t, u.ID, foundByEmail.ID)

	// 3. Update
	u.Name = "Updated CRUD User"
	err = repo.Update(ctx, u)
	require.NoError(t, err)

	// 验证更新
	updated, err := repo.FindByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated CRUD User", updated.Name)

	// 4. List
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, users, 1)

	// 5. Delete
	err = repo.Delete(ctx, u.ID)
	require.NoError(t, err)

	// 验证删除
	_, err = repo.FindByID(ctx, u.ID)
	assert.Equal(t, user.ErrUserNotFound, err)
}

// TestEntUserRepository_ConcurrentAccess 测试并发访问
func TestEntUserRepository_ConcurrentAccess(t *testing.T) {
	repo, client := setupEntUserRepository(t)
	defer client.Close()

	ctx := context.Background()

	// 并发创建用户
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			u := user.NewUser("concurrent"+string(rune('0'+n))+"@example.com", "Concurrent User")
			err := repo.Create(ctx, u)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有用户都已创建
	users, err := repo.List(ctx, 20, 0)
	require.NoError(t, err)
	assert.Len(t, users, 10)
}

// BenchmarkEntUserRepository_Create 基准测试：创建用户
func BenchmarkEntUserRepository_Create(b *testing.B) {
	client := enttest.Open(b, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntUserRepository(client)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := user.NewUser("bench"+string(rune(i))+"@example.com", "Benchmark User")
		repo.Create(ctx, u)
	}
}

// BenchmarkEntUserRepository_FindByID 基准测试：根据 ID 查找
func BenchmarkEntUserRepository_FindByID(b *testing.B) {
	client := enttest.Open(b, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntUserRepository(client)
	ctx := context.Background()

	// 准备数据
	u := user.NewUser("benchfind@example.com", "Benchmark Find User")
	err := repo.Create(ctx, u)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByID(ctx, u.ID)
	}
}

// BenchmarkEntUserRepository_FindByEmail 基准测试：根据邮箱查找
func BenchmarkEntUserRepository_FindByEmail(b *testing.B) {
	client := enttest.Open(b, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntUserRepository(client)
	ctx := context.Background()

	// 准备数据
	u := user.NewUser("benchemail@example.com", "Benchmark Email User")
	err := repo.Create(ctx, u)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByEmail(ctx, "benchemail@example.com")
	}
}


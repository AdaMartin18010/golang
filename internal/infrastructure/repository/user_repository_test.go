package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
)

// InMemoryUserRepository 内存用户仓储（用于测试）
type InMemoryUserRepository struct {
	users map[string]*user.User
	emails map[string]string // email -> id mapping
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[string]*user.User),
		emails: make(map[string]string),
	}
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if id, ok := r.emails[email]; ok {
		return r.users[id], nil
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) Save(ctx context.Context, u *user.User) error {
	if _, exists := r.users[u.ID]; exists {
		return ErrUserAlreadyExists
	}
	if _, exists := r.emails[u.Email]; exists {
		return ErrEmailAlreadyExists
	}

	r.users[u.ID] = u
	r.emails[u.Email] = u.ID
	return nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, u *user.User) error {
	if _, exists := r.users[u.ID]; !exists {
		return ErrUserNotFound
	}

	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	u, exists := r.users[id]
	if !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.emails, u.Email)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	start := offset
	if start > len(users) {
		return []*user.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// 错误定义
var (
	ErrUserNotFound = assert.AnError
	ErrUserAlreadyExists = assert.AnError
	ErrEmailAlreadyExists = assert.AnError
)

// TestInMemoryUserRepository_Save 测试保存用户
func TestInMemoryUserRepository_Save(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	testUser := user.NewUser("test@example.com", "Test User")

	err := repo.Save(ctx, testUser)
	require.NoError(t, err)

	// 验证保存成功
	retrieved, err := repo.FindByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrieved.ID)
	assert.Equal(t, testUser.Email, retrieved.Email)
}

// TestInMemoryUserRepository_Save_Duplicate 测试保存重复用户
func TestInMemoryUserRepository_Save_Duplicate(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	testUser := user.NewUser("test@example.com", "Test User")
	require.NoError(t, repo.Save(ctx, testUser))

	// 尝试保存相同ID的用户
	err := repo.Save(ctx, testUser)
	assert.Error(t, err)
}

// TestInMemoryUserRepository_FindByEmail 测试通过邮箱查找
func TestInMemoryUserRepository_FindByEmail(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	testUser := user.NewUser("test@example.com", "Test User")
	require.NoError(t, repo.Save(ctx, testUser))

	retrieved, err := repo.FindByEmail(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrieved.ID)
}

// TestInMemoryUserRepository_Update 测试更新用户
func TestInMemoryUserRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	testUser := user.NewUser("test@example.com", "Old Name")
	require.NoError(t, repo.Save(ctx, testUser))

	testUser.UpdateName("New Name")
	err := repo.Update(ctx, testUser)
	require.NoError(t, err)

	retrieved, err := repo.FindByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Name", retrieved.Name)
}

// TestInMemoryUserRepository_Delete 测试删除用户
func TestInMemoryUserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	testUser := user.NewUser("test@example.com", "Test User")
	require.NoError(t, repo.Save(ctx, testUser))

	err := repo.Delete(ctx, testUser.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, testUser.ID)
	assert.Error(t, err)
}

// TestInMemoryUserRepository_List 测试列出用户
func TestInMemoryUserRepository_List(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	// 添加多个用户
	for i := 0; i < 5; i++ {
		u := user.NewUser("user"+string(rune('0'+i))+"@example.com", "User")
		require.NoError(t, repo.Save(ctx, u))
	}

	// 列出前3个
	users, err := repo.List(ctx, 3, 0)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// 列出接下来的2个
	users, err = repo.List(ctx, 3, 3)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

// BenchmarkInMemoryUserRepository_Save 性能测试
func BenchmarkInMemoryUserRepository_Save(b *testing.B) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := user.NewUser("user@example.com", "User")
		repo.Save(ctx, u)
	}
}

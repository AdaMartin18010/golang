package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewUser 测试创建新用户
func TestNewUser(t *testing.T) {
	email := "test@example.com"
	name := "Test User"

	user := NewUser(email, name)

	assert.NotEmpty(t, user.ID, "User ID should not be empty")
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.False(t, user.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, user.UpdatedAt.IsZero(), "UpdatedAt should be set")
}

// TestUser_UpdateName 测试更新用户名称
func TestUser_UpdateName(t *testing.T) {
	user := NewUser("test@example.com", "Old Name")
	oldUpdatedAt := user.UpdatedAt

	// 等待一小段时间确保时间戳变化
	time.Sleep(10 * time.Millisecond)

	user.UpdateName("New Name")

	assert.Equal(t, "New Name", user.Name)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
}

// TestUser_UpdateEmail 测试更新用户邮箱
func TestUser_UpdateEmail(t *testing.T) {
	user := NewUser("old@example.com", "Test User")
	oldUpdatedAt := user.UpdatedAt

	// 等待一小段时间确保时间戳变化
	time.Sleep(10 * time.Millisecond)

	user.UpdateEmail("new@example.com")

	assert.Equal(t, "new@example.com", user.Email)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
}

// TestUser_IsValid 测试用户验证
func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				ID:    "test-id",
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			user: &User{
				ID:    "",
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			user: &User{
				ID:   "test-id",
				Name: "Test User",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			user: &User{
				ID:    "test-id",
				Email: "invalid-email",
				Name:  "Test User",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			user: &User{
				ID:    "test-id",
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "name too short",
			user: &User{
				ID:    "test-id",
				Email: "test@example.com",
				Name:  "A",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkNewUser 性能基准测试 - 创建用户
func BenchmarkNewUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUser("test@example.com", "Test User")
	}
}

// BenchmarkUser_IsValid 性能基准测试 - 验证用户
func BenchmarkUser_IsValid(b *testing.B) {
	user := NewUser("test@example.com", "Test User")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.IsValid()
	}
}

package user

import (
	"testing"
	"time"

	domain "github.com/yourusername/golang/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

// TestNewUser 测试创建新用户
func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		userName  string
		wantEmail string
		wantName  string
	}{
		{
			name:      "valid user",
			email:     "test@example.com",
			userName:  "Test User",
			wantEmail: "test@example.com",
			wantName:  "Test User",
		},
		{
			name:      "user with long name",
			email:     "longname@example.com",
			userName:  "Very Long User Name That Should Still Work",
			wantEmail: "longname@example.com",
			wantName:  "Very Long User Name That Should Still Work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := domain.NewUser(tt.email, tt.userName)

			assert.NotNil(t, user)
			assert.NotEmpty(t, user.ID, "ID should be generated")
			assert.Equal(t, tt.wantEmail, user.Email)
			assert.Equal(t, tt.wantName, user.Name)
			assert.False(t, user.CreatedAt.IsZero(), "CreatedAt should be set")
			assert.False(t, user.UpdatedAt.IsZero(), "UpdatedAt should be set")
			assert.Equal(t, user.CreatedAt, user.UpdatedAt, "CreatedAt and UpdatedAt should be equal for new user")
		})
	}
}

// TestUser_UpdateName 测试更新用户名称
func TestUser_UpdateName(t *testing.T) {
	tests := []struct {
		name         string
		initialName  string
		newName      string
		wantName     string
		wantUpdated  bool
	}{
		{
			name:        "update name",
			initialName: "Old Name",
			newName:     "New Name",
			wantName:    "New Name",
			wantUpdated: true,
		},
		{
			name:        "update to same name",
			initialName: "Same Name",
			newName:     "Same Name",
			wantName:    "Same Name",
			wantUpdated: true, // UpdatedAt should still be updated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := domain.NewUser("test@example.com", tt.initialName)
			oldUpdatedAt := user.UpdatedAt

			time.Sleep(10 * time.Millisecond) // 确保时间不同
			user.UpdateName(tt.newName)

			assert.Equal(t, tt.wantName, user.Name)
			if tt.wantUpdated {
				assert.True(t, user.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
			}
		})
	}
}

// TestUser_UpdateEmail 测试更新用户邮箱
func TestUser_UpdateEmail(t *testing.T) {
	tests := []struct {
		name         string
		initialEmail string
		newEmail     string
		wantEmail    string
		wantUpdated  bool
	}{
		{
			name:         "update email",
			initialEmail: "old@example.com",
			newEmail:     "new@example.com",
			wantEmail:    "new@example.com",
			wantUpdated:  true,
		},
		{
			name:         "update to same email",
			initialEmail: "same@example.com",
			newEmail:     "same@example.com",
			wantEmail:    "same@example.com",
			wantUpdated:  true, // UpdatedAt should still be updated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := domain.NewUser(tt.initialEmail, "Test User")
			oldUpdatedAt := user.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			user.UpdateEmail(tt.newEmail)

			assert.Equal(t, tt.wantEmail, user.Email)
			if tt.wantUpdated {
				assert.True(t, user.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
			}
		})
	}
}

// TestUser_IsValid 测试用户验证
func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		user    *domain.User
		wantErr error
	}{
		{
			name: "valid user",
			user: domain.NewUser("test@example.com", "Test User"),
			wantErr: nil,
		},
		{
			name: "empty ID",
			user: &domain.User{
				ID:        "",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrInvalidUserID,
		},
		{
			name: "empty email",
			user: &domain.User{
				ID:        "test-id",
				Email:     "",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrEmailRequired,
		},
		{
			name: "invalid email format - no @",
			user: &domain.User{
				ID:        "test-id",
				Email:     "invalidemail.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrInvalidEmailFormat,
		},
		{
			name: "invalid email format - no dot",
			user: &domain.User{
				ID:        "test-id",
				Email:     "invalid@email",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrInvalidEmailFormat,
		},
		{
			name: "invalid email format - @ at start",
			user: &domain.User{
				ID:        "test-id",
				Email:     "@example.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrInvalidEmailFormat,
		},
		{
			name: "invalid email format - @ at end",
			user: &domain.User{
				ID:        "test-id",
				Email:     "test@",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrInvalidEmailFormat,
		},
		{
			name: "empty name",
			user: &domain.User{
				ID:        "test-id",
				Email:     "test@example.com",
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrNameRequired,
		},
		{
			name: "name too short",
			user: &domain.User{
				ID:        "test-id",
				Email:     "test@example.com",
				Name:      "A",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: domain.ErrNameTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

// TestUser_ID_Uniqueness 测试用户 ID 的唯一性
func TestUser_ID_Uniqueness(t *testing.T) {
	user1 := domain.NewUser("user1@example.com", "User 1")
	user2 := domain.NewUser("user2@example.com", "User 2")

	assert.NotEqual(t, user1.ID, user2.ID, "User IDs should be unique")
}

// TestUser_Timestamps 测试时间戳设置
func TestUser_Timestamps(t *testing.T) {
	before := time.Now()
	user := domain.NewUser("test@example.com", "Test User")
	after := time.Now()

	assert.True(t, user.CreatedAt.After(before) || user.CreatedAt.Equal(before))
	assert.True(t, user.CreatedAt.Before(after) || user.CreatedAt.Equal(after))
	assert.True(t, user.UpdatedAt.After(before) || user.UpdatedAt.Equal(before))
	assert.True(t, user.UpdatedAt.Before(after) || user.UpdatedAt.Equal(after))
}

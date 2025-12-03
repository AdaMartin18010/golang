package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUser_Validate 测试用户实体验证
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: &User{
				ID:        "123",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty email",
			user: &User{
				ID:        "123",
				Email:     "",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid email format",
			user: &User{
				ID:        "123",
				Email:     "invalid-email",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "empty name",
			user: &User{
				ID:        "123",
				Email:     "test@example.com",
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr {
				require.Error(t, err, "Expected error but got nil")
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
			}
		})
	}
}

// TestUser_IsValid 测试用户有效性检查
func TestUser_IsValid(t *testing.T) {
	validUser := &User{
		ID:        "123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.True(t, validUser.IsValid(), "Valid user should return true")

	invalidUser := &User{
		ID:    "123",
		Email: "",
	}

	assert.False(t, invalidUser.IsValid(), "Invalid user should return false")
}

// BenchmarkUser_Validate 性能基准测试
func BenchmarkUser_Validate(b *testing.B) {
	user := &User{
		ID:        "123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

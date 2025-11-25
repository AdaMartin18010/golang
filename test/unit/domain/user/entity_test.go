package user

import (
	"testing"
	"time"

	domain "github.com/yourusername/golang/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	name := "Test User"

	user := domain.NewUser(email, name)

	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUser_UpdateName(t *testing.T) {
	user := domain.NewUser("test@example.com", "Old Name")
	oldUpdatedAt := user.UpdatedAt

	time.Sleep(10 * time.Millisecond) // 确保时间不同
	user.UpdateName("New Name")

	assert.Equal(t, "New Name", user.Name)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestUser_UpdateEmail(t *testing.T) {
	user := domain.NewUser("old@example.com", "Test User")
	oldUpdatedAt := user.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	user.UpdateEmail("new@example.com")

	assert.Equal(t, "new@example.com", user.Email)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

package user

import (
	"time"
	"github.com/google/uuid"
)

// User 用户实体
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 创建新用户
func NewUser(email, name string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateName 更新用户名
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}

// IsValid 验证用户是否有效
func (u *User) IsValid() bool {
	return u.Email != "" && u.Name != "" && u.ID != ""
}

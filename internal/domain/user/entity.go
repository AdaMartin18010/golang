package user

import "time"

// User 用户实体（领域模型）
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

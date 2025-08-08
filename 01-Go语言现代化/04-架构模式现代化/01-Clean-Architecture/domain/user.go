package domain

import (
	"errors"
	"time"
)

// User 是核心业务实体，包含用户的基本信息和业务规则
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser 创建新的用户实例
func NewUser(email, name string, age int) (*User, error) {
	user := &User{
		Email:     email,
		Name:      name,
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate 验证用户数据的有效性
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	if u.Name == "" {
		return errors.New("name cannot be empty")
	}

	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if u.Age < 0 || u.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(name string, age int) error {
	originalName := u.Name
	originalAge := u.Age

	u.Name = name
	u.Age = age
	u.UpdatedAt = time.Now()

	if err := u.Validate(); err != nil {
		// 回滚更改
		u.Name = originalName
		u.Age = originalAge
		u.UpdatedAt = time.Now()
		return err
	}

	return nil
}

// IsAdult 检查用户是否成年
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

// GetDisplayName 获取用户的显示名称
func (u *User) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

// 简单的邮箱验证函数
func isValidEmail(email string) bool {
	// 这里使用简单的验证逻辑
	// 在实际项目中，可以使用更复杂的正则表达式或专门的验证库
	return len(email) > 3 && len(email) < 255
}

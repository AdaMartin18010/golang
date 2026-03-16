package user

import (
	"time"

	"github.com/google/uuid"
)

// User 是用户领域实体，包含用户的核心业务逻辑。
//
// 设计原则：
// 1. 实体包含业务规则和业务逻辑
// 2. 实体是不可变的（通过方法修改）
// 3. 实体负责自己的验证
//
// 字段说明：
// - ID: 用户唯一标识（UUID）
// - Email: 用户邮箱（唯一）
// - Name: 用户名称
// - CreatedAt: 创建时间
// - UpdatedAt: 更新时间
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 创建新的用户实体。
//
// 功能说明：
// - 生成用户 ID（UUID）
// - 设置创建时间和更新时间
// - 不进行验证（验证应该在业务逻辑层进行）
//
// 参数：
// - email: 用户邮箱
// - name: 用户名称
//
// 返回：
// - *User: 新创建的用户实体
//
// 使用示例：
//
//	user := user.NewUser("test@example.com", "Test User")
//	// 在保存前进行验证
//	if err := validateUser(user); err != nil {
//	    return err
//	}
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

// UpdateName 更新用户名称。
//
// 功能说明：
// - 更新用户名称
// - 自动更新 UpdatedAt 时间戳
//
// 参数：
// - name: 新的用户名称
//
// 业务规则：
// - 名称不能为空（应该在业务逻辑层验证）
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// UpdateEmail 更新用户邮箱。
//
// 功能说明：
// - 更新用户邮箱
// - 自动更新 UpdatedAt 时间戳
//
// 参数：
// - email: 新的用户邮箱
//
// 业务规则：
// - 邮箱不能为空（应该在业务逻辑层验证）
// - 邮箱格式应该有效（应该在业务逻辑层验证）
// - 邮箱应该唯一（应该在业务逻辑层验证）
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}

// IsValid 验证用户实体的有效性。
//
// 功能说明：
// - 验证用户实体的基本业务规则
// - 返回验证错误列表
//
// 返回：
// - error: 验证失败时返回错误
//
// 验证规则：
// - ID 不能为空
// - Email 不能为空
// - Email 格式应该有效
// - Name 不能为空
// - Name 长度应该 >= 2
func (u *User) IsValid() error {
	if u.ID == "" {
		return ErrInvalidUserID
	}

	if u.Email == "" {
		return ErrEmailRequired
	}

	if !isValidEmailFormat(u.Email) {
		return ErrInvalidEmailFormat
	}

	if u.Name == "" {
		return ErrNameRequired
	}

	if len(u.Name) < 2 {
		return ErrNameTooShort
	}

	return nil
}

// isValidEmailFormat 验证邮箱格式（简单验证）。
//
// 功能说明：
// - 简单的邮箱格式验证
// - 生产环境应该使用更严格的验证
//
// 参数：
// - email: 邮箱地址
//
// 返回：
// - bool: 格式有效返回 true
func isValidEmailFormat(email string) bool {
	if len(email) < 3 || len(email) > 255 {
		return false
	}

	hasAt := false
	hasDot := false
	for i, char := range email {
		if char == '@' {
			if hasAt || i == 0 || i == len(email)-1 {
				return false
			}
			hasAt = true
		}
		if char == '.' && hasAt {
			hasDot = true
		}
	}

	return hasAt && hasDot
}

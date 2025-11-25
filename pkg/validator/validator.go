package validator

import (
	"regexp"
	"strings"
)

// Validator 验证器
type Validator struct{}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateName 验证名称
func ValidateName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	if len(name) < 2 || len(name) > 100 {
		return false
	}
	return true
}

// ValidateRequired 验证必填字段
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Struct 验证结构体（使用 go-playground/validator）
func (v *Validator) Struct(s interface{}) error {
	// TODO: 实现结构体验证
	// 可以使用 github.com/go-playground/validator/v10
	return nil
}

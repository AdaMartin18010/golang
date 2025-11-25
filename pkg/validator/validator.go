package validator

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidateName 验证用户名
func ValidateName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 2 && len(name) <= 100
}

// ValidateRequired 验证必填字段
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

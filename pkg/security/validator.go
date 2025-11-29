package security

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	// ErrInvalidInput 无效输入
	ErrInvalidInput = errors.New("invalid input")
	// ErrInputTooLong 输入过长
	ErrInputTooLong = errors.New("input too long")
	// ErrInputTooShort 输入过短
	ErrInputTooShort = errors.New("input too short")
	// ErrInvalidFormat 无效格式
	ErrInvalidFormat = errors.New("invalid format")
)

// InputValidator 输入验证器
type InputValidator struct {
	maxLength int
	minLength int
	pattern   *regexp.Regexp
	required  bool
}

// InputValidatorConfig 输入验证器配置
type InputValidatorConfig struct {
	MaxLength int
	MinLength int
	Pattern   string
	Required  bool
}

// NewInputValidator 创建输入验证器
func NewInputValidator(config InputValidatorConfig) *InputValidator {
	validator := &InputValidator{
		maxLength: config.MaxLength,
		minLength: config.MinLength,
		required:  config.Required,
	}

	if config.Pattern != "" {
		validator.pattern = regexp.MustCompile(config.Pattern)
	}

	return validator
}

// Validate 验证输入
func (v *InputValidator) Validate(input string) error {
	if v.required && input == "" {
		return ErrInvalidInput
	}

	if input == "" {
		return nil // 空值在非必需时允许
	}

	if v.minLength > 0 && len(input) < v.minLength {
		return ErrInputTooShort
	}

	if v.maxLength > 0 && len(input) > v.maxLength {
		return ErrInputTooLong
	}

	if v.pattern != nil && !v.pattern.MatchString(input) {
		return ErrInvalidFormat
	}

	return nil
}

// EmailValidator 邮箱验证器
type EmailValidator struct {
	*InputValidator
}

// NewEmailValidator 创建邮箱验证器
func NewEmailValidator() *EmailValidator {
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return &EmailValidator{
		InputValidator: NewInputValidator(InputValidatorConfig{
			MaxLength: 254,
			MinLength: 5,
			Pattern:   emailPattern,
			Required:  false,
		}),
	}
}

// ValidateEmail 验证邮箱
func (v *EmailValidator) ValidateEmail(email string) error {
	return v.Validate(email)
}

// URLValidator URL 验证器
type URLValidator struct {
	*InputValidator
	allowedSchemes []string
}

// NewURLValidator 创建 URL 验证器
func NewURLValidator(allowedSchemes []string) *URLValidator {
	if allowedSchemes == nil {
		allowedSchemes = []string{"http", "https"}
	}

	return &URLValidator{
		InputValidator: NewInputValidator(InputValidatorConfig{
			MaxLength: 2048,
			MinLength: 4,
			Required:  false,
		}),
		allowedSchemes: allowedSchemes,
	}
}

// ValidateURL 验证 URL
func (v *URLValidator) ValidateURL(urlStr string) error {
	if err := v.Validate(urlStr); err != nil {
		return err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ErrInvalidFormat
	}

	// 检查协议
	validScheme := false
	for _, scheme := range v.allowedSchemes {
		if parsedURL.Scheme == scheme {
			validScheme = true
			break
		}
	}

	if !validScheme {
		return fmt.Errorf("scheme must be one of: %v", v.allowedSchemes)
	}

	// 检查主机
	if parsedURL.Host == "" {
		return ErrInvalidFormat
	}

	return nil
}

// PhoneValidator 手机号验证器
type PhoneValidator struct {
	*InputValidator
}

// NewPhoneValidator 创建手机号验证器
func NewPhoneValidator() *PhoneValidator {
	// 中国手机号格式
	phonePattern := `^1[3-9]\d{9}$`
	return &PhoneValidator{
		InputValidator: NewInputValidator(InputValidatorConfig{
			MaxLength: 11,
			MinLength: 11,
			Pattern:   phonePattern,
			Required:  false,
		}),
	}
}

// ValidatePhone 验证手机号
func (v *PhoneValidator) ValidatePhone(phone string) error {
	return v.Validate(phone)
}

// StringSanitizer 字符串清理器
type StringSanitizer struct {
	trimWhitespace bool
	removeNewlines  bool
	removeTabs       bool
}

// NewStringSanitizer 创建字符串清理器
func NewStringSanitizer() *StringSanitizer {
	return &StringSanitizer{
		trimWhitespace: true,
		removeNewlines: true,
		removeTabs:     true,
	}
}

// Sanitize 清理字符串
func (s *StringSanitizer) Sanitize(input string) string {
	if s.trimWhitespace {
		input = strings.TrimSpace(input)
	}

	if s.removeNewlines {
		input = strings.ReplaceAll(input, "\n", " ")
		input = strings.ReplaceAll(input, "\r", " ")
	}

	if s.removeTabs {
		input = strings.ReplaceAll(input, "\t", " ")
	}

	// 移除多个连续空格
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	return input
}

// RemoveSpecialChars 移除特殊字符
func (s *StringSanitizer) RemoveSpecialChars(input string, allowed string) string {
	var result strings.Builder
	allowedSet := make(map[rune]bool)
	for _, r := range allowed {
		allowedSet[r] = true
	}

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || allowedSet[r] {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// NormalizeWhitespace 规范化空白字符
func (s *StringSanitizer) NormalizeWhitespace(input string) string {
	// 移除所有空白字符，只保留单个空格
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	return strings.TrimSpace(input)
}

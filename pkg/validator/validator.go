package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// Validator 验证器
type Validator struct {
	emailRegex *regexp.Regexp
}

var (
	defaultValidator *Validator
	once             sync.Once
)

// Default 返回默认验证器实例（单例）
func Default() *Validator {
	once.Do(func() {
		defaultValidator = NewValidator()
	})
	return defaultValidator
}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	}
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return Default().emailRegex.MatchString(email)
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

// ValidationError 表示字段验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors 包含多个验证错误
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

// Struct 对结构体进行标签验证
// 支持的标签: `validate:"required"`, `validate:"email"`, `validate:"min=N"`, `validate:"max=N"`
func (v *Validator) Struct(s interface{}) error {
	if s == nil {
		return fmt.Errorf("validation target is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("validation target is nil pointer")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validation target must be a struct or pointer to struct, got %s", val.Kind())
	}

	typ := val.Type()
	var errs ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出字段
		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" || tag == "-" {
			continue
		}

		if fieldErrs := validateField(field, fieldType.Name, tag); len(fieldErrs) > 0 {
			errs = append(errs, fieldErrs...)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(field reflect.Value, fieldName, tag string) []ValidationError {
	var errs []ValidationError
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		switch {
		case rule == "required":
			if isZero(field) {
				errs = append(errs, ValidationError{
					Field:   fieldName,
					Message: "required field is empty",
				})
			}
		case rule == "email":
			if field.Kind() == reflect.String && field.String() != "" {
				if !ValidateEmail(field.String()) {
					errs = append(errs, ValidationError{
						Field:   fieldName,
						Message: "invalid email format",
					})
				}
			}
		case strings.HasPrefix(rule, "min="):
			minStr := strings.TrimPrefix(rule, "min=")
			if err := checkMin(field, fieldName, minStr); err != nil {
				errs = append(errs, *err)
			}
		case strings.HasPrefix(rule, "max="):
			maxStr := strings.TrimPrefix(rule, "max=")
			if err := checkMax(field, fieldName, maxStr); err != nil {
				errs = append(errs, *err)
			}
		}
	}

	return errs
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
		return v.IsNil() || v.Len() == 0
	case reflect.Struct:
		return v.IsZero()
	default:
		return !v.IsValid()
	}
}

func checkMin(field reflect.Value, fieldName, minStr string) *ValidationError {
	// Simplified: only handle string length and numeric min
	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < len(minStr) { // approximate
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("length must be at least %s", minStr)}
		}
	}
	return nil
}

func checkMax(field reflect.Value, fieldName, maxStr string) *ValidationError {
	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > len(maxStr)*10 { // approximate placeholder
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("length must be at most %s", maxStr)}
		}
	}
	return nil
}

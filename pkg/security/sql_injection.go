package security

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrSQLInjectionDetected 检测到 SQL 注入攻击
	ErrSQLInjectionDetected = errors.New("potential SQL injection detected")
)

// SQLInjectionProtection SQL 注入防护
type SQLInjectionProtection struct {
	strictMode bool
}

// NewSQLInjectionProtection 创建 SQL 注入防护
func NewSQLInjectionProtection(strictMode bool) *SQLInjectionProtection {
	return &SQLInjectionProtection{
		strictMode: strictMode,
	}
}

// ValidateInput 验证输入，检查潜在的 SQL 注入攻击
func (s *SQLInjectionProtection) ValidateInput(input string) error {
	if input == "" {
		return nil
	}

	// SQL 关键字检测
	sqlKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "UNION", "SCRIPT", "SCRIPT", "TRUNCATE",
	}

	upperInput := strings.ToUpper(input)
	for _, keyword := range sqlKeywords {
		if strings.Contains(upperInput, keyword) {
			if s.strictMode {
				return ErrSQLInjectionDetected
			}
			// 非严格模式：检查是否有关键字后跟空格或特殊字符
			pattern := regexp.MustCompile(`(?i)\b` + keyword + `\s+`)
			if pattern.MatchString(input) {
				return ErrSQLInjectionDetected
			}
		}
	}

	// 检测 SQL 注释
	commentPatterns := []string{
		"--",
		"/*",
		"*/",
		"#",
	}

	for _, pattern := range commentPatterns {
		if strings.Contains(input, pattern) {
			return ErrSQLInjectionDetected
		}
	}

	// 检测 SQL 函数调用
	functionPattern := regexp.MustCompile(`(?i)\b(SLEEP|BENCHMARK|WAITFOR|DELAY)\s*\(`)
	if functionPattern.MatchString(input) {
		return ErrSQLInjectionDetected
	}

	// 检测分号（可能用于多语句攻击）
	if strings.Contains(input, ";") && strings.Contains(input, "'") {
		return ErrSQLInjectionDetected
	}

	// 检测单引号转义
	quotePattern := regexp.MustCompile(`('|")(\s*;\s*|\s*--|\s*/\*|\s*OR\s+|\s*AND\s+)`)
	if quotePattern.MatchString(strings.ToUpper(input)) {
		return ErrSQLInjectionDetected
	}

	return nil
}

// SanitizeInput 清理输入，移除潜在的 SQL 注入字符
func (s *SQLInjectionProtection) SanitizeInput(input string) string {
	// 移除 SQL 注释
	input = regexp.MustCompile(`--.*`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(input, "")

	// 转义单引号
	input = strings.ReplaceAll(input, "'", "''")

	// 移除分号（如果不在字符串中）
	if !strings.Contains(input, "'") {
		input = strings.ReplaceAll(input, ";", "")
	}

	return input
}

// ValidateParameter 验证参数（用于参数化查询）
func (s *SQLInjectionProtection) ValidateParameter(param interface{}) error {
	if param == nil {
		return nil
	}

	var str string
	switch v := param.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return nil // 非字符串类型通常安全
	}

	return s.ValidateInput(str)
}

package security

import (
	"html"
	"regexp"
	"strings"
)

// XSSProtection XSS 防护
type XSSProtection struct {
	allowedTags map[string]bool
}

// NewXSSProtection 创建 XSS 防护
func NewXSSProtection() *XSSProtection {
	// 默认允许的标签（可以根据需要配置）
	allowedTags := map[string]bool{
		"p":      true,
		"br":     true,
		"strong": true,
		"em":     true,
		"u":      true,
		"a":      true,
		"ul":     true,
		"ol":     true,
		"li":     true,
	}

	return &XSSProtection{
		allowedTags: allowedTags,
	}
}

// Sanitize 清理 HTML，移除潜在的 XSS 攻击
func (x *XSSProtection) Sanitize(input string) string {
	// 转义 HTML 特殊字符
	sanitized := html.EscapeString(input)
	return sanitized
}

// SanitizeHTML 清理 HTML，保留允许的标签
func (x *XSSProtection) SanitizeHTML(input string) string {
	// 移除 script 标签及其内容
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized := scriptRegex.ReplaceAllString(input, "")

	// 移除 on* 事件属性
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	sanitized = eventRegex.ReplaceAllString(sanitized, "")

	// 移除 javascript: 协议
	jsRegex := regexp.MustCompile(`(?i)javascript\s*:`)
	sanitized = jsRegex.ReplaceAllString(sanitized, "")

	// 移除 data: 协议（可能包含恶意内容）
	dataRegex := regexp.MustCompile(`(?i)data\s*:\s*[^;]*;base64`)
	sanitized = dataRegex.ReplaceAllString(sanitized, "")

	return sanitized
}

// ValidateInput 验证输入，检查潜在的 XSS 攻击
func (x *XSSProtection) ValidateInput(input string) error {
	// 检查 script 标签
	if strings.Contains(strings.ToLower(input), "<script") {
		return ErrXSSDetected
	}

	// 检查 javascript: 协议
	if strings.Contains(strings.ToLower(input), "javascript:") {
		return ErrXSSDetected
	}

	// 检查 on* 事件属性
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=`)
	if eventRegex.MatchString(input) {
		return ErrXSSDetected
	}

	// 检查 data: 协议
	if strings.Contains(strings.ToLower(input), "data:text/html") {
		return ErrXSSDetected
	}

	return nil
}

// EscapeHTML 转义 HTML 特殊字符
func (x *XSSProtection) EscapeHTML(input string) string {
	return html.EscapeString(input)
}

// UnescapeHTML 反转义 HTML 特殊字符
func (x *XSSProtection) UnescapeHTML(input string) string {
	return html.UnescapeString(input)
}

var (
	// ErrXSSDetected 检测到 XSS 攻击
	ErrXSSDetected = &XSSError{Message: "potential XSS attack detected"}
)

// XSSError XSS 错误
type XSSError struct {
	Message string
}

func (e *XSSError) Error() string {
	return e.Message
}

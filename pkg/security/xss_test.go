package security

import (
	"testing"
)

func TestXSSProtection_Sanitize(t *testing.T) {
	xss := NewXSSProtection()

	tests := []struct {
		input    string
		expected string
	}{
		{"<script>alert('XSS')</script>", "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;"},
		{"Hello <b>World</b>", "Hello &lt;b&gt;World&lt;/b&gt;"},
		{"javascript:alert('XSS')", "javascript:alert(&#39;XSS&#39;)"},
	}

	for _, tt := range tests {
		result := xss.Sanitize(tt.input)
		if result != tt.expected {
			t.Errorf("Sanitize(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestXSSProtection_SanitizeHTML(t *testing.T) {
	xss := NewXSSProtection()

	tests := []struct {
		input    string
		contains string
	}{
		{"<script>alert('XSS')</script>", ""},
		{"<p>Hello</p>", "<p>Hello</p>"},
		{"<div onclick='alert(1)'>Click</div>", ""},
		{"<a href='javascript:alert(1)'>Link</a>", ""},
	}

	for _, tt := range tests {
		result := xss.SanitizeHTML(tt.input)
		if tt.contains != "" && !contains(result, tt.contains) {
			t.Errorf("SanitizeHTML(%q) should contain %q, got %q", tt.input, tt.contains, result)
		}
		if tt.contains == "" && contains(result, "script") {
			t.Errorf("SanitizeHTML(%q) should not contain script, got %q", tt.input, result)
		}
	}
}

func TestXSSProtection_ValidateInput(t *testing.T) {
	xss := NewXSSProtection()

	tests := []struct {
		input string
		valid bool
	}{
		{"<script>alert('XSS')</script>", false},
		{"Hello World", true},
		{"javascript:alert('XSS')", false},
		{"<div onclick='alert(1)'>", false},
		{"data:text/html,<script>alert(1)</script>", false},
	}

	for _, tt := range tests {
		err := xss.ValidateInput(tt.input)
		if tt.valid && err != nil {
			t.Errorf("ValidateInput(%q) should be valid, got error: %v", tt.input, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateInput(%q) should be invalid", tt.input)
		}
	}
}

func TestXSSProtection_EscapeHTML(t *testing.T) {
	xss := NewXSSProtection()

	input := "<script>alert('XSS')</script>"
	escaped := xss.EscapeHTML(input)

	if escaped == input {
		t.Error("HTML should be escaped")
	}

	// 验证转义后的内容
	if !contains(escaped, "&lt;") {
		t.Error("Escaped HTML should contain &lt;")
	}
}

func TestXSSProtection_UnescapeHTML(t *testing.T) {
	xss := NewXSSProtection()

	escaped := "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;"
	unescaped := xss.UnescapeHTML(escaped)

	if unescaped == escaped {
		t.Error("HTML should be unescaped")
	}

	// 验证反转义后的内容
	if !contains(unescaped, "<script>") {
		t.Error("Unescaped HTML should contain <script>")
	}
}

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && 
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

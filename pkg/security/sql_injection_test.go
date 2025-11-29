package security

import (
	"testing"
)

func TestSQLInjectionProtection_ValidateInput(t *testing.T) {
	protection := NewSQLInjectionProtection(true)

	tests := []struct {
		input string
		valid bool
	}{
		{"SELECT * FROM users", false},
		{"'; DROP TABLE users; --", false},
		{"' OR '1'='1", false},
		{"admin'--", false},
		{"1' UNION SELECT * FROM users--", false},
		{"Hello World", true},
		{"user@example.com", true},
		{"12345", true},
		{"normal text", true},
	}

	for _, tt := range tests {
		err := protection.ValidateInput(tt.input)
		if tt.valid && err != nil {
			t.Errorf("ValidateInput(%q) should be valid, got error: %v", tt.input, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateInput(%q) should be invalid", tt.input)
		}
	}
}

func TestSQLInjectionProtection_SanitizeInput(t *testing.T) {
	protection := NewSQLInjectionProtection(false)

	tests := []struct {
		input    string
		expected string
	}{
		{"test'value", "test''value"},
		{"test--comment", "test"},
		{"test/*comment*/value", "testvalue"},
	}

	for _, tt := range tests {
		result := protection.SanitizeInput(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeInput(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestSQLInjectionProtection_ValidateParameter(t *testing.T) {
	protection := NewSQLInjectionProtection(true)

	// 测试字符串参数
	err := protection.ValidateParameter("SELECT * FROM users")
	if err == nil {
		t.Error("Should detect SQL injection in string parameter")
	}

	// 测试非字符串参数（应该安全）
	err = protection.ValidateParameter(123)
	if err != nil {
		t.Errorf("Non-string parameter should be safe, got error: %v", err)
	}

	err = protection.ValidateParameter([]byte("normal data"))
	if err != nil {
		t.Errorf("Byte slice with normal data should be safe, got error: %v", err)
	}
}

package security

import (
	"testing"
)

func TestInputValidator_Validate(t *testing.T) {
	validator := NewInputValidator(InputValidatorConfig{
		MinLength: 5,
		MaxLength: 10,
		Required:  true,
	})

	tests := []struct {
		input string
		valid bool
	}{
		{"short", false},      // 太短
		{"valid123", true},    // 有效
		{"toolongstring", false}, // 太长
		{"", false},           // 空值（必需）
	}

	for _, tt := range tests {
		err := validator.Validate(tt.input)
		if tt.valid && err != nil {
			t.Errorf("Validate(%q) should be valid, got error: %v", tt.input, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Validate(%q) should be invalid", tt.input)
		}
	}
}

func TestEmailValidator_ValidateEmail(t *testing.T) {
	validator := NewEmailValidator()

	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", true}, // 非必需，空值允许
	}

	for _, tt := range tests {
		err := validator.ValidateEmail(tt.email)
		if tt.valid && err != nil {
			t.Errorf("ValidateEmail(%q) should be valid, got error: %v", tt.email, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateEmail(%q) should be invalid", tt.email)
		}
	}
}

func TestURLValidator_ValidateURL(t *testing.T) {
	validator := NewURLValidator([]string{"http", "https"})

	tests := []struct {
		url   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"ftp://example.com", false}, // 不允许的协议
		{"invalid", false},
		{"", true}, // 非必需，空值允许
	}

	for _, tt := range tests {
		err := validator.ValidateURL(tt.url)
		if tt.valid && err != nil {
			t.Errorf("ValidateURL(%q) should be valid, got error: %v", tt.url, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("ValidateURL(%q) should be invalid", tt.url)
		}
	}
}

func TestStringSanitizer_Sanitize(t *testing.T) {
	sanitizer := NewStringSanitizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  world  ", "hello world"},
		{"hello\nworld", "hello world"},
		{"hello\tworld", "hello world"},
		{"hello   world", "hello world"},
	}

	for _, tt := range tests {
		result := sanitizer.Sanitize(tt.input)
		if result != tt.expected {
			t.Errorf("Sanitize(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

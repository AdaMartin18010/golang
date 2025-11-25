package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"invalid email", "invalid", false},
		{"empty email", "", false},
		{"missing @", "testexample.com", false},
		{"missing domain", "test@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"valid name", "John Doe", true},
		{"too short", "A", false},
		{"too long", string(make([]byte, 101)), false},
		{"empty", "", false},
		{"whitespace only", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateName(tt.val)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"valid", "value", true},
		{"empty", "", false},
		{"whitespace only", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateRequired(tt.val)
			assert.Equal(t, tt.want, got)
		})
	}
}

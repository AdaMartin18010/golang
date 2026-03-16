package strings

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"hello", false},
		{"  hello  ", false},
	}

	for _, tt := range tests {
		result := IsEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmpty(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"hi", 2, "hi"},
	}

	for _, tt := range tests {
		result := Truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, expected %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestContainsAny(t *testing.T) {
	if !ContainsAny("hello world", "world", "foo") {
		t.Error("Expected ContainsAny to return true")
	}
	if ContainsAny("hello", "foo", "bar") {
		t.Error("Expected ContainsAny to return false")
	}
}

func TestContainsAll(t *testing.T) {
	if !ContainsAll("hello world", "hello", "world") {
		t.Error("Expected ContainsAll to return true")
	}
	if ContainsAll("hello", "hello", "world") {
		t.Error("Expected ContainsAll to return false")
	}
}

func TestRemoveWhitespace(t *testing.T) {
	result := RemoveWhitespace("hello world")
	if result != "helloworld" {
		t.Errorf("Expected 'helloworld', got %q", result)
	}
}

func TestReverse(t *testing.T) {
	result := Reverse("hello")
	if result != "olleh" {
		t.Errorf("Expected 'olleh', got %q", result)
	}
}

func TestPadLeft(t *testing.T) {
	result := PadLeft("hello", 10, '0')
	if result != "00000hello" {
		t.Errorf("Expected '00000hello', got %q", result)
	}
}

func TestPadRight(t *testing.T) {
	result := PadRight("hello", 10, '0')
	if result != "hello00000" {
		t.Errorf("Expected 'hello00000', got %q", result)
	}
}

func TestCamelToSnake(t *testing.T) {
	result := CamelToSnake("HelloWorld")
	if result != "hello_world" {
		t.Errorf("Expected 'hello_world', got %q", result)
	}
}

func TestSnakeToCamel(t *testing.T) {
	result := SnakeToCamel("hello_world")
	if result != "helloWorld" {
		t.Errorf("Expected 'helloWorld', got %q", result)
	}
}

func TestFirstUpper(t *testing.T) {
	result := FirstUpper("hello")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got %q", result)
	}
}

func TestFirstLower(t *testing.T) {
	result := FirstLower("Hello")
	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}
}

func TestMask(t *testing.T) {
	result := Mask("1234567890", 3, 7, '*')
	if result != "123***7890" {
		t.Errorf("Expected '123***7890', got %q", result)
	}
}

func TestMaskEmail(t *testing.T) {
	result := MaskEmail("test@example.com")
	if result != "t***t@example.com" {
		t.Errorf("Expected 't***t@example.com', got %q", result)
	}
}

func TestMaskPhone(t *testing.T) {
	result := MaskPhone("13812345678")
	if result != "138****5678" {
		t.Errorf("Expected '138****5678', got %q", result)
	}
}

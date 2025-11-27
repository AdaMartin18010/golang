package converter

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultConverter_ToString(t *testing.T) {
	conv := NewConverter()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", 123, "123"},
		{"string", "hello", "hello"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float64", 123.45, "123.45"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := conv.ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultConverter_ToInt(t *testing.T) {
	conv := NewConverter()

	tests := []struct {
		name     string
		input    interface{}
		expected int
		wantErr  bool
	}{
		{"int", 123, 123, false},
		{"string", "123", 123, false},
		{"float64", 123.45, 123, false},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := conv.ToInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ToInt(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultConverter_ToJSON(t *testing.T) {
	conv := NewConverter()

	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	jsonStr, err := conv.ToJSON(data)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// 验证可以解析回来
	var result map[string]interface{}
	if err := conv.FromJSON(jsonStr, &result); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if result["name"] != "John" {
		t.Errorf("Expected name=John, got %v", result["name"])
	}
}

func TestDefaultConverter_ToMap(t *testing.T) {
	conv := NewConverter()

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}

	m, err := conv.ToMap(user)
	if err != nil {
		t.Fatalf("ToMap failed: %v", err)
	}

	if m["id"] != 1 {
		t.Errorf("Expected id=1, got %v", m["id"])
	}
	if m["name"] != "John" {
		t.Errorf("Expected name=John, got %v", m["name"])
	}
}

func TestDefaultConverter_ToSlice(t *testing.T) {
	conv := NewConverter()

	slice := []int{1, 2, 3}
	result, err := conv.ToSlice(slice)
	if err != nil {
		t.Fatalf("ToSlice failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
}

func TestDefaultConverter_Convert(t *testing.T) {
	conv := NewConverter()

	targetType := reflect.TypeOf(int64(0))
	result, err := conv.Convert("123", targetType)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if result.(int64) != 123 {
		t.Errorf("Expected 123, got %v", result)
	}
}


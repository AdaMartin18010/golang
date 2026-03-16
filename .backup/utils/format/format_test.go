package format

import (
	"testing"
	"time"
)

func TestFormatNumber(t *testing.T) {
	if FormatNumber(1234567) != "1,234,567" {
		t.Errorf("Expected '1,234,567', got %s", FormatNumber(1234567))
	}
}

func TestFormatFloat(t *testing.T) {
	result := FormatFloat(1234567.89, 2)
	if result != "1,234,567.89" {
		t.Errorf("Expected '1,234,567.89', got %s", result)
	}
}

func TestFormatPercent(t *testing.T) {
	result := FormatPercent(25, 100)
	if result != "25.00%" {
		t.Errorf("Expected '25.00%%', got %s", result)
	}
}

func TestFormatDuration(t *testing.T) {
	d := 2*time.Hour + 30*time.Minute + 45*time.Second
	result := FormatDuration(d)
	if result == "" {
		t.Error("Expected non-empty duration string")
	}
}

func TestFormatBytes(t *testing.T) {
	result := FormatBytes(1024 * 1024)
	if result == "" {
		t.Error("Expected non-empty bytes string")
	}
}

func TestFormatTimeHuman(t *testing.T) {
	now := time.Now()
	result := FormatTimeHuman(now.Add(-5 * time.Minute))
	if result == "" {
		t.Error("Expected non-empty time string")
	}
}

func TestFormatPhone(t *testing.T) {
	result := FormatPhone("13800138000")
	if result == "" {
		t.Error("Expected non-empty phone string")
	}
}

func TestFormatMask(t *testing.T) {
	result := FormatMask("1234567890", 3, 7, '*')
	if result != "123***7890" {
		t.Errorf("Expected '123***7890', got %s", result)
	}
}

func TestFormatMaskPhone(t *testing.T) {
	result := FormatMaskPhone("13800138000")
	if result == "" {
		t.Error("Expected non-empty masked phone string")
	}
}

func TestFormatPlural(t *testing.T) {
	result := FormatPlural(1, "item", "items")
	if result != "1 item" {
		t.Errorf("Expected '1 item', got %s", result)
	}
	result = FormatPlural(2, "item", "items")
	if result != "2 items" {
		t.Errorf("Expected '2 items', got %s", result)
	}
}

func TestFormatListWithAnd(t *testing.T) {
	items := []string{"apple", "banana", "orange"}
	result := FormatListWithAnd(items)
	if result == "" {
		t.Error("Expected non-empty list string")
	}
}

func TestFormatTruncate(t *testing.T) {
	result := FormatTruncate("hello world", 8, "...")
	if result != "hello..." {
		t.Errorf("Expected 'hello...', got %s", result)
	}
}

func TestFormatPadLeft(t *testing.T) {
	result := FormatPadLeft("123", 5, '0')
	if result != "00123" {
		t.Errorf("Expected '00123', got %s", result)
	}
}

func TestFormatPadRight(t *testing.T) {
	result := FormatPadRight("123", 5, '0')
	if result != "12300" {
		t.Errorf("Expected '12300', got %s", result)
	}
}

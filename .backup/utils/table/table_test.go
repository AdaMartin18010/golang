package table

import (
	"strings"
	"testing"
)

func TestNewTable(t *testing.T) {
	tbl := NewTable("Name", "Age", "City")
	if len(tbl.headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(tbl.headers))
	}
}

func TestAddRow(t *testing.T) {
	tbl := NewTable("Name", "Age")
	tbl.AddRow("Alice", "30")
	if len(tbl.rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(tbl.rows))
	}
}

func TestRender(t *testing.T) {
	tbl := NewTable("Name", "Age")
	tbl.AddRow("Alice", "30")
	tbl.AddRow("Bob", "25")
	result := tbl.Render()
	if !strings.Contains(result, "Name") {
		t.Error("Expected table to contain 'Name'")
	}
	if !strings.Contains(result, "Alice") {
		t.Error("Expected table to contain 'Alice'")
	}
}

func TestSimpleTable(t *testing.T) {
	st := NewSimpleTable("Name", "Age")
	st.AddRow("Alice", "30")
	result := st.Render()
	if !strings.Contains(result, "Name") {
		t.Error("Expected table to contain 'Name'")
	}
}

func TestFormatTable(t *testing.T) {
	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	}
	result := FormatTable(headers, rows)
	if result == "" {
		t.Error("Expected non-empty table")
	}
}

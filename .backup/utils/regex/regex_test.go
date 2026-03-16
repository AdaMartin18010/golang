package regex

import (
	"testing"
)

func TestMatch(t *testing.T) {
	matched, err := Match(`\d+`, "123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !matched {
		t.Error("Expected match")
	}
}

func TestMatchString(t *testing.T) {
	if !MatchString(`\d+`, "123") {
		t.Error("Expected match")
	}
}

func TestFind(t *testing.T) {
	match, err := Find(`\d+`, "abc123def")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match != "123" {
		t.Errorf("Expected '123', got %s", match)
	}
}

func TestFindAll(t *testing.T) {
	matches, err := FindAll(`\d+`, "abc123def456", -1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}
}

func TestReplace(t *testing.T) {
	result, err := Replace(`\d+`, "abc123def", "XXX")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "abcXXXdef" {
		t.Errorf("Expected 'abcXXXdef', got %s", result)
	}
}

func TestSplit(t *testing.T) {
	result, err := Split(`\s+`, "a b c", -1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 parts, got %d", len(result))
	}
}

func TestCount(t *testing.T) {
	count, err := Count(`\d+`, "abc123def456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid(`\d+`) {
		t.Error("Expected valid pattern")
	}
	if IsValid(`[`) {
		t.Error("Expected invalid pattern")
	}
}

func TestEscape(t *testing.T) {
	escaped := Escape(`.*+?^${}()|[]\`)
	if escaped == "" {
		t.Error("Expected escaped string")
	}
}

func TestMatchEmail(t *testing.T) {
	if !MatchEmail("test@example.com") {
		t.Error("Expected valid email")
	}
	if MatchEmail("invalid") {
		t.Error("Expected invalid email")
	}
}

func TestMatchPhone(t *testing.T) {
	if !MatchPhone("13800138000") {
		t.Error("Expected valid phone")
	}
	if MatchPhone("12345678901") {
		t.Error("Expected invalid phone")
	}
}

func TestMatchURL(t *testing.T) {
	if !MatchURL("https://example.com") {
		t.Error("Expected valid URL")
	}
	if MatchURL("invalid") {
		t.Error("Expected invalid URL")
	}
}

func TestExtractGroups(t *testing.T) {
	pattern := `(?P<name>\w+)\s+(?P<age>\d+)`
	s := "John 30"
	groups, err := ExtractGroups(pattern, s)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if groups["name"] != "John" {
		t.Errorf("Expected 'John', got %s", groups["name"])
	}
	if groups["age"] != "30" {
		t.Errorf("Expected '30', got %s", groups["age"])
	}
}

func TestRemove(t *testing.T) {
	result, err := Remove(`\d+`, "abc123def")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "abcdef" {
		t.Errorf("Expected 'abcdef', got %s", result)
	}
}

func TestReplaceN(t *testing.T) {
	result, err := ReplaceN(`\d+`, "abc123def456", "XXX", 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !MatchString(`abcXXXdef456`, result) {
		t.Errorf("Expected 'abcXXXdef456', got %s", result)
	}
}

func TestGetMatches(t *testing.T) {
	matches, err := GetMatches(`\d+`, "abc123def456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}
}

func TestCompileWithOptions(t *testing.T) {
	re, err := CompileWithOptions(`abc`, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !re.MatchString("ABC") {
		t.Error("Expected case-insensitive match")
	}
}

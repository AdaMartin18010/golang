package url

import (
	"strings"
	"testing"
)

func TestBuildURL(t *testing.T) {
	baseURL := "https://api.example.com"
	path := "/users"
	params := map[string]string{
		"page":  "1",
		"limit": "10",
	}

	result, err := BuildURL(baseURL, path, params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "/users") {
		t.Error("Expected URL to contain /users")
	}
	if !strings.Contains(result, "page=1") {
		t.Error("Expected URL to contain page=1")
	}
}

func TestAddQuery(t *testing.T) {
	rawURL := "https://api.example.com/users"
	result, err := AddQuery(rawURL, "page", "1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "page=1") {
		t.Error("Expected URL to contain page=1")
	}
}

func TestRemoveQuery(t *testing.T) {
	rawURL := "https://api.example.com/users?page=1&limit=10"
	result, err := RemoveQuery(rawURL, "page")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if strings.Contains(result, "page=1") {
		t.Error("Expected URL not to contain page=1")
	}
}

func TestGetQuery(t *testing.T) {
	rawURL := "https://api.example.com/users?page=1"
	value, err := GetQuery(rawURL, "page")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if value != "1" {
		t.Errorf("Expected '1', got %s", value)
	}
}

func TestEncodeDecode(t *testing.T) {
	original := "hello world"
	encoded := Encode(original)
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if decoded != original {
		t.Errorf("Expected %s, got %s", original, decoded)
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("https://api.example.com") {
		t.Error("Expected valid URL")
	}

	if IsValid("not a url") {
		t.Error("Expected invalid URL")
	}
}

func TestIsAbsolute(t *testing.T) {
	if !IsAbsolute("https://api.example.com") {
		t.Error("Expected absolute URL")
	}

	if IsAbsolute("/relative/path") {
		t.Error("Expected relative URL")
	}
}

func TestGetDomain(t *testing.T) {
	rawURL := "https://api.example.com:8080/users"
	domain, err := GetDomain(rawURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if domain != "api.example.com" {
		t.Errorf("Expected 'api.example.com', got %s", domain)
	}
}

func TestGetPort(t *testing.T) {
	rawURL := "https://api.example.com:8080/users"
	port, err := GetPort(rawURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if port != "8080" {
		t.Errorf("Expected '8080', got %s", port)
	}
}

func TestBuildQueryString(t *testing.T) {
	params := map[string]string{
		"page":  "1",
		"limit": "10",
	}

	queryString := BuildQueryString(params)
	if !strings.Contains(queryString, "page=1") {
		t.Error("Expected query string to contain page=1")
	}
}

func TestParseQueryString(t *testing.T) {
	queryString := "page=1&limit=10"
	params, err := ParseQueryString(queryString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if params["page"] != "1" {
		t.Errorf("Expected '1', got %s", params["page"])
	}
}

func TestIsHTTPS(t *testing.T) {
	if !IsHTTPS("https://api.example.com") {
		t.Error("Expected HTTPS URL")
	}

	if IsHTTPS("http://api.example.com") {
		t.Error("Expected not HTTPS URL")
	}
}

func TestMaskURL(t *testing.T) {
	rawURL := "https://user:pass@api.example.com/users?token=secret123"
	masked, err := MaskURL(rawURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if strings.Contains(masked, "pass") {
		t.Error("Expected password to be masked")
	}
	if strings.Contains(masked, "secret123") {
		t.Error("Expected token to be masked")
	}
}

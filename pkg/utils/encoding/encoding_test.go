package encoding

import (
	"strings"
	"testing"
)

func TestBase64EncodeDecode(t *testing.T) {
	data := []byte("hello world")
	encoded := Base64Encode(data)
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decoded))
	}
}

func TestBase64URLEncodeDecode(t *testing.T) {
	data := []byte("hello world")
	encoded := Base64URLEncode(data)
	decoded, err := Base64URLDecode(encoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decoded))
	}
}

func TestHexEncodeDecode(t *testing.T) {
	data := []byte("hello world")
	encoded := HexEncode(data)
	decoded, err := HexDecode(encoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decoded))
	}
}

func TestStringToInt(t *testing.T) {
	result, err := StringToInt("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}
}

func TestIntToString(t *testing.T) {
	result := IntToString(123)
	if result != "123" {
		t.Errorf("Expected '123', got %s", result)
	}
}

func TestFloat64ToString(t *testing.T) {
	result := Float64ToString(123.456)
	if result != "123.456" {
		t.Errorf("Expected '123.456', got %s", result)
	}
}

func TestStringToFloat64(t *testing.T) {
	result, err := StringToFloat64("123.456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != 123.456 {
		t.Errorf("Expected 123.456, got %f", result)
	}
}

func TestBoolToString(t *testing.T) {
	if BoolToString(true) != "true" {
		t.Error("Expected 'true'")
	}
	if BoolToString(false) != "false" {
		t.Error("Expected 'false'")
	}
}

func TestStringToBool(t *testing.T) {
	result, err := StringToBool("true")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !result {
		t.Error("Expected true")
	}
}

func TestJSONEncodeDecode(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	encoded, err := JSONEncode(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var decoded map[string]interface{}
	err = JSONDecode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if decoded["name"] != "test" {
		t.Errorf("Expected 'test', got %v", decoded["name"])
	}
}

func TestEscapeString(t *testing.T) {
	result := EscapeString("<script>alert('xss')</script>")
	if !strings.Contains(result, "&lt;") {
		t.Error("Expected escaped string")
	}
}

func TestUnescapeString(t *testing.T) {
	escaped := "&lt;script&gt;"
	result := UnescapeString(escaped)
	if result != "<script>" {
		t.Errorf("Expected '<script>', got %s", result)
	}
}

func TestIsBase64(t *testing.T) {
	if !IsBase64(Base64Encode([]byte("test"))) {
		t.Error("Expected valid Base64")
	}
	if IsBase64("invalid!") {
		t.Error("Expected invalid Base64")
	}
}

func TestIsHex(t *testing.T) {
	if !IsHex(HexEncode([]byte("test"))) {
		t.Error("Expected valid Hex")
	}
	if IsHex("invalid!") {
		t.Error("Expected invalid Hex")
	}
}

func TestIsJSON(t *testing.T) {
	if !IsJSON(`{"name":"test"}`) {
		t.Error("Expected valid JSON")
	}
	if IsJSON("invalid json") {
		t.Error("Expected invalid JSON")
	}
}

func TestBase64EncodeString(t *testing.T) {
	result := Base64EncodeString("hello")
	decoded, err := Base64DecodeString(result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if decoded != "hello" {
		t.Errorf("Expected 'hello', got %s", decoded)
	}
}

func TestHexEncodeString(t *testing.T) {
	result := HexEncodeString("hello")
	decoded, err := HexDecodeString(result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if decoded != "hello" {
		t.Errorf("Expected 'hello', got %s", decoded)
	}
}

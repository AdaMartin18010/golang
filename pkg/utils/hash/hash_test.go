package hash

import (
	"testing"
)

func TestMD5(t *testing.T) {
	data := []byte("hello world")
	hash := MD5(data)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestMD5String(t *testing.T) {
	hash := MD5String("hello world")
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestSHA1(t *testing.T) {
	data := []byte("hello world")
	hash := SHA1(data)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestSHA256(t *testing.T) {
	data := []byte("hello world")
	hash := SHA256(data)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestSHA512(t *testing.T) {
	data := []byte("hello world")
	hash := SHA512(data)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestCRC32(t *testing.T) {
	data := []byte("hello world")
	hash := CRC32(data)
	if hash == 0 {
		t.Error("Expected non-zero hash")
	}
}

func TestCRC64(t *testing.T) {
	data := []byte("hello world")
	hash := CRC64(data)
	if hash == 0 {
		t.Error("Expected non-zero hash")
	}
}

func TestFNV32(t *testing.T) {
	data := []byte("hello world")
	hash := FNV32(data)
	if hash == 0 {
		t.Error("Expected non-zero hash")
	}
}

func TestFNV64(t *testing.T) {
	data := []byte("hello world")
	hash := FNV64(data)
	if hash == 0 {
		t.Error("Expected non-zero hash")
	}
}

func TestHash(t *testing.T) {
	data := []byte("hello world")
	hash, err := Hash(data, "md5")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestHashString(t *testing.T) {
	hash, err := HashString("hello world", "sha256")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestCompareHash(t *testing.T) {
	data := []byte("hello world")
	hash1 := MD5(data)
	hash2 := MD5(data)
	if !CompareHash(hash1, hash2) {
		t.Error("Expected hashes to be equal")
	}
}

func TestVerifyHash(t *testing.T) {
	data := []byte("hello world")
	expectedHash := MD5(data)
	valid, err := VerifyHash(data, "md5", expectedHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected hash to be valid")
	}
}

func TestVerifyHashString(t *testing.T) {
	s := "hello world"
	expectedHash := MD5String(s)
	valid, err := VerifyHashString(s, "md5", expectedHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected hash to be valid")
	}
}

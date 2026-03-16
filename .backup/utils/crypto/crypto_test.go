package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
	if hash == password {
		t.Error("Hash should not equal password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword"
	hash, _ := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Error("Expected password check to succeed")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Error("Expected password check to fail")
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	plaintext := "test message"
	ciphertext, err := AESEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decrypted, err := AESDecrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestGenerateAESKey(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if key == "" {
		t.Error("Expected non-empty key")
	}
}

func TestDeriveKey(t *testing.T) {
	password := "testpassword"
	salt := "testsalt"
	key := DeriveKey(password, salt, 10000)
	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}
}

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt(16)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if salt == "" {
		t.Error("Expected non-empty salt")
	}
}

func TestSHA256Hash(t *testing.T) {
	data := "test data"
	hash := SHA256Hash(data)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}
}

func TestHMAC(t *testing.T) {
	key := "testkey"
	message := "test message"
	hmac := HMAC(key, message)
	if hmac == "" {
		t.Error("Expected non-empty HMAC")
	}
}

func TestRandomBytes(t *testing.T) {
	bytes, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(bytes) != 16 {
		t.Errorf("Expected length 16, got %d", len(bytes))
	}
}

func TestRandomString(t *testing.T) {
	str, err := RandomString(16)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(str) != 16 {
		t.Errorf("Expected length 16, got %d", len(str))
	}
}

func TestAESEncryptor(t *testing.T) {
	key, _ := GenerateAESKey()
	encryptor, err := NewAESEncryptor(key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	plaintext := "test message"
	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

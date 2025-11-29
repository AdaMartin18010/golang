package security

import (
	"testing"
)

func TestAES256Encryptor_EncryptDecrypt(t *testing.T) {
	// 生成测试密钥（32 字节）
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := NewAES256Encryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := "Hello, World! This is a test message."

	// 加密
	ciphertext, err := encryptor.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Error("Ciphertext should not be empty")
	}

	// 解密
	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Expected '%s', got '%s'", plaintext, string(decrypted))
	}
}

func TestAES256Encryptor_EncryptDecryptString(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := NewAES256Encryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := "Hello, World!"

	// 加密字符串
	ciphertext, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt string: %v", err)
	}

	if ciphertext == "" {
		t.Error("Ciphertext should not be empty")
	}

	// 解密字符串
	decrypted, err := encryptor.DecryptString(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt string: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestAES256Encryptor_InvalidKey(t *testing.T) {
	// 测试无效密钥长度
	invalidKey := make([]byte, 16) // 应该是 32 字节
	_, err := NewAES256Encryptor(invalidKey)
	if err == nil {
		t.Error("Should return error for invalid key length")
	}
}

func TestAES256EncryptorFromString(t *testing.T) {
	keyString := "my-secret-key-12345"
	encryptor, err := NewAES256EncryptorFromString(keyString)
	if err != nil {
		t.Fatalf("Failed to create encryptor from string: %v", err)
	}

	plaintext := "Test message"
	ciphertext, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := encryptor.DecryptString(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestFieldEncryptor(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, _ := NewAES256Encryptor(key)
	fieldEncryptor := NewFieldEncryptor(encryptor)

	value := "sensitive-data"
	encrypted, err := fieldEncryptor.EncryptField(value)
	if err != nil {
		t.Fatalf("Failed to encrypt field: %v", err)
	}

	decrypted, err := fieldEncryptor.DecryptField(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt field: %v", err)
	}

	if decrypted != value {
		t.Errorf("Expected '%s', got '%s'", value, decrypted)
	}
}

func TestDataMasker_MaskEmail(t *testing.T) {
	masker := NewDataMasker()

	tests := []struct {
		input    string
		expected string
	}{
		{"test@example.com", "t***t@***.com"},
		{"ab@example.com", "***@***"},
		{"", ""},
	}

	for _, tt := range tests {
		result := masker.MaskEmail(tt.input)
		if result != tt.expected {
			t.Errorf("MaskEmail(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestDataMasker_MaskPhone(t *testing.T) {
	masker := NewDataMasker()

	tests := []struct {
		input    string
		expected string
	}{
		{"13812345678", "138****5678"},
		{"1234", "****"},
		{"", ""},
	}

	for _, tt := range tests {
		result := masker.MaskPhone(tt.input)
		if result != tt.expected {
			t.Errorf("MaskPhone(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestDataMasker_MaskIDCard(t *testing.T) {
	masker := NewDataMasker()

	tests := []struct {
		input    string
		expected string
	}{
		{"123456789012345678", "1234********5678"},
		{"1234", "********"},
		{"", ""},
	}

	for _, tt := range tests {
		result := masker.MaskIDCard(tt.input)
		if result != tt.expected {
			t.Errorf("MaskIDCard(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestDataMasker_MaskName(t *testing.T) {
	masker := NewDataMasker()

	tests := []struct {
		input    string
		expected string
	}{
		{"张三", "张*"},
		{"李四", "李*"},
		{"王五", "王*"},
		{"", ""},
		{"A", "*"},
		{"AB", "A*"},
		{"ABC", "A**C"},
	}

	for _, tt := range tests {
		result := masker.MaskName(tt.input)
		if result != tt.expected {
			t.Errorf("MaskName(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

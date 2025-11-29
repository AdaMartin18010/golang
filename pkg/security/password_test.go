package security

import (
	"testing"
)

func TestPasswordHasher_Hash(t *testing.T) {
	hasher := NewPasswordHasher(DefaultPasswordHashConfig())

	password := "test-password-123"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should be different from password")
	}
}

func TestPasswordHasher_Verify(t *testing.T) {
	hasher := NewPasswordHasher(DefaultPasswordHashConfig())

	password := "test-password-123"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 验证正确密码
	valid, err := hasher.Verify(password, hash)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if !valid {
		t.Error("Password verification should succeed")
	}

	// 验证错误密码
	valid, err = hasher.Verify("wrong-password", hash)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if valid {
		t.Error("Password verification should fail for wrong password")
	}
}

func TestPasswordHasher_NeedsRehash(t *testing.T) {
	hasher1 := NewPasswordHasher(DefaultPasswordHashConfig())
	hasher2 := NewPasswordHasher(PasswordHashConfig{
		Memory:      128 * 1024, // 不同的配置
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})

	password := "test-password"
	hash, err := hasher1.Hash(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 使用不同配置检查
	needsRehash := hasher2.NeedsRehash(hash)
	if !needsRehash {
		t.Error("Should need rehash with different config")
	}

	// 使用相同配置检查
	needsRehash = hasher1.NeedsRehash(hash)
	if needsRehash {
		t.Error("Should not need rehash with same config")
	}
}

func TestPasswordValidator_Validate(t *testing.T) {
	validator := NewPasswordValidator(DefaultPasswordValidatorConfig())

	tests := []struct {
		password string
		valid    bool
	}{
		{"ValidPass123", true},
		{"short", false},           // 太短
		{"nouppercase123", false},  // 没有大写
		{"NOLOWERCASE123", false},  // 没有小写
		{"NoDigits", false},         // 没有数字
		{"ValidPass123", true},
	}

	for _, tt := range tests {
		err := validator.Validate(tt.password)
		if tt.valid && err != nil {
			t.Errorf("Password '%s' should be valid, got error: %v", tt.password, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Password '%s' should be invalid", tt.password)
		}
	}
}

func TestPasswordValidator_MinLength(t *testing.T) {
	validator := NewPasswordValidator(PasswordValidatorConfig{
		MinLength: 10,
		MaxLength: 128,
	})

	err := validator.Validate("Short1")
	if err == nil {
		t.Error("Should fail for password shorter than min length")
	}

	err = validator.Validate("LongEnough1")
	if err != nil {
		t.Errorf("Should pass for password meeting min length: %v", err)
	}
}

func TestPasswordValidator_MaxLength(t *testing.T) {
	validator := NewPasswordValidator(PasswordValidatorConfig{
		MinLength: 8,
		MaxLength: 10,
	})

	longPassword := "VeryLongPassword123"
	err := validator.Validate(longPassword)
	if err == nil {
		t.Error("Should fail for password longer than max length")
	}
}

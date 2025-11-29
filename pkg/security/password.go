package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash 无效的哈希
	ErrInvalidHash = errors.New("invalid hash format")
	// ErrIncompatibleVersion 不兼容的版本
	ErrIncompatibleVersion = errors.New("incompatible version")
)

// PasswordHasher 密码哈希器
type PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// PasswordHashConfig 密码哈希配置
type PasswordHashConfig struct {
	Memory      uint32 // 内存使用（KB）
	Iterations  uint32 // 迭代次数
	Parallelism uint8  // 并行度
	SaltLength  uint32 // 盐长度
	KeyLength   uint32 // 密钥长度
}

// DefaultPasswordHashConfig 默认密码哈希配置
func DefaultPasswordHashConfig() PasswordHashConfig {
	return PasswordHashConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// NewPasswordHasher 创建密码哈希器
func NewPasswordHasher(config PasswordHashConfig) *PasswordHasher {
	if config.Memory == 0 {
		config = DefaultPasswordHashConfig()
	}

	return &PasswordHasher{
		memory:      config.Memory,
		iterations:  config.Iterations,
		parallelism: config.Parallelism,
		saltLength:  config.SaltLength,
		keyLength:   config.KeyLength,
	}
}

// Hash 哈希密码
func (h *PasswordHasher) Hash(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// 使用 Argon2id 哈希密码
	hash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallelism, h.keyLength)

	// 编码为 base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 返回格式：argon2id$v=19$m=65536,t=3,p=2$salt$hash
	return fmt.Sprintf("argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.memory,
		h.iterations,
		h.parallelism,
		b64Salt,
		b64Hash,
	), nil
}

// Verify 验证密码
func (h *PasswordHasher) Verify(password, encodedHash string) (bool, error) {
	// 解析哈希字符串
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, ErrInvalidHash
	}

	// 检查算法
	if parts[0] != "argon2id" {
		return false, ErrInvalidHash
	}

	// 解析版本
	var version int
	if _, err := fmt.Sscanf(parts[1], "v=%d", &version); err != nil {
		return false, ErrInvalidHash
	}
	if version != argon2.Version {
		return false, ErrIncompatibleVersion
	}

	// 解析参数
	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, ErrInvalidHash
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, ErrInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, ErrInvalidHash
	}

	// 计算密码哈希
	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	// 使用常量时间比较
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// NeedsRehash 检查是否需要重新哈希
func (h *PasswordHasher) NeedsRehash(encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return true
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return true
	}

	// 检查参数是否匹配
	return memory != h.memory || iterations != h.iterations || parallelism != h.parallelism
}

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength    int
	maxLength    int
	requireUpper bool
	requireLower bool
	requireDigit bool
	requireSpecial bool
}

// PasswordValidatorConfig 密码验证器配置
type PasswordValidatorConfig struct {
	MinLength     int
	MaxLength     int
	RequireUpper  bool
	RequireLower  bool
	RequireDigit  bool
	RequireSpecial bool
}

// DefaultPasswordValidatorConfig 默认密码验证器配置
func DefaultPasswordValidatorConfig() PasswordValidatorConfig {
	return PasswordValidatorConfig{
		MinLength:     8,
		MaxLength:     128,
		RequireUpper:  true,
		RequireLower:  true,
		RequireDigit:  true,
		RequireSpecial: false,
	}
}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator(config PasswordValidatorConfig) *PasswordValidator {
	if config.MinLength == 0 {
		config = DefaultPasswordValidatorConfig()
	}

	return &PasswordValidator{
		minLength:     config.MinLength,
		maxLength:     config.MaxLength,
		requireUpper:  config.RequireUpper,
		requireLower:  config.RequireLower,
		requireDigit:  config.RequireDigit,
		requireSpecial: config.RequireSpecial,
	}
}

// Validate 验证密码
func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.minLength {
		return fmt.Errorf("password must be at least %d characters", v.minLength)
	}

	if len(password) > v.maxLength {
		return fmt.Errorf("password must be at most %d characters", v.maxLength)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if v.requireUpper && !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if v.requireLower && !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if v.requireDigit && !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	if v.requireSpecial && !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

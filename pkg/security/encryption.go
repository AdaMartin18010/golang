package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrInvalidKey 无效的密钥
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrInvalidCiphertext 无效的密文
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
)

// AES256Encryptor AES-256 加密器
type AES256Encryptor struct {
	key []byte
}

// NewAES256Encryptor 创建 AES-256 加密器
func NewAES256Encryptor(key []byte) (*AES256Encryptor, error) {
	// 密钥必须是 32 字节（256 位）
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	return &AES256Encryptor{key: key}, nil
}

// NewAES256EncryptorFromString 从字符串创建加密器（使用 SHA-256 哈希）
func NewAES256EncryptorFromString(keyString string) (*AES256Encryptor, error) {
	hash := sha256.Sum256([]byte(keyString))
	return NewAES256Encryptor(hash[:])
}

// Encrypt 加密数据
func (e *AES256Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密数据
func (e *AES256Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 提取 nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString 加密字符串（返回 Base64 编码）
func (e *AES256Encryptor) EncryptString(plaintext string) (string, error) {
	ciphertext, err := e.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString 解密字符串（从 Base64 解码）
func (e *AES256Encryptor) DecryptString(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := e.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// FieldEncryptor 字段级加密器
type FieldEncryptor struct {
	encryptor *AES256Encryptor
}

// NewFieldEncryptor 创建字段级加密器
func NewFieldEncryptor(encryptor *AES256Encryptor) *FieldEncryptor {
	return &FieldEncryptor{encryptor: encryptor}
}

// EncryptField 加密字段
func (f *FieldEncryptor) EncryptField(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	return f.encryptor.EncryptString(value)
}

// DecryptField 解密字段
func (f *FieldEncryptor) DecryptField(encryptedValue string) (string, error) {
	if encryptedValue == "" {
		return "", nil
	}
	return f.encryptor.DecryptString(encryptedValue)
}

// DataMasker 数据脱敏器
type DataMasker struct{}

// NewDataMasker 创建数据脱敏器
func NewDataMasker() *DataMasker {
	return &DataMasker{}
}

// MaskEmail 脱敏邮箱
func (m *DataMasker) MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := splitEmail(email)
	if len(parts) != 2 {
		return "***@***"
	}

	local, domain := parts[0], parts[1]
	if len(local) <= 2 {
		return "***@***"
	}

	maskedLocal := local[:1] + "***" + local[len(local)-1:]
	return maskedLocal + "@" + maskDomain(domain)
}

// MaskPhone 脱敏手机号
func (m *DataMasker) MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}

	if len(phone) <= 4 {
		return "****"
	}

	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskIDCard 脱敏身份证号
func (m *DataMasker) MaskIDCard(idCard string) string {
	if idCard == "" {
		return ""
	}

	if len(idCard) <= 8 {
		return "********"
	}

	return idCard[:4] + "********" + idCard[len(idCard)-4:]
}

// MaskName 脱敏姓名
func (m *DataMasker) MaskName(name string) string {
	if name == "" {
		return ""
	}

	if len(name) <= 1 {
		return "*"
	}

	if len(name) == 2 {
		return name[:1] + "*"
	}

	return name[:1] + "**" + name[len(name)-1:]
}

// splitEmail 分割邮箱地址
func splitEmail(email string) []string {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return []string{email[:i], email[i+1:]}
		}
	}
	return nil
}

// maskDomain 脱敏域名
func maskDomain(domain string) string {
	if domain == "" {
		return "***"
	}

	parts := splitDomain(domain)
	if len(parts) == 0 {
		return "***"
	}

	if len(parts) == 1 {
		return maskString(parts[0])
	}

	// 保留最后一个部分（通常是 .com, .org 等）
	lastPart := parts[len(parts)-1]
	maskedParts := make([]string, len(parts)-1)
	for i := 0; i < len(parts)-1; i++ {
		maskedParts[i] = maskString(parts[i])
	}
	maskedParts = append(maskedParts, lastPart)

	result := ""
	for i, part := range maskedParts {
		if i > 0 {
			result += "."
		}
		result += part
	}
	return result
}

// splitDomain 分割域名
func splitDomain(domain string) []string {
	parts := make([]string, 0)
	start := 0
	for i := 0; i < len(domain); i++ {
		if domain[i] == '.' {
			if i > start {
				parts = append(parts, domain[start:i])
			}
			start = i + 1
		}
	}
	if start < len(domain) {
		parts = append(parts, domain[start:])
	}
	return parts
}

// maskString 脱敏字符串
func maskString(s string) string {
	if len(s) <= 2 {
		return "***"
	}
	return s[:1] + "***" + s[len(s)-1:]
}

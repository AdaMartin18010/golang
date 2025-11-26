package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

// HashPassword 使用bcrypt哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AESEncrypt AES加密
func AESEncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES解密
func AESDecrypt(ciphertext, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateAESKey 生成AES密钥（32字节，256位）
func GenerateAESKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// DeriveKey 从密码派生密钥（PBKDF2）
func DeriveKey(password, salt string, iterations int) []byte {
	if iterations == 0 {
		iterations = 10000
	}
	return pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)
}

// GenerateSalt 生成随机盐值
func GenerateSalt(length int) (string, error) {
	if length == 0 {
		length = 16
	}
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// SHA256Hash SHA256哈希
func SHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// SHA256HashBytes SHA256哈希（返回字节）
func SHA256HashBytes(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HMAC 计算HMAC（简化版，使用SHA256）
func HMAC(key, message string) string {
	h := sha256.New()
	h.Write([]byte(key + message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// RandomBytes 生成随机字节
func RandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// RandomString 生成随机字符串（Base64编码）
func RandomString(length int) (string, error) {
	bytes, err := RandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// Encryptor 加密器接口
type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// AESEncryptor AES加密器
type AESEncryptor struct {
	key string
}

// NewAESEncryptor 创建AES加密器
func NewAESEncryptor(key string) (*AESEncryptor, error) {
	// 验证密钥长度（AES-256需要32字节）
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("key must be 32 bytes (256 bits)")
	}
	return &AESEncryptor{key: key}, nil
}

// Encrypt 加密
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	return AESEncrypt(plaintext, e.key)
}

// Decrypt 解密
func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
	return AESDecrypt(ciphertext, e.key)
}

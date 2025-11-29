package security

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrInvalidConfig 无效配置
	ErrInvalidConfig = errors.New("invalid security configuration")
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	// 加密配置
	Encryption EncryptionConfig

	// 密钥管理配置
	KeyManagement KeyManagementConfig

	// 会话配置
	Session SessionConfig

	// CSRF 配置
	CSRF CSRFConfig

	// 速率限制配置
	RateLimit RateLimiterConfig

	// 密码配置
	Password PasswordHashConfig

	// 审计日志配置
	Audit AuditLogConfig

	// TLS 配置
	TLS TLSConfig
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Algorithm string // 加密算法
	KeySize   int    // 密钥大小
}

// KeyManagementConfig 密钥管理配置
type KeyManagementConfig struct {
	Provider     string        // 密钥提供者（memory, vault）
	RotationInterval time.Duration // 轮换间隔
	MaxKeyAge    time.Duration // 最大密钥年龄
}

// AuditLogConfig 审计日志配置
type AuditLogConfig struct {
	Enabled     bool          // 是否启用
	Retention   time.Duration // 保留时间
	StorageType string        // 存储类型（memory, database, file）
}

// TLSConfig TLS 配置
type TLSConfig struct {
	Enabled            bool
	MinVersion         uint16
	MaxVersion         uint16
	CipherSuites       []uint16
	InsecureSkipVerify bool
	CertFile           string
	KeyFile            string
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		Encryption: EncryptionConfig{
			Algorithm: "AES-256-GCM",
			KeySize:   256,
		},
		KeyManagement: KeyManagementConfig{
			Provider:        "memory",
			RotationInterval: 30 * 24 * time.Hour,
			MaxKeyAge:       90 * 24 * time.Hour,
		},
		Session: DefaultSessionConfig(),
		CSRF:    DefaultCSRFConfig(),
		RateLimit: RateLimiterConfig{
			Limit:  100,
			Window: 1 * time.Minute,
		},
		Password: DefaultPasswordHashConfig(),
		Audit: AuditLogConfig{
			Enabled:     true,
			Retention:   90 * 24 * time.Hour,
			StorageType: "memory",
		},
		TLS: TLSConfig{
			Enabled:    false,
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		},
	}
}

// Validate 验证配置
func (c *SecurityConfig) Validate() error {
	// 验证加密配置
	if c.Encryption.Algorithm == "" {
		return fmt.Errorf("%w: encryption algorithm is required", ErrInvalidConfig)
	}

	if c.Encryption.KeySize != 128 && c.Encryption.KeySize != 192 && c.Encryption.KeySize != 256 {
		return fmt.Errorf("%w: invalid encryption key size", ErrInvalidConfig)
	}

	// 验证密钥管理配置
	if c.KeyManagement.Provider == "" {
		return fmt.Errorf("%w: key management provider is required", ErrInvalidConfig)
	}

	// 验证 TLS 配置
	if c.TLS.Enabled {
		if c.TLS.CertFile == "" || c.TLS.KeyFile == "" {
			return fmt.Errorf("%w: TLS certificate and key files are required", ErrInvalidConfig)
		}
	}

	return nil
}

// SecurityConfigManager 安全配置管理器
type SecurityConfigManager struct {
	config *SecurityConfig
	mu     sync.RWMutex
}

// NewSecurityConfigManager 创建安全配置管理器
func NewSecurityConfigManager(config SecurityConfig) (*SecurityConfigManager, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &SecurityConfigManager{
		config: &config,
	}, nil
}

// GetConfig 获取配置
func (m *SecurityConfigManager) GetConfig() SecurityConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return *m.config
}

// UpdateConfig 更新配置
func (m *SecurityConfigManager) UpdateConfig(config SecurityConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}

	m.mu.Lock()
	m.config = &config
	m.mu.Unlock()

	return nil
}

// GetEncryptionConfig 获取加密配置
func (m *SecurityConfigManager) GetEncryptionConfig() EncryptionConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.Encryption
}

// GetKeyManagementConfig 获取密钥管理配置
func (m *SecurityConfigManager) GetKeyManagementConfig() KeyManagementConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.KeyManagement
}

// GetSessionConfig 获取会话配置
func (m *SecurityConfigManager) GetSessionConfig() SessionConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.Session
}

// GetCSRFConfig 获取 CSRF 配置
func (m *SecurityConfigManager) GetCSRFConfig() CSRFConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.CSRF
}

// GetRateLimitConfig 获取速率限制配置
func (m *SecurityConfigManager) GetRateLimitConfig() RateLimiterConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.RateLimit
}

// GetPasswordConfig 获取密码配置
func (m *SecurityConfigManager) GetPasswordConfig() PasswordHashConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.Password
}

// GetAuditConfig 获取审计日志配置
func (m *SecurityConfigManager) GetAuditConfig() AuditLogConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.Audit
}

// GetTLSConfig 获取 TLS 配置
func (m *SecurityConfigManager) GetTLSConfig() TLSConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.config.TLS
}

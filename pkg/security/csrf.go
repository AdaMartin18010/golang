package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrInvalidCSRFToken 无效的 CSRF 令牌
	ErrInvalidCSRFToken = errors.New("invalid CSRF token")
	// ErrCSRFTokenExpired CSRF 令牌已过期
	ErrCSRFTokenExpired = errors.New("CSRF token expired")
)

// CSRFProtection CSRF 防护
type CSRFProtection struct {
	tokens      map[string]*CSRFToken
	mu          sync.RWMutex
	secret      []byte
	tokenLength int
	expiry      time.Duration
	cleanup     *time.Ticker
	stopCleanup chan struct{}
}

// CSRFToken CSRF 令牌
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFConfig CSRF 配置
type CSRFConfig struct {
	Secret      []byte        // 密钥
	TokenLength int           // 令牌长度
	Expiry      time.Duration // 过期时间
}

// DefaultCSRFConfig 默认 CSRF 配置
func DefaultCSRFConfig() CSRFConfig {
	secret := make([]byte, 32)
	rand.Read(secret)

	return CSRFConfig{
		Secret:      secret,
		TokenLength: 32,
		Expiry:      24 * time.Hour,
	}
}

// NewCSRFProtection 创建 CSRF 防护
func NewCSRFProtection(config CSRFConfig) *CSRFProtection {
	if config.Secret == nil {
		config = DefaultCSRFConfig()
	}
	if config.TokenLength == 0 {
		config.TokenLength = 32
	}
	if config.Expiry == 0 {
		config.Expiry = 24 * time.Hour
	}

	csrf := &CSRFProtection{
		tokens:      make(map[string]*CSRFToken),
		secret:      config.Secret,
		tokenLength: config.TokenLength,
		expiry:      config.Expiry,
		stopCleanup: make(chan struct{}),
	}

	// 启动清理协程
	csrf.cleanup = time.NewTicker(1 * time.Hour)
	go csrf.cleanupExpired()

	return csrf
}

// GenerateToken 生成 CSRF 令牌
func (c *CSRFProtection) GenerateToken(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("session ID is required")
	}

	// 生成随机令牌
	tokenBytes := make([]byte, c.tokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// 存储令牌
	c.mu.Lock()
	c.tokens[sessionID] = &CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(c.expiry),
	}
	c.mu.Unlock()

	return token, nil
}

// ValidateToken 验证 CSRF 令牌
func (c *CSRFProtection) ValidateToken(sessionID, token string) error {
	if sessionID == "" || token == "" {
		return ErrInvalidCSRFToken
	}

	c.mu.RLock()
	storedToken, exists := c.tokens[sessionID]
	c.mu.RUnlock()

	if !exists {
		return ErrInvalidCSRFToken
	}

	if storedToken.Token != token {
		return ErrInvalidCSRFToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		// 删除过期令牌
		c.mu.Lock()
		delete(c.tokens, sessionID)
		c.mu.Unlock()
		return ErrCSRFTokenExpired
	}

	return nil
}

// RevokeToken 撤销 CSRF 令牌
func (c *CSRFProtection) RevokeToken(sessionID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.tokens, sessionID)
	return nil
}

// cleanupExpired 清理过期令牌
func (c *CSRFProtection) cleanupExpired() {
	for {
		select {
		case <-c.cleanup.C:
			c.mu.Lock()
			now := time.Now()
			for sessionID, token := range c.tokens {
				if now.After(token.ExpiresAt) {
					delete(c.tokens, sessionID)
				}
			}
			c.mu.Unlock()

		case <-c.stopCleanup:
			c.cleanup.Stop()
			return
		}
	}
}

// Shutdown 关闭 CSRF 防护
func (c *CSRFProtection) Shutdown() error {
	close(c.stopCleanup)
	return nil
}

package oauth2

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryTokenStore 内存令牌存储（用于测试和开发）
type MemoryTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*Token
	// 按刷新令牌索引
	refreshTokens map[string]string
}

// NewMemoryTokenStore 创建内存令牌存储
func NewMemoryTokenStore() *MemoryTokenStore {
	store := &MemoryTokenStore{
		tokens:        make(map[string]*Token),
		refreshTokens: make(map[string]string),
	}

	// 启动清理协程
	go store.cleanup()

	return store
}

// Save 保存令牌
func (s *MemoryTokenStore) Save(ctx context.Context, token *Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[token.AccessToken] = token
	if token.RefreshToken != "" {
		s.refreshTokens[token.RefreshToken] = token.AccessToken
	}

	return nil
}

// Get 获取令牌
func (s *MemoryTokenStore) Get(ctx context.Context, accessToken string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, ok := s.tokens[accessToken]
	if !ok {
		return nil, ErrTokenNotFound
	}

	return token, nil
}

// Delete 删除令牌
func (s *MemoryTokenStore) Delete(ctx context.Context, accessToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	token, ok := s.tokens[accessToken]
	if ok {
		delete(s.tokens, accessToken)
		if token.RefreshToken != "" {
			delete(s.refreshTokens, token.RefreshToken)
		}
	}

	return nil
}

// DeleteByRefreshToken 通过刷新令牌删除
func (s *MemoryTokenStore) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	accessToken, ok := s.refreshTokens[refreshToken]
	if ok {
		delete(s.tokens, accessToken)
		delete(s.refreshTokens, refreshToken)
	}

	return nil
}

// cleanup 清理过期令牌
func (s *MemoryTokenStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for accessToken, token := range s.tokens {
			if now.After(token.ExpiresAt) {
				delete(s.tokens, accessToken)
				if token.RefreshToken != "" {
					delete(s.refreshTokens, token.RefreshToken)
				}
			}
		}
		s.mu.Unlock()
	}
}

// MemoryClientStore 内存客户端存储
type MemoryClientStore struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

// NewMemoryClientStore 创建内存客户端存储
func NewMemoryClientStore() *MemoryClientStore {
	return &MemoryClientStore{
		clients: make(map[string]*Client),
	}
}

// Get 获取客户端
func (s *MemoryClientStore) Get(ctx context.Context, clientID string) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, ok := s.clients[clientID]
	if !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}

// ValidateSecret 验证客户端密钥
func (s *MemoryClientStore) ValidateSecret(ctx context.Context, clientID, secret string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, ok := s.clients[clientID]
	if !ok {
		return ErrClientNotFound
	}

	if client.Secret != secret {
		return ErrInvalidClientSecret
	}

	return nil
}

// Save 保存客户端
func (s *MemoryClientStore) Save(ctx context.Context, client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.ID] = client
	return nil
}

// MemoryCodeStore 内存授权码存储
type MemoryCodeStore struct {
	mu   sync.RWMutex
	codes map[string]*AuthCode
}

// NewMemoryCodeStore 创建内存授权码存储
func NewMemoryCodeStore() *MemoryCodeStore {
	store := &MemoryCodeStore{
		codes: make(map[string]*AuthCode),
	}

	// 启动清理协程
	go store.cleanup()

	return store
}

// Save 保存授权码
func (s *MemoryCodeStore) Save(ctx context.Context, code *AuthCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.codes[code.Code] = code
	return nil
}

// Get 获取授权码
func (s *MemoryCodeStore) Get(ctx context.Context, code string) (*AuthCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	authCode, ok := s.codes[code]
	if !ok {
		return nil, ErrCodeNotFound
	}

	return authCode, nil
}

// Delete 删除授权码
func (s *MemoryCodeStore) Delete(ctx context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.codes, code)
	return nil
}

// cleanup 清理过期授权码
func (s *MemoryCodeStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for code, authCode := range s.codes {
			if now.After(authCode.ExpiresAt) {
				delete(s.codes, code)
			}
		}
		s.mu.Unlock()
	}
}

// 错误定义
var (
	ErrTokenNotFound      = &Error{Code: "token_not_found", Message: "Token not found"}
	ErrClientNotFound     = &Error{Code: "client_not_found", Message: "Client not found"}
	ErrInvalidClientSecret = &Error{Code: "invalid_client_secret", Message: "Invalid client secret"}
	ErrCodeNotFound       = &Error{Code: "code_not_found", Message: "Authorization code not found"}
)

// Error OAuth2 错误
type Error struct {
	Code        string
	Message     string
	Description string
}

func (e *Error) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

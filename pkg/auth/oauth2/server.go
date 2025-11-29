package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Server OAuth2 服务器
// 提供 OAuth2 授权服务器功能
type Server struct {
	clients      map[string]*Client
	tokens       map[string]*Token
	authCodes    map[string]*AuthCode
	config       *ServerConfig
	tokenStore   TokenStore
	clientStore  ClientStore
	codeStore    CodeStore
}

// ServerConfig OAuth2 服务器配置
type ServerConfig struct {
	// AccessTokenLifetime 访问令牌生命周期（秒）
	AccessTokenLifetime time.Duration
	// RefreshTokenLifetime 刷新令牌生命周期（秒）
	RefreshTokenLifetime time.Duration
	// AuthCodeLifetime 授权码生命周期（秒）
	AuthCodeLifetime time.Duration
	// TokenType 令牌类型（通常是 "Bearer"）
	TokenType string
	// AllowedGrantTypes 允许的授权类型
	AllowedGrantTypes []GrantType
	// AllowedScopes 允许的作用域
	AllowedScopes []string
}

// DefaultServerConfig 返回默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		AccessTokenLifetime:  3600 * time.Second,  // 1 小时
		RefreshTokenLifetime: 86400 * time.Second, // 24 小时
		AuthCodeLifetime:     600 * time.Second,   // 10 分钟
		TokenType:            "Bearer",
		AllowedGrantTypes: []GrantType{
			GrantTypeAuthorizationCode,
			GrantTypeClientCredentials,
			GrantTypeRefreshToken,
		},
		AllowedScopes: []string{"read", "write", "admin"},
	}
}

// GrantType 授权类型
type GrantType string

const (
	// GrantTypeAuthorizationCode 授权码模式
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	// GrantTypeClientCredentials 客户端凭证模式
	GrantTypeClientCredentials GrantType = "client_credentials"
	// GrantTypeRefreshToken 刷新令牌模式
	GrantTypeRefreshToken GrantType = "refresh_token"
	// GrantTypeImplicit 隐式模式（不推荐，已废弃）
	GrantTypeImplicit GrantType = "implicit"
)

// Client OAuth2 客户端
type Client struct {
	ID           string
	Secret       string
	RedirectURIs []string
	GrantTypes   []GrantType
	Scopes       []string
	CreatedAt    time.Time
}

// Token 访问令牌
type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
	Scope        string
	ClientID     string
	UserID       string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

// AuthCode 授权码
type AuthCode struct {
	Code        string
	ClientID    string
	RedirectURI string
	Scope       string
	UserID      string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// TokenStore 令牌存储接口
type TokenStore interface {
	Save(ctx context.Context, token *Token) error
	Get(ctx context.Context, accessToken string) (*Token, error)
	Delete(ctx context.Context, accessToken string) error
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
}

// ClientStore 客户端存储接口
type ClientStore interface {
	Get(ctx context.Context, clientID string) (*Client, error)
	ValidateSecret(ctx context.Context, clientID, secret string) error
}

// CodeStore 授权码存储接口
type CodeStore interface {
	Save(ctx context.Context, code *AuthCode) error
	Get(ctx context.Context, code string) (*AuthCode, error)
	Delete(ctx context.Context, code string) error
}

// NewServer 创建新的 OAuth2 服务器
func NewServer(config *ServerConfig) *Server {
	if config == nil {
		config = DefaultServerConfig()
	}

	return &Server{
		clients:    make(map[string]*Client),
		tokens:     make(map[string]*Token),
		authCodes:  make(map[string]*AuthCode),
		config:     config,
		tokenStore: NewMemoryTokenStore(),
		clientStore: NewMemoryClientStore(),
		codeStore:   NewMemoryCodeStore(),
	}
}

// SetTokenStore 设置令牌存储
func (s *Server) SetTokenStore(store TokenStore) {
	s.tokenStore = store
}

// SetClientStore 设置客户端存储
func (s *Server) SetClientStore(store ClientStore) {
	s.clientStore = store
}

// SetCodeStore 设置授权码存储
func (s *Server) SetCodeStore(store CodeStore) {
	s.codeStore = store
}

// RegisterClient 注册客户端
func (s *Server) RegisterClient(client *Client) error {
	if client.ID == "" {
		return errors.New("client ID is required")
	}
	if client.Secret == "" {
		return errors.New("client secret is required")
	}

	s.clients[client.ID] = client
	return nil
}

// GenerateAuthCode 生成授权码
func (s *Server) GenerateAuthCode(ctx context.Context, clientID, redirectURI, scope, userID string) (string, error) {
	// 验证客户端
	client, err := s.clientStore.Get(ctx, clientID)
	if err != nil {
		return "", fmt.Errorf("invalid client: %w", err)
	}

	// 验证重定向 URI
	if !s.validateRedirectURI(client, redirectURI) {
		return "", errors.New("invalid redirect URI")
	}

	// 生成授权码
	code := s.generateRandomString(32)
	authCode := &AuthCode{
		Code:        code,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		UserID:      userID,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.AuthCodeLifetime),
	}

	// 保存授权码
	if err := s.codeStore.Save(ctx, authCode); err != nil {
		return "", fmt.Errorf("failed to save auth code: %w", err)
	}

	return code, nil
}

// ExchangeAuthCode 交换授权码获取令牌
func (s *Server) ExchangeAuthCode(ctx context.Context, code, clientID, clientSecret, redirectURI string) (*Token, error) {
	// 验证客户端
	if err := s.clientStore.ValidateSecret(ctx, clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("invalid client credentials: %w", err)
	}

	// 获取授权码
	authCode, err := s.codeStore.Get(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid authorization code: %w", err)
	}

	// 验证授权码是否过期
	if time.Now().After(authCode.ExpiresAt) {
		s.codeStore.Delete(ctx, code)
		return nil, errors.New("authorization code expired")
	}

	// 验证客户端 ID
	if authCode.ClientID != clientID {
		return nil, errors.New("client ID mismatch")
	}

	// 验证重定向 URI
	if authCode.RedirectURI != redirectURI {
		return nil, errors.New("redirect URI mismatch")
	}

	// 删除已使用的授权码
	s.codeStore.Delete(ctx, code)

	// 生成令牌
	return s.generateToken(ctx, clientID, authCode.UserID, authCode.Scope)
}

// GenerateClientCredentialsToken 生成客户端凭证令牌
func (s *Server) GenerateClientCredentialsToken(ctx context.Context, clientID, clientSecret, scope string) (*Token, error) {
	// 验证客户端
	if err := s.clientStore.ValidateSecret(ctx, clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("invalid client credentials: %w", err)
	}

	// 生成令牌（客户端凭证模式没有用户 ID）
	return s.generateToken(ctx, clientID, "", scope)
}

// RefreshToken 刷新访问令牌
func (s *Server) RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret string) (*Token, error) {
	// 验证客户端
	if err := s.clientStore.ValidateSecret(ctx, clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("invalid client credentials: %w", err)
	}

	// 查找令牌
	token, err := s.tokenStore.Get(ctx, refreshToken)
	if err != nil {
		// 尝试通过刷新令牌查找
		tokens := s.tokens
		for _, t := range tokens {
			if t.RefreshToken == refreshToken {
				token = t
				break
			}
		}
		if token == nil {
			return nil, errors.New("invalid refresh token")
		}
	}

	// 验证客户端 ID
	if token.ClientID != clientID {
		return nil, errors.New("client ID mismatch")
	}

	// 验证令牌是否过期
	if time.Now().After(token.ExpiresAt) {
		s.tokenStore.DeleteByRefreshToken(ctx, refreshToken)
		return nil, errors.New("refresh token expired")
	}

	// 删除旧令牌
	s.tokenStore.DeleteByRefreshToken(ctx, refreshToken)

	// 生成新令牌
	return s.generateToken(ctx, clientID, token.UserID, token.Scope)
}

// ValidateToken 验证访问令牌
func (s *Server) ValidateToken(ctx context.Context, accessToken string) (*Token, error) {
	token, err := s.tokenStore.Get(ctx, accessToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// 验证令牌是否过期
	if time.Now().After(token.ExpiresAt) {
		s.tokenStore.Delete(ctx, accessToken)
		return nil, errors.New("token expired")
	}

	return token, nil
}

// RevokeToken 撤销令牌
func (s *Server) RevokeToken(ctx context.Context, token string) error {
	return s.tokenStore.Delete(ctx, token)
}

// generateToken 生成令牌
func (s *Server) generateToken(ctx context.Context, clientID, userID, scope string) (*Token, error) {
	accessToken := s.generateRandomString(32)
	refreshToken := s.generateRandomString(32)

	now := time.Now()
	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    s.config.TokenType,
		ExpiresIn:    int64(s.config.AccessTokenLifetime.Seconds()),
		Scope:         scope,
		ClientID:      clientID,
		UserID:        userID,
		CreatedAt:     now,
		ExpiresAt:     now.Add(s.config.AccessTokenLifetime),
	}

	// 保存令牌
	if err := s.tokenStore.Save(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	return token, nil
}

// generateRandomString 生成随机字符串
func (s *Server) generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用 UUID 作为后备
		return uuid.New().String()
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// validateRedirectURI 验证重定向 URI
func (s *Server) validateRedirectURI(client *Client, redirectURI string) bool {
	for _, uri := range client.RedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

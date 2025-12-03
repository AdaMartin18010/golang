package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// Provider OAuth2 提供者接口
// 支持多种 OAuth2 流程
type Provider interface {
	// AuthorizationURL 生成授权 URL
	AuthorizationURL(state string, opts ...oauth2.AuthCodeOption) string

	// Exchange 交换授权码获取令牌
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)

	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error)

	// ValidateToken 验证令牌
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)

	// RevokeToken 撤销令牌
	RevokeToken(ctx context.Context, token string) error
}

// TokenInfo 令牌信息
type TokenInfo struct {
	UserID    string
	Username  string
	Email     string
	Scope     []string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

// Config OAuth2 配置
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
	// PKCE 支持（推荐用于公共客户端）
	UsePKCE bool
}

// StandardProvider 标准 OAuth2 提供者实现
type StandardProvider struct {
	config *oauth2.Config
	usePKCE bool
}

// NewStandardProvider 创建标准 OAuth2 提供者
func NewStandardProvider(cfg Config) *StandardProvider {
	return &StandardProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     cfg.Endpoint,
		},
		usePKCE: cfg.UsePKCE,
	}
}

// AuthorizationURL 生成授权 URL
func (p *StandardProvider) AuthorizationURL(state string, opts ...oauth2.AuthCodeOption) string {
	if p.usePKCE {
		// 生成 PKCE challenge
		verifier := generateCodeVerifier()
		challenge := generateCodeChallenge(verifier)

		// 注意：实际实现需要存储 verifier 以便后续 Exchange 使用
		// 可以使用 session 或 cache 存储

		opts = append(opts,
			oauth2.SetAuthURLParam("code_challenge", challenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)
	}

	return p.config.AuthCodeURL(state, opts...)
}

// Exchange 交换授权码获取令牌
func (p *StandardProvider) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("authorization code is required")
	}

	if p.usePKCE {
		// 注意：实际实现需要从存储中获取 verifier
		// verifier := getStoredVerifier(state)
		// opts = append(opts, oauth2.SetAuthURLParam("code_verifier", verifier))
	}

	token, err := p.config.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// RefreshToken 刷新令牌
func (p *StandardProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := p.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// ValidateToken 验证令牌
func (p *StandardProvider) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	// 注意：实际实现需要调用 OAuth2 提供者的 introspection endpoint
	// 或验证 JWT 签名（如果是 JWT 令牌）

	// 示例实现：
	// 1. 解析 JWT
	// 2. 验证签名
	// 3. 检查过期时间
	// 4. 提取用户信息

	return nil, errors.New("not implemented")
}

// RevokeToken 撤销令牌
func (p *StandardProvider) RevokeToken(ctx context.Context, token string) error {
	// 注意：实际实现需要调用 OAuth2 提供者的 revocation endpoint
	return errors.New("not implemented")
}

// PKCE 辅助函数

// generateCodeVerifier 生成 PKCE code verifier
func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// generateCodeChallenge 生成 PKCE code challenge
func generateCodeChallenge(verifier string) string {
	// 使用 S256 方法（SHA256）
	// 实际实现需要：
	// hash := sha256.Sum256([]byte(verifier))
	// return base64.RawURLEncoding.EncodeToString(hash[:])
	return verifier // 占位实现
}

// 预定义的 OAuth2 Endpoints

// GoogleEndpoint Google OAuth2 endpoint
var GoogleEndpoint = oauth2.Endpoint{
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://oauth2.googleapis.com/token",
}

// GitHubEndpoint GitHub OAuth2 endpoint
var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

// MicrosoftEndpoint Microsoft OAuth2 endpoint
var MicrosoftEndpoint = oauth2.Endpoint{
	AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
	TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
}

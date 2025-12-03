package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCProvider OpenID Connect 提供者
// 扩展 OAuth2，添加 ID Token 支持
type OIDCProvider struct {
	oauth2Provider *StandardProvider
	oidcProvider   *oidc.Provider
	oidcVerifier   *oidc.IDTokenVerifier
	config         OIDCConfig
}

// OIDCConfig OIDC 配置
type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	IssuerURL    string // OIDC Issuer URL (e.g., https://accounts.google.com)
	// SkipIssuerCheck 跳过 Issuer 验证（仅用于测试）
	SkipIssuerCheck bool
	// SkipExpiryCheck 跳过过期检查（仅用于测试）
	SkipExpiryCheck bool
}

// NewOIDCProvider 创建 OIDC 提供者
func NewOIDCProvider(ctx context.Context, cfg OIDCConfig) (*OIDCProvider, error) {
	// 发现 OIDC 配置
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// 确保包含 openid scope
	scopes := cfg.Scopes
	hasOpenID := false
	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			hasOpenID = true
			break
		}
	}
	if !hasOpenID {
		scopes = append([]string{oidc.ScopeOpenID}, scopes...)
	}

	// 创建 OAuth2 配置
	oauth2Config := Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       scopes,
		Endpoint:     provider.Endpoint(),
		UsePKCE:      true, // OIDC 推荐使用 PKCE
	}

	// 创建 ID Token 验证器
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          cfg.ClientID,
		SkipIssuerCheck:   cfg.SkipIssuerCheck,
		SkipExpiryCheck:   cfg.SkipExpiryCheck,
	})

	return &OIDCProvider{
		oauth2Provider: NewStandardProvider(oauth2Config),
		oidcProvider:   provider,
		oidcVerifier:   verifier,
		config:         cfg,
	}, nil
}

// AuthorizationURL 生成授权 URL
func (p *OIDCProvider) AuthorizationURL(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2Provider.AuthorizationURL(state, opts...)
}

// Exchange 交换授权码获取令牌（包含 ID Token）
func (p *OIDCProvider) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*OIDCToken, error) {
	// 交换授权码
	token, err := p.oauth2Provider.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, err
	}

	// 提取 ID Token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token not found in token response")
	}

	// 验证 ID Token
	idToken, err := p.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// 解析 Claims
	var claims IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse ID token claims: %w", err)
	}

	return &OIDCToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
		IDToken:      rawIDToken,
		Claims:       claims,
	}, nil
}

// RefreshToken 刷新令牌
func (p *OIDCProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return p.oauth2Provider.RefreshToken(ctx, refreshToken)
}

// ValidateToken 验证访问令牌
func (p *OIDCProvider) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	return p.oauth2Provider.ValidateToken(ctx, token)
}

// VerifyIDToken 验证 ID Token
func (p *OIDCProvider) VerifyIDToken(ctx context.Context, rawIDToken string) (*IDTokenClaims, error) {
	idToken, err := p.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}

// GetUserInfo 获取用户信息
func (p *OIDCProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfo, err := p.oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	var info UserInfo
	if err := userInfo.Claims(&info); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &info, nil
}

// RevokeToken 撤销令牌
func (p *OIDCProvider) RevokeToken(ctx context.Context, token string) error {
	return p.oauth2Provider.RevokeToken(ctx, token)
}

// OIDCToken OIDC 令牌（包含 ID Token）
type OIDCToken struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Expiry       time.Time
	IDToken      string
	Claims       IDTokenClaims
}

// IDTokenClaims ID Token 标准声明
type IDTokenClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	Expiry    int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	AuthTime  int64  `json:"auth_time,omitempty"`
	Nonce     string `json:"nonce,omitempty"`

	// 标准声明
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	PhoneVerified     bool   `json:"phone_number_verified,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

// UserInfo 用户信息（从 UserInfo endpoint 获取）
type UserInfo struct {
	Subject           string `json:"sub"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	PhoneVerified     bool   `json:"phone_number_verified,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

// 预定义的 OIDC 提供者

// NewGoogleOIDCProvider 创建 Google OIDC 提供者
func NewGoogleOIDCProvider(ctx context.Context, clientID, clientSecret, redirectURL string) (*OIDCProvider, error) {
	return NewOIDCProvider(ctx, OIDCConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		IssuerURL:    "https://accounts.google.com",
	})
}

// NewMicrosoftOIDCProvider 创建 Microsoft OIDC 提供者
func NewMicrosoftOIDCProvider(ctx context.Context, clientID, clientSecret, redirectURL, tenant string) (*OIDCProvider, error) {
	if tenant == "" {
		tenant = "common"
	}

	return NewOIDCProvider(ctx, OIDCConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		IssuerURL:    fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenant),
	})
}

// NewAuth0OIDCProvider 创建 Auth0 OIDC 提供者
func NewAuth0OIDCProvider(ctx context.Context, clientID, clientSecret, redirectURL, domain string) (*OIDCProvider, error) {
	return NewOIDCProvider(ctx, OIDCConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		IssuerURL:    fmt.Sprintf("https://%s/", domain),
	})
}

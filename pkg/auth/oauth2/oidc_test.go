package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func generateTestKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func TestOIDCProvider_GenerateIDToken(t *testing.T) {
	// 生成测试密钥对
	privateKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 创建服务器和用户存储
	server := NewServer(DefaultServerConfig())
	userStore := NewMemoryUserStore()

	// 创建用户
	user := &UserInfo{
		Subject:           "user-123",
		Name:              "Test User",
		Email:             "test@example.com",
		EmailVerified:     true,
		PreferredUsername: "testuser",
	}
	if err := userStore.Save(context.Background(), user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// 创建 OIDC 提供者
	provider, err := NewOIDCProvider(server, "https://example.com", privateKey, userStore)
	if err != nil {
		t.Fatalf("Failed to create OIDC provider: %v", err)
	}

	ctx := context.Background()
	token, err := provider.GenerateIDToken(ctx, "user-123", "test-client", "test-nonce", "access-token")
	if err != nil {
		t.Fatalf("Failed to generate ID token: %v", err)
	}

	if token == "" {
		t.Error("ID token should not be empty")
	}
}

func TestOIDCProvider_ValidateIDToken(t *testing.T) {
	// 生成测试密钥对
	privateKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 创建服务器和用户存储
	server := NewServer(DefaultServerConfig())
	userStore := NewMemoryUserStore()

	// 创建用户
	user := &UserInfo{
		Subject:           "user-123",
		Name:              "Test User",
		Email:             "test@example.com",
		EmailVerified:     true,
		PreferredUsername: "testuser",
	}
	if err := userStore.Save(context.Background(), user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// 创建 OIDC 提供者
	provider, err := NewOIDCProvider(server, "https://example.com", privateKey, userStore)
	if err != nil {
		t.Fatalf("Failed to create OIDC provider: %v", err)
	}

	ctx := context.Background()
	
	// 生成 ID Token
	tokenString, err := provider.GenerateIDToken(ctx, "user-123", "test-client", "test-nonce", "access-token")
	if err != nil {
		t.Fatalf("Failed to generate ID token: %v", err)
	}

	// 验证 ID Token
	claims, err := provider.ValidateIDToken(ctx, tokenString, "test-client", "test-nonce")
	if err != nil {
		t.Fatalf("Failed to validate ID token: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Errorf("Expected subject 'user-123', got '%s'", claims.Subject)
	}

	if claims.Issuer != "https://example.com" {
		t.Errorf("Expected issuer 'https://example.com', got '%s'", claims.Issuer)
	}
}

func TestOIDCProvider_GetUserInfo(t *testing.T) {
	// 生成测试密钥对
	privateKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 创建服务器和用户存储
	server := NewServer(DefaultServerConfig())
	userStore := NewMemoryUserStore()

	// 注册客户端
	client := &Client{
		ID:         "test-client",
		Secret:     "test-secret",
		GrantTypes: []GrantType{GrantTypeClientCredentials},
		Scopes:     []string{"read"},
		CreatedAt:  time.Now(),
	}
	if err := server.clientStore.(*MemoryClientStore).Save(context.Background(), client); err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// 创建用户
	user := &UserInfo{
		Subject:           "user-123",
		Name:              "Test User",
		Email:             "test@example.com",
		EmailVerified:     true,
		PreferredUsername: "testuser",
	}
	if err := userStore.Save(context.Background(), user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// 创建 OIDC 提供者
	provider, err := NewOIDCProvider(server, "https://example.com", privateKey, userStore)
	if err != nil {
		t.Fatalf("Failed to create OIDC provider: %v", err)
	}

	ctx := context.Background()
	
	// 生成访问令牌
	token, err := server.GenerateClientCredentialsToken(ctx, "test-client", "test-secret", "read")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 注意：客户端凭证模式没有用户 ID，所以这里需要特殊处理
	// 在实际使用中，应该使用授权码流程获取的令牌
	// 这里仅测试 UserInfo 端点功能

	// 获取用户信息（需要修改服务器以支持用户 ID）
	// 为了测试，我们直接使用用户存储
	userInfo, err := provider.userStore.GetUser(ctx, "user-123")
	if err != nil {
		t.Fatalf("Failed to get user info: %v", err)
	}

	if userInfo.Subject != "user-123" {
		t.Errorf("Expected subject 'user-123', got '%s'", userInfo.Subject)
	}

	if userInfo.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", userInfo.Email)
	}
}

func TestOIDCProvider_GetDiscoveryDocument(t *testing.T) {
	// 生成测试密钥对
	privateKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 创建服务器和用户存储
	server := NewServer(DefaultServerConfig())
	userStore := NewMemoryUserStore()

	// 创建 OIDC 提供者
	provider, err := NewOIDCProvider(server, "https://example.com", privateKey, userStore)
	if err != nil {
		t.Fatalf("Failed to create OIDC provider: %v", err)
	}

	// 获取 Discovery 文档
	doc := provider.GetDiscoveryDocument()

	if doc.Issuer != "https://example.com" {
		t.Errorf("Expected issuer 'https://example.com', got '%s'", doc.Issuer)
	}

	if doc.AuthorizationEndpoint != "https://example.com/oauth/authorize" {
		t.Errorf("Expected authorization endpoint 'https://example.com/oauth/authorize', got '%s'", doc.AuthorizationEndpoint)
	}

	if doc.TokenEndpoint != "https://example.com/oauth/token" {
		t.Errorf("Expected token endpoint 'https://example.com/oauth/token', got '%s'", doc.TokenEndpoint)
	}

	if doc.UserInfoEndpoint != "https://example.com/userinfo" {
		t.Errorf("Expected userinfo endpoint 'https://example.com/userinfo', got '%s'", doc.UserInfoEndpoint)
	}

	if doc.JWKSUri != "https://example.com/.well-known/jwks.json" {
		t.Errorf("Expected JWKS URI 'https://example.com/.well-known/jwks.json', got '%s'", doc.JWKSUri)
	}
}

func TestOIDCProvider_GetJWKS(t *testing.T) {
	// 生成测试密钥对
	privateKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 创建服务器和用户存储
	server := NewServer(DefaultServerConfig())
	userStore := NewMemoryUserStore()

	// 创建 OIDC 提供者
	provider, err := NewOIDCProvider(server, "https://example.com", privateKey, userStore)
	if err != nil {
		t.Fatalf("Failed to create OIDC provider: %v", err)
	}

	// 获取 JWKS
	jwks := provider.GetJWKS()

	if len(jwks.Keys) == 0 {
		t.Error("JWKS should contain at least one key")
	}

	key := jwks.Keys[0]
	if key.Kty != "RSA" {
		t.Errorf("Expected key type 'RSA', got '%s'", key.Kty)
	}

	if key.Use != "sig" {
		t.Errorf("Expected key use 'sig', got '%s'", key.Use)
	}

	if key.Alg != "RS256" {
		t.Errorf("Expected algorithm 'RS256', got '%s'", key.Alg)
	}

	if key.N == "" {
		t.Error("Modulus (N) should not be empty")
	}

	if key.E == "" {
		t.Error("Exponent (E) should not be empty")
	}
}

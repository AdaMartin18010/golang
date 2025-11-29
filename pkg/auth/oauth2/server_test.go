package oauth2

import (
	"context"
	"testing"
	"time"
)

func TestServer_GenerateAuthCode(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	// 注册客户端
	client := &Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []GrantType{GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now(),
	}
	
	if err := server.clientStore.(*MemoryClientStore).Save(context.Background(), client); err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	ctx := context.Background()
	code, err := server.GenerateAuthCode(ctx, "test-client", "http://localhost:8080/callback", "read write", "user-123")
	if err != nil {
		t.Fatalf("Failed to generate auth code: %v", err)
	}

	if code == "" {
		t.Error("Auth code should not be empty")
	}

	// 验证授权码可以获取
	authCode, err := server.codeStore.Get(ctx, code)
	if err != nil {
		t.Fatalf("Failed to get auth code: %v", err)
	}

	if authCode.ClientID != "test-client" {
		t.Errorf("Expected client ID 'test-client', got '%s'", authCode.ClientID)
	}

	if authCode.UserID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", authCode.UserID)
	}
}

func TestServer_ExchangeAuthCode(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	// 注册客户端
	client := &Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []GrantType{GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now(),
	}
	
	if err := server.clientStore.(*MemoryClientStore).Save(context.Background(), client); err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	ctx := context.Background()
	
	// 生成授权码
	code, err := server.GenerateAuthCode(ctx, "test-client", "http://localhost:8080/callback", "read write", "user-123")
	if err != nil {
		t.Fatalf("Failed to generate auth code: %v", err)
	}

	// 交换授权码获取令牌
	token, err := server.ExchangeAuthCode(ctx, code, "test-client", "test-secret", "http://localhost:8080/callback")
	if err != nil {
		t.Fatalf("Failed to exchange auth code: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	if token.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	if token.ClientID != "test-client" {
		t.Errorf("Expected client ID 'test-client', got '%s'", token.ClientID)
	}

	if token.UserID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", token.UserID)
	}

	// 验证授权码已被删除
	_, err = server.codeStore.Get(ctx, code)
	if err == nil {
		t.Error("Auth code should be deleted after exchange")
	}
}

func TestServer_GenerateClientCredentialsToken(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	// 注册客户端
	client := &Client{
		ID:         "test-client",
		Secret:     "test-secret",
		GrantTypes: []GrantType{GrantTypeClientCredentials},
		Scopes:     []string{"read", "write"},
		CreatedAt:  time.Now(),
	}
	
	if err := server.clientStore.(*MemoryClientStore).Save(context.Background(), client); err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	ctx := context.Background()
	token, err := server.GenerateClientCredentialsToken(ctx, "test-client", "test-secret", "read write")
	if err != nil {
		t.Fatalf("Failed to generate client credentials token: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	if token.ClientID != "test-client" {
		t.Errorf("Expected client ID 'test-client', got '%s'", token.ClientID)
	}

	// 客户端凭证模式没有用户 ID
	if token.UserID != "" {
		t.Errorf("Expected empty user ID, got '%s'", token.UserID)
	}
}

func TestServer_RefreshToken(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
	// 注册客户端
	client := &Client{
		ID:         "test-client",
		Secret:     "test-secret",
		GrantTypes: []GrantType{GrantTypeRefreshToken},
		Scopes:     []string{"read", "write"},
		CreatedAt:  time.Now(),
	}
	
	if err := server.clientStore.(*MemoryClientStore).Save(context.Background(), client); err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	ctx := context.Background()
	
	// 生成初始令牌
	originalToken, err := server.GenerateClientCredentialsToken(ctx, "test-client", "test-secret", "read write")
	if err != nil {
		t.Fatalf("Failed to generate initial token: %v", err)
	}

	// 刷新令牌
	newToken, err := server.RefreshToken(ctx, originalToken.RefreshToken, "test-client", "test-secret")
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if newToken.AccessToken == originalToken.AccessToken {
		t.Error("New access token should be different from original")
	}

	if newToken.RefreshToken == originalToken.RefreshToken {
		t.Error("New refresh token should be different from original")
	}

	// 验证旧令牌已被删除
	_, err = server.tokenStore.Get(ctx, originalToken.AccessToken)
	if err == nil {
		t.Error("Original token should be deleted after refresh")
	}
}

func TestServer_ValidateToken(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
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

	ctx := context.Background()
	token, err := server.GenerateClientCredentialsToken(ctx, "test-client", "test-secret", "read")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 验证令牌
	validatedToken, err := server.ValidateToken(ctx, token.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedToken.AccessToken != token.AccessToken {
		t.Error("Validated token should match original token")
	}
}

func TestServer_RevokeToken(t *testing.T) {
	server := NewServer(DefaultServerConfig())
	
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

	ctx := context.Background()
	token, err := server.GenerateClientCredentialsToken(ctx, "test-client", "test-secret", "read")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 撤销令牌
	if err := server.RevokeToken(ctx, token.AccessToken); err != nil {
		t.Fatalf("Failed to revoke token: %v", err)
	}

	// 验证令牌已被删除
	_, err = server.tokenStore.Get(ctx, token.AccessToken)
	if err == nil {
		t.Error("Token should be deleted after revocation")
	}
}

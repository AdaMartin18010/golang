package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/golang/pkg/auth/oauth2"
)

// TestOAuth2AuthorizationCodeFlow 测试 OAuth2 授权码流程
func TestOAuth2AuthorizationCodeFlow(t *testing.T) {
	// 创建测试框架
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建 OAuth2 服务器
	server := oauth2.NewServer(oauth2.DefaultServerConfig())

	// 创建测试客户端
	client := &oauth2.Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []oauth2.GrantType{oauth2.GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
	}

	ctx := context.Background()
	err = server.RegisterClient(client)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// 1. 生成授权码
	code, err := server.GenerateAuthCode(ctx, client.ID, "http://localhost:8080/callback", "read write", "test-user")
	if err != nil {
		t.Fatalf("Failed to generate authorization code: %v", err)
	}

	if code == "" {
		t.Error("Authorization code should not be empty")
	}

	// 2. 交换访问令牌
	token, err := server.ExchangeAuthCode(ctx, code, client.ID, client.Secret, "http://localhost:8080/callback")
	if err != nil {
		t.Fatalf("Failed to exchange authorization code: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	// 3. 验证令牌
	validToken, err := server.ValidateToken(ctx, token.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validToken.UserID != "test-user" {
		t.Errorf("Expected user ID 'test-user', got '%s'", validToken.UserID)
	}
}

// TestOAuth2ClientCredentialsFlow 测试 OAuth2 客户端凭证流程
func TestOAuth2ClientCredentialsFlow(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	ctx := context.Background()

	// 创建 OAuth2 服务器
	server := oauth2.NewServer(oauth2.DefaultServerConfig())

	// 创建测试客户端
	client := &oauth2.Client{
		ID:         "test-client",
		Secret:     "test-secret",
		GrantTypes: []oauth2.GrantType{oauth2.GrantTypeClientCredentials},
		Scopes:     []string{"read", "write"},
	}

	err = server.RegisterClient(client)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// 生成客户端凭证令牌
	token, err := server.GenerateClientCredentialsToken(ctx, client.ID, client.Secret, "read write")
	if err != nil {
		t.Fatalf("Failed to generate client credentials token: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("Access token should not be empty")
	}

	// 验证令牌
	validToken, err := server.ValidateToken(ctx, token.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validToken.ClientID != client.ID {
		t.Errorf("Expected client ID '%s', got '%s'", client.ID, validToken.ClientID)
	}
}

// TestOAuth2RefreshTokenFlow 测试 OAuth2 刷新令牌流程
func TestOAuth2RefreshTokenFlow(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	ctx := context.Background()

	// 创建 OAuth2 服务器
	server := oauth2.NewServer(oauth2.DefaultServerConfig())

	// 创建测试客户端
	client := &oauth2.Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []oauth2.GrantType{oauth2.GrantTypeAuthorizationCode, oauth2.GrantTypeRefreshToken},
		Scopes:       []string{"read", "write"},
	}

	err = server.RegisterClient(client)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// 1. 生成初始令牌（带刷新令牌）
	code, _ := server.GenerateAuthCode(ctx, client.ID, "http://localhost:8080/callback", "read write", "test-user")
	token, _ := server.ExchangeAuthCode(ctx, code, client.ID, client.Secret, "http://localhost:8080/callback")

	if token.RefreshToken == "" {
		t.Fatal("Refresh token should be generated")
	}

	// 2. 使用刷新令牌获取新令牌
	newToken, err := server.RefreshToken(ctx, token.RefreshToken, client.ID, client.Secret)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if newToken.AccessToken == "" {
		t.Error("New access token should not be empty")
	}

	if newToken.AccessToken == token.AccessToken {
		t.Error("New access token should be different from old one")
	}
}

// TestOAuth2HTTPHandler 测试 OAuth2 HTTP 处理器
func TestOAuth2HTTPHandler(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建 OAuth2 服务器
	server := oauth2.NewServer(oauth2.DefaultServerConfig())

	// 创建测试客户端
	client := &oauth2.Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []oauth2.GrantType{oauth2.GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
	}

	err = server.RegisterClient(client)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟授权端点
		if r.URL.Path == "/oauth2/authorize" {
			// 处理授权请求
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Authorization page"))
		}
	}))
	defer ts.Close()

	// 测试授权端点
	req, _ := http.NewRequest("GET", ts.URL+"/oauth2/authorize?client_id=test-client&redirect_uri=http://localhost:8080/callback&response_type=code&scope=read+write", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

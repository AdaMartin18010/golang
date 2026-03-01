package oauth2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestStandardProvider_Exchange_InvalidCode(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://invalid.example.com/token",
		},
	}

	provider := NewStandardProvider(cfg)

	// 使用无效授权码应该返回错误
	ctx := context.Background()
	_, err := provider.Exchange(ctx, "invalid-code")
	require.Error(t, err)
}

func TestStandardProvider_RefreshToken_InvalidToken(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://invalid.example.com/token",
		},
	}

	provider := NewStandardProvider(cfg)

	ctx := context.Background()
	// 使用无效刷新令牌应该返回错误
	_, err := provider.RefreshToken(ctx, "invalid-refresh-token")
	require.Error(t, err)
}

func TestStandardProvider_ValidateToken_NotImplemented(t *testing.T) {
	cfg := Config{
		ClientID: "test-client",
		Endpoint: GoogleEndpoint,
	}
	provider := NewStandardProvider(cfg)

	ctx := context.Background()
	_, err := provider.ValidateToken(ctx, "some-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestStandardProvider_RevokeToken_NotImplemented(t *testing.T) {
	cfg := Config{
		ClientID: "test-client",
		Endpoint: GoogleEndpoint,
	}
	provider := NewStandardProvider(cfg)

	ctx := context.Background()
	err := provider.RevokeToken(ctx, "some-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "complete config",
			config: Config{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				RedirectURL:  "http://localhost/callback",
				Scopes:       []string{"openid", "email"},
				Endpoint:     GoogleEndpoint,
				UsePKCE:      true,
			},
			valid: true,
		},
		{
			name: "minimal config",
			config: Config{
				ClientID: "client-id",
				Endpoint: GoogleEndpoint,
			},
			valid: true,
		},
		{
			name: "empty config",
			config: Config{},
			valid: true, // 可以创建，但可能无法正常工作
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewStandardProvider(tt.config)
			if tt.valid {
				assert.NotNil(t, provider)
			}
		})
	}
}

func TestTokenInfo_Empty(t *testing.T) {
	info := TokenInfo{}
	assert.Empty(t, info.UserID)
	assert.Empty(t, info.Username)
	assert.Empty(t, info.Email)
	assert.Empty(t, info.Scope)
	assert.True(t, info.ExpiresAt.IsZero())
	assert.True(t, info.IssuedAt.IsZero())
}

func TestAuthorizationURL_Variations(t *testing.T) {
	tests := []struct {
		name  string
		state string
		opts  []oauth2.AuthCodeOption
	}{
		{"simple state", "state123", nil},
		{"empty state", "", nil},
		{"with options", "state456", []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}},
		{"long state", "very-long-state-string-with-many-characters-123456789", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Scopes:       []string{"openid"},
				Endpoint:     GoogleEndpoint,
				UsePKCE:      false,
			}
			provider := NewStandardProvider(cfg)

			url := provider.AuthorizationURL(tt.state, tt.opts...)
			assert.NotEmpty(t, url)
			assert.Contains(t, url, "client_id=test-client")
			if tt.state != "" {
				assert.Contains(t, url, "state="+tt.state)
			}
		})
	}
}

func TestPKCE_Flow(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid"},
		Endpoint:     GoogleEndpoint,
		UsePKCE:      true,
	}
	provider := NewStandardProvider(cfg)

	// 生成授权 URL（应包含 PKCE 参数）
	url := provider.AuthorizationURL("state123")
	assert.Contains(t, url, "code_challenge=")
	assert.Contains(t, url, "code_challenge_method=S256")

	// 验证 code_verifier 格式
	verifier := generateCodeVerifier()
	assert.Len(t, verifier, 43) // base64 URL encoding of 32 bytes
	assert.Regexp(t, `^[A-Za-z0-9_-]+$`, verifier)
}

func TestProvider_Interface(t *testing.T) {
	// 验证 StandardProvider 实现了 Provider 接口
	var _ Provider = (*StandardProvider)(nil)
}

func TestEndpoints_Constants(t *testing.T) {
	tests := []struct {
		name     string
		endpoint oauth2.Endpoint
		authURL  string
		tokenURL string
	}{
		{"Google", GoogleEndpoint, "https://accounts.google.com/o/oauth2/auth", "https://oauth2.googleapis.com/token"},
		{"GitHub", GitHubEndpoint, "https://github.com/login/oauth/authorize", "https://github.com/login/oauth/access_token"},
		{"Microsoft", MicrosoftEndpoint, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize", "https://login.microsoftonline.com/common/oauth2/v2.0/token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.authURL, tt.endpoint.AuthURL)
			assert.Equal(t, tt.tokenURL, tt.endpoint.TokenURL)
		})
	}
}



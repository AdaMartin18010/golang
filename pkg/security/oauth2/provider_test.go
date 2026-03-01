package oauth2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// TestNewStandardProvider 测试创建标准 OAuth2 提供者
func TestNewStandardProvider(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     GoogleEndpoint,
		UsePKCE:      true,
	}

	provider := NewStandardProvider(cfg)

	require.NotNil(t, provider)
	assert.NotNil(t, provider.config)
	assert.True(t, provider.usePKCE)
}

// TestStandardProvider_AuthorizationURL 测试生成授权 URL
func TestStandardProvider_AuthorizationURL(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		state    string
		wantHost string
	}{
		{
			name: "Google OAuth2",
			config: Config{
				ClientID:     "google-client-id",
				ClientSecret: "google-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Scopes:       []string{"openid", "email"},
				Endpoint:     GoogleEndpoint,
				UsePKCE:      false,
			},
			state:    "random-state-123",
			wantHost: "accounts.google.com",
		},
		{
			name: "GitHub OAuth2",
			config: Config{
				ClientID:     "github-client-id",
				ClientSecret: "github-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Scopes:       []string{"user:email"},
				Endpoint:     GitHubEndpoint,
				UsePKCE:      false,
			},
			state:    "state-xyz",
			wantHost: "github.com",
		},
		{
			name: "Microsoft OAuth2",
			config: Config{
				ClientID:     "microsoft-client-id",
				ClientSecret: "microsoft-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Scopes:       []string{"openid", "profile"},
				Endpoint:     MicrosoftEndpoint,
				UsePKCE:      true,
			},
			state:    "state-abc",
			wantHost: "login.microsoftonline.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewStandardProvider(tt.config)
			url := provider.AuthorizationURL(tt.state)

			assert.Contains(t, url, tt.wantHost)
			assert.Contains(t, url, "client_id="+tt.config.ClientID)
			assert.Contains(t, url, "redirect_uri=")
			assert.Contains(t, url, "state="+tt.state)
			assert.Contains(t, url, "response_type=code")
		})
	}
}

// TestStandardProvider_AuthorizationURLWithPKCE 测试 PKCE
func TestStandardProvider_AuthorizationURLWithPKCE(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid"},
		Endpoint:     GoogleEndpoint,
		UsePKCE:      true,
	}

	provider := NewStandardProvider(cfg)
	url := provider.AuthorizationURL("state-123")

	// 验证 PKCE 参数
	assert.Contains(t, url, "code_challenge=")
	assert.Contains(t, url, "code_challenge_method=S256")
}

// TestStandardProvider_Exchange 测试交换授权码
func TestStandardProvider_Exchange(t *testing.T) {
	// 注意：这是单元测试，不涉及真实的 OAuth2 服务器
	// 实际测试需要 mock OAuth2 服务器

	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://example.com/oauth/authorize",
			TokenURL: "https://example.com/oauth/token",
		},
	}

	provider := NewStandardProvider(cfg)

	t.Run("空授权码", func(t *testing.T) {
		ctx := context.Background()
		_, err := provider.Exchange(ctx, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "authorization code is required")
	})

	t.Run("无效授权码", func(t *testing.T) {
		ctx := context.Background()
		// 这会失败，因为没有真实的 OAuth2 服务器
		_, err := provider.Exchange(ctx, "invalid-code")
		// 期望错误，因为没有真实的 token endpoint
		require.Error(t, err)
	})
}

// TestStandardProvider_RefreshToken 测试刷新令牌
func TestStandardProvider_RefreshToken(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://example.com/oauth/token",
		},
	}

	provider := NewStandardProvider(cfg)

	t.Run("空刷新令牌", func(t *testing.T) {
		ctx := context.Background()
		_, err := provider.RefreshToken(ctx, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "refresh token is required")
	})
}

// TestStandardProvider_ValidateToken 测试验证令牌
func TestStandardProvider_ValidateToken(t *testing.T) {
	cfg := Config{
		ClientID: "test-client",
		Endpoint: GoogleEndpoint,
	}
	provider := NewStandardProvider(cfg)

	t.Run("未实现", func(t *testing.T) {
		ctx := context.Background()
		_, err := provider.ValidateToken(ctx, "some-token")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})
}

// TestStandardProvider_RevokeToken 测试撤销令牌
func TestStandardProvider_RevokeToken(t *testing.T) {
	cfg := Config{
		ClientID: "test-client",
		Endpoint: GoogleEndpoint,
	}
	provider := NewStandardProvider(cfg)

	t.Run("未实现", func(t *testing.T) {
		ctx := context.Background()
		err := provider.RevokeToken(ctx, "some-token")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})
}

// TestGenerateCodeVerifier 测试生成 PKCE code verifier
func TestGenerateCodeVerifier(t *testing.T) {
	verifier1 := generateCodeVerifier()
	verifier2 := generateCodeVerifier()

	// 验证长度（base64 URL encoding of 32 bytes = 43 characters）
	assert.Len(t, verifier1, 43)
	assert.Len(t, verifier2, 43)

	// 验证每次生成不同
	assert.NotEqual(t, verifier1, verifier2)

	// 验证只包含 base64 URL safe 字符
	assert.Regexp(t, `^[A-Za-z0-9_-]+$`, verifier1)
}

// TestGenerateCodeChallenge 测试生成 PKCE code challenge
func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-string"
	challenge := generateCodeChallenge(verifier)

	// 当前是占位实现，返回 verifier 本身
	assert.Equal(t, verifier, challenge)
}

// TestTokenInfo 测试 TokenInfo 结构
func TestTokenInfo(t *testing.T) {
	info := TokenInfo{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Scope:    []string{"read", "write"},
	}

	assert.Equal(t, "user-123", info.UserID)
	assert.Equal(t, "testuser", info.Username)
	assert.Equal(t, "test@example.com", info.Email)
	assert.Equal(t, []string{"read", "write"}, info.Scope)
}

// TestConfig 测试 Config 结构
func TestConfig(t *testing.T) {
	cfg := Config{
		ClientID:     "client-123",
		ClientSecret: "secret-456",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     GoogleEndpoint,
		UsePKCE:      true,
	}

	assert.Equal(t, "client-123", cfg.ClientID)
	assert.Equal(t, "secret-456", cfg.ClientSecret)
	assert.Equal(t, "http://localhost/callback", cfg.RedirectURL)
	assert.Equal(t, []string{"openid", "profile", "email"}, cfg.Scopes)
	assert.True(t, cfg.UsePKCE)
}

// TestEndpoints 测试预定义的 Endpoints
func TestEndpoints(t *testing.T) {
	t.Run("Google Endpoint", func(t *testing.T) {
		assert.Equal(t, "https://accounts.google.com/o/oauth2/auth", GoogleEndpoint.AuthURL)
		assert.Equal(t, "https://oauth2.googleapis.com/token", GoogleEndpoint.TokenURL)
	})

	t.Run("GitHub Endpoint", func(t *testing.T) {
		assert.Equal(t, "https://github.com/login/oauth/authorize", GitHubEndpoint.AuthURL)
		assert.Equal(t, "https://github.com/login/oauth/access_token", GitHubEndpoint.TokenURL)
	})

	t.Run("Microsoft Endpoint", func(t *testing.T) {
		assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize", MicrosoftEndpoint.AuthURL)
		assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/token", MicrosoftEndpoint.TokenURL)
	})
}

// BenchmarkGenerateCodeVerifier 性能测试
func BenchmarkGenerateCodeVerifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateCodeVerifier()
	}
}

// BenchmarkAuthorizationURL 性能测试
func BenchmarkAuthorizationURL(b *testing.B) {
	cfg := Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid"},
		Endpoint:     GoogleEndpoint,
		UsePKCE:      false,
	}
	provider := NewStandardProvider(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.AuthorizationURL("state-123")
	}
}

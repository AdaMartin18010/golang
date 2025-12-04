package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenManager_InvalidSignature 测试无效签名
func TestTokenManager_InvalidSignature(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	// 生成令牌
	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
	require.NoError(t, err)

	// 篡改令牌（修改最后几个字符）
	tamperedToken := token[:len(token)-10] + "TAMPERED12"

	// 验证应该失败
	_, err = tm.ValidateToken(tamperedToken)
	assert.Error(t, err)
}

// TestTokenManager_MalformedToken 测试格式错误的令牌
func TestTokenManager_MalformedToken(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "not.a.jwt"},
		{"only two parts", "header.payload"},
		{"random string", "random-string-not-jwt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tm.ValidateToken(tt.token)
			assert.Error(t, err)
		})
	}
}

// TestTokenManager_CustomClaims 测试自定义声明
func TestTokenManager_CustomClaims(t *testing.T) {
	tm, err := NewTokenManager(Config{
		Issuer: "custom-issuer",
	})
	require.NoError(t, err)

	token, err := tm.GenerateAccessToken(
		"user-456",
		"jane.doe",
		"jane@example.com",
		[]string{"admin", "moderator"},
	)
	require.NoError(t, err)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)

	assert.Equal(t, "user-456", claims.UserID)
	assert.Equal(t, "jane.doe", claims.Username)
	assert.Equal(t, "jane@example.com", claims.Email)
	assert.Contains(t, claims.Roles, "admin")
	assert.Contains(t, claims.Roles, "moderator")
	assert.Equal(t, "custom-issuer", claims.Issuer)
}

// TestTokenManager_MultipleRoles 测试多角色
func TestTokenManager_MultipleRoles(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	roles := []string{"user", "admin", "moderator", "editor"}
	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", roles)
	require.NoError(t, err)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)

	assert.Len(t, claims.Roles, 4)
	for _, role := range roles {
		assert.Contains(t, claims.Roles, role)
	}
}

// TestTokenManager_EmptyRoles 测试空角色
func TestTokenManager_EmptyRoles(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{})
	require.NoError(t, err)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)

	assert.Empty(t, claims.Roles)
}

// TestTokenManager_LongLivedToken 测试长期令牌
func TestTokenManager_LongLivedToken(t *testing.T) {
	tm, err := NewTokenManager(Config{
		AccessTokenTTL: 24 * time.Hour,
	})
	require.NoError(t, err)

	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
	require.NoError(t, err)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)

	// 验证过期时间
	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	expectedExpiry := time.Now().Add(24 * time.Hour)

	// 允许1分钟的误差
	assert.WithinDuration(t, expectedExpiry, expiresAt, 1*time.Minute)
}

// TestTokenManager_ConcurrentGeneration 测试并发生成令牌
func TestTokenManager_ConcurrentGeneration(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	const numGoroutines = 100
	tokens := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			token, err := tm.GenerateAccessToken(
				"user-"+string(rune(id)),
				"user",
				"user@example.com",
				[]string{"user"},
			)
			if err != nil {
				errors <- err
				return
			}
			tokens <- token
		}(i)
	}

	// 收集结果
	for i := 0; i < numGoroutines; i++ {
		select {
		case token := <-tokens:
			assert.NotEmpty(t, token)
		case err := <-errors:
			t.Errorf("Error generating token: %v", err)
		}
	}
}

// BenchmarkTokenManager_GenerateAndValidate 性能测试 - 生成和验证
func BenchmarkTokenManager_GenerateAndValidate(b *testing.B) {
	tm, _ := NewTokenManager(Config{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _ := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
		tm.ValidateToken(token)
	}
}

// BenchmarkTokenManager_GenerateTokenPair 性能测试 - 生成令牌对
func BenchmarkTokenManager_GenerateTokenPair(b *testing.B) {
	tm, _ := NewTokenManager(Config{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.GenerateTokenPair("user-123", "john", "john@example.com", []string{"user"})
	}
}

package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenManager(t *testing.T) {
	tm, err := NewTokenManager(Config{
		Issuer:         "test-issuer",
		AccessTokenTTL: 15 * time.Minute,
	})
	require.NoError(t, err)
	assert.NotNil(t, tm)
	assert.NotNil(t, tm.privateKey)
	assert.NotNil(t, tm.publicKey)
}

func TestTokenManager_GenerateAccessToken(t *testing.T) {
	tm, err := NewTokenManager(Config{
		Issuer:         "test-issuer",
		AccessTokenTTL: 15 * time.Minute,
	})
	require.NoError(t, err)

	token, err := tm.GenerateAccessToken(
		"user-123",
		"john.doe",
		"john@example.com",
		[]string{"user", "admin"},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTokenManager_ValidateToken(t *testing.T) {
	tm, err := NewTokenManager(Config{
		Issuer:         "test-issuer",
		AccessTokenTTL: 15 * time.Minute,
	})
	require.NoError(t, err)

	// 生成令牌
	token, err := tm.GenerateAccessToken(
		"user-123",
		"john.doe",
		"john@example.com",
		[]string{"user"},
	)
	require.NoError(t, err)

	// 验证令牌
	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "john.doe", claims.Username)
	assert.Equal(t, "john@example.com", claims.Email)
	assert.Contains(t, claims.Roles, "user")
}

func TestTokenManager_ValidateToken_Invalid(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	// 验证无效令牌
	_, err = tm.ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestTokenManager_GenerateTokenPair(t *testing.T) {
	tm, err := NewTokenManager(Config{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	require.NoError(t, err)

	tokenPair, err := tm.GenerateTokenPair(
		"user-123",
		"john.doe",
		"john@example.com",
		[]string{"user"},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Greater(t, tokenPair.ExpiresIn, int64(0))
}

func TestTokenManager_RefreshAccessToken(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	// 生成初始令牌对
	tokenPair, err := tm.GenerateTokenPair(
		"user-123",
		"john.doe",
		"john@example.com",
		[]string{"user"},
	)
	require.NoError(t, err)

	// 使用刷新令牌获取新的令牌对
	newTokenPair, err := tm.RefreshAccessToken(tokenPair.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.NotEqual(t, tokenPair.AccessToken, newTokenPair.AccessToken)
}

func TestTokenManager_ExpiredToken(t *testing.T) {
	tm, err := NewTokenManager(Config{
		AccessTokenTTL: 1 * time.Millisecond, // 极短的过期时间
	})
	require.NoError(t, err)

	// 生成令牌
	token, err := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
	require.NoError(t, err)

	// 等待令牌过期
	time.Sleep(10 * time.Millisecond)

	// 验证过期令牌
	_, err = tm.ValidateToken(token)
	assert.Error(t, err)
}

func TestTokenManager_ExportKeys(t *testing.T) {
	tm, err := NewTokenManager(Config{})
	require.NoError(t, err)

	privateKeyPEM, publicKeyPEM, err := tm.ExportKeys()
	require.NoError(t, err)
	assert.NotEmpty(t, privateKeyPEM)
	assert.NotEmpty(t, publicKeyPEM)
	assert.Contains(t, string(privateKeyPEM), "RSA PRIVATE KEY")
	assert.Contains(t, string(publicKeyPEM), "RSA PUBLIC KEY")
}

func BenchmarkTokenManager_GenerateToken(b *testing.B) {
	tm, _ := NewTokenManager(Config{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})
	}
}

func BenchmarkTokenManager_ValidateToken(b *testing.B) {
	tm, _ := NewTokenManager(Config{})
	token, _ := tm.GenerateAccessToken("user-123", "john", "john@example.com", []string{"user"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.ValidateToken(token)
	}
}

package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestJWT_GenerateAccessToken(t *testing.T) {
	config := Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token, got empty string")
	}
}

func TestJWT_ValidateToken(t *testing.T) {
	config := Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := j.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("Expected UserID 'user-123', got '%s'", claims.UserID)
	}

	if claims.Username != "john" {
		t.Errorf("Expected Username 'john', got '%s'", claims.Username)
	}
}

func TestJWT_ValidateToken_Expired(t *testing.T) {
	config := Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: -1 * time.Hour, // 已过期
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = j.ValidateToken(token)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}

func TestJWT_RefreshToken(t *testing.T) {
	config := Config{
		SecretKey:       "test-secret-key",
		SigningMethod:   "HS256",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	refreshToken, err := j.GenerateRefreshToken("user-123")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	accessToken, newRefreshToken, err := j.RefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if accessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if newRefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}
}

func TestJWT_RS256(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	config := Config{
		PrivateKey:     privateKey,
		PublicKey:      &privateKey.PublicKey,
		SigningMethod:  "RS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := j.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("Expected UserID 'user-123', got '%s'", claims.UserID)
	}
}

func TestJWT_ExtractUserID(t *testing.T) {
	config := Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	token, err := j.GenerateAccessToken("user-123", "john", []string{"user"}, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	userID, err := j.ExtractUserID(token)
	if err != nil {
		t.Fatalf("Failed to extract user ID: %v", err)
	}

	if userID != "user-123" {
		t.Errorf("Expected UserID 'user-123', got '%s'", userID)
	}
}

func TestJWT_ExtractRoles(t *testing.T) {
	config := Config{
		SecretKey:      "test-secret-key",
		SigningMethod:  "HS256",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	j, err := NewJWT(config)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	roles := []string{"user", "admin"}
	token, err := j.GenerateAccessToken("user-123", "john", roles, "john@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	extractedRoles, err := j.ExtractRoles(token)
	if err != nil {
		t.Fatalf("Failed to extract roles: %v", err)
	}

	if len(extractedRoles) != len(roles) {
		t.Errorf("Expected %d roles, got %d", len(roles), len(extractedRoles))
	}
}

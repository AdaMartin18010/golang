package security

import (
	"testing"
	"time"
)

func TestCSRFProtection_GenerateToken(t *testing.T) {
	csrf := NewCSRFProtection(DefaultCSRFConfig())
	defer csrf.Shutdown()

	sessionID := "session-123"
	token, err := csrf.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestCSRFProtection_ValidateToken(t *testing.T) {
	csrf := NewCSRFProtection(DefaultCSRFConfig())
	defer csrf.Shutdown()

	sessionID := "session-123"
	token, err := csrf.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 验证正确令牌
	if err := csrf.ValidateToken(sessionID, token); err != nil {
		t.Errorf("Token validation should succeed, got error: %v", err)
	}

	// 验证错误令牌
	if err := csrf.ValidateToken(sessionID, "wrong-token"); err == nil {
		t.Error("Token validation should fail for wrong token")
	}

	// 验证错误会话 ID
	if err := csrf.ValidateToken("wrong-session", token); err == nil {
		t.Error("Token validation should fail for wrong session")
	}
}

func TestCSRFProtection_ExpiredToken(t *testing.T) {
	config := DefaultCSRFConfig()
	config.Expiry = 100 * time.Millisecond

	csrf := NewCSRFProtection(config)
	defer csrf.Shutdown()

	sessionID := "session-123"
	token, err := csrf.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 验证过期令牌
	if err := csrf.ValidateToken(sessionID, token); err != ErrCSRFTokenExpired {
		t.Errorf("Expected ErrCSRFTokenExpired, got %v", err)
	}
}

func TestCSRFProtection_RevokeToken(t *testing.T) {
	csrf := NewCSRFProtection(DefaultCSRFConfig())
	defer csrf.Shutdown()

	sessionID := "session-123"
	token, err := csrf.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 撤销令牌
	if err := csrf.RevokeToken(sessionID); err != nil {
		t.Fatalf("Failed to revoke token: %v", err)
	}

	// 验证令牌已撤销
	if err := csrf.ValidateToken(sessionID, token); err == nil {
		t.Error("Token validation should fail after revocation")
	}
}

package security

import (
	"testing"
	"time"
)

func TestSessionManager_CreateSession(t *testing.T) {
	sm := NewSessionManager(DefaultSessionConfig())
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", map[string]interface{}{
		"name": "Test User",
	})
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}

	if session.UserID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", session.UserID)
	}
}

func TestSessionManager_GetSession(t *testing.T) {
	sm := NewSessionManager(DefaultSessionConfig())
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	retrieved, err := sm.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrieved.ID)
	}
}

func TestSessionManager_ExpiredSession(t *testing.T) {
	config := DefaultSessionConfig()
	config.DefaultTTL = 100 * time.Millisecond

	sm := NewSessionManager(config)
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	_, err = sm.GetSession(session.ID)
	if err != ErrSessionExpired {
		t.Errorf("Expected ErrSessionExpired, got %v", err)
	}
}

func TestSessionManager_UpdateSession(t *testing.T) {
	sm := NewSessionManager(DefaultSessionConfig())
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", map[string]interface{}{
		"key1": "value1",
	})
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	err = sm.UpdateSession(session.ID, map[string]interface{}{
		"key2": "value2",
	})
	if err != nil {
		t.Fatalf("Failed to update session: %v", err)
	}

	updated, err := sm.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	if updated.Data["key1"] != "value1" {
		t.Error("Original data should be preserved")
	}

	if updated.Data["key2"] != "value2" {
		t.Error("New data should be added")
	}
}

func TestSessionManager_DeleteSession(t *testing.T) {
	sm := NewSessionManager(DefaultSessionConfig())
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	err = sm.DeleteSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	_, err = sm.GetSession(session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound, got %v", err)
	}
}

func TestSessionManager_RefreshSession(t *testing.T) {
	config := DefaultSessionConfig()
	config.DefaultTTL = 1 * time.Hour

	sm := NewSessionManager(config)
	defer sm.Shutdown()

	session, err := sm.CreateSession("user-123", nil)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	originalExpiry := session.ExpiresAt

	// 刷新会话
	err = sm.RefreshSession(session.ID, 2*time.Hour)
	if err != nil {
		t.Fatalf("Failed to refresh session: %v", err)
	}

	refreshed, err := sm.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get refreshed session: %v", err)
	}

	if !refreshed.ExpiresAt.After(originalExpiry) {
		t.Error("Session expiry should be extended")
	}
}

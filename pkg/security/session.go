package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrSessionNotFound 会话未找到
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired 会话已过期
	ErrSessionExpired = errors.New("session expired")
	// ErrInvalidSessionID 无效的会话 ID
	ErrInvalidSessionID = errors.New("invalid session ID")
)

// SessionManager 会话管理器
type SessionManager struct {
	sessions    map[string]*Session
	mu          sync.RWMutex
	secret      []byte
	idLength    int
	defaultTTL  time.Duration
	cleanup     *time.Ticker
	stopCleanup chan struct{}
}

// Session 会话
type Session struct {
	ID        string
	UserID    string
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	LastAccess time.Time
}

// SessionConfig 会话配置
type SessionConfig struct {
	Secret     []byte        // 密钥
	IDLength   int           // 会话 ID 长度
	DefaultTTL time.Duration // 默认过期时间
}

// DefaultSessionConfig 默认会话配置
func DefaultSessionConfig() SessionConfig {
	secret := make([]byte, 32)
	rand.Read(secret)

	return SessionConfig{
		Secret:     secret,
		IDLength:   32,
		DefaultTTL: 24 * time.Hour,
	}
}

// NewSessionManager 创建会话管理器
func NewSessionManager(config SessionConfig) *SessionManager {
	if config.Secret == nil {
		config = DefaultSessionConfig()
	}
	if config.IDLength == 0 {
		config.IDLength = 32
	}
	if config.DefaultTTL == 0 {
		config.DefaultTTL = 24 * time.Hour
	}

	sm := &SessionManager{
		sessions:    make(map[string]*Session),
		secret:      config.Secret,
		idLength:    config.IDLength,
		defaultTTL:  config.DefaultTTL,
		stopCleanup: make(chan struct{}),
	}

	// 启动清理协程
	sm.cleanup = time.NewTicker(1 * time.Hour)
	go sm.cleanupExpired()

	return sm
}

// CreateSession 创建会话
func (sm *SessionManager) CreateSession(userID string, data map[string]interface{}) (*Session, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	sessionID := sm.generateSessionID()
	now := time.Now()

	session := &Session{
		ID:         sessionID,
		UserID:     userID,
		Data:       data,
		CreatedAt:  now,
		ExpiresAt:  now.Add(sm.defaultTTL),
		LastAccess: now,
	}

	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}

	sm.mu.Lock()
	sm.sessions[sessionID] = session
	sm.mu.Unlock()

	return session, nil
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	if sessionID == "" {
		return nil, ErrInvalidSessionID
	}

	sm.mu.RLock()
	session, exists := sm.sessions[sessionID]
	sm.mu.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
	}

	// 检查是否过期
	if time.Now().After(session.ExpiresAt) {
		sm.mu.Lock()
		delete(sm.sessions, sessionID)
		sm.mu.Unlock()
		return nil, ErrSessionExpired
	}

	// 更新最后访问时间
	sm.mu.Lock()
	session.LastAccess = time.Now()
	sm.mu.Unlock()

	return session, nil
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(sessionID string, data map[string]interface{}) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	sm.mu.Lock()
	if data != nil {
		for k, v := range data {
			session.Data[k] = v
		}
	}
	session.LastAccess = time.Now()
	sm.mu.Unlock()

	return nil
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, sessionID)
	return nil
}

// RefreshSession 刷新会话（延长过期时间）
func (sm *SessionManager) RefreshSession(sessionID string, ttl time.Duration) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	if ttl == 0 {
		ttl = sm.defaultTTL
	}

	sm.mu.Lock()
	session.ExpiresAt = time.Now().Add(ttl)
	session.LastAccess = time.Now()
	sm.mu.Unlock()

	return nil
}

// ListSessions 列出所有会话
func (sm *SessionManager) ListSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}

// generateSessionID 生成会话 ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, sm.idLength)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// cleanupExpired 清理过期会话
func (sm *SessionManager) cleanupExpired() {
	for {
		select {
		case <-sm.cleanup.C:
			sm.mu.Lock()
			now := time.Now()
			for sessionID, session := range sm.sessions {
				if now.After(session.ExpiresAt) {
					delete(sm.sessions, sessionID)
				}
			}
			sm.mu.Unlock()

		case <-sm.stopCleanup:
			sm.cleanup.Stop()
			return
		}
	}
}

// Shutdown 关闭会话管理器
func (sm *SessionManager) Shutdown() error {
	close(sm.stopCleanup)
	return nil
}

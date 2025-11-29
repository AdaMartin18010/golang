package oauth2

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	// 注意：需要安装 PostgreSQL 驱动
	// go get github.com/lib/pq
	// _ "github.com/lib/pq"
)

// 注意：这个测试需要 PostgreSQL 数据库
// 可以通过环境变量 POSTGRES_DSN 配置数据库连接
// 例如：POSTGRES_DSN=postgres://user:pass@localhost/testdb?sslmode=disable

func getTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_DSN not set, skipping PostgreSQL tests")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

func TestPostgresTokenStore(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store, err := NewPostgresTokenStore(db)
	if err != nil {
		t.Fatalf("Failed to create PostgresTokenStore: %v", err)
	}

	ctx := context.Background()

	// 测试保存令牌
	token := &Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:         "read write",
		ClientID:      "test-client",
		UserID:        "test-user",
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(1 * time.Hour),
	}

	if err := store.Save(ctx, token); err != nil {
		t.Fatalf("Failed to save token: %v", err)
	}

	// 测试获取令牌
	retrievedToken, err := store.Get(ctx, "test-access-token")
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	if retrievedToken.AccessToken != token.AccessToken {
		t.Errorf("Expected access token %s, got %s", token.AccessToken, retrievedToken.AccessToken)
	}

	// 测试删除令牌
	if err := store.Delete(ctx, "test-access-token"); err != nil {
		t.Fatalf("Failed to delete token: %v", err)
	}

	// 验证令牌已删除
	_, err = store.Get(ctx, "test-access-token")
	if err != ErrTokenNotFound {
		t.Error("Token should be deleted")
	}
}

func TestPostgresClientStore(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store, err := NewPostgresClientStore(db)
	if err != nil {
		t.Fatalf("Failed to create PostgresClientStore: %v", err)
	}

	ctx := context.Background()

	// 创建客户端
	client := &Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []GrantType{GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now(),
	}

	if err := store.Save(ctx, client); err != nil {
		t.Fatalf("Failed to save client: %v", err)
	}

	// 测试获取客户端
	retrievedClient, err := store.Get(ctx, "test-client")
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	if retrievedClient.ID != client.ID {
		t.Errorf("Expected client ID %s, got %s", client.ID, retrievedClient.ID)
	}

	// 测试验证密钥
	if err := store.ValidateSecret(ctx, "test-client", "test-secret"); err != nil {
		t.Errorf("Secret validation should succeed, got error: %v", err)
	}

	if err := store.ValidateSecret(ctx, "test-client", "wrong-secret"); err == nil {
		t.Error("Secret validation should fail for wrong secret")
	}
}

func TestPostgresCodeStore(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store, err := NewPostgresCodeStore(db)
	if err != nil {
		t.Fatalf("Failed to create PostgresCodeStore: %v", err)
	}

	ctx := context.Background()

	// 测试保存授权码
	code := &AuthCode{
		Code:        "test-code",
		ClientID:    "test-client",
		RedirectURI: "http://localhost:8080/callback",
		Scope:       "read write",
		UserID:      "test-user",
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	if err := store.Save(ctx, code); err != nil {
		t.Fatalf("Failed to save code: %v", err)
	}

	// 测试获取授权码
	retrievedCode, err := store.Get(ctx, "test-code")
	if err != nil {
		t.Fatalf("Failed to get code: %v", err)
	}

	if retrievedCode.Code != code.Code {
		t.Errorf("Expected code %s, got %s", code.Code, retrievedCode.Code)
	}

	// 测试删除授权码
	if err := store.Delete(ctx, "test-code"); err != nil {
		t.Fatalf("Failed to delete code: %v", err)
	}

	// 验证授权码已删除
	_, err = store.Get(ctx, "test-code")
	if err != ErrCodeNotFound {
		t.Error("Code should be deleted")
	}
}

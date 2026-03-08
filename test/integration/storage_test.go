package integration

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/golang/pkg/auth/oauth2"
)

// TestPostgresStorageIntegration 测试PostgreSQL存储集成
func TestPostgresStorageIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 检查PostgreSQL是否可用
	if tf.GetDB() == nil {
		t.Skip("PostgreSQL not available")
	}

	// 创建PostgreSQL存储
	tokenStore, err := oauth2.NewPostgresTokenStore(tf.GetDB())
	if err != nil {
		t.Fatalf("Failed to create token store: %v", err)
	}
	clientStore, err := oauth2.NewPostgresClientStore(tf.GetDB())
	if err != nil {
		t.Fatalf("Failed to create client store: %v", err)
	}

	ctx := context.Background()

	// 测试客户端存储
	client := &oauth2.Client{
		ID:           "test-client-123",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:3000/callback"},
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now(),
	}

	err = clientStore.Save(ctx, client)
	if err != nil {
		t.Fatalf("Failed to save client: %v", err)
	}

	// 读取客户端
	retrievedClient, err := clientStore.Get(ctx, client.ID)
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	if retrievedClient.ID != client.ID {
		t.Errorf("Expected client ID %s, got %s", client.ID, retrievedClient.ID)
	}

	// 测试令牌存储
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "test-refresh-token",
		Scope:        "read write",
		ClientID:     client.ID,
		UserID:       "user-123",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	err = tokenStore.Save(ctx, token)
	if err != nil {
		t.Fatalf("Failed to save token: %v", err)
	}

	// 读取令牌
	retrievedToken, err := tokenStore.Get(ctx, token.AccessToken)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	if retrievedToken.AccessToken != token.AccessToken {
		t.Errorf("Expected access token %s, got %s", token.AccessToken, retrievedToken.AccessToken)
	}

	// 清理
	_ = tokenStore.Delete(ctx, token.AccessToken)
	// 注意：PostgresClientStore 没有 Delete 方法，跳过客户端清理
}

// TestRedisStorageIntegration 测试Redis存储集成
func TestRedisStorageIntegration(t *testing.T) {
	// Redis 暂未配置，跳过此测试
	t.Skip("Redis not configured - skipping Redis integration test")
}

// TestStoragePerformance 测试存储性能
func TestStoragePerformance(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	if tf.GetDB() == nil {
		t.Skip("PostgreSQL not available")
	}

	tokenStore, err := oauth2.NewPostgresTokenStore(tf.GetDB())
	if err != nil {
		t.Fatalf("Failed to create token store: %v", err)
	}
	ctx := context.Background()

	// 性能测试：写入100个令牌
	start := time.Now()
	for i := 0; i < 100; i++ {
		token := &oauth2.Token{
			AccessToken:  "test-token-" + string(rune(i)),
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			ClientID:     "test-client",
			UserID:       "user-123",
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = tokenStore.Save(ctx, token)
	}
	elapsed := time.Since(start)

	t.Logf("Wrote 100 tokens in %v (avg: %v per token)", elapsed, elapsed/100)

	// 性能测试：读取100个令牌
	start = time.Now()
	for i := 0; i < 100; i++ {
		_, _ = tokenStore.Get(ctx, "test-token-"+string(rune(i)))
	}
	elapsed = time.Since(start)

	t.Logf("Read 100 tokens in %v (avg: %v per token)", elapsed, elapsed/100)
}

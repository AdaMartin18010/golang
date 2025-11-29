package oauth2

import (
	"context"
	"os"
	"testing"
	"time"

	// 注意：需要安装 Redis 客户端
	// go get github.com/redis/go-redis/v9
	// "github.com/redis/go-redis/v9"
)

// 注意：这个测试需要 Redis 服务器
// 可以通过环境变量 REDIS_ADDR 配置 Redis 地址
// 例如：REDIS_ADDR=localhost:6379

func getTestRedisClient(t *testing.T) *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available at %s, skipping Redis tests: %v", addr, err)
	}

	return client
}

func TestRedisTokenStore(t *testing.T) {
	client := getTestRedisClient(t)
	defer client.Close()

	store, err := NewRedisTokenStore(client, "test:")
	if err != nil {
		t.Fatalf("Failed to create RedisTokenStore: %v", err)
	}

	ctx := context.Background()

	// 清理测试数据
	defer func() {
		client.FlushDB(ctx)
	}()

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

func TestRedisClientStore(t *testing.T) {
	client := getTestRedisClient(t)
	defer client.Close()

	store, err := NewRedisClientStore(client, "test:")
	if err != nil {
		t.Fatalf("Failed to create RedisClientStore: %v", err)
	}

	ctx := context.Background()

	// 清理测试数据
	defer func() {
		client.FlushDB(ctx)
	}()

	// 创建客户端
	clientObj := &Client{
		ID:           "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []GrantType{GrantTypeAuthorizationCode},
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now(),
	}

	if err := store.Save(ctx, clientObj); err != nil {
		t.Fatalf("Failed to save client: %v", err)
	}

	// 测试获取客户端
	retrievedClient, err := store.Get(ctx, "test-client")
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	if retrievedClient.ID != clientObj.ID {
		t.Errorf("Expected client ID %s, got %s", clientObj.ID, retrievedClient.ID)
	}

	// 测试验证密钥
	if err := store.ValidateSecret(ctx, "test-client", "test-secret"); err != nil {
		t.Errorf("Secret validation should succeed, got error: %v", err)
	}

	if err := store.ValidateSecret(ctx, "test-client", "wrong-secret"); err == nil {
		t.Error("Secret validation should fail for wrong secret")
	}
}

func TestRedisCodeStore(t *testing.T) {
	client := getTestRedisClient(t)
	defer client.Close()

	store, err := NewRedisCodeStore(client, "test:")
	if err != nil {
		t.Fatalf("Failed to create RedisCodeStore: %v", err)
	}

	ctx := context.Background()

	// 清理测试数据
	defer func() {
		client.FlushDB(ctx)
	}()

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

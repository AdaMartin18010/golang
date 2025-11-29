package security

import (
	"context"
	"testing"
	"time"
)

func TestKeyManager_GenerateAESKey(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	// 测试生成 256 位密钥
	key, err := km.GenerateAESKey(ctx, "test-key", 256)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	if key.ID == "" {
		t.Error("Key ID should not be empty")
	}

	if key.Name != "test-key" {
		t.Errorf("Expected key name 'test-key', got '%s'", key.Name)
	}

	if key.Type != KeyTypeAES {
		t.Errorf("Expected key type AES, got '%s'", key.Type)
	}

	if len(key.Data) != 32 { // 256 bits = 32 bytes
		t.Errorf("Expected key length 32, got %d", len(key.Data))
	}
}

func TestKeyManager_GenerateRSAKeyPair(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	privateKey, publicKey, err := km.GenerateRSAKeyPair(ctx, "test-rsa", 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	if privateKey.ID == "" {
		t.Error("Private key ID should not be empty")
	}

	if publicKey.ID == "" {
		t.Error("Public key ID should not be empty")
	}

	if privateKey.Type != KeyTypeRSA {
		t.Errorf("Expected private key type RSA, got '%s'", privateKey.Type)
	}

	if publicKey.Type != KeyTypeRSA {
		t.Errorf("Expected public key type RSA, got '%s'", publicKey.Type)
	}
}

func TestKeyManager_SaveAndGetKey(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	key, err := km.GenerateAESKey(ctx, "test-key", 256)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// 获取密钥
	retrievedKey, err := km.GetKey(ctx, key.ID)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if retrievedKey.ID != key.ID {
		t.Errorf("Expected key ID %s, got %s", key.ID, retrievedKey.ID)
	}

	if len(retrievedKey.Data) != len(key.Data) {
		t.Errorf("Expected key data length %d, got %d", len(key.Data), len(retrievedKey.Data))
	}
}

func TestKeyManager_DeleteKey(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	key, err := km.GenerateAESKey(ctx, "test-key", 256)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// 删除密钥
	if err := km.DeleteKey(ctx, key.ID); err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// 验证密钥已删除
	_, err = km.GetKey(ctx, key.ID)
	if err != ErrKeyNotFound {
		t.Error("Key should be deleted")
	}
}

func TestKeyManager_ListKeys(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	// 生成多个密钥
	_, err := km.GenerateAESKey(ctx, "key1", 256)
	if err != nil {
		t.Fatalf("Failed to generate key1: %v", err)
	}

	_, err = km.GenerateAESKey(ctx, "key2", 256)
	if err != nil {
		t.Fatalf("Failed to generate key2: %v", err)
	}

	// 列出所有密钥
	keys, err := km.ListKeys(ctx)
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}

	if len(keys) < 2 {
		t.Errorf("Expected at least 2 keys, got %d", len(keys))
	}
}

func TestKeyManager_RotateKey(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	// 生成原始密钥
	originalKey, err := km.GenerateAESKey(ctx, "test-key", 256)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// 生成新密钥数据
	newKeyData := make([]byte, 32)
	for i := range newKeyData {
		newKeyData[i] = byte(i + 10)
	}

	// 轮换密钥
	newKey, err := km.RotateKey(ctx, originalKey.ID, newKeyData)
	if err != nil {
		t.Fatalf("Failed to rotate key: %v", err)
	}

	if newKey.Version != originalKey.Version+1 {
		t.Errorf("Expected version %d, got %d", originalKey.Version+1, newKey.Version)
	}

	if newKey.Metadata["previous_key_id"] != originalKey.ID {
		t.Error("Metadata should contain previous key ID")
	}
}

func TestKeyManager_ExpiredKey(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	key, err := km.GenerateAESKey(ctx, "test-key", 256)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// 设置过期时间（过去）
	expiredTime := time.Now().Add(-1 * time.Hour)
	key.ExpiresAt = &expiredTime
	km.SaveKey(ctx, key)

	// 尝试获取过期密钥
	_, err = km.GetKey(ctx, key.ID)
	if err != ErrKeyExpired {
		t.Errorf("Expected ErrKeyExpired, got %v", err)
	}
}

func TestKeyManager_InvalidKeySize(t *testing.T) {
	store := NewMemoryKeyStore()
	km := NewKeyManager(store)

	ctx := context.Background()

	// 测试无效的密钥大小
	_, err := km.GenerateAESKey(ctx, "test-key", 512)
	if err == nil {
		t.Error("Should return error for invalid key size")
	}
}

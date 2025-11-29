package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	// ErrKeyNotFound 密钥未找到
	ErrKeyNotFound = errors.New("key not found")
	// ErrKeyExpired 密钥已过期
	ErrKeyExpired = errors.New("key expired")
)

// KeyManager 密钥管理器
type KeyManager struct {
	keys      map[string]*Key
	mu        sync.RWMutex
	keyStore  KeyStore
}

// Key 密钥
type Key struct {
	ID        string
	Name      string
	Type      KeyType
	Data      []byte
	Version   int
	CreatedAt time.Time
	ExpiresAt *time.Time
	Metadata  map[string]string
}

// KeyType 密钥类型
type KeyType string

const (
	// KeyTypeAES AES 密钥
	KeyTypeAES KeyType = "aes"
	// KeyTypeRSA RSA 密钥对
	KeyTypeRSA KeyType = "rsa"
	// KeyTypeHMAC HMAC 密钥
	KeyTypeHMAC KeyType = "hmac"
)

// KeyStore 密钥存储接口
type KeyStore interface {
	Save(ctx context.Context, key *Key) error
	Get(ctx context.Context, keyID string) (*Key, error)
	Delete(ctx context.Context, keyID string) error
	List(ctx context.Context) ([]*Key, error)
}

// NewKeyManager 创建密钥管理器
func NewKeyManager(keyStore KeyStore) *KeyManager {
	return &KeyManager{
		keys:     make(map[string]*Key),
		keyStore: keyStore,
	}
}

// GenerateAESKey 生成 AES 密钥
func (km *KeyManager) GenerateAESKey(ctx context.Context, name string, size int) (*Key, error) {
	if size != 128 && size != 192 && size != 256 {
		return nil, fmt.Errorf("invalid key size: %d (must be 128, 192, or 256)", size)
	}

	keyBytes := make([]byte, size/8)
	if _, err := io.ReadFull(rand.Reader, keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	key := &Key{
		ID:        generateKeyID(),
		Name:      name,
		Type:      KeyTypeAES,
		Data:      keyBytes,
		Version:   1,
		CreatedAt: time.Now(),
		Metadata:  make(map[string]string),
	}

	if err := km.SaveKey(ctx, key); err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateRSAKeyPair 生成 RSA 密钥对
func (km *KeyManager) GenerateRSAKeyPair(ctx context.Context, name string, bits int) (*Key, *Key, error) {
	if bits < 2048 {
		return nil, nil, fmt.Errorf("RSA key size must be at least 2048 bits")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// 编码私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 编码公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	privateKeyObj := &Key{
		ID:        generateKeyID(),
		Name:      name + "-private",
		Type:      KeyTypeRSA,
		Data:      privateKeyPEM,
		Version:   1,
		CreatedAt: time.Now(),
		Metadata: map[string]string{
			"bits": fmt.Sprintf("%d", bits),
		},
	}

	publicKeyObj := &Key{
		ID:        generateKeyID(),
		Name:      name + "-public",
		Type:      KeyTypeRSA,
		Data:      publicKeyPEM,
		Version:   1,
		CreatedAt: time.Now(),
		Metadata: map[string]string{
			"bits": fmt.Sprintf("%d", bits),
		},
	}

	if err := km.SaveKey(ctx, privateKeyObj); err != nil {
		return nil, nil, err
	}

	if err := km.SaveKey(ctx, publicKeyObj); err != nil {
		return nil, nil, err
	}

	return privateKeyObj, publicKeyObj, nil
}

// SaveKey 保存密钥
func (km *KeyManager) SaveKey(ctx context.Context, key *Key) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	km.keys[key.ID] = key

	if km.keyStore != nil {
		return km.keyStore.Save(ctx, key)
	}

	return nil
}

// GetKey 获取密钥
func (km *KeyManager) GetKey(ctx context.Context, keyID string) (*Key, error) {
	km.mu.RLock()
	key, exists := km.keys[keyID]
	km.mu.RUnlock()

	if exists {
		// 检查是否过期
		if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
			return nil, ErrKeyExpired
		}
		return key, nil
	}

	// 从存储中获取
	if km.keyStore != nil {
		key, err := km.keyStore.Get(ctx, keyID)
		if err != nil {
			return nil, err
		}

		// 检查是否过期
		if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
			return nil, ErrKeyExpired
		}

		// 缓存到内存
		km.mu.Lock()
		km.keys[keyID] = key
		km.mu.Unlock()

		return key, nil
	}

	return nil, ErrKeyNotFound
}

// DeleteKey 删除密钥
func (km *KeyManager) DeleteKey(ctx context.Context, keyID string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	delete(km.keys, keyID)

	if km.keyStore != nil {
		return km.keyStore.Delete(ctx, keyID)
	}

	return nil
}

// ListKeys 列出所有密钥
func (km *KeyManager) ListKeys(ctx context.Context) ([]*Key, error) {
	if km.keyStore != nil {
		return km.keyStore.List(ctx)
	}

	km.mu.RLock()
	defer km.mu.RUnlock()

	keys := make([]*Key, 0, len(km.keys))
	for _, key := range km.keys {
		keys = append(keys, key)
	}

	return keys, nil
}

// RotateKey 轮换密钥（创建新版本）
func (km *KeyManager) RotateKey(ctx context.Context, keyID string, newKeyData []byte) (*Key, error) {
	oldKey, err := km.GetKey(ctx, keyID)
	if err != nil {
		return nil, err
	}

	newKey := &Key{
		ID:        generateKeyID(),
		Name:      oldKey.Name,
		Type:      oldKey.Type,
		Data:      newKeyData,
		Version:   oldKey.Version + 1,
		CreatedAt: time.Now(),
		Metadata:  make(map[string]string),
	}

	// 复制元数据
	for k, v := range oldKey.Metadata {
		newKey.Metadata[k] = v
	}
	newKey.Metadata["previous_key_id"] = oldKey.ID

	if err := km.SaveKey(ctx, newKey); err != nil {
		return nil, err
	}

	return newKey, nil
}

// generateKeyID 生成密钥 ID
func generateKeyID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("key_%x", b)
}

// MemoryKeyStore 内存密钥存储（用于测试）
type MemoryKeyStore struct {
	keys map[string]*Key
	mu   sync.RWMutex
}

// NewMemoryKeyStore 创建内存密钥存储
func NewMemoryKeyStore() *MemoryKeyStore {
	return &MemoryKeyStore{
		keys: make(map[string]*Key),
	}
}

// Save 保存密钥
func (s *MemoryKeyStore) Save(ctx context.Context, key *Key) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keys[key.ID] = key
	return nil
}

// Get 获取密钥
func (s *MemoryKeyStore) Get(ctx context.Context, keyID string) (*Key, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key, exists := s.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	return key, nil
}

// Delete 删除密钥
func (s *MemoryKeyStore) Delete(ctx context.Context, keyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.keys, keyID)
	return nil
}

// List 列出所有密钥
func (s *MemoryKeyStore) List(ctx context.Context) ([]*Key, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]*Key, 0, len(s.keys))
	for _, key := range s.keys {
		keys = append(keys, key)
	}

	return keys, nil
}

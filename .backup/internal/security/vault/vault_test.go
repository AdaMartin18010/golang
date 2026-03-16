// Package vault provides HashiCorp Vault integration for secure secret management.
//
// vault_test.go 包含 Vault 客户端的单元测试，使用 mock 实现。
//
// 测试覆盖：
// 1. 客户端创建和配置
// 2. 密钥管理操作（KV v1/v2）
// 3. 动态数据库凭据
// 4. 加密/解密服务
// 5. 密钥轮换
// 6. 错误处理和重试机制
package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockClient 是 Vault 客户端的 mock 实现
type MockClient struct {
	connected       bool
	kv              *MockKVManager
	db              *MockDBManager
	encryption      *MockEncryptionService
	rotation        *MockRotationService
	healthFunc      func(ctx context.Context) (*HealthStatus, error)
	closeFunc       func() error
	isConnectedFunc func() bool
}

// NewMockClient 创建新的 mock 客户端
func NewMockClient() *MockClient {
	return &MockClient{
		connected:  true,
		kv:         NewMockKVManager(),
		db:         NewMockDBManager(),
		encryption: NewMockEncryptionService(),
		rotation:   NewMockRotationService(),
	}
}

func (m *MockClient) KV() KVManager                 { return m.kv }
func (m *MockClient) DB() DBManager                 { return m.db }
func (m *MockClient) Encryption() EncryptionService { return m.encryption }
func (m *MockClient) Rotation() RotationService     { return m.rotation }
func (m *MockClient) Health(ctx context.Context) (*HealthStatus, error) {
	if m.healthFunc != nil {
		return m.healthFunc(ctx)
	}
	return &HealthStatus{
		Initialized: true,
		Sealed:      false,
		Standby:     false,
		Version:     "1.15.0",
		CheckedAt:   time.Now(),
	}, nil
}
func (m *MockClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	m.connected = false
	return nil
}
func (m *MockClient) IsConnected() bool {
	if m.isConnectedFunc != nil {
		return m.isConnectedFunc()
	}
	return m.connected
}

// MockKVManager 是 KVManager 的 mock 实现
type MockKVManager struct {
	secrets    map[string]*Secret
	versions   map[string][]SecretVersion
	metadata   map[string]*KVMetadata
	isV2Flag   bool
	getFunc    func(ctx context.Context, path string, version ...int) (*Secret, error)
	putFunc    func(ctx context.Context, path string, data map[string]interface{}) error
	deleteFunc func(ctx context.Context, path string) error
	listFunc   func(ctx context.Context, path string) ([]string, error)
}

// NewMockKVManager 创建新的 mock KV 管理器
func NewMockKVManager() *MockKVManager {
	return &MockKVManager{
		secrets:  make(map[string]*Secret),
		versions: make(map[string][]SecretVersion),
		metadata: make(map[string]*KVMetadata),
		isV2Flag: true,
	}
}

func (m *MockKVManager) Get(ctx context.Context, path string, version ...int) (*Secret, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, path, version...)
	}
	if secret, ok := m.secrets[path]; ok {
		return secret, nil
	}
	return nil, errors.New("secret not found")
}

func (m *MockKVManager) GetRaw(ctx context.Context, path string) (map[string]interface{}, error) {
	secret, err := m.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}

func (m *MockKVManager) Put(ctx context.Context, path string, data map[string]interface{}) error {
	if m.putFunc != nil {
		return m.putFunc(ctx, path, data)
	}
	version := 1
	if m.isV2Flag {
		if versions, ok := m.versions[path]; ok {
			version = len(versions) + 1
		}
	}
	m.secrets[path] = &Secret{
		Data:    data,
		Version: version,
		Metadata: &SecretMetadata{
			Version:     version,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
	}
	m.versions[path] = append(m.versions[path], SecretVersion{
		Version:     version,
		CreatedTime: time.Now(),
	})
	return nil
}

func (m *MockKVManager) Delete(ctx context.Context, path string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, path)
	}
	delete(m.secrets, path)
	return nil
}

func (m *MockKVManager) DeleteVersion(ctx context.Context, path string, version int) error {
	if !m.isV2Flag {
		return errors.New("not supported in v1")
	}
	return nil
}

func (m *MockKVManager) DestroyVersion(ctx context.Context, path string, version int) error {
	if !m.isV2Flag {
		return errors.New("not supported in v1")
	}
	if versions, ok := m.versions[path]; ok {
		for i := range versions {
			if versions[i].Version == version {
				versions[i].Destroyed = true
				break
			}
		}
	}
	return nil
}

func (m *MockKVManager) List(ctx context.Context, path string) ([]string, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, path)
	}
	keys := make([]string, 0)
	for key := range m.secrets {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *MockKVManager) GetVersions(ctx context.Context, path string) ([]SecretVersion, error) {
	if !m.isV2Flag {
		return nil, errors.New("not supported in v1")
	}
	if versions, ok := m.versions[path]; ok {
		return versions, nil
	}
	return []SecretVersion{}, nil
}

func (m *MockKVManager) Rollback(ctx context.Context, path string, version int) error {
	if !m.isV2Flag {
		return errors.New("not supported in v1")
	}
	return nil
}

func (m *MockKVManager) GetMetadata(ctx context.Context, path string) (*KVMetadata, error) {
	if !m.isV2Flag {
		return nil, errors.New("not supported in v1")
	}
	if metadata, ok := m.metadata[path]; ok {
		return metadata, nil
	}
	return &KVMetadata{}, nil
}

func (m *MockKVManager) UpdateMetadata(ctx context.Context, path string, metadata *KVMetadata) error {
	if !m.isV2Flag {
		return errors.New("not supported in v1")
	}
	m.metadata[path] = metadata
	return nil
}

func (m *MockKVManager) IsV2() bool { return m.isV2Flag }

// MockDBManager 是 DBManager 的 mock 实现
type MockDBManager struct {
	credentials  map[string]*DBCredentials
	getCredsFunc func(ctx context.Context, role string) (*DBCredentials, error)
	renewFunc    func(ctx context.Context, leaseID string, increment int) error
	revokeFunc   func(ctx context.Context, leaseID string) error
}

// NewMockDBManager 创建新的 mock 数据库凭据管理器
func NewMockDBManager() *MockDBManager {
	return &MockDBManager{credentials: make(map[string]*DBCredentials)}
}

func (m *MockDBManager) GetCredentials(ctx context.Context, role string) (*DBCredentials, error) {
	if m.getCredsFunc != nil {
		return m.getCredsFunc(ctx, role)
	}
	if creds, ok := m.credentials[role]; ok {
		return creds, nil
	}
	return &DBCredentials{
		Username:      "test_user",
		Password:      "test_password",
		LeaseID:       "lease-123",
		LeaseDuration: 3600,
		Renewable:     true,
	}, nil
}

func (m *MockDBManager) GetCredentialsFromPath(ctx context.Context, mount, role string) (*DBCredentials, error) {
	return m.GetCredentials(ctx, role)
}

func (m *MockDBManager) RenewLease(ctx context.Context, leaseID string, increment int) error {
	if m.renewFunc != nil {
		return m.renewFunc(ctx, leaseID, increment)
	}
	return nil
}

func (m *MockDBManager) RevokeLease(ctx context.Context, leaseID string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, leaseID)
	}
	return nil
}

// MockEncryptionService 是 EncryptionService 的 mock 实现
type MockEncryptionService struct {
	keys        map[string]*KeyInfo
	encryptFunc func(ctx context.Context, keyName string, plaintext []byte, opts *EncryptOptions) (*EncryptResult, error)
	decryptFunc func(ctx context.Context, keyName string, ciphertext string, opts *DecryptOptions) ([]byte, error)
	signFunc    func(ctx context.Context, keyName string, input []byte, opts *SignOptions) (*SignResult, error)
	verifyFunc  func(ctx context.Context, keyName string, input []byte, signature string, opts *VerifyOptions) (bool, error)
}

// NewMockEncryptionService 创建新的 mock 加密服务
func NewMockEncryptionService() *MockEncryptionService {
	return &MockEncryptionService{keys: make(map[string]*KeyInfo)}
}

func (m *MockEncryptionService) CreateKey(ctx context.Context, name string, keyType KeyType, opts *KeyOptions) error {
	m.keys[name] = &KeyInfo{
		Name:               name,
		Type:               string(keyType),
		LatestVersion:      1,
		SupportsEncryption: true,
		SupportsDecryption: true,
		SupportsSigning:    keyType == KeyTypeECDSAP256 || keyType == KeyTypeECDSAP384 || keyType == KeyTypeED25519,
	}
	return nil
}

func (m *MockEncryptionService) DeleteKey(ctx context.Context, name string) error {
	delete(m.keys, name)
	return nil
}

func (m *MockEncryptionService) ListKeys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0, len(m.keys))
	for name := range m.keys {
		keys = append(keys, name)
	}
	return keys, nil
}

func (m *MockEncryptionService) GetKeyInfo(ctx context.Context, name string) (*KeyInfo, error) {
	if key, ok := m.keys[name]; ok {
		return key, nil
	}
	return nil, errors.New("key not found")
}

func (m *MockEncryptionService) RotateKey(ctx context.Context, name string) error {
	if key, ok := m.keys[name]; ok {
		key.LatestVersion++
	}
	return nil
}

func (m *MockEncryptionService) ExportKey(ctx context.Context, name string, keyType string) (map[string]string, error) {
	return map[string]string{"1": "exported-key-1"}, nil
}

func (m *MockEncryptionService) Encrypt(ctx context.Context, keyName string, plaintext []byte, opts *EncryptOptions) (*EncryptResult, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(ctx, keyName, plaintext, opts)
	}
	return &EncryptResult{Ciphertext: "vault:v1:" + string(plaintext), KeyVersion: 1}, nil
}

func (m *MockEncryptionService) EncryptBatch(ctx context.Context, keyName string, items []*BatchEncryptItem) ([]*EncryptResult, error) {
	results := make([]*EncryptResult, len(items))
	for i, item := range items {
		results[i], _ = m.Encrypt(ctx, keyName, item.Plaintext, nil)
	}
	return results, nil
}

func (m *MockEncryptionService) Decrypt(ctx context.Context, keyName string, ciphertext string, opts *DecryptOptions) ([]byte, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(ctx, keyName, ciphertext, opts)
	}
	if len(ciphertext) > 9 {
		return []byte(ciphertext[9:]), nil
	}
	return []byte(ciphertext), nil
}

func (m *MockEncryptionService) DecryptBatch(ctx context.Context, keyName string, items []*BatchDecryptItem) ([][]byte, error) {
	results := make([][]byte, len(items))
	for i, item := range items {
		results[i], _ = m.Decrypt(ctx, keyName, item.Ciphertext, nil)
	}
	return results, nil
}

func (m *MockEncryptionService) Rewrap(ctx context.Context, keyName string, ciphertext string) (*EncryptResult, error) {
	plaintext, err := m.Decrypt(ctx, keyName, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return m.Encrypt(ctx, keyName, plaintext, nil)
}

func (m *MockEncryptionService) Sign(ctx context.Context, keyName string, input []byte, opts *SignOptions) (*SignResult, error) {
	if m.signFunc != nil {
		return m.signFunc(ctx, keyName, input, opts)
	}
	return &SignResult{Signature: "signature-" + string(input), KeyVersion: 1}, nil
}

func (m *MockEncryptionService) Verify(ctx context.Context, keyName string, input []byte, signature string, opts *VerifyOptions) (bool, error) {
	if m.verifyFunc != nil {
		return m.verifyFunc(ctx, keyName, input, signature, opts)
	}
	expectedSig, _ := m.Sign(ctx, keyName, input, nil)
	return expectedSig.Signature == signature, nil
}

func (m *MockEncryptionService) Hash(ctx context.Context, input []byte, algorithm string) (string, error) {
	return "hash-" + string(input), nil
}

func (m *MockEncryptionService) GenerateHMAC(ctx context.Context, keyName string, input []byte, algorithm string) (string, error) {
	return "hmac-" + string(input), nil
}

func (m *MockEncryptionService) VerifyHMAC(ctx context.Context, keyName string, input []byte, hmac string, algorithm string) (bool, error) {
	expectedHMAC, _ := m.GenerateHMAC(ctx, keyName, input, algorithm)
	return expectedHMAC == hmac, nil
}

// MockRotationService 是 RotationService 的 mock 实现
type MockRotationService struct {
	policies     map[string]*RotationPolicy
	history      []RotationHistory
	running      bool
	registerFunc func(policy *RotationPolicy) error
	rotateFunc   func(ctx context.Context, path string) error
}

// NewMockRotationService 创建新的 mock 轮换服务
func NewMockRotationService() *MockRotationService {
	return &MockRotationService{
		policies: make(map[string]*RotationPolicy),
		history:  make([]RotationHistory, 0),
	}
}

func (m *MockRotationService) RegisterPolicy(policy *RotationPolicy) error {
	if m.registerFunc != nil {
		return m.registerFunc(policy)
	}
	if policy.ID == "" {
		policy.ID = "policy-" + time.Now().String()
	}
	m.policies[policy.ID] = policy
	return nil
}

func (m *MockRotationService) UnregisterPolicy(policyID string) error {
	delete(m.policies, policyID)
	return nil
}

func (m *MockRotationService) GetPolicy(policyID string) (*RotationPolicy, error) {
	if policy, ok := m.policies[policyID]; ok {
		return policy, nil
	}
	return nil, errors.New("policy not found")
}

func (m *MockRotationService) ListPolicies() []*RotationPolicy {
	policies := make([]*RotationPolicy, 0, len(m.policies))
	for _, policy := range m.policies {
		policies = append(policies, policy)
	}
	return policies
}

func (m *MockRotationService) Rotate(ctx context.Context, path string) error {
	if m.rotateFunc != nil {
		return m.rotateFunc(ctx, path)
	}
	m.history = append(m.history, RotationHistory{
		ID:        "hist-" + time.Now().String(),
		Path:      path,
		RotatedAt: time.Now(),
		RotatedBy: "manual",
		Status:    "success",
	})
	return nil
}

func (m *MockRotationService) RotateWithPolicy(ctx context.Context, policyID string) error {
	policy, err := m.GetPolicy(policyID)
	if err != nil {
		return err
	}
	return m.Rotate(ctx, policy.Path)
}

func (m *MockRotationService) StartAutoRotation() error {
	m.running = true
	return nil
}

func (m *MockRotationService) StopAutoRotation() error {
	m.running = false
	return nil
}

func (m *MockRotationService) IsAutoRotationRunning() bool { return m.running }

func (m *MockRotationService) GetRotationHistory(path string) []RotationHistory {
	history := make([]RotationHistory, 0)
	for _, h := range m.history {
		if h.Path == path {
			history = append(history, h)
		}
	}
	return history
}

func (m *MockRotationService) GetAllRotationHistory() []RotationHistory { return m.history }

func (m *MockRotationService) CleanupOldVersions(ctx context.Context, path string, keepVersions int) error {
	return nil
}

// TestKVManager 测试 KV 管理器
func TestKVManager(t *testing.T) {
	ctx := context.Background()
	kv := NewMockKVManager()

	t.Run("Put and Get", func(t *testing.T) {
		path := "secret/test"
		data := map[string]interface{}{"username": "admin", "password": "secret123"}

		if err := kv.Put(ctx, path, data); err != nil {
			t.Errorf("Put() error = %v", err)
			return
		}

		secret, err := kv.Get(ctx, path)
		if err != nil {
			t.Errorf("Get() error = %v", err)
			return
		}

		if secret.Data["username"] != "admin" {
			t.Errorf("expected username = admin, got %v", secret.Data["username"])
		}
	})

	t.Run("List secrets", func(t *testing.T) {
		kv.Put(ctx, "secret/app1", map[string]interface{}{"key": "value1"})
		kv.Put(ctx, "secret/app2", map[string]interface{}{"key": "value2"})

		keys, err := kv.List(ctx, "secret/")
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}
		if len(keys) < 2 {
			t.Errorf("expected at least 2 keys, got %d", len(keys))
		}
	})

	t.Run("Delete secret", func(t *testing.T) {
		path := "secret/to-delete"
		kv.Put(ctx, path, map[string]interface{}{"key": "value"})

		if err := kv.Delete(ctx, path); err != nil {
			t.Errorf("Delete() error = %v", err)
			return
		}

		if _, err := kv.Get(ctx, path); err == nil {
			t.Error("expected error after delete, got nil")
		}
	})

	t.Run("Version management", func(t *testing.T) {
		path := "secret/versioned"
		kv.Put(ctx, path, map[string]interface{}{"version": "1"})
		kv.Put(ctx, path, map[string]interface{}{"version": "2"})
		kv.Put(ctx, path, map[string]interface{}{"version": "3"})

		versions, err := kv.GetVersions(ctx, path)
		if err != nil {
			t.Errorf("GetVersions() error = %v", err)
			return
		}
		if len(versions) != 3 {
			t.Errorf("expected 3 versions, got %d", len(versions))
		}
	})
}

// TestDBManager 测试数据库凭据管理器
func TestDBManager(t *testing.T) {
	ctx := context.Background()
	db := NewMockDBManager()

	t.Run("GetCredentials", func(t *testing.T) {
		creds, err := db.GetCredentials(ctx, "db-role")
		if err != nil {
			t.Errorf("GetCredentials() error = %v", err)
			return
		}
		if creds.Username == "" {
			t.Error("expected non-empty username")
		}
		if creds.Password == "" {
			t.Error("expected non-empty password")
		}
	})

	t.Run("RenewLease", func(t *testing.T) {
		if err := db.RenewLease(ctx, "lease-123", 3600); err != nil {
			t.Errorf("RenewLease() error = %v", err)
		}
	})

	t.Run("RevokeLease", func(t *testing.T) {
		if err := db.RevokeLease(ctx, "lease-123"); err != nil {
			t.Errorf("RevokeLease() error = %v", err)
		}
	})
}

// TestEncryptionService 测试加密服务
func TestEncryptionService(t *testing.T) {
	ctx := context.Background()
	enc := NewMockEncryptionService()

	t.Run("CreateKey and GetKeyInfo", func(t *testing.T) {
		if err := enc.CreateKey(ctx, "my-key", KeyTypeAES, nil); err != nil {
			t.Errorf("CreateKey() error = %v", err)
			return
		}
		info, err := enc.GetKeyInfo(ctx, "my-key")
		if err != nil {
			t.Errorf("GetKeyInfo() error = %v", err)
			return
		}
		if info.Name != "my-key" {
			t.Errorf("expected name = my-key, got %s", info.Name)
		}
	})

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		enc.CreateKey(ctx, "encrypt-key", KeyTypeAES, nil)
		plaintext := []byte("hello world")

		result, err := enc.Encrypt(ctx, "encrypt-key", plaintext, nil)
		if err != nil {
			t.Errorf("Encrypt() error = %v", err)
			return
		}

		decrypted, err := enc.Decrypt(ctx, "encrypt-key", result.Ciphertext, nil)
		if err != nil {
			t.Errorf("Decrypt() error = %v", err)
			return
		}

		if string(decrypted) != string(plaintext) {
			t.Errorf("expected %s, got %s", plaintext, decrypted)
		}
	})

	t.Run("Sign and Verify", func(t *testing.T) {
		enc.CreateKey(ctx, "sign-key", KeyTypeED25519, nil)
		data := []byte("data to sign")

		signResult, err := enc.Sign(ctx, "sign-key", data, nil)
		if err != nil {
			t.Errorf("Sign() error = %v", err)
			return
		}

		valid, err := enc.Verify(ctx, "sign-key", data, signResult.Signature, nil)
		if err != nil {
			t.Errorf("Verify() error = %v", err)
			return
		}
		if !valid {
			t.Error("expected signature to be valid")
		}
	})

	t.Run("Batch Encrypt", func(t *testing.T) {
		enc.CreateKey(ctx, "batch-key", KeyTypeAES, nil)
		items := []*BatchEncryptItem{
			{ID: "1", Plaintext: []byte("data1")},
			{ID: "2", Plaintext: []byte("data2")},
			{ID: "3", Plaintext: []byte("data3")},
		}

		results, err := enc.EncryptBatch(ctx, "batch-key", items)
		if err != nil {
			t.Errorf("EncryptBatch() error = %v", err)
			return
		}
		if len(results) != len(items) {
			t.Errorf("expected %d results, got %d", len(items), len(results))
		}
	})

	t.Run("RotateKey", func(t *testing.T) {
		enc.CreateKey(ctx, "rotate-key", KeyTypeAES, nil)
		enc.RotateKey(ctx, "rotate-key")

		info, err := enc.GetKeyInfo(ctx, "rotate-key")
		if err != nil {
			t.Errorf("GetKeyInfo() error = %v", err)
			return
		}
		if info.LatestVersion != 2 {
			t.Errorf("expected version 2 after rotation, got %d", info.LatestVersion)
		}
	})
}

// TestRotationService 测试密钥轮换服务
func TestRotationService(t *testing.T) {
	ctx := context.Background()
	rotation := NewMockRotationService()

	t.Run("Register and Get Policy", func(t *testing.T) {
		policy := &RotationPolicy{
			Name:           "test-policy",
			Path:           "secret/test",
			Interval:       24 * time.Hour,
			Enabled:        true,
			AutoRotate:     true,
			RetainVersions: 10,
		}

		if err := rotation.RegisterPolicy(policy); err != nil {
			t.Errorf("RegisterPolicy() error = %v", err)
			return
		}

		retrieved, err := rotation.GetPolicy(policy.ID)
		if err != nil {
			t.Errorf("GetPolicy() error = %v", err)
			return
		}
		if retrieved.Name != "test-policy" {
			t.Errorf("expected name = test-policy, got %s", retrieved.Name)
		}
	})

	t.Run("List Policies", func(t *testing.T) {
		rotation.RegisterPolicy(&RotationPolicy{Name: "policy-1", Path: "secret/app1"})
		rotation.RegisterPolicy(&RotationPolicy{Name: "policy-2", Path: "secret/app2"})

		policies := rotation.ListPolicies()
		if len(policies) < 2 {
			t.Errorf("expected at least 2 policies, got %d", len(policies))
		}
	})

	t.Run("Manual Rotate", func(t *testing.T) {
		if err := rotation.Rotate(ctx, "secret/manual-rotate"); err != nil {
			t.Errorf("Rotate() error = %v", err)
			return
		}

		history := rotation.GetRotationHistory("secret/manual-rotate")
		if len(history) == 0 {
			t.Error("expected rotation history")
		}
		if history[0].Status != "success" {
			t.Errorf("expected status = success, got %s", history[0].Status)
		}
	})

	t.Run("Auto Rotation", func(t *testing.T) {
		if err := rotation.StartAutoRotation(); err != nil {
			t.Errorf("StartAutoRotation() error = %v", err)
			return
		}
		if !rotation.IsAutoRotationRunning() {
			t.Error("expected auto rotation to be running")
		}

		if err := rotation.StopAutoRotation(); err != nil {
			t.Errorf("StopAutoRotation() error = %v", err)
			return
		}
		if rotation.IsAutoRotationRunning() {
			t.Error("expected auto rotation to be stopped")
		}
	})

	t.Run("Unregister Policy", func(t *testing.T) {
		policy := &RotationPolicy{Name: "to-unregister", Path: "secret/temp"}
		rotation.RegisterPolicy(policy)

		if err := rotation.UnregisterPolicy(policy.ID); err != nil {
			t.Errorf("UnregisterPolicy() error = %v", err)
			return
		}
		if _, err := rotation.GetPolicy(policy.ID); err == nil {
			t.Error("expected error after unregister, got nil")
		}
	})
}

// TestWithRetry 测试重试机制
func TestWithRetry(t *testing.T) {
	t.Run("successful on first attempt", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		result, err := withRetry(ctx, 3, 100*time.Millisecond, func() (string, error) {
			callCount++
			return "success", nil
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "success" {
			t.Errorf("expected 'success', got '%s'", result)
		}
		if callCount != 1 {
			t.Errorf("expected 1 call, got %d", callCount)
		}
	})

	t.Run("successful after retries", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		result, err := withRetry(ctx, 3, 100*time.Millisecond, func() (string, error) {
			callCount++
			if callCount < 3 {
				return "", errors.New("temporary error")
			}
			return "success", nil
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "success" {
			t.Errorf("expected 'success', got '%s'", result)
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got %d", callCount)
		}
	})

	t.Run("failed after max retries", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		_, err := withRetry(ctx, 2, 10*time.Millisecond, func() (string, error) {
			callCount++
			return "", errors.New("persistent error")
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got %d", callCount)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := withRetry(ctx, 3, 100*time.Millisecond, func() (string, error) {
			return "", errors.New("error")
		})

		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})
}

// TestHealthStatus 测试健康检查
func TestHealthStatus(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	status, err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() error = %v", err)
		return
	}
	if !status.Initialized {
		t.Error("expected initialized = true")
	}
	if status.Sealed {
		t.Error("expected sealed = false")
	}
}

// TestClose 测试客户端关闭
func TestClose(t *testing.T) {
	client := NewMockClient()

	if !client.IsConnected() {
		t.Error("expected client to be connected")
	}

	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if client.IsConnected() {
		t.Error("expected client to be disconnected after Close()")
	}
}

// BenchmarkEncrypt 加密性能基准测试
func BenchmarkEncrypt(b *testing.B) {
	ctx := context.Background()
	enc := NewMockEncryptionService()
	enc.CreateKey(ctx, "bench-key", KeyTypeAES, nil)
	plaintext := []byte("benchmark data for encryption")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.Encrypt(ctx, "bench-key", plaintext, nil)
	}
}

// BenchmarkDecrypt 解密性能基准测试
func BenchmarkDecrypt(b *testing.B) {
	ctx := context.Background()
	enc := NewMockEncryptionService()
	enc.CreateKey(ctx, "bench-key", KeyTypeAES, nil)
	plaintext := []byte("benchmark data for decryption")
	result, _ := enc.Encrypt(ctx, "bench-key", plaintext, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.Decrypt(ctx, "bench-key", result.Ciphertext, nil)
	}
}

// TestIntegration 集成测试
func TestIntegration(t *testing.T) {
	ctx := context.Background()
	client := NewMockClient()

	t.Run("complete workflow", func(t *testing.T) {
		// 1. 存储密钥
		err := client.KV().Put(ctx, "secret/app/config", map[string]interface{}{
			"api_key":     "secret-api-key",
			"db_password": "secret-db-password",
		})
		if err != nil {
			t.Fatalf("failed to store secret: %v", err)
		}

		// 2. 读取密钥
		secret, err := client.KV().Get(ctx, "secret/app/config")
		if err != nil {
			t.Fatalf("failed to get secret: %v", err)
		}
		if secret.Data["api_key"] != "secret-api-key" {
			t.Errorf("unexpected api_key: %v", secret.Data["api_key"])
		}

		// 3. 获取数据库凭据
		creds, err := client.DB().GetCredentials(ctx, "db-readonly")
		if err != nil {
			t.Fatalf("failed to get db credentials: %v", err)
		}
		if creds.Username == "" {
			t.Error("expected non-empty db username")
		}

		// 4. 创建加密密钥
		err = client.Encryption().CreateKey(ctx, "app-encryption-key", KeyTypeAES, nil)
		if err != nil {
			t.Fatalf("failed to create encryption key: %v", err)
		}

		// 5. 加密敏感数据
		sensitiveData := []byte("sensitive information")
		encrypted, err := client.Encryption().Encrypt(ctx, "app-encryption-key", sensitiveData, nil)
		if err != nil {
			t.Fatalf("failed to encrypt data: %v", err)
		}

		// 6. 解密数据
		decrypted, err := client.Encryption().Decrypt(ctx, "app-encryption-key", encrypted.Ciphertext, nil)
		if err != nil {
			t.Fatalf("failed to decrypt data: %v", err)
		}
		if string(decrypted) != string(sensitiveData) {
			t.Error("decrypted data doesn't match original")
		}

		// 7. 注册轮换策略
		policy := &RotationPolicy{
			Name:           "app-key-rotation",
			Path:           "secret/app/config",
			Interval:       24 * time.Hour,
			Enabled:        true,
			AutoRotate:     false,
			RetainVersions: 5,
		}
		err = client.Rotation().RegisterPolicy(policy)
		if err != nil {
			t.Fatalf("failed to register rotation policy: %v", err)
		}

		// 8. 手动轮换
		err = client.Rotation().Rotate(ctx, "secret/app/config")
		if err != nil {
			t.Fatalf("failed to rotate secret: %v", err)
		}

		// 9. 检查轮换历史
		history := client.Rotation().GetRotationHistory("secret/app/config")
		if len(history) == 0 {
			t.Error("expected rotation history to be recorded")
		}

		// 10. 健康检查
		status, err := client.Health(ctx)
		if err != nil {
			t.Fatalf("health check failed: %v", err)
		}
		if !status.Initialized {
			t.Error("vault not initialized")
		}

		t.Log("Complete workflow test passed!")
	})
}

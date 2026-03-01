// Package vault provides HashiCorp Vault integration for secure secret management.
//
// encryption.go 提供加密服务，包括：
// 1. 数据加密/解密
// 2. 密钥派生
// 3. 签名和验证
// 4. 批量加密操作
//
// 设计原则：
// 1. 支持 Transit 密钥引擎
// 2. 统一的加密接口
// 3. 支持上下文控制
// 4. 提供清晰的错误信息
//
// 使用场景：
// - 应用数据的加密存储
// - 敏感信息的加密传输
// - 数据签名和验证
// - 密钥派生
package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// EncryptionAlgorithm 定义加密算法类型
type EncryptionAlgorithm string

const (
	// AlgorithmAES256GCM AES-256-GCM 算法
	AlgorithmAES256GCM EncryptionAlgorithm = "aes256-gcm"
	// AlgorithmChaCha20Poly1305 ChaCha20-Poly1305 算法
	AlgorithmChaCha20Poly1305 EncryptionAlgorithm = "chacha20-poly1305"
	// AlgorithmRSA2048 RSA-2048 算法
	AlgorithmRSA2048 EncryptionAlgorithm = "rsa-2048"
	// AlgorithmRSA4096 RSA-4096 算法
	AlgorithmRSA4096 EncryptionAlgorithm = "rsa-4096"
	// AlgorithmECDSAP256 ECDSA P-256 算法
	AlgorithmECDSAP256 EncryptionAlgorithm = "ecdsa-p256"
	// AlgorithmECDSAP384 ECDSA P-384 算法
	AlgorithmECDSAP384 EncryptionAlgorithm = "ecdsa-p384"
	// AlgorithmED25519 Ed25519 算法
	AlgorithmED25519 EncryptionAlgorithm = "ed25519"
)

// KeyType 定义密钥类型
type KeyType string

const (
	// KeyTypeAES AES 对称密钥
	KeyTypeAES KeyType = "aes256-gcm96"
	// KeyTypeChaCha20 ChaCha20 对称密钥
	KeyTypeChaCha20 KeyType = "chacha20-poly1305"
	// KeyTypeRSA2048 RSA-2048 非对称密钥
	KeyTypeRSA2048 KeyType = "rsa-2048"
	// KeyTypeRSA4096 RSA-4096 非对称密钥
	KeyTypeRSA4096 KeyType = "rsa-4096"
	// KeyTypeECDSAP256 ECDSA P-256 非对称密钥
	KeyTypeECDSAP256 KeyType = "ecdsa-p256"
	// KeyTypeECDSAP384 ECDSA P-384 非对称密钥
	KeyTypeECDSAP384 KeyType = "ecdsa-p384"
	// KeyTypeECDSAP521 ECDSA P-521 非对称密钥
	KeyTypeECDSAP521 KeyType = "ecdsa-p521"
	// KeyTypeED25519 Ed25519 非对称密钥
	KeyTypeED25519 KeyType = "ed25519"
)

// EncryptionService 是加密服务接口。
//
// 提供基于 Vault Transit 密钥引擎的加密/解密、签名/验证功能。
// 所有加密操作都使用 Vault 管理的密钥，密钥不会离开 Vault。
type EncryptionService interface {
	// CreateKey 创建新的加密密钥
	CreateKey(ctx context.Context, name string, keyType KeyType, opts *KeyOptions) error

	// DeleteKey 删除加密密钥
	DeleteKey(ctx context.Context, name string) error

	// ListKeys 列出所有加密密钥
	ListKeys(ctx context.Context) ([]string, error)

	// GetKeyInfo 获取密钥信息
	GetKeyInfo(ctx context.Context, name string) (*KeyInfo, error)

	// RotateKey 轮换加密密钥
	RotateKey(ctx context.Context, name string) error

	// ExportKey 导出密钥（如果允许）
	ExportKey(ctx context.Context, name string, keyType string) (map[string]string, error)

	// Encrypt 加密数据
	Encrypt(ctx context.Context, keyName string, plaintext []byte, opts *EncryptOptions) (*EncryptResult, error)

	// EncryptBatch 批量加密数据
	EncryptBatch(ctx context.Context, keyName string, items []*BatchEncryptItem) ([]*EncryptResult, error)

	// Decrypt 解密数据
	Decrypt(ctx context.Context, keyName string, ciphertext string, opts *DecryptOptions) ([]byte, error)

	// DecryptBatch 批量解密数据
	DecryptBatch(ctx context.Context, keyName string, items []*BatchDecryptItem) ([][]byte, error)

	// Rewrap 使用新版本密钥重新加密数据
	Rewrap(ctx context.Context, keyName string, ciphertext string) (*EncryptResult, error)

	// Sign 使用私钥签名数据
	Sign(ctx context.Context, keyName string, input []byte, opts *SignOptions) (*SignResult, error)

	// Verify 使用公钥验证签名
	Verify(ctx context.Context, keyName string, input []byte, signature string, opts *VerifyOptions) (bool, error)

	// Hash 计算数据哈希
	Hash(ctx context.Context, input []byte, algorithm string) (string, error)

	// GenerateHMAC 生成 HMAC
	GenerateHMAC(ctx context.Context, keyName string, input []byte, algorithm string) (string, error)

	// VerifyHMAC 验证 HMAC
	VerifyHMAC(ctx context.Context, keyName string, input []byte, hmac string, algorithm string) (bool, error)
}

// KeyOptions 是创建密钥的选项
type KeyOptions struct {
	// Derived 表示是否启用密钥派生
	Derived bool `json:"derived,omitempty"`
	// Exportable 表示密钥是否可导出
	Exportable bool `json:"exportable,omitempty"`
	// AllowPlaintextBackup 是否允许明文备份
	AllowPlaintextBackup bool `json:"allow_plaintext_backup,omitempty"`
	// AutoRotatePeriod 自动轮换周期（0 表示禁用）
	AutoRotatePeriod time.Duration `json:"auto_rotate_period,omitempty"`
	// ConvergentEncryption 是否启用收敛加密
	ConvergentEncryption bool `json:"convergent_encryption,omitempty"`
}

// KeyInfo 包含密钥信息
type KeyInfo struct {
	// Name 是密钥名称
	Name string `json:"name"`
	// Type 是密钥类型
	Type string `json:"type"`
	// Keys 是所有密钥版本（版本号 -> 公钥）
	Keys map[string]interface{} `json:"keys"`
	// LatestVersion 是最新版本号
	LatestVersion int `json:"latest_version"`
	// MinDecryptionVersion 是最小解密版本
	MinDecryptionVersion int `json:"min_decryption_version"`
	// MinEncryptionVersion 是最小加密版本
	MinEncryptionVersion int `json:"min_encryption_version"`
	// SupportsEncryption 是否支持加密
	SupportsEncryption bool `json:"supports_encryption"`
	// SupportsDecryption 是否支持解密
	SupportsDecryption bool `json:"supports_decryption"`
	// SupportsSigning 是否支持签名
	SupportsSigning bool `json:"supports_signing"`
	// SupportsDerivation 是否支持派生
	SupportsDerivation bool `json:"supports_derivation"`
	// Derived 是否启用派生
	Derived bool `json:"derived"`
	// Exportable 是否可导出
	Exportable bool `json:"exportable"`
	// AllowPlaintextBackup 是否允许明文备份
	AllowPlaintextBackup bool `json:"allow_plaintext_backup"`
	// AutoRotatePeriod 自动轮换周期
	AutoRotatePeriod time.Duration `json:"auto_rotate_period"`
}

// EncryptOptions 是加密选项
type EncryptOptions struct {
	// Context 是加密上下文（用于派生密钥）
	Context []byte `json:"context,omitempty"`
	// Nonce 是非ce（某些算法需要）
	Nonce []byte `json:"nonce,omitempty"`
	// KeyVersion 指定使用的密钥版本（0 表示最新版本）
	KeyVersion int `json:"key_version,omitempty"`
	// AssociatedData 是附加认证数据（AEAD 算法）
	AssociatedData []byte `json:"associated_data,omitempty"`
}

// EncryptResult 包含加密结果
type EncryptResult struct {
	// Ciphertext 是加密后的数据
	Ciphertext string `json:"ciphertext"`
	// KeyVersion 是使用的密钥版本
	KeyVersion int `json:"key_version"`
}

// BatchEncryptItem 是批量加密的单个项目
type BatchEncryptItem struct {
	// ID 是项目标识
	ID string `json:"id"`
	// Plaintext 是要加密的明文
	Plaintext []byte `json:"plaintext"`
	// Context 是加密上下文
	Context []byte `json:"context,omitempty"`
	// Nonce 是非ce
	Nonce []byte `json:"nonce,omitempty"`
}

// DecryptOptions 是解密选项
type DecryptOptions struct {
	// Context 是解密上下文（用于派生密钥）
	Context []byte `json:"context,omitempty"`
	// Nonce 是非ce（某些算法需要）
	Nonce []byte `json:"nonce,omitempty"`
	// AssociatedData 是附加认证数据（AEAD 算法）
	AssociatedData []byte `json:"associated_data,omitempty"`
}

// BatchDecryptItem 是批量解密的单个项目
type BatchDecryptItem struct {
	// ID 是项目标识
	ID string `json:"id"`
	// Ciphertext 是要解密的密文
	Ciphertext string `json:"ciphertext"`
	// Context 是解密上下文
	Context []byte `json:"context,omitempty"`
	// Nonce 是非ce
	Nonce []byte `json:"nonce,omitempty"`
}

// SignOptions 是签名选项
type SignOptions struct {
	// HashAlgorithm 是哈希算法（sha2-256、sha2-384、sha2-512 等）
	HashAlgorithm string `json:"hash_algorithm,omitempty"`
	// KeyVersion 指定使用的密钥版本（0 表示最新版本）
	KeyVersion int `json:"key_version,omitempty"`
	// Prehashed 表示输入是否已哈希
	Prehashed bool `json:"prehashed,omitempty"`
	// SignatureAlgorithm 是签名算法（RSA-PSS、RSA-PKCS1v15 等）
	SignatureAlgorithm string `json:"signature_algorithm,omitempty"`
	// MarshalingAlgorithm 是编组算法（asn1、jws）
	MarshalingAlgorithm string `json:"marshaling_algorithm,omitempty"`
}

// SignResult 包含签名结果
type SignResult struct {
	// Signature 是签名值
	Signature string `json:"signature"`
	// KeyVersion 是使用的密钥版本
	KeyVersion int `json:"key_version"`
}

// VerifyOptions 是验证选项
type VerifyOptions struct {
	// HashAlgorithm 是哈希算法
	HashAlgorithm string `json:"hash_algorithm,omitempty"`
	// SignatureAlgorithm 是签名算法
	SignatureAlgorithm string `json:"signature_algorithm,omitempty"`
	// MarshalingAlgorithm 是编组算法
	MarshalingAlgorithm string `json:"marshaling_algorithm,omitempty"`
}

// encryptionService 是 EncryptionService 的实现
type encryptionService struct {
	client      *api.Client
	transitPath string
}

// newEncryptionService 创建新的加密服务
func newEncryptionService(client *api.Client) EncryptionService {
	return &encryptionService{
		client:      client,
		transitPath: "transit",
	}
}

// CreateKey 创建新的加密密钥
func (e *encryptionService) CreateKey(ctx context.Context, name string, keyType KeyType, opts *KeyOptions) error {
	if opts == nil {
		opts = &KeyOptions{}
	}

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		path := fmt.Sprintf("%s/keys/%s", e.transitPath, name)
		data := map[string]interface{}{
			"type": keyType,
		}

		if opts.Derived {
			data["derived"] = true
		}
		if opts.Exportable {
			data["exportable"] = true
		}
		if opts.AllowPlaintextBackup {
			data["allow_plaintext_backup"] = true
		}
		if opts.AutoRotatePeriod > 0 {
			data["auto_rotate_period"] = opts.AutoRotatePeriod.String()
		}
		if opts.ConvergentEncryption {
			data["convergent_encryption"] = true
		}

		_, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to create key %s: %w", name, err)
		}

		return struct{}{}, nil
	})
}

// DeleteKey 删除加密密钥
func (e *encryptionService) DeleteKey(ctx context.Context, name string) error {
	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		path := fmt.Sprintf("%s/keys/%s", e.transitPath, name)
		_, err := e.client.Logical().DeleteWithContext(ctx, path)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to delete key %s: %w", name, err)
		}

		return struct{}{}, nil
	})
}

// ListKeys 列出所有加密密钥
func (e *encryptionService) ListKeys(ctx context.Context) ([]string, error) {
	return withRetry(ctx, 3, 1*time.Second, func() ([]string, error) {
		path := fmt.Sprintf("%s/keys", e.transitPath)
		secret, err := e.client.Logical().ListWithContext(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to list keys: %w", err)
		}

		if secret == nil {
			return []string{}, nil
		}

		keys, ok := secret.Data["keys"].([]interface{})
		if !ok {
			return []string{}, nil
		}

		result := make([]string, len(keys))
		for i, key := range keys {
			result[i] = key.(string)
		}

		return result, nil
	})
}

// GetKeyInfo 获取密钥信息
func (e *encryptionService) GetKeyInfo(ctx context.Context, name string) (*KeyInfo, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (*KeyInfo, error) {
		path := fmt.Sprintf("%s/keys/%s", e.transitPath, name)
		secret, err := e.client.Logical().ReadWithContext(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to get key info %s: %w", name, err)
		}

		if secret == nil {
			return nil, fmt.Errorf("key not found: %s", name)
		}

		data := secret.Data
		info := &KeyInfo{
			Name:                 getString(data, "name"),
			Type:                 getString(data, "type"),
			Keys:                 getMapInterface(data, "keys"),
			LatestVersion:        getInt(data, "latest_version"),
			MinDecryptionVersion: getInt(data, "min_decryption_version"),
			MinEncryptionVersion: getInt(data, "min_encryption_version"),
			Derived:              getBool(data, "derived"),
			Exportable:           getBool(data, "exportable"),
			AllowPlaintextBackup: getBool(data, "allow_plaintext_backup"),
		}

		// 获取功能支持信息
		supports := getMapInterface(data, "supports_encryption")
		if val, ok := supports["encryption"].(bool); ok {
			info.SupportsEncryption = val
		}
		if val, ok := supports["decryption"].(bool); ok {
			info.SupportsDecryption = val
		}
		if val, ok := supports["signing"].(bool); ok {
			info.SupportsSigning = val
		}
		if val, ok := supports["derivation"].(bool); ok {
			info.SupportsDerivation = val
		}

		return info, nil
	})
}

// RotateKey 轮换加密密钥
func (e *encryptionService) RotateKey(ctx context.Context, name string) error {
	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		path := fmt.Sprintf("%s/keys/%s/rotate", e.transitPath, name)
		_, err := e.client.Logical().WriteWithContext(ctx, path, nil)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to rotate key %s: %w", name, err)
		}

		return struct{}{}, nil
	})
}

// ExportKey 导出密钥（如果允许）
func (e *encryptionService) ExportKey(ctx context.Context, name string, keyType string) (map[string]string, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (map[string]string, error) {
		path := fmt.Sprintf("%s/export/%s/%s", e.transitPath, keyType, name)
		secret, err := e.client.Logical().ReadWithContext(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to export key %s: %w", name, err)
		}

		if secret == nil {
			return nil, fmt.Errorf("key not found or not exportable: %s", name)
		}

		result := make(map[string]string)
		if keys, ok := secret.Data["keys"].(map[string]interface{}); ok {
			for version, key := range keys {
				if keyStr, ok := key.(string); ok {
					result[version] = keyStr
				}
			}
		}

		return result, nil
	})
}

// Encrypt 加密数据
func (e *encryptionService) Encrypt(ctx context.Context, keyName string, plaintext []byte, opts *EncryptOptions) (*EncryptResult, error) {
	if opts == nil {
		opts = &EncryptOptions{}
	}

	return withRetry(ctx, 3, 1*time.Second, func() (*EncryptResult, error) {
		path := fmt.Sprintf("%s/encrypt/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString(plaintext),
		}

		if len(opts.Context) > 0 {
			data["context"] = base64.StdEncoding.EncodeToString(opts.Context)
		}
		if len(opts.Nonce) > 0 {
			data["nonce"] = base64.StdEncoding.EncodeToString(opts.Nonce)
		}
		if opts.KeyVersion > 0 {
			data["key_version"] = opts.KeyVersion
		}
		if len(opts.AssociatedData) > 0 {
			data["associated_data"] = base64.StdEncoding.EncodeToString(opts.AssociatedData)
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", err)
		}

		ciphertext := getString(secret.Data, "ciphertext")
		keyVersion := getInt(secret.Data, "key_version")

		return &EncryptResult{
			Ciphertext: ciphertext,
			KeyVersion: keyVersion,
		}, nil
	})
}

// EncryptBatch 批量加密数据
func (e *encryptionService) EncryptBatch(ctx context.Context, keyName string, items []*BatchEncryptItem) ([]*EncryptResult, error) {
	return withRetry(ctx, 3, 1*time.Second, func() ([]*EncryptResult, error) {
		path := fmt.Sprintf("%s/encrypt/%s", e.transitPath, keyName)
		batchInput := make([]map[string]interface{}, len(items))

		for i, item := range items {
			batchItem := map[string]interface{}{
				"plaintext": base64.StdEncoding.EncodeToString(item.Plaintext),
			}
			if len(item.Context) > 0 {
				batchItem["context"] = base64.StdEncoding.EncodeToString(item.Context)
			}
			if len(item.Nonce) > 0 {
				batchItem["nonce"] = base64.StdEncoding.EncodeToString(item.Nonce)
			}
			batchInput[i] = batchItem
		}

		data := map[string]interface{}{
			"batch_input": batchInput,
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to batch encrypt: %w", err)
		}

		batchResults, ok := secret.Data["batch_results"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid batch response format")
		}

		results := make([]*EncryptResult, len(batchResults))
		for i, result := range batchResults {
			if resultMap, ok := result.(map[string]interface{}); ok {
				results[i] = &EncryptResult{
					Ciphertext: getString(resultMap, "ciphertext"),
					KeyVersion: getInt(resultMap, "key_version"),
				}
			}
		}

		return results, nil
	})
}

// Decrypt 解密数据
func (e *encryptionService) Decrypt(ctx context.Context, keyName string, ciphertext string, opts *DecryptOptions) ([]byte, error) {
	if opts == nil {
		opts = &DecryptOptions{}
	}

	return withRetry(ctx, 3, 1*time.Second, func() ([]byte, error) {
		path := fmt.Sprintf("%s/decrypt/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"ciphertext": ciphertext,
		}

		if len(opts.Context) > 0 {
			data["context"] = base64.StdEncoding.EncodeToString(opts.Context)
		}
		if len(opts.Nonce) > 0 {
			data["nonce"] = base64.StdEncoding.EncodeToString(opts.Nonce)
		}
		if len(opts.AssociatedData) > 0 {
			data["associated_data"] = base64.StdEncoding.EncodeToString(opts.AssociatedData)
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}

		plaintextB64 := getString(secret.Data, "plaintext")
		plaintext, err := base64.StdEncoding.DecodeString(plaintextB64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode plaintext: %w", err)
		}

		return plaintext, nil
	})
}

// DecryptBatch 批量解密数据
func (e *encryptionService) DecryptBatch(ctx context.Context, keyName string, items []*BatchDecryptItem) ([][]byte, error) {
	return withRetry(ctx, 3, 1*time.Second, func() ([][]byte, error) {
		path := fmt.Sprintf("%s/decrypt/%s", e.transitPath, keyName)
		batchInput := make([]map[string]interface{}, len(items))

		for i, item := range items {
			batchItem := map[string]interface{}{
				"ciphertext": item.Ciphertext,
			}
			if len(item.Context) > 0 {
				batchItem["context"] = base64.StdEncoding.EncodeToString(item.Context)
			}
			if len(item.Nonce) > 0 {
				batchItem["nonce"] = base64.StdEncoding.EncodeToString(item.Nonce)
			}
			batchInput[i] = batchItem
		}

		data := map[string]interface{}{
			"batch_input": batchInput,
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to batch decrypt: %w", err)
		}

		batchResults, ok := secret.Data["batch_results"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid batch response format")
		}

		results := make([][]byte, len(batchResults))
		for i, result := range batchResults {
			if resultMap, ok := result.(map[string]interface{}); ok {
				plaintextB64 := getString(resultMap, "plaintext")
				plaintext, err := base64.StdEncoding.DecodeString(plaintextB64)
				if err != nil {
					results[i] = nil
				} else {
					results[i] = plaintext
				}
			}
		}

		return results, nil
	})
}

// Rewrap 使用新版本密钥重新加密数据
func (e *encryptionService) Rewrap(ctx context.Context, keyName string, ciphertext string) (*EncryptResult, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (*EncryptResult, error) {
		path := fmt.Sprintf("%s/rewrap/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"ciphertext": ciphertext,
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to rewrap data: %w", err)
		}

		return &EncryptResult{
			Ciphertext: getString(secret.Data, "ciphertext"),
			KeyVersion: getInt(secret.Data, "key_version"),
		}, nil
	})
}

// Sign 使用私钥签名数据
func (e *encryptionService) Sign(ctx context.Context, keyName string, input []byte, opts *SignOptions) (*SignResult, error) {
	if opts == nil {
		opts = &SignOptions{}
	}

	return withRetry(ctx, 3, 1*time.Second, func() (*SignResult, error) {
		path := fmt.Sprintf("%s/sign/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"input": base64.StdEncoding.EncodeToString(input),
		}

		if opts.HashAlgorithm != "" {
			data["hash_algorithm"] = opts.HashAlgorithm
		}
		if opts.KeyVersion > 0 {
			data["key_version"] = opts.KeyVersion
		}
		if opts.Prehashed {
			data["prehashed"] = true
		}
		if opts.SignatureAlgorithm != "" {
			data["signature_algorithm"] = opts.SignatureAlgorithm
		}
		if opts.MarshalingAlgorithm != "" {
			data["marshaling_algorithm"] = opts.MarshalingAlgorithm
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return nil, fmt.Errorf("failed to sign data: %w", err)
		}

		return &SignResult{
			Signature:  getString(secret.Data, "signature"),
			KeyVersion: getInt(secret.Data, "key_version"),
		}, nil
	})
}

// Verify 使用公钥验证签名
func (e *encryptionService) Verify(ctx context.Context, keyName string, input []byte, signature string, opts *VerifyOptions) (bool, error) {
	if opts == nil {
		opts = &VerifyOptions{}
	}

	return withRetry(ctx, 3, 1*time.Second, func() (bool, error) {
		path := fmt.Sprintf("%s/verify/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"input":     base64.StdEncoding.EncodeToString(input),
			"signature": signature,
		}

		if opts.HashAlgorithm != "" {
			data["hash_algorithm"] = opts.HashAlgorithm
		}
		if opts.SignatureAlgorithm != "" {
			data["signature_algorithm"] = opts.SignatureAlgorithm
		}
		if opts.MarshalingAlgorithm != "" {
			data["marshaling_algorithm"] = opts.MarshalingAlgorithm
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return false, fmt.Errorf("failed to verify signature: %w", err)
		}

		return getBool(secret.Data, "valid"), nil
	})
}

// Hash 计算数据哈希
func (e *encryptionService) Hash(ctx context.Context, input []byte, algorithm string) (string, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (string, error) {
		path := fmt.Sprintf("%s/hash", e.transitPath)
		data := map[string]interface{}{
			"input":     base64.StdEncoding.EncodeToString(input),
			"algorithm": algorithm,
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return "", fmt.Errorf("failed to hash data: %w", err)
		}

		return getString(secret.Data, "sum"), nil
	})
}

// GenerateHMAC 生成 HMAC
func (e *encryptionService) GenerateHMAC(ctx context.Context, keyName string, input []byte, algorithm string) (string, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (string, error) {
		path := fmt.Sprintf("%s/hmac/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"input": base64.StdEncoding.EncodeToString(input),
		}

		if algorithm != "" {
			data["algorithm"] = algorithm
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return "", fmt.Errorf("failed to generate HMAC: %w", err)
		}

		return getString(secret.Data, "hmac"), nil
	})
}

// VerifyHMAC 验证 HMAC
func (e *encryptionService) VerifyHMAC(ctx context.Context, keyName string, input []byte, hmac string, algorithm string) (bool, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (bool, error) {
		path := fmt.Sprintf("%s/verify/%s", e.transitPath, keyName)
		data := map[string]interface{}{
			"input": base64.StdEncoding.EncodeToString(input),
			"hmac":  hmac,
		}

		if algorithm != "" {
			data["algorithm"] = algorithm
		}

		secret, err := e.client.Logical().WriteWithContext(ctx, path, data)
		if err != nil {
			return false, fmt.Errorf("failed to verify HMAC: %w", err)
		}

		return getBool(secret.Data, "valid"), nil
	})
}

// 辅助函数

// getInt 从 map 中获取整数值
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key].(int); ok {
		return val
	}
	if val, ok := data[key].(json.Number); ok {
		i, _ := val.Int64()
		return int(i)
	}
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

// getBool 从 map 中获取布尔值
func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

// getMapInterface 从 map 中获取 map 值
func getMapInterface(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}

// json 用于解析 json.Number
type json interface {
	Number
}

// Number 是 JSON 数字类型
type Number string

// Int64 转换为 int64
func (n Number) Int64() (int64, error) {
	var i int64
	_, err := fmt.Sscanf(string(n), "%d", &i)
	return i, err
}

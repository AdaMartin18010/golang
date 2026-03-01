// Package vault provides HashiCorp Vault integration for secure secret management.
//
// secret.go 提供密钥管理功能，包括：
// 1. KV v1/v2 密钥引擎支持
// 2. 动态数据库凭据管理
// 3. 密钥版本管理（KV v2）
// 4. 密钥元数据管理
//
// 设计原则：
// 1. 支持 KV v1 和 v2 两种密钥引擎版本
// 2. 统一的接口，隐藏版本差异
// 3. 支持上下文控制
// 4. 提供清晰的错误信息
//
// 使用场景：
// - 存储应用配置（API 密钥、数据库连接字符串等）
// - 动态获取数据库凭据
// - 密钥版本控制和回滚
// - 密钥元数据管理（TTL、描述等）
package vault

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// Secret 表示从 Vault 获取的密钥
type Secret struct {
	// Data 是密钥的数据内容
	Data map[string]interface{} `json:"data"`
	// Metadata 是密钥的元数据（仅 KV v2）
	Metadata *SecretMetadata `json:"metadata,omitempty"`
	// Version 是密钥版本（仅 KV v2）
	Version int `json:"version,omitempty"`
}

// SecretMetadata 包含密钥的元数据信息（仅 KV v2）
type SecretMetadata struct {
	// CreatedTime 是创建时间
	CreatedTime time.Time `json:"created_time"`
	// UpdatedTime 是最后更新时间
	UpdatedTime time.Time `json:"updated_time"`
	// DeletionTime 是删除时间（如果已删除）
	DeletionTime *time.Time `json:"deletion_time,omitempty"`
	// Destroyed 表示密钥是否已被销毁
	Destroyed bool `json:"destroyed"`
	// Version 是密钥版本号
	Version int `json:"version"`
}

// SecretVersion 表示密钥的特定版本
type SecretVersion struct {
	// Version 是版本号
	Version int `json:"version"`
	// CreatedTime 是创建时间
	CreatedTime time.Time `json:"created_time"`
	// DeletionTime 是删除时间（如果已删除）
	DeletionTime *time.Time `json:"deletion_time,omitempty"`
	// Destroyed 表示密钥是否已被销毁
	Destroyed bool `json:"destroyed"`
}

// DBCredentials 表示动态数据库凭据
type DBCredentials struct {
	// Username 是数据库用户名
	Username string `json:"username"`
	// Password 是数据库密码
	Password string `json:"password"`
	// LeaseID 是租约 ID
	LeaseID string `json:"lease_id"`
	// LeaseDuration 是租约有效期（秒）
	LeaseDuration int `json:"lease_duration"`
	// Renewable 表示租约是否可续期
	Renewable bool `json:"renewable"`
}

// KVManager 是 KV 密钥管理器接口。
//
// 支持 KV v1 和 v2 两种密钥引擎版本，提供统一的接口。
// KV v2 额外支持版本控制和元数据管理功能。
type KVManager interface {
	// Get 获取指定路径的密钥
	// 对于 KV v2，可以指定版本，默认为最新版本
	Get(ctx context.Context, path string, version ...int) (*Secret, error)

	// GetRaw 获取原始密钥数据（不包含元数据）
	GetRaw(ctx context.Context, path string) (map[string]interface{}, error)

	// Put 存储密钥到指定路径
	// 对于 KV v2，会自动创建新版本
	Put(ctx context.Context, path string, data map[string]interface{}) error

	// Delete 删除指定路径的密钥
	Delete(ctx context.Context, path string) error

	// DeleteVersion 删除指定版本的密钥（仅 KV v2）
	DeleteVersion(ctx context.Context, path string, version int) error

	// DestroyVersion 销毁指定版本的密钥（无法恢复，仅 KV v2）
	DestroyVersion(ctx context.Context, path string, version int) error

	// List 列出指定路径下的所有密钥
	List(ctx context.Context, path string) ([]string, error)

	// GetVersions 获取密钥的所有版本信息（仅 KV v2）
	GetVersions(ctx context.Context, path string) ([]SecretVersion, error)

	// Rollback 回滚到指定版本（仅 KV v2）
	// 实际上是复制指定版本的内容为新版本
	Rollback(ctx context.Context, path string, version int) error

	// GetMetadata 获取密钥的元数据（仅 KV v2）
	GetMetadata(ctx context.Context, path string) (*KVMetadata, error)

	// UpdateMetadata 更新密钥的元数据（仅 KV v2）
	UpdateMetadata(ctx context.Context, path string, metadata *KVMetadata) error

	// IsV2 返回是否使用 KV v2 引擎
	IsV2() bool
}

// KVMetadata 是 KV v2 密钥的元数据
type KVMetadata struct {
	// MaxVersions 是最大版本数（0 表示无限制）
	MaxVersions int `json:"max_versions,omitempty"`
	// CasRequired 表示是否需要 CAS（Check-And-Set）
	CasRequired bool `json:"cas_required,omitempty"`
	// DeleteVersionAfter 是自动删除版本的时间
	DeleteVersionAfter string `json:"delete_version_after,omitempty"`
	// CustomMetadata 是自定义元数据
	CustomMetadata map[string]string `json:"custom_metadata,omitempty"`
}

// DBManager 是动态数据库凭据管理器接口。
//
// 用于从 Vault 动态获取数据库凭据，支持自动租约续期。
type DBManager interface {
	// GetCredentials 从指定角色获取动态数据库凭据
	GetCredentials(ctx context.Context, role string) (*DBCredentials, error)

	// GetCredentialsFromPath 从指定路径获取动态数据库凭据
	GetCredentialsFromPath(ctx context.Context, mount, role string) (*DBCredentials, error)

	// RenewLease 续期租约
	RenewLease(ctx context.Context, leaseID string, increment int) error

	// RevokeLease 撤销租约
	RevokeLease(ctx context.Context, leaseID string) error
}

// kvManager 是 KVManager 的实现
type kvManager struct {
	client       *api.Client
	version      int
	dataPath     string
	metadataPath string
}

// newKVManager 创建新的 KV 管理器
func newKVManager(client *api.Client, version int) KVManager {
	return &kvManager{
		client:       client,
		version:      version,
		dataPath:     "secret/data",
		metadataPath: "secret/metadata",
	}
}

// normalizePath 规范化路径
func (k *kvManager) normalizePath(path string) string {
	path = strings.TrimPrefix(path, "secret/")
	path = strings.TrimPrefix(path, "data/")
	path = strings.TrimPrefix(path, "metadata/")
	return path
}

// Get 获取密钥
func (k *kvManager) Get(ctx context.Context, secretPath string, version ...int) (*Secret, error) {
	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (*Secret, error) {
		var secret *api.Secret
		var err error

		if k.version == 2 {
			fullPath := path.Join(k.dataPath, secretPath)
			if len(version) > 0 {
				secret, err = k.client.KVv2("secret").GetVersion(ctx, version[0], secretPath)
			} else {
				secret, err = k.client.KVv2("secret").Get(ctx, secretPath)
			}
			if err != nil {
				return nil, fmt.Errorf("failed to get secret from path %s: %w", fullPath, err)
			}

			return k.parseV2Secret(secret), nil
		}

		// KV v1
		fullPath := path.Join("secret", secretPath)
		secret, err = k.client.KVv1("secret").Get(ctx, secretPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get secret from path %s: %w", fullPath, err)
		}

		return &Secret{
			Data: secret.Data,
		}, nil
	})
}

// GetRaw 获取原始密钥数据
func (k *kvManager) GetRaw(ctx context.Context, secretPath string) (map[string]interface{}, error) {
	secret, err := k.Get(ctx, secretPath)
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}

// parseV2Secret 解析 KV v2 密钥
func (k *kvManager) parseV2Secret(secret *api.Secret) *Secret {
	if secret == nil {
		return nil
	}

	result := &Secret{
		Data: make(map[string]interface{}),
	}

	// 提取数据
	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		result.Data = data
	}

	// 提取元数据
	if metadata, ok := secret.Data["metadata"].(map[string]interface{}); ok {
		result.Metadata = &SecretMetadata{}

		if createdTime, ok := metadata["created_time"].(string); ok {
			result.Metadata.CreatedTime, _ = time.Parse(time.RFC3339, createdTime)
		}
		if updatedTime, ok := metadata["updated_time"].(string); ok {
			result.Metadata.UpdatedTime, _ = time.Parse(time.RFC3339, updatedTime)
		}
		if deletionTime, ok := metadata["deletion_time"].(string); ok && deletionTime != "" {
			t, _ := time.Parse(time.RFC3339, deletionTime)
			result.Metadata.DeletionTime = &t
		}
		if destroyed, ok := metadata["destroyed"].(bool); ok {
			result.Metadata.Destroyed = destroyed
		}
		if version, ok := metadata["version"].(json.Number); ok {
			v, _ := version.Int64()
			result.Metadata.Version = int(v)
			result.Version = int(v)
		}
	}

	return result
}

// Put 存储密钥
func (k *kvManager) Put(ctx context.Context, secretPath string, data map[string]interface{}) error {
	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		if k.version == 2 {
			_, err := k.client.KVv2("secret").Put(ctx, secretPath, data)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to put secret to path %s: %w", secretPath, err)
			}
		} else {
			err := k.client.KVv1("secret").Put(ctx, secretPath, data)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to put secret to path %s: %w", secretPath, err)
			}
		}
		return struct{}{}, nil
	})
}

// Delete 删除密钥
func (k *kvManager) Delete(ctx context.Context, secretPath string) error {
	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		if k.version == 2 {
			err := k.client.KVv2("secret").DeleteLatestVersion(ctx, secretPath)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to delete secret at path %s: %w", secretPath, err)
			}
		} else {
			err := k.client.KVv1("secret").Delete(ctx, secretPath)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to delete secret at path %s: %w", secretPath, err)
			}
		}
		return struct{}{}, nil
	})
}

// DeleteVersion 删除指定版本（仅 KV v2）
func (k *kvManager) DeleteVersion(ctx context.Context, secretPath string, version int) error {
	if k.version != 2 {
		return fmt.Errorf("delete version is only supported in KV v2")
	}

	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		err := k.client.KVv2("secret").DeleteVersions(ctx, secretPath, []int{version})
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to delete version %d of secret at path %s: %w", version, secretPath, err)
		}
		return struct{}{}, nil
	})
}

// DestroyVersion 销毁指定版本（仅 KV v2）
func (k *kvManager) DestroyVersion(ctx context.Context, secretPath string, version int) error {
	if k.version != 2 {
		return fmt.Errorf("destroy version is only supported in KV v2")
	}

	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		err := k.client.KVv2("secret").DestroyVersions(ctx, secretPath, []int{version})
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to destroy version %d of secret at path %s: %w", version, secretPath, err)
		}
		return struct{}{}, nil
	})
}

// List 列出密钥
func (k *kvManager) List(ctx context.Context, secretPath string) ([]string, error) {
	secretPath = strings.TrimSuffix(secretPath, "/")

	return withRetry(ctx, 3, 1*time.Second, func() ([]string, error) {
		var keys []string
		var err error

		if k.version == 2 {
			keys, err = k.client.KVv2("secret").List(ctx, secretPath)
		} else {
			keys, err = k.client.KVv1("secret").List(ctx, secretPath)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list secrets at path %s: %w", secretPath, err)
		}

		return keys, nil
	})
}

// GetVersions 获取所有版本（仅 KV v2）
func (k *kvManager) GetVersions(ctx context.Context, secretPath string) ([]SecretVersion, error) {
	if k.version != 2 {
		return nil, fmt.Errorf("get versions is only supported in KV v2")
	}

	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() ([]SecretVersion, error) {
		fullPath := path.Join(k.metadataPath, secretPath)
		secret, err := k.client.Logical().ReadWithContext(ctx, fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get versions for secret at path %s: %w", secretPath, err)
		}

		if secret == nil {
			return nil, fmt.Errorf("secret not found at path %s", secretPath)
		}

		return k.parseVersions(secret.Data)
	})
}

// parseVersions 解析版本信息
func (k *kvManager) parseVersions(data map[string]interface{}) ([]SecretVersion, error) {
	versionsData, ok := data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid versions data format")
	}

	versions := make([]SecretVersion, 0, len(versionsData))
	for versionStr, versionData := range versionsData {
		versionNum := 0
		fmt.Sscanf(versionStr, "%d", &versionNum)

		versionInfo, ok := versionData.(map[string]interface{})
		if !ok {
			continue
		}

		version := SecretVersion{Version: versionNum}

		if createdTime, ok := versionInfo["created_time"].(string); ok {
			version.CreatedTime, _ = time.Parse(time.RFC3339, createdTime)
		}
		if deletionTime, ok := versionInfo["deletion_time"].(string); ok && deletionTime != "" {
			t, _ := time.Parse(time.RFC3339, deletionTime)
			version.DeletionTime = &t
		}
		if destroyed, ok := versionInfo["destroyed"].(bool); ok {
			version.Destroyed = destroyed
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// Rollback 回滚到指定版本（仅 KV v2）
func (k *kvManager) Rollback(ctx context.Context, secretPath string, version int) error {
	if k.version != 2 {
		return fmt.Errorf("rollback is only supported in KV v2")
	}

	// 获取指定版本的数据
	oldSecret, err := k.Get(ctx, secretPath, version)
	if err != nil {
		return fmt.Errorf("failed to get version %d for rollback: %w", version, err)
	}

	// 存储为新版本
	if err := k.Put(ctx, secretPath, oldSecret.Data); err != nil {
		return fmt.Errorf("failed to rollback to version %d: %w", version, err)
	}

	return nil
}

// GetMetadata 获取元数据（仅 KV v2）
func (k *kvManager) GetMetadata(ctx context.Context, secretPath string) (*KVMetadata, error) {
	if k.version != 2 {
		return nil, fmt.Errorf("metadata is only supported in KV v2")
	}

	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (*KVMetadata, error) {
		metadata, err := k.client.KVv2("secret").GetMetadata(ctx, secretPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for secret at path %s: %w", secretPath, err)
		}

		return &KVMetadata{
			MaxVersions:        metadata.MaxVersions,
			CasRequired:        metadata.CASRequired,
			DeleteVersionAfter: metadata.DeleteVersionAfter,
			CustomMetadata:     metadata.CustomMetadata,
		}, nil
	})
}

// UpdateMetadata 更新元数据（仅 KV v2）
func (k *kvManager) UpdateMetadata(ctx context.Context, secretPath string, metadata *KVMetadata) error {
	if k.version != 2 {
		return fmt.Errorf("metadata is only supported in KV v2")
	}

	secretPath = k.normalizePath(secretPath)

	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		apiMetadata := &api.KVMetadata{
			MaxVersions:        metadata.MaxVersions,
			CASRequired:        metadata.CasRequired,
			DeleteVersionAfter: metadata.DeleteVersionAfter,
			CustomMetadata:     metadata.CustomMetadata,
		}

		err := k.client.KVv2("secret").UpdateMetadata(ctx, secretPath, apiMetadata)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to update metadata for secret at path %s: %w", secretPath, err)
		}

		return struct{}{}, nil
	})
}

// IsV2 返回是否使用 KV v2 引擎
func (k *kvManager) IsV2() bool {
	return k.version == 2
}

// dbManager 是 DBManager 的实现
type dbManager struct {
	client *api.Client
}

// newDBManager 创建新的数据库凭据管理器
func newDBManager(client *api.Client) DBManager {
	return &dbManager{client: client}
}

// GetCredentials 获取动态数据库凭据
func (d *dbManager) GetCredentials(ctx context.Context, role string) (*DBCredentials, error) {
	return d.GetCredentialsFromPath(ctx, "database", role)
}

// GetCredentialsFromPath 从指定路径获取动态数据库凭据
func (d *dbManager) GetCredentialsFromPath(ctx context.Context, mount, role string) (*DBCredentials, error) {
	return withRetry(ctx, 3, 1*time.Second, func() (*DBCredentials, error) {
		path := fmt.Sprintf("%s/creds/%s", mount, role)
		secret, err := d.client.Logical().ReadWithContext(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to get credentials from %s: %w", path, err)
		}

		if secret == nil {
			return nil, fmt.Errorf("no credentials found for role %s", role)
		}

		return &DBCredentials{
			Username:      getString(secret.Data, "username"),
			Password:      getString(secret.Data, "password"),
			LeaseID:       secret.LeaseID,
			LeaseDuration: secret.LeaseDuration,
			Renewable:     secret.Renewable,
		}, nil
	})
}

// RenewLease 续期租约
func (d *dbManager) RenewLease(ctx context.Context, leaseID string, increment int) error {
	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		_, err := d.client.Sys().RenewWithContext(ctx, leaseID, increment)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to renew lease %s: %w", leaseID, err)
		}
		return struct{}{}, nil
	})
}

// RevokeLease 撤销租约
func (d *dbManager) RevokeLease(ctx context.Context, leaseID string) error {
	return withRetry(ctx, 3, 1*time.Second, func() (struct{}, error) {
		err := d.client.Sys().RevokeWithContext(ctx, leaseID)
		if err != nil {
			return struct{}{}, fmt.Errorf("failed to revoke lease %s: %w", leaseID, err)
		}
		return struct{}{}, nil
	})
}

// getString 从 map 中获取字符串值
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

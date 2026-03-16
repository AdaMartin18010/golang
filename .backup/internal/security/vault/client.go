// Package vault provides HashiCorp Vault integration for secure secret management.
//
// 功能说明：
// 1. Vault 客户端连接管理（支持多种认证方式）
// 2. 连接重试和超时配置
// 3. 健康检查和连接池管理
// 4. 上下文支持，支持取消和超时
//
// 设计原则：
// 1. 使用接口抽象 Vault 客户端，便于测试和 mock
// 2. 支持上下文控制，便于超时和取消操作
// 3. 提供连接重试机制，提高可用性
// 4. 统一错误处理，提供清晰的错误信息
//
// 使用场景：
// - 应用启动时初始化 Vault 客户端
// - 运行时动态获取密钥
// - 数据库凭据的动态获取
// - 加密/解密操作
//
// 认证方式：
// - Token 认证
// - Kubernetes 认证
// - AppRole 认证
// - TLS 证书认证
//
// 示例：
//
//	// 创建配置
//	cfg := &vault.Config{
//	    Address:     "https://vault.example.com:8200",
//	    AuthMethod:  vault.AuthMethodToken,
//	    Token:       "your-vault-token",
//	    MaxRetries:  3,
//	    Timeout:     30 * time.Second,
//	}
//
//	// 创建客户端
//	client, err := vault.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 使用客户端
//	secret, err := client.KV().Get(ctx, "secret/data/myapp")
package vault

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/kubernetes"
)

// AuthMethod 定义 Vault 认证方式类型
type AuthMethod string

const (
	// AuthMethodToken 使用 Token 认证
	AuthMethodToken AuthMethod = "token"
	// AuthMethodKubernetes 使用 Kubernetes 认证
	AuthMethodKubernetes AuthMethod = "kubernetes"
	// AuthMethodAppRole 使用 AppRole 认证
	AuthMethodAppRole AuthMethod = "approle"
	// AuthMethodCert 使用 TLS 证书认证
	AuthMethodCert AuthMethod = "cert"
)

// Config 是 Vault 客户端的配置结构。
//
// 字段说明：
// - Address: Vault 服务器地址（格式：https://vault.example.com:8200）
// - AuthMethod: 认证方式（token、kubernetes、approle、cert）
// - Token: Token 认证方式的令牌（仅当 AuthMethod=token 时使用）
// - Namespace: Vault 命名空间（企业版功能，可选）
// - MaxRetries: 最大重试次数（默认：3）
// - Timeout: 请求超时时间（默认：30s）
// - RetryDelay: 重试间隔时间（默认：1s）
// - KVVersion: KV 密钥引擎版本（1 或 2，默认：2）
// - TLSCertPath: TLS 证书路径（用于 cert 认证）
// - TLSKeyPath: TLS 私钥路径（用于 cert 认证）
// - CACertPath: CA 证书路径（用于 TLS 验证）
// - InsecureSkipVerify: 是否跳过 TLS 证书验证（开发环境使用，默认：false）
// - K8sRole: Kubernetes 角色名称（用于 kubernetes 认证）
// - K8sTokenPath: Kubernetes 服务账户令牌路径（默认：/var/run/secrets/kubernetes.io/serviceaccount/token）
// - AppRoleID: AppRole ID（用于 approle 认证）
// - AppSecretID: AppRole Secret ID（用于 approle 认证）
type Config struct {
	Address            string        `mapstructure:"address"`
	AuthMethod         AuthMethod    `mapstructure:"auth_method"`
	Token              string        `mapstructure:"token"`
	Namespace          string        `mapstructure:"namespace"`
	MaxRetries         int           `mapstructure:"max_retries"`
	Timeout            time.Duration `mapstructure:"timeout"`
	RetryDelay         time.Duration `mapstructure:"retry_delay"`
	KVVersion          int           `mapstructure:"kv_version"`
	TLSCertPath        string        `mapstructure:"tls_cert_path"`
	TLSKeyPath         string        `mapstructure:"tls_key_path"`
	CACertPath         string        `mapstructure:"ca_cert_path"`
	InsecureSkipVerify bool          `mapstructure:"insecure_skip_verify"`
	K8sRole            string        `mapstructure:"k8s_role"`
	K8sTokenPath       string        `mapstructure:"k8s_token_path"`
	AppRoleID          string        `mapstructure:"app_role_id"`
	AppSecretID        string        `mapstructure:"app_secret_id"`
}

// Client 是 Vault 客户端接口。
//
// 提供统一的 Vault 操作接口，包括：
// - 密钥管理（KV v1/v2）
// - 动态密钥（数据库凭据等）
// - 加密/解密操作
// - 密钥轮换
// - 健康检查
//
// 所有方法都支持上下文控制，可以设置超时和取消操作。
type Client interface {
	// KV 返回 KV 密钥管理器（支持 v1 和 v2）
	KV() KVManager
	// DB 返回动态数据库凭据管理器
	DB() DBManager
	// Encryption 返回加密服务
	Encryption() EncryptionService
	// Rotation 返回密钥轮换服务
	Rotation() RotationService
	// Health 检查 Vault 连接状态
	Health(ctx context.Context) (*HealthStatus, error)
	// Close 关闭客户端连接
	Close() error
	// IsConnected 检查客户端是否已连接
	IsConnected() bool
}

// HealthStatus 表示 Vault 健康状态
type HealthStatus struct {
	Initialized                bool      `json:"initialized"`
	Sealed                     bool      `json:"sealed"`
	Standby                    bool      `json:"standby"`
	PerformanceStandby         bool      `json:"performance_standby,omitempty"`
	ReplicationPerformanceMode string    `json:"replication_performance_mode,omitempty"`
	ReplicationDRMode          string    `json:"replication_dr_mode,omitempty"`
	ServerTimeUTC              int64     `json:"server_time_utc"`
	Version                    string    `json:"version"`
	ClusterName                string    `json:"cluster_name,omitempty"`
	ClusterID                  string    `json:"cluster_id,omitempty"`
	CheckedAt                  time.Time `json:"checked_at"`
}

// vaultClient 是 Client 接口的实现
type vaultClient struct {
	client     *api.Client
	config     *Config
	mu         sync.RWMutex
	connected  bool
	kv         KVManager
	db         DBManager
	encryption EncryptionService
	rotation   RotationService
}

// NewClient 创建新的 Vault 客户端。
//
// 参数：
//   - config: Vault 配置，必须包含 Address 和认证信息
//
// 返回：
//   - Client: Vault 客户端实例
//   - error: 创建失败时返回错误
//
// 示例：
//
//	cfg := &vault.Config{
//	    Address:    "https://vault.example.com:8200",
//	    AuthMethod: vault.AuthMethodToken,
//	    Token:      os.Getenv("VAULT_TOKEN"),
//	}
//	client, err := vault.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClient(config *Config) (Client, error) {
	if config == nil {
		return nil, fmt.Errorf("vault config is required")
	}

	// 设置默认值
	if err := setDefaultConfig(config); err != nil {
		return nil, fmt.Errorf("failed to set default config: %w", err)
	}

	// 创建 Vault API 客户端配置
	apiConfig := api.DefaultConfig()
	apiConfig.Address = config.Address
	apiConfig.Timeout = config.Timeout

	// 配置 TLS
	if err := configureTLS(apiConfig, config); err != nil {
		return nil, fmt.Errorf("failed to configure TLS: %w", err)
	}

	// 创建 Vault 客户端
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// 设置命名空间（企业版功能）
	if config.Namespace != "" {
		client.SetNamespace(config.Namespace)
	}

	// 执行认证
	if err := authenticate(client, config); err != nil {
		return nil, fmt.Errorf("failed to authenticate with vault: %w", err)
	}

	vc := &vaultClient{
		client:    client,
		config:    config,
		connected: true,
	}

	// 初始化子模块
	vc.kv = newKVManager(client, config.KVVersion)
	vc.db = newDBManager(client)
	vc.encryption = newEncryptionService(client)
	vc.rotation = newRotationService(client, vc.kv)

	return vc, nil
}

// setDefaultConfig 设置配置默认值
func setDefaultConfig(config *Config) error {
	if config.Address == "" {
		config.Address = "http://127.0.0.1:8200"
	}

	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	if config.KVVersion == 0 {
		config.KVVersion = 2
	}

	if config.K8sTokenPath == "" {
		config.K8sTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	return nil
}

// configureTLS 配置 TLS
func configureTLS(apiConfig *api.Config, config *Config) error {
	tlsConfig := &api.TLSConfig{}

	if config.CACertPath != "" {
		tlsConfig.CACert = config.CACertPath
	}

	if config.TLSCertPath != "" && config.TLSKeyPath != "" {
		tlsConfig.ClientCert = config.TLSCertPath
		tlsConfig.ClientKey = config.TLSKeyPath
	}

	tlsConfig.Insecure = config.InsecureSkipVerify

	if err := apiConfig.ConfigureTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to configure TLS: %w", err)
	}

	return nil
}

// authenticate 执行 Vault 认证
func authenticate(client *api.Client, config *Config) error {
	switch config.AuthMethod {
	case AuthMethodToken:
		if config.Token == "" {
			return fmt.Errorf("token is required for token authentication")
		}
		client.SetToken(config.Token)

	case AuthMethodKubernetes:
		if config.K8sRole == "" {
			return fmt.Errorf("k8s_role is required for kubernetes authentication")
		}
		if err := authenticateWithKubernetes(client, config); err != nil {
			return err
		}

	case AuthMethodAppRole:
		if config.AppRoleID == "" || config.AppSecretID == "" {
			return fmt.Errorf("app_role_id and app_secret_id are required for approle authentication")
		}
		if err := authenticateWithAppRole(client, config); err != nil {
			return err
		}

	case AuthMethodCert:
		if config.TLSCertPath == "" || config.TLSKeyPath == "" {
			return fmt.Errorf("tls_cert_path and tls_key_path are required for cert authentication")
		}
		// TLS 证书认证在 configureTLS 中已经配置

	default:
		return fmt.Errorf("unsupported authentication method: %s", config.AuthMethod)
	}

	return nil
}

// authenticateWithKubernetes 使用 Kubernetes 认证
func authenticateWithKubernetes(client *api.Client, config *Config) error {
	k8sAuth, err := kubernetes.NewKubernetesAuth(
		config.K8sRole,
		kubernetes.WithServiceAccountTokenPath(config.K8sTokenPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes auth: %w", err)
	}

	_, err = client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		return fmt.Errorf("failed to login with kubernetes auth: %w", err)
	}

	return nil
}

// authenticateWithAppRole 使用 AppRole 认证
func authenticateWithAppRole(client *api.Client, config *Config) error {
	appRoleAuth, err := approle.NewAppRoleAuth(
		config.AppRoleID,
		&approle.SecretID{FromString: config.AppSecretID},
	)
	if err != nil {
		return fmt.Errorf("failed to create approle auth: %w", err)
	}

	_, err = client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return fmt.Errorf("failed to login with approle auth: %w", err)
	}

	return nil
}

// KV 返回 KV 密钥管理器
func (c *vaultClient) KV() KVManager {
	return c.kv
}

// DB 返回动态数据库凭据管理器
func (c *vaultClient) DB() DBManager {
	return c.db
}

// Encryption 返回加密服务
func (c *vaultClient) Encryption() EncryptionService {
	return c.encryption
}

// Rotation 返回密钥轮换服务
func (c *vaultClient) Rotation() RotationService {
	return c.rotation
}

// Health 检查 Vault 健康状态
func (c *vaultClient) Health(ctx context.Context) (*HealthStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("vault client is not connected")
	}

	health, err := c.client.Sys().Health()
	if err != nil {
		return nil, fmt.Errorf("failed to check vault health: %w", err)
	}

	return &HealthStatus{
		Initialized:                health.Initialized,
		Sealed:                     health.Sealed,
		Standby:                    health.Standby,
		PerformanceStandby:         health.PerformanceStandby,
		ReplicationPerformanceMode: health.ReplicationPerformanceMode,
		ReplicationDRMode:          health.ReplicationDRMode,
		ServerTimeUTC:              health.ServerTimeUTC,
		Version:                    health.Version,
		ClusterName:                health.ClusterName,
		ClusterID:                  health.ClusterID,
		CheckedAt:                  time.Now(),
	}, nil
}

// Close 关闭客户端连接
func (c *vaultClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected = false
	// Vault API 客户端没有显式关闭方法
	// 这里可以清理相关资源
	return nil
}

// IsConnected 检查客户端是否已连接
func (c *vaultClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// withRetry 执行带重试的操作
func withRetry[T any](ctx context.Context, maxRetries int, retryDelay time.Duration, operation func() (T, error)) (T, error) {
	var result T
	var err error

	for i := 0; i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		result, err = operation()
		if err == nil {
			return result, nil
		}

		if i < maxRetries {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(retryDelay):
				// 继续重试
			}
		}
	}

	return result, fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}

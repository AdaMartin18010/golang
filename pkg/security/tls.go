package security

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	// ErrInvalidTLSConfig 无效的 TLS 配置
	ErrInvalidTLSConfig = errors.New("invalid TLS configuration")
)

// TLSManager TLS 管理器
type TLSManager struct {
	config *tls.Config
}

// NewTLSManager 创建 TLS 管理器
func NewTLSManager(config TLSConfig) (*TLSManager, error) {
	if !config.Enabled {
		return &TLSManager{config: nil}, nil
	}

	if config.CertFile == "" || config.KeyFile == "" {
		return nil, fmt.Errorf("%w: certificate and key files are required", ErrInvalidTLSConfig)
	}

	// 加载证书和密钥
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   config.MinVersion,
		MaxVersion:   config.MaxVersion,
		CipherSuites: config.CipherSuites,
	}

	// 设置默认值
	if tlsConfig.MinVersion == 0 {
		tlsConfig.MinVersion = tls.VersionTLS12
	}
	if tlsConfig.MaxVersion == 0 {
		tlsConfig.MaxVersion = tls.VersionTLS13
	}

	// 配置密码套件（如果未指定，使用安全默认值）
	if len(tlsConfig.CipherSuites) == 0 {
		tlsConfig.CipherSuites = defaultCipherSuites()
	}

	// 安全配置
	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.InsecureSkipVerify = config.InsecureSkipVerify

	return &TLSManager{config: tlsConfig}, nil
}

// GetConfig 获取 TLS 配置
func (m *TLSManager) GetConfig() *tls.Config {
	if m.config == nil {
		return nil
	}

	// 返回配置的副本，避免并发修改
	config := m.config.Clone()
	return config
}

// LoadCA 加载 CA 证书
func (m *TLSManager) LoadCA(caFile string) error {
	if caFile == "" {
		return nil
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	if m.config == nil {
		m.config = &tls.Config{}
	}

	m.config.RootCAs = caCertPool
	return nil
}

// defaultCipherSuites 返回默认的安全密码套件
func defaultCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
}

// ValidateCertificate 验证证书
func (m *TLSManager) ValidateCertificate(certFile, keyFile string) error {
	_, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("invalid certificate: %w", err)
	}

	// 读取证书文件
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	// 解析证书
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// 检查证书是否过期
	now := time.Now()
	if now.After(cert.NotAfter) {
		return fmt.Errorf("certificate expired on %v", cert.NotAfter)
	}

	if now.Before(cert.NotBefore) {
		return fmt.Errorf("certificate not valid until %v", cert.NotBefore)
	}

	return nil
}

package security

import (
	"crypto/tls"
	"testing"
)

func TestTLSManager_NewTLSManager(t *testing.T) {
	// 测试禁用 TLS
	config := TLSConfig{
		Enabled: false,
	}

	manager, err := NewTLSManager(config)
	if err != nil {
		t.Fatalf("Failed to create TLS manager: %v", err)
	}

	if manager.config != nil {
		t.Error("TLS config should be nil when disabled")
	}
}

func TestTLSManager_InvalidConfig(t *testing.T) {
	config := TLSConfig{
		Enabled:  true,
		CertFile: "",
		KeyFile:  "",
	}

	_, err := NewTLSManager(config)
	if err == nil {
		t.Error("Should return error for invalid config")
	}
}

func TestDefaultCipherSuites(t *testing.T) {
	suites := defaultCipherSuites()
	if len(suites) == 0 {
		t.Error("Default cipher suites should not be empty")
	}

	// 验证所有套件都是安全的
	for _, suite := range suites {
		if suite == 0 {
			t.Error("Cipher suite should not be zero")
		}
	}
}

func TestTLSConfig_Validate(t *testing.T) {
	config := DefaultSecurityConfig()

	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
}

func TestTLSConfig_InvalidTLS(t *testing.T) {
	config := DefaultSecurityConfig()
	config.TLS.Enabled = true
	config.TLS.CertFile = ""
	config.TLS.KeyFile = ""

	if err := config.Validate(); err == nil {
		t.Error("TLS config without cert/key should be invalid")
	}
}

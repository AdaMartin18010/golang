package security

import (
	"crypto/tls"
	"testing"
	"time"
)

func TestSecurityConfig_Validate(t *testing.T) {
	config := DefaultSecurityConfig()

	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
}

func TestSecurityConfig_ValidateInvalid(t *testing.T) {
	config := DefaultSecurityConfig()
	config.Encryption.Algorithm = ""

	if err := config.Validate(); err == nil {
		t.Error("Config with empty algorithm should be invalid")
	}
}

func TestSecurityConfig_ValidateInvalidKeySize(t *testing.T) {
	config := DefaultSecurityConfig()
	config.Encryption.KeySize = 512

	if err := config.Validate(); err == nil {
		t.Error("Config with invalid key size should be invalid")
	}
}

func TestSecurityConfig_ValidateTLS(t *testing.T) {
	config := DefaultSecurityConfig()
	config.TLS.Enabled = true
	config.TLS.CertFile = ""
	config.TLS.KeyFile = ""

	if err := config.Validate(); err == nil {
		t.Error("TLS config without cert/key files should be invalid")
	}
}

func TestSecurityConfigManager(t *testing.T) {
	config := DefaultSecurityConfig()
	manager, err := NewSecurityConfigManager(config)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	retrieved := manager.GetConfig()
	if retrieved.Encryption.Algorithm != config.Encryption.Algorithm {
		t.Error("Retrieved config should match original")
	}
}

func TestSecurityConfigManager_UpdateConfig(t *testing.T) {
	config := DefaultSecurityConfig()
	manager, err := NewSecurityConfigManager(config)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	newConfig := DefaultSecurityConfig()
	newConfig.Encryption.KeySize = 192

	err = manager.UpdateConfig(newConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	updated := manager.GetConfig()
	if updated.Encryption.KeySize != 192 {
		t.Error("Config should be updated")
	}
}

func TestSecurityConfigManager_GetEncryptionConfig(t *testing.T) {
	config := DefaultSecurityConfig()
	manager, _ := NewSecurityConfigManager(config)

	encConfig := manager.GetEncryptionConfig()
	if encConfig.Algorithm != config.Encryption.Algorithm {
		t.Error("Should return correct encryption config")
	}
}

func TestDefaultSecurityConfig(t *testing.T) {
	config := DefaultSecurityConfig()

	if config.Encryption.Algorithm != "AES-256-GCM" {
		t.Error("Default encryption algorithm should be AES-256-GCM")
	}

	if config.TLS.MinVersion != tls.VersionTLS12 {
		t.Error("Default TLS min version should be TLS 1.2")
	}

	if config.Audit.Enabled != true {
		t.Error("Audit should be enabled by default")
	}
}

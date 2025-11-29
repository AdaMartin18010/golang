package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/golang/pkg/security"
)

// TestSecurityMiddlewareIntegration 测试安全中间件集成
func TestSecurityMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建安全中间件
	smConfig := security.SecurityMiddlewareConfig{
		SecurityHeaders: func() *security.SecurityHeadersConfig {
			cfg := security.DefaultSecurityHeadersConfig()
			return &cfg
		}(),
		RateLimit: &security.RateLimiterConfig{
			Limit:  10,
			Window: 1 * time.Minute,
		},
		CSRF: func() *security.CSRFConfig {
			cfg := security.DefaultCSRFConfig()
			return &cfg
		}(),
		EnableXSS: true,
	}

	middleware := security.NewSecurityMiddleware(smConfig)
	defer middleware.Shutdown()

	// 创建测试处理器
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// 测试请求
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// 验证安全头部
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header should be set")
	}

	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("X-Content-Type-Options header should be set")
	}
}

// TestRateLimitIntegration 测试速率限制集成
func TestRateLimitIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建速率限制器
	limiter := security.NewIPRateLimiter(security.RateLimiterConfig{
		Limit:  2,
		Window: 1 * time.Minute,
	})
	defer limiter.Shutdown(context.Background())

	ctx := context.Background()
	ip := "192.168.1.1"

	// 前两次请求应该成功
	for i := 0; i < 2; i++ {
		allowed, err := limiter.AllowIP(ctx, ip)
		if err != nil {
			t.Fatalf("Failed to allow request: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 第三次请求应该被限制
	allowed, err := limiter.AllowIP(ctx, ip)
	if err == nil {
		t.Error("Should return error for rate limit exceeded")
	}
	if allowed {
		t.Error("Request should be denied")
	}
}

// TestCSRFProtectionIntegration 测试 CSRF 防护集成
func TestCSRFProtectionIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建 CSRF 防护
	csrf := security.NewCSRFProtection(security.DefaultCSRFConfig())
	defer csrf.Shutdown()

	sessionID := "test-session-123"
	token, err := csrf.GenerateToken(sessionID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 验证有效令牌
	if err := csrf.ValidateToken(sessionID, token); err != nil {
		t.Errorf("Valid token should pass validation: %v", err)
	}

	// 验证无效令牌
	if err := csrf.ValidateToken(sessionID, "invalid-token"); err == nil {
		t.Error("Invalid token should fail validation")
	}
}

// TestPasswordHashingIntegration 测试密码哈希集成
func TestPasswordHashingIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建密码哈希器
	hasher := security.NewPasswordHasher(security.DefaultPasswordHashConfig())

	password := "test-password-123"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 验证正确密码
	valid, err := hasher.Verify(password, hash)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if !valid {
		t.Error("Password verification should succeed")
	}

	// 验证错误密码
	valid, err = hasher.Verify("wrong-password", hash)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if valid {
		t.Error("Password verification should fail for wrong password")
	}
}

// TestEncryptionIntegration 测试加密集成
func TestEncryptionIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建加密器
	encryptor, err := security.NewAES256EncryptorFromString("test-secret-key-12345")
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := "sensitive data"
	ciphertext, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := encryptor.DecryptString(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected '%s', got '%s'", plaintext, decrypted)
	}
}


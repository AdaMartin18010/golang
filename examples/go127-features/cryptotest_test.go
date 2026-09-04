package main

import (
	"crypto/rand"
	"testing"
	"testing/cryptotest"
)

// TestCryptoTestSetGlobalRandom 演示 testing/cryptotest.SetGlobalRandom（Go 1.27）：
// 为单个测试设置全局确定性的加密随机源，影响 crypto/rand 及 crypto/... 包的
// 隐式随机源，使涉及随机性的测试可重现。
func TestCryptoTestSetGlobalRandom(t *testing.T) {
	cryptotest.SetGlobalRandom(t, 42)

	a := make([]byte, 16)
	if _, err := rand.Read(a); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	cryptotest.SetGlobalRandom(t, 42) // 重置后序列应重现
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}

	if string(a) != string(b) {
		t.Errorf("same seed should reproduce the same randomness: %x vs %x", a, b)
	}
	t.Logf("deterministic random bytes: %x", a)
}

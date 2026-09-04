package main

import (
	"bytes"
	"crypto/mldsa"
	"crypto/rand"
	"fmt"
	"hash/maphash"
	"math/big"
)

// ============================================================================
// Feature 9: crypto/mldsa — 后量子签名（FIPS 204，Go 1.27）
// ============================================================================

// demonstrateMLDSA shows post-quantum ML-DSA signatures.
// 注意：若使用 FIPS 140-3 Go Cryptographic Module v1.0.0，本包不可用（返回错误）；
// v1.26.0+ 的 FIPS module 可用。
func demonstrateMLDSA() {
	fmt.Println("=== Feature 9: crypto/mldsa (FIPS 204 post-quantum signatures) ===")

	priv, err := mldsa.GenerateKey(mldsa.MLDSA44())
	if err != nil {
		fmt.Printf("GenerateKey error (FIPS module 模式下预期): %v\n", err)
		fmt.Println()
		return
	}

	msg := []byte("message signed with ML-DSA-44")
	sig, err := priv.Sign(rand.Reader, msg, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MLDSA44 signature size: %d bytes\n", len(sig))

	pub := priv.PublicKey()
	if err := mldsa.Verify(pub, msg, sig, nil); err != nil {
		panic(err)
	}
	fmt.Println("Verify: signature valid")

	tampered := bytes.Clone(msg)
	tampered[0] ^= 0xff
	if err := mldsa.Verify(pub, tampered, sig, nil); err != nil {
		fmt.Printf("Verify tampered: rejected (%v)\n", err)
	}
	fmt.Println()
}

// ============================================================================
// Feature 10: math/big.Int.Divide — 四种舍入模式的整数除法（Go 1.27）
// ============================================================================

// demonstrateBigDivide shows Euclidean division with explicit rounding modes.
func demonstrateBigDivide() {
	fmt.Println("=== Feature 10: math/big.Int.Divide ===")

	x, y := big.NewInt(-7), big.NewInt(2)
	for _, mode := range []struct {
		name string
		m    big.RoundingMode
	}{
		{"Trunc", big.Trunc}, // 向零截断: q = -3, r = -1
		{"Floor", big.Floor}, // 向下取整: q = -4, r =  1（欧几里得余数恒非负）
		{"Round", big.Round}, // 四舍五入: q = -4, r =  1
		{"Ceil", big.Ceil},   // 向上取整: q = -3, r = -1
	} {
		q, r := new(big.Int).Divide(x, y, new(big.Int), mode.m)
		fmt.Printf("Divide(-7, 2, %s) = q:%d r:%d\n", mode.name, q, r)
	}
	fmt.Println()
}

// ============================================================================
// Feature 11: hash/maphash.Hasher[T] — 类型化哈希器（Go 1.27）
// ============================================================================

// demonstrateMaphash shows the new generic Hasher interfaces.
// ComparableHasher[T] 的 Equal 与 == 一致；Hasher[T] 允许不可比较类型做哈希键。
func demonstrateMaphash() {
	fmt.Println("=== Feature 11: hash/maphash.Hasher[T] / ComparableHasher[T] ===")

	var hasher maphash.ComparableHasher[string] // 零值可用
	var h1, h2 maphash.Hash

	// 两个 Hash 需使用相同 Seed 才能比较（零值 Hash 各自持有随机 Seed）
	seed := maphash.MakeSeed()
	h1.SetSeed(seed)
	h2.SetSeed(seed)

	hasher.Hash(&h1, "go-1.27")
	hasher.Hash(&h2, "go-1.27")
	fmt.Printf("same key: Sum64 equal = %t, hasher.Equal = %t\n", h1.Sum64() == h2.Sum64(), hasher.Equal("go-1.27", "go-1.27"))

	hasher.Hash(&h2, "other")
	fmt.Printf("diff key: Sum64 equal = %t\n", h1.Sum64() == h2.Sum64())
	fmt.Println()
}

// ============================================================================
// Feature 12: simd / archsimd — 实验性 SIMD（Go 1.27，需 GOEXPERIMENT=simd）
// ============================================================================

// demonstrateSIMD 说明实验性 SIMD 包的使用方式。
// simd/archsimd 需要 GOEXPERIMENT=simd 编译，默认工具链下包不可见，
// 因此这里仅打印说明；启用方式：
//
//	GOEXPERIMENT=simd go run .
//
// 1.27 提供 amd64（SSE/AVX 修订）、arm64（Neon）、wasm（128-bit）后端；
// 1.27.1 修复了 arm64 的 SIGILL（#81110）与 stub 参数问题（#81109）。
func demonstrateSIMD() {
	fmt.Println("=== Feature 12: simd/archsimd (experimental, GOEXPERIMENT=simd) ===")
	fmt.Println("需以 GOEXPERIMENT=simd 编译；amd64/arm64(Neon)/wasm 后端。")
	fmt.Println("详细 API 见 go doc simd/archsimd（实验包无稳定性保证）。")
	fmt.Println()
}

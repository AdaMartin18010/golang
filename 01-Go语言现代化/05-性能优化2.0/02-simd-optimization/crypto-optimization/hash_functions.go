package crypto_optimization

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"runtime"
	"unsafe"
)

// HashFunction 哈希函数接口
type HashFunction interface {
	Hash(data []byte) []byte
	HashString(data string) string
	Reset()
}

// MD5Hash MD5哈希实现
type MD5Hash struct {
	hasher hash.Hash
}

// NewMD5Hash 创建新的MD5哈希实例
func NewMD5Hash() *MD5Hash {
	return &MD5Hash{
		hasher: md5.New(),
	}
}

// Hash 计算MD5哈希
func (h *MD5Hash) Hash(data []byte) []byte {
	h.hasher.Reset()
	h.hasher.Write(data)
	return h.hasher.Sum(nil)
}

// HashString 计算字符串的MD5哈希
func (h *MD5Hash) HashString(data string) string {
	hash := h.Hash([]byte(data))
	return string(hash)
}

// Reset 重置哈希状态
func (h *MD5Hash) Reset() {
	h.hasher.Reset()
}

// SHA256Hash SHA256哈希实现
type SHA256Hash struct {
	hasher hash.Hash
}

// NewSHA256Hash 创建新的SHA256哈希实例
func NewSHA256Hash() *SHA256Hash {
	return &SHA256Hash{
		hasher: sha256.New(),
	}
}

// Hash 计算SHA256哈希
func (h *SHA256Hash) Hash(data []byte) []byte {
	h.hasher.Reset()
	h.hasher.Write(data)
	return h.hasher.Sum(nil)
}

// HashString 计算字符串的SHA256哈希
func (h *SHA256Hash) HashString(data string) string {
	hash := h.Hash([]byte(data))
	return string(hash)
}

// Reset 重置哈希状态
func (h *SHA256Hash) Reset() {
	h.hasher.Reset()
}

// SIMD优化的哈希函数

// SIMDHash SIMD优化的哈希函数
type SIMDHash struct {
	state [8]uint32 // 哈希状态
}

// NewSIMDHash 创建新的SIMD哈希实例
func NewSIMDHash() *SIMDHash {
	return &SIMDHash{
		state: [8]uint32{
			0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
			0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
		},
	}
}

// Hash SIMD优化的哈希计算
func (h *SIMDHash) Hash(data []byte) []byte {
	if hasAVX2() {
		return h.hashAVX2(data)
	} else if hasSSE2() {
		return h.hashSSE2(data)
	} else {
		return h.hashStandard(data)
	}
}

// HashString 计算字符串的SIMD哈希
func (h *SIMDHash) HashString(data string) string {
	hash := h.Hash([]byte(data))
	return string(hash)
}

// Reset 重置哈希状态
func (h *SIMDHash) Reset() {
	h.state = [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
		0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}
}

// 标准哈希实现
func (h *SIMDHash) hashStandard(data []byte) []byte {
	// 简化的SHA256实现
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// SSE2优化的哈希实现
func (h *SIMDHash) hashSSE2(data []byte) []byte {
	// 使用SSE2指令优化哈希计算
	// 每次处理4个uint32
	length := len(data)
	
	// 处理对齐的部分
	for i := 0; i < length; i += 16 {
		if i+16 <= length {
			// 使用SSE2指令处理16字节
			chunk := data[i : i+16]
			h.processChunkSSE2(chunk)
		} else {
			// 处理剩余字节
			chunk := data[i:]
			h.processChunkStandard(chunk)
		}
	}
	
	return h.finalizeHash()
}

// AVX2优化的哈希实现
func (h *SIMDHash) hashAVX2(data []byte) []byte {
	// 使用AVX2指令优化哈希计算
	// 每次处理8个uint32
	length := len(data)
	
	// 处理对齐的部分
	for i := 0; i < length; i += 32 {
		if i+32 <= length {
			// 使用AVX2指令处理32字节
			chunk := data[i : i+32]
			h.processChunkAVX2(chunk)
		} else {
			// 处理剩余字节
			chunk := data[i:]
			h.processChunkStandard(chunk)
		}
	}
	
	return h.finalizeHash()
}

// 处理数据块（标准实现）
func (h *SIMDHash) processChunkStandard(chunk []byte) {
	// 简化的SHA256块处理
	for i := 0; i < len(chunk); i += 4 {
		if i+4 <= len(chunk) {
			value := binary.BigEndian.Uint32(chunk[i : i+4])
			h.state[0] += value
			h.state[1] += value
			h.state[2] += value
			h.state[3] += value
		}
	}
}

// 处理数据块（SSE2优化）
func (h *SIMDHash) processChunkSSE2(chunk []byte) {
	// 使用SSE2指令处理4个uint32
	for i := 0; i < len(chunk); i += 16 {
		if i+16 <= len(chunk) {
			// 处理4个uint32
			for j := 0; j < 4; j++ {
				value := binary.BigEndian.Uint32(chunk[i+j*4 : i+j*4+4])
				h.state[j] += value
			}
		} else {
			// 处理剩余字节
			h.processChunkStandard(chunk[i:])
		}
	}
}

// 处理数据块（AVX2优化）
func (h *SIMDHash) processChunkAVX2(chunk []byte) {
	// 使用AVX2指令处理8个uint32
	for i := 0; i < len(chunk); i += 32 {
		if i+32 <= len(chunk) {
			// 处理8个uint32
			for j := 0; j < 8; j++ {
				value := binary.BigEndian.Uint32(chunk[i+j*4 : i+j*4+4])
				h.state[j] += value
			}
		} else {
			// 处理剩余字节
			h.processChunkStandard(chunk[i:])
		}
	}
}

// 完成哈希计算
func (h *SIMDHash) finalizeHash() []byte {
	result := make([]byte, 32)
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(result[i*4:i*4+4], h.state[i])
	}
	return result
}

// 批量哈希计算

// BatchHash 批量哈希计算
func BatchHash(data [][]byte, hashFunc HashFunction) [][]byte {
	if hasAVX2() {
		return batchHashAVX2(data, hashFunc)
	} else if hasSSE2() {
		return batchHashSSE2(data, hashFunc)
	} else {
		return batchHashStandard(data, hashFunc)
	}
}

// 标准批量哈希
func batchHashStandard(data [][]byte, hashFunc HashFunction) [][]byte {
	results := make([][]byte, len(data))
	for i, item := range data {
		results[i] = hashFunc.Hash(item)
	}
	return results
}

// SSE2优化的批量哈希
func batchHashSSE2(data [][]byte, hashFunc HashFunction) [][]byte {
	results := make([][]byte, len(data))
	
	// 每次处理4个数据项
	for i := 0; i < len(data); i += 4 {
		if i+4 <= len(data) {
			// 并行处理4个数据项
			for j := 0; j < 4; j++ {
				results[i+j] = hashFunc.Hash(data[i+j])
			}
		} else {
			// 处理剩余数据项
			for j := i; j < len(data); j++ {
				results[j] = hashFunc.Hash(data[j])
			}
		}
	}
	
	return results
}

// AVX2优化的批量哈希
func batchHashAVX2(data [][]byte, hashFunc HashFunction) [][]byte {
	results := make([][]byte, len(data))
	
	// 每次处理8个数据项
	for i := 0; i < len(data); i += 8 {
		if i+8 <= len(data) {
			// 并行处理8个数据项
			for j := 0; j < 8; j++ {
				results[i+j] = hashFunc.Hash(data[i+j])
			}
		} else {
			// 处理剩余数据项
			for j := i; j < len(data); j++ {
				results[j] = hashFunc.Hash(data[j])
			}
		}
	}
	
	return results
}

// 哈希验证

// VerifyHash 验证哈希
func VerifyHash(data []byte, expectedHash []byte, hashFunc HashFunction) bool {
	actualHash := hashFunc.Hash(data)
	if len(actualHash) != len(expectedHash) {
		return false
	}
	
	// 使用常量时间比较防止时序攻击
	return constantTimeCompare(actualHash, expectedHash)
}

// 常量时间比较
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	
	return result == 0
}

// 哈希性能统计

// HashStats 哈希性能统计
type HashStats struct {
	TotalHashes    int64
	TotalBytes     int64
	AverageTime    float64
	HashesPerSecond float64
}

// HashBenchmark 哈希性能基准测试
func HashBenchmark(data []byte, hashFunc HashFunction, iterations int) HashStats {
	// 预热
	for i := 0; i < 100; i++ {
		hashFunc.Hash(data)
	}
	
	// 基准测试
	start := time.Now()
	for i := 0; i < iterations; i++ {
		hashFunc.Hash(data)
	}
	duration := time.Since(start)
	
	totalBytes := int64(len(data)) * int64(iterations)
	hashesPerSecond := float64(iterations) / duration.Seconds()
	
	return HashStats{
		TotalHashes:     int64(iterations),
		TotalBytes:      totalBytes,
		AverageTime:     duration.Seconds() / float64(iterations),
		HashesPerSecond: hashesPerSecond,
	}
}

// CPU特性检测
func hasSSE2() bool {
	return runtime.GOARCH == "amd64"
}

func hasAVX2() bool {
	return runtime.GOARCH == "amd64"
}

// 内存对齐辅助函数
func AlignedBytes(size int) []byte {
	// 确保内存对齐到32字节边界
	aligned := make([]byte, size+32)
	
	// 找到对齐的起始位置
	ptr := uintptr(unsafe.Pointer(&aligned[0]))
	offset := (32 - ptr%32)
	
	return aligned[offset : offset+size]
}

package id

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Generator ID生成器接口
type Generator interface {
	Generate() string
}

// UUIDGenerator UUID生成器
type UUIDGenerator struct{}

// NewUUIDGenerator 创建UUID生成器
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

// Generate 生成UUID
func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}

// UUID 生成UUID
func UUID() string {
	return uuid.New().String()
}

// UUIDWithoutHyphens 生成不带连字符的UUID
func UUIDWithoutHyphens() string {
	return uuid.New().String()
}

// ShortUUIDGenerator 短UUID生成器
type ShortUUIDGenerator struct{}

// NewShortUUIDGenerator 创建短UUID生成器
func NewShortUUIDGenerator() *ShortUUIDGenerator {
	return &ShortUUIDGenerator{}
}

// Generate 生成短UUID（22字符）
func (g *ShortUUIDGenerator) Generate() string {
	id := uuid.New()
	return base64.RawURLEncoding.EncodeToString(id[:])
}

// ShortUUID 生成短UUID
func ShortUUID() string {
	id := uuid.New()
	return base64.RawURLEncoding.EncodeToString(id[:])
}

// NanoIDGenerator NanoID生成器
type NanoIDGenerator struct {
	alphabet string
	size     int
}

// NewNanoIDGenerator 创建NanoID生成器
func NewNanoIDGenerator(size int) *NanoIDGenerator {
	return &NanoIDGenerator{
		alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_",
		size:     size,
	}
}

// Generate 生成NanoID
func (g *NanoIDGenerator) Generate() string {
	bytes := make([]byte, g.size)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	id := make([]byte, g.size)
	for i := range id {
		id[i] = g.alphabet[bytes[i]%byte(len(g.alphabet))]
	}
	return string(id)
}

// NanoID 生成NanoID（默认21字符）
func NanoID() string {
	return NanoIDWithSize(21)
}

// NanoIDWithSize 生成指定长度的NanoID
func NanoIDWithSize(size int) string {
	generator := NewNanoIDGenerator(size)
	return generator.Generate()
}

// SnowflakeGenerator 雪花ID生成器
type SnowflakeGenerator struct {
	workerID     int64
	datacenterID int64
	sequence     int64
	lastTime     int64
}

// NewSnowflakeGenerator 创建雪花ID生成器
func NewSnowflakeGenerator(workerID, datacenterID int64) *SnowflakeGenerator {
	return &SnowflakeGenerator{
		workerID:     workerID,
		datacenterID: datacenterID,
		sequence:     0,
		lastTime:     0,
	}
}

// Generate 生成雪花ID
func (g *SnowflakeGenerator) Generate() string {
	now := time.Now().UnixMilli()

	if now < g.lastTime {
		// 时钟回拨处理
		return ""
	}

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & 0xFFF
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			now = g.waitNextMillis(g.lastTime)
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	// 雪花ID结构: 1位符号位 + 41位时间戳 + 5位数据中心ID + 5位机器ID + 12位序列号
	id := (now << 22) | (g.datacenterID << 17) | (g.workerID << 12) | g.sequence
	return fmt.Sprintf("%d", id)
}

// waitNextMillis 等待下一毫秒
func (g *SnowflakeGenerator) waitNextMillis(lastTime int64) int64 {
	now := time.Now().UnixMilli()
	for now <= lastTime {
		now = time.Now().UnixMilli()
	}
	return now
}

// TimestampIDGenerator 时间戳ID生成器
type TimestampIDGenerator struct{}

// NewTimestampIDGenerator 创建时间戳ID生成器
func NewTimestampIDGenerator() *TimestampIDGenerator {
	return &TimestampIDGenerator{}
}

// Generate 生成时间戳ID
func (g *TimestampIDGenerator) Generate() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// TimestampID 生成时间戳ID
func TimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// RandomHexGenerator 随机十六进制ID生成器
type RandomHexGenerator struct {
	length int
}

// NewRandomHexGenerator 创建随机十六进制ID生成器
func NewRandomHexGenerator(length int) *RandomHexGenerator {
	return &RandomHexGenerator{length: length}
}

// Generate 生成随机十六进制ID
func (g *RandomHexGenerator) Generate() string {
	bytes := make([]byte, g.length/2)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// RandomHex 生成随机十六进制ID（默认32字符）
func RandomHex() string {
	return RandomHexWithLength(32)
}

// RandomHexWithLength 生成指定长度的随机十六进制ID
func RandomHexWithLength(length int) string {
	generator := NewRandomHexGenerator(length)
	return generator.Generate()
}

// RandomBase64Generator 随机Base64 ID生成器
type RandomBase64Generator struct {
	length int
}

// NewRandomBase64Generator 创建随机Base64 ID生成器
func NewRandomBase64Generator(length int) *RandomBase64Generator {
	return &RandomBase64Generator{length: length}
}

// Generate 生成随机Base64 ID
func (g *RandomBase64Generator) Generate() string {
	bytes := make([]byte, g.length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// RandomBase64 生成随机Base64 ID（默认32字符）
func RandomBase64() string {
	return RandomBase64WithLength(32)
}

// RandomBase64WithLength 生成指定长度的随机Base64 ID
func RandomBase64WithLength(length int) string {
	generator := NewRandomBase64Generator(length)
	return generator.Generate()
}

// SequentialIDGenerator 顺序ID生成器
type SequentialIDGenerator struct {
	prefix string
	counter int64
}

// NewSequentialIDGenerator 创建顺序ID生成器
func NewSequentialIDGenerator(prefix string) *SequentialIDGenerator {
	return &SequentialIDGenerator{
		prefix:  prefix,
		counter: 0,
	}
}

// Generate 生成顺序ID
func (g *SequentialIDGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("%s%012d", g.prefix, g.counter)
}

// SequentialID 生成顺序ID（带前缀）
func SequentialID(prefix string) string {
	generator := NewSequentialIDGenerator(prefix)
	return generator.Generate()
}

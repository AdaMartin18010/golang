package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"
	"os"
)

// MD5 MD5哈希
func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// MD5String MD5哈希字符串
func MD5String(s string) string {
	return MD5([]byte(s))
}

// MD5File MD5哈希文件
func MD5File(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SHA1 SHA1哈希
func SHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1String SHA1哈希字符串
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// SHA1File SHA1哈希文件
func SHA1File(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SHA256 SHA256哈希
func SHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256String SHA256哈希字符串
func SHA256String(s string) string {
	return SHA256([]byte(s))
}

// SHA256File SHA256哈希文件
func SHA256File(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SHA512 SHA512哈希
func SHA512(data []byte) string {
	h := sha512.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA512String SHA512哈希字符串
func SHA512String(s string) string {
	return SHA512([]byte(s))
}

// SHA512File SHA512哈希文件
func SHA512File(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha512.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// CRC32 CRC32校验和
func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// CRC32String CRC32校验和字符串
func CRC32String(s string) uint32 {
	return CRC32([]byte(s))
}

// CRC32File CRC32校验和文件
func CRC32File(filename string) (uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	table := crc32.MakeTable(crc32.IEEE)
	hash := crc32.New(table)
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}

// CRC64 CRC64校验和
func CRC64(data []byte) uint64 {
	table := crc64.MakeTable(crc64.ISO)
	return crc64.Checksum(data, table)
}

// CRC64String CRC64校验和字符串
func CRC64String(s string) uint64 {
	return CRC64([]byte(s))
}

// CRC64File CRC64校验和文件
func CRC64File(filename string) (uint64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.New(table)
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}
	return hash.Sum64(), nil
}

// FNV32 FNV32哈希
func FNV32(data []byte) uint32 {
	h := fnv.New32()
	h.Write(data)
	return h.Sum32()
}

// FNV32String FNV32哈希字符串
func FNV32String(s string) uint32 {
	return FNV32([]byte(s))
}

// FNV32a FNV32a哈希
func FNV32a(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

// FNV32aString FNV32a哈希字符串
func FNV32aString(s string) uint32 {
	return FNV32a([]byte(s))
}

// FNV64 FNV64哈希
func FNV64(data []byte) uint64 {
	h := fnv.New64()
	h.Write(data)
	return h.Sum64()
}

// FNV64String FNV64哈希字符串
func FNV64String(s string) uint64 {
	return FNV64([]byte(s))
}

// FNV64a FNV64a哈希
func FNV64a(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}

// FNV64aString FNV64a哈希字符串
func FNV64aString(s string) uint64 {
	return FNV64a([]byte(s))
}

// FNV128 FNV128哈希
func FNV128(data []byte) string {
	h := fnv.New128()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// FNV128String FNV128哈希字符串
func FNV128String(s string) string {
	return FNV128([]byte(s))
}

// FNV128a FNV128a哈希
func FNV128a(data []byte) string {
	h := fnv.New128a()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// FNV128aString FNV128a哈希字符串
func FNV128aString(s string) string {
	return FNV128a([]byte(s))
}

// Hash 通用哈希函数
func Hash(data []byte, algorithm string) (string, error) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashString 通用哈希函数（字符串）
func HashString(s string, algorithm string) (string, error) {
	return Hash([]byte(s), algorithm)
}

// HashFile 通用哈希函数（文件）
func HashFile(filename string, algorithm string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// CompareHash 比较哈希值
func CompareHash(hash1, hash2 string) bool {
	return hash1 == hash2
}

// VerifyHash 验证哈希值
func VerifyHash(data []byte, algorithm string, expectedHash string) (bool, error) {
	hash, err := Hash(data, algorithm)
	if err != nil {
		return false, err
	}
	return CompareHash(hash, expectedHash), nil
}

// VerifyHashString 验证哈希值（字符串）
func VerifyHashString(s string, algorithm string, expectedHash string) (bool, error) {
	return VerifyHash([]byte(s), algorithm, expectedHash)
}

// VerifyHashFile 验证哈希值（文件）
func VerifyHashFile(filename string, algorithm string, expectedHash string) (bool, error) {
	hash, err := HashFile(filename, algorithm)
	if err != nil {
		return false, err
	}
	return CompareHash(hash, expectedHash), nil
}

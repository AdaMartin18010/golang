package encoding

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64URLEncode Base64 URL安全编码
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URLDecode Base64 URL安全解码
func Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Base64RawStdEncode Base64原始标准编码（无填充）
func Base64RawStdEncode(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

// Base64RawStdDecode Base64原始标准解码（无填充）
func Base64RawStdDecode(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}

// Base64RawURLEncode Base64原始URL安全编码（无填充）
func Base64RawURLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64RawURLDecode Base64原始URL安全解码（无填充）
func Base64RawURLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// HexEncode 十六进制编码
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode 十六进制解码
func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// HexEncodeUpper 十六进制编码（大写）
func HexEncodeUpper(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

// HexDecodeUpper 十六进制解码（大写）
func HexDecodeUpper(s string) ([]byte, error) {
	return hex.DecodeString(strings.ToUpper(s))
}

// StringToBytes 字符串转字节数组
func StringToBytes(s string) []byte {
	return []byte(s)
}

// BytesToString 字节数组转字符串
func BytesToString(b []byte) string {
	return string(b)
}

// IntToString 整数转字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int64ToString 64位整数转字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt 字符串转整数
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 字符串转64位整数
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Float64ToString 浮点数转字符串
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// StringToFloat64 字符串转浮点数
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// BoolToString 布尔值转字符串
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// StringToBool 字符串转布尔值
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// JSONEncode JSON编码
func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONEncodePretty JSON编码（格式化）
func JSONEncodePretty(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// JSONDecode JSON解码
func JSONDecode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// JSONEncodeString JSON编码为字符串
func JSONEncodeString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONEncodePrettyString JSON编码为字符串（格式化）
func JSONEncodePrettyString(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONDecodeString JSON解码字符串
func JSONDecodeString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// EscapeString 转义字符串（HTML实体）
func EscapeString(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	).Replace(s)
}

// UnescapeString 反转义字符串（HTML实体）
func UnescapeString(s string) string {
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
	).Replace(s)
}

// EscapeURL 转义URL
func EscapeURL(s string) string {
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '~' {
			result.WriteRune(r)
		} else {
			result.WriteString(fmt.Sprintf("%%%02X", r))
		}
	}
	return result.String()
}

// UnescapeURL 反转义URL
func UnescapeURL(s string) (string, error) {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			if i+2 >= len(s) {
				return "", errors.New("invalid URL encoding")
			}
			hexStr := s[i+1 : i+3]
			val, err := strconv.ParseUint(hexStr, 16, 8)
			if err != nil {
				return "", err
			}
			result.WriteByte(byte(val))
			i += 2
		} else {
			result.WriteByte(s[i])
		}
	}
	return result.String(), nil
}

// RuneToBytes 字符转字节数组
func RuneToBytes(r rune) []byte {
	return []byte(string(r))
}

// BytesToRunes 字节数组转字符数组
func BytesToRunes(b []byte) []rune {
	return []rune(string(b))
}

// RunesToString 字符数组转字符串
func RunesToString(runes []rune) string {
	return string(runes)
}

// StringToRunes 字符串转字符数组
func StringToRunes(s string) []rune {
	return []rune(s)
}

// IsBase64 检查字符串是否为有效的Base64编码
func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// IsHex 检查字符串是否为有效的十六进制编码
func IsHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

// IsJSON 检查字符串是否为有效的JSON
func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// Base64EncodeString Base64编码字符串
func Base64EncodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64DecodeString Base64解码字符串
func Base64DecodeString(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HexEncodeString 十六进制编码字符串
func HexEncodeString(s string) string {
	return hex.EncodeToString([]byte(s))
}

// HexDecodeString 十六进制解码字符串
func HexDecodeString(s string) (string, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EncodeBase64URL Base64 URL编码（别名）
func EncodeBase64URL(data []byte) string {
	return Base64URLEncode(data)
}

// DecodeBase64URL Base64 URL解码（别名）
func DecodeBase64URL(s string) ([]byte, error) {
	return Base64URLDecode(s)
}

// EncodeHex 十六进制编码（别名）
func EncodeHex(data []byte) string {
	return HexEncode(data)
}

// DecodeHex 十六进制解码（别名）
func DecodeHex(s string) ([]byte, error) {
	return HexDecode(s)
}

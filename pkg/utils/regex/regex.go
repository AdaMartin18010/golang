package regex

import (
	"regexp"
	"strings"
)

// Match 检查字符串是否匹配正则表达式
func Match(pattern, s string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

// MatchString 检查字符串是否匹配正则表达式（忽略错误）
func MatchString(pattern, s string) bool {
	matched, _ := Match(pattern, s)
	return matched
}

// Find 查找第一个匹配的子串
func Find(pattern, s string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	match := re.FindString(s)
	return match, nil
}

// FindString 查找第一个匹配的子串（忽略错误）
func FindString(pattern, s string) string {
	match, _ := Find(pattern, s)
	return match
}

// FindAll 查找所有匹配的子串
func FindAll(pattern, s string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		n = -1
	}
	return re.FindAllString(s, n), nil
}

// FindAllString 查找所有匹配的子串（忽略错误）
func FindAllString(pattern, s string, n int) []string {
	matches, _ := FindAll(pattern, s, n)
	return matches
}

// FindSubmatch 查找第一个匹配的子串和子组
func FindSubmatch(pattern, s string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.FindStringSubmatch(s), nil
}

// FindAllSubmatch 查找所有匹配的子串和子组
func FindAllSubmatch(pattern, s string, n int) ([][]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		n = -1
	}
	return re.FindAllStringSubmatch(s, n), nil
}

// Replace 替换匹配的子串
func Replace(pattern, src, repl string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(src, repl), nil
}

// ReplaceString 替换匹配的子串（忽略错误）
func ReplaceString(pattern, src, repl string) string {
	result, _ := Replace(pattern, src, repl)
	return result
}

// ReplaceFunc 使用函数替换匹配的子串
func ReplaceFunc(pattern, src string, repl func(string) string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllStringFunc(src, repl), nil
}

// ReplaceAll 替换所有匹配的子串
func ReplaceAll(pattern, src, repl string) (string, error) {
	return Replace(pattern, src, repl)
}

// Split 使用正则表达式分割字符串
func Split(pattern, s string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		n = -1
	}
	return re.Split(s, n), nil
}

// SplitString 使用正则表达式分割字符串（忽略错误）
func SplitString(pattern, s string, n int) []string {
	result, _ := Split(pattern, s, n)
	return result
}

// Extract 提取匹配的子串
func Extract(pattern, s string) ([]string, error) {
	return FindAll(pattern, s, -1)
}

// ExtractString 提取匹配的子串（忽略错误）
func ExtractString(pattern, s string) []string {
	return FindAllString(pattern, s, -1)
}

// ExtractGroups 提取匹配的子组
func ExtractGroups(pattern, s string) (map[string]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, nil
	}
	groupNames := re.SubexpNames()
	result := make(map[string]string)
	for i, name := range groupNames {
		if i > 0 && name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}
	return result, nil
}

// Count 统计匹配的数量
func Count(pattern, s string) (int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}
	matches := re.FindAllString(s, -1)
	return len(matches), nil
}

// CountString 统计匹配的数量（忽略错误）
func CountString(pattern, s string) int {
	count, _ := Count(pattern, s)
	return count
}

// IsValid 检查正则表达式是否有效
func IsValid(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

// Escape 转义正则表达式特殊字符
func Escape(s string) string {
	return regexp.QuoteMeta(s)
}

// Compile 编译正则表达式
func Compile(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile(pattern)
}

// MustCompile 编译正则表达式（失败则panic）
func MustCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

// Common patterns
var (
	// EmailPattern 邮箱正则表达式
	EmailPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// PhonePattern 手机号正则表达式（中国）
	PhonePattern = `^1[3-9]\d{9}$`

	// URLPattern URL正则表达式
	URLPattern = `^(https?|ftp)://[^\s/$.?#].[^\s]*$`

	// IPv4Pattern IPv4地址正则表达式
	IPv4Pattern = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// IPv6Pattern IPv6地址正则表达式
	IPv6Pattern = `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`

	// UUIDPattern UUID正则表达式
	UUIDPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

	// DatePattern 日期正则表达式（YYYY-MM-DD）
	DatePattern = `^\d{4}-\d{2}-\d{2}$`

	// TimePattern 时间正则表达式（HH:MM:SS）
	TimePattern = `^\d{2}:\d{2}:\d{2}$`

	// DateTimePattern 日期时间正则表达式（YYYY-MM-DD HH:MM:SS）
	DateTimePattern = `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`

	// ChinesePattern 中文字符正则表达式
	ChinesePattern = `[\u4e00-\u9fa5]`

	// NumberPattern 数字正则表达式
	NumberPattern = `^\d+$`

	// LetterPattern 字母正则表达式
	LetterPattern = `^[a-zA-Z]+$`

	// AlphanumericPattern 字母数字正则表达式
	AlphanumericPattern = `^[a-zA-Z0-9]+$`
)

// MatchEmail 检查是否为邮箱
func MatchEmail(s string) bool {
	return MatchString(EmailPattern, s)
}

// MatchPhone 检查是否为手机号（中国）
func MatchPhone(s string) bool {
	return MatchString(PhonePattern, s)
}

// MatchURL 检查是否为URL
func MatchURL(s string) bool {
	return MatchString(URLPattern, s)
}

// MatchIPv4 检查是否为IPv4地址
func MatchIPv4(s string) bool {
	return MatchString(IPv4Pattern, s)
}

// MatchIPv6 检查是否为IPv6地址
func MatchIPv6(s string) bool {
	return MatchString(IPv6Pattern, s)
}

// MatchUUID 检查是否为UUID
func MatchUUID(s string) bool {
	return MatchString(UUIDPattern, s)
}

// MatchDate 检查是否为日期格式
func MatchDate(s string) bool {
	return MatchString(DatePattern, s)
}

// MatchTime 检查是否为时间格式
func MatchTime(s string) bool {
	return MatchString(TimePattern, s)
}

// MatchDateTime 检查是否为日期时间格式
func MatchDateTime(s string) bool {
	return MatchString(DateTimePattern, s)
}

// MatchChinese 检查是否包含中文字符
func MatchChinese(s string) bool {
	return MatchString(ChinesePattern, s)
}

// MatchNumber 检查是否为数字
func MatchNumber(s string) bool {
	return MatchString(NumberPattern, s)
}

// MatchLetter 检查是否为字母
func MatchLetter(s string) bool {
	return MatchString(LetterPattern, s)
}

// MatchAlphanumeric 检查是否为字母数字
func MatchAlphanumeric(s string) bool {
	return MatchString(AlphanumericPattern, s)
}

// Remove 移除匹配的子串
func Remove(pattern, s string) (string, error) {
	return Replace(pattern, s, "")
}

// RemoveString 移除匹配的子串（忽略错误）
func RemoveString(pattern, s string) string {
	return ReplaceString(pattern, s, "")
}

// ExtractFirst 提取第一个匹配的子串
func ExtractFirst(pattern, s string) (string, error) {
	return Find(pattern, s)
}

// ExtractFirstString 提取第一个匹配的子串（忽略错误）
func ExtractFirstString(pattern, s string) string {
	return FindString(pattern, s)
}

// ExtractLast 提取最后一个匹配的子串
func ExtractLast(pattern, s string) (string, error) {
	matches, err := FindAll(pattern, s, -1)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", nil
	}
	return matches[len(matches)-1], nil
}

// ExtractLastString 提取最后一个匹配的子串（忽略错误）
func ExtractLastString(pattern, s string) string {
	match, _ := ExtractLast(pattern, s)
	return match
}

// ReplaceN 替换前n个匹配的子串
func ReplaceN(pattern, src, repl string, n int) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllStringFunc(src, func(match string) string {
		if n > 0 {
			n--
			return repl
		}
		return match
	}), nil
}

// ReplaceNString 替换前n个匹配的子串（忽略错误）
func ReplaceNString(pattern, src, repl string, n int) string {
	result, _ := ReplaceN(pattern, src, repl, n)
	return result
}

// HasMatch 检查是否有匹配
func HasMatch(pattern, s string) (bool, error) {
	return Match(pattern, s)
}

// HasMatchString 检查是否有匹配（忽略错误）
func HasMatchString(pattern, s string) bool {
	return MatchString(pattern, s)
}

// GetMatches 获取所有匹配的位置
func GetMatches(pattern, s string) ([][]int, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.FindAllStringIndex(s, -1), nil
}

// GetMatchPositions 获取所有匹配的位置（返回开始和结束位置）
func GetMatchPositions(pattern, s string) ([]struct {
	Start int
	End   int
	Match string
}, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringSubmatchIndex(s, -1)
	result := make([]struct {
		Start int
		End   int
		Match string
	}, 0, len(matches))
	for _, match := range matches {
		if len(match) >= 2 {
			result = append(result, struct {
				Start int
				End   int
				Match string
			}{
				Start: match[0],
				End:   match[1],
				Match: s[match[0]:match[1]],
			})
		}
	}
	return result, nil
}

// ReplaceWithCallback 使用回调函数替换匹配的子串
func ReplaceWithCallback(pattern, src string, callback func([]string) string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllStringFunc(src, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		return callback(submatches)
	}), nil
}

// Validate 验证字符串是否匹配正则表达式
func Validate(pattern, s string) error {
	matched, err := Match(pattern, s)
	if err != nil {
		return err
	}
	if !matched {
		return &ValidationError{Pattern: pattern, String: s}
	}
	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Pattern string
	String  string
}

func (e *ValidationError) Error() string {
	return "string does not match pattern"
}

// CompileWithOptions 编译正则表达式（带选项）
func CompileWithOptions(pattern string, caseSensitive bool) (*regexp.Regexp, error) {
	if !caseSensitive {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

// MatchWithOptions 检查字符串是否匹配正则表达式（带选项）
func MatchWithOptions(pattern, s string, caseSensitive bool) (bool, error) {
	re, err := CompileWithOptions(pattern, caseSensitive)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

// FindWithOptions 查找第一个匹配的子串（带选项）
func FindWithOptions(pattern, s string, caseSensitive bool) (string, error) {
	re, err := CompileWithOptions(pattern, caseSensitive)
	if err != nil {
		return "", err
	}
	return re.FindString(s), nil
}

// ReplaceWithOptions 替换匹配的子串（带选项）
func ReplaceWithOptions(pattern, src, repl string, caseSensitive bool) (string, error) {
	re, err := CompileWithOptions(pattern, caseSensitive)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(src, repl), nil
}

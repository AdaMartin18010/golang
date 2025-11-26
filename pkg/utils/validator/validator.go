package validator

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// EmailRegex 邮箱正则表达式
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// PhoneRegex 手机号正则表达式（中国）
	PhoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

	// IDCardRegex 身份证号正则表达式（中国）
	IDCardRegex = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)

	// URLRegex URL正则表达式
	URLRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

	// IPv4Regex IPv4地址正则表达式
	IPv4Regex = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

	// IPv6Regex IPv6地址正则表达式
	IPv6Regex = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`)

	// CreditCardRegex 信用卡号正则表达式
	CreditCardRegex = regexp.MustCompile(`^\d{13,19}$`)

	// UUIDRegex UUID正则表达式
	UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// IsEmail 验证邮箱
func IsEmail(email string) bool {
	if email == "" {
		return false
	}
	return EmailRegex.MatchString(email)
}

// IsPhone 验证手机号（中国）
func IsPhone(phone string) bool {
	if phone == "" {
		return false
	}
	return PhoneRegex.MatchString(phone)
}

// IsIDCard 验证身份证号（中国）
func IsIDCard(idCard string) bool {
	if idCard == "" {
		return false
	}
	if !IDCardRegex.MatchString(idCard) {
		return false
	}
	return validateIDCardChecksum(idCard)
}

// validateIDCardChecksum 验证身份证校验位
func validateIDCardChecksum(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}

	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	for i := 0; i < 17; i++ {
		if idCard[i] < '0' || idCard[i] > '9' {
			return false
		}
		sum += int(idCard[i]-'0') * weights[i]
	}

	checkCode := checkCodes[sum%11]
	lastChar := idCard[17]
	if lastChar >= 'a' && lastChar <= 'z' {
		lastChar -= 32
	}

	return lastChar == checkCode
}

// IsURL 验证URL
func IsURL(url string) bool {
	if url == "" {
		return false
	}
	return URLRegex.MatchString(url)
}

// IsIPv4 验证IPv4地址
func IsIPv4(ip string) bool {
	if ip == "" {
		return false
	}
	return IPv4Regex.MatchString(ip)
}

// IsIPv6 验证IPv6地址
func IsIPv6(ip string) bool {
	if ip == "" {
		return false
	}
	return IPv6Regex.MatchString(ip)
}

// IsIP 验证IP地址（IPv4或IPv6）
func IsIP(ip string) bool {
	return IsIPv4(ip) || IsIPv6(ip)
}

// IsCreditCard 验证信用卡号
func IsCreditCard(card string) bool {
	if card == "" {
		return false
	}
	if !CreditCardRegex.MatchString(card) {
		return false
	}
	return validateLuhn(card)
}

// validateLuhn Luhn算法验证
func validateLuhn(card string) bool {
	sum := 0
	alternate := false

	for i := len(card) - 1; i >= 0; i-- {
		digit := int(card[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// IsUUID 验证UUID
func IsUUID(uuid string) bool {
	if uuid == "" {
		return false
	}
	return UUIDRegex.MatchString(strings.ToLower(uuid))
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsNumeric 检查字符串是否为数字
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 检查字符串是否只包含字母
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查字符串是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsLower 检查字符串是否全为小写
func IsLower(s string) bool {
	if s == "" {
		return false
	}
	return s == strings.ToLower(s) && IsAlpha(s)
}

// IsUpper 检查字符串是否全为大写
func IsUpper(s string) bool {
	if s == "" {
		return false
	}
	return s == strings.ToUpper(s) && IsAlpha(s)
}

// HasMinLength 检查字符串长度是否大于等于最小值
func HasMinLength(s string, min int) bool {
	return len(s) >= min
}

// HasMaxLength 检查字符串长度是否小于等于最大值
func HasMaxLength(s string, max int) bool {
	return len(s) <= max
}

// HasLength 检查字符串长度是否在指定范围内
func HasLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// Contains 检查字符串是否包含子串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// StartsWith 检查字符串是否以指定前缀开始
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith 检查字符串是否以指定后缀结束
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// IsIn 检查值是否在切片中
func IsIn[T comparable](value T, slice []T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// IsNotIn 检查值是否不在切片中
func IsNotIn[T comparable](value T, slice []T) bool {
	return !IsIn(value, slice)
}

// IsBetween 检查数值是否在指定范围内
func IsBetween[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value, min, max T) bool {
	return value >= min && value <= max
}

// IsPositive 检查数值是否为正数
func IsPositive[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value T) bool {
	return value > 0
}

// IsNegative 检查数值是否为负数
func IsNegative[T int | int8 | int16 | int32 | int64 | float32 | float64](value T) bool {
	return value < 0
}

// IsZero 检查数值是否为零
func IsZero[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value T) bool {
	return value == 0
}

// IsNonZero 检查数值是否非零
func IsNonZero[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value T) bool {
	return value != 0
}

// Matches 检查字符串是否匹配正则表达式
func Matches(s string, pattern string) bool {
	if s == "" || pattern == "" {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// IsStrongPassword 检查是否为强密码
// 要求：至少8位，包含大小写字母、数字和特殊字符
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// IsWeakPassword 检查是否为弱密码
func IsWeakPassword(password string) bool {
	return !IsStrongPassword(password)
}

// IsChinese 检查字符串是否只包含中文字符
func IsChinese(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

// HasChinese 检查字符串是否包含中文字符
func HasChinese(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// IsDate 验证日期格式（YYYY-MM-DD）
func IsDate(date string) bool {
	if date == "" {
		return false
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

// IsTime 验证时间格式（HH:MM:SS）
func IsTime(time string) bool {
	if time == "" {
		return false
	}
	timeRegex := regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
	return timeRegex.MatchString(time)
}

// IsDateTime 验证日期时间格式（YYYY-MM-DD HH:MM:SS）
func IsDateTime(datetime string) bool {
	if datetime == "" {
		return false
	}
	datetimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	return datetimeRegex.MatchString(datetime)
}

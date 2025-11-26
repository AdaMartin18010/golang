package format

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatNumber 格式化数字（添加千分位分隔符）
func FormatNumber(n int64) string {
	s := strconv.FormatInt(n, 10)
	return formatNumberString(s)
}

// FormatFloat 格式化浮点数（添加千分位分隔符）
func FormatFloat(f float64, precision int) string {
	s := strconv.FormatFloat(f, 'f', precision, 64)
	parts := strings.Split(s, ".")
	integerPart := formatNumberString(parts[0])
	if len(parts) == 2 {
		return integerPart + "." + parts[1]
	}
	return integerPart
}

// formatNumberString 格式化数字字符串（添加千分位分隔符）
func formatNumberString(s string) string {
	if len(s) <= 3 {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(r)
	}
	return result.String()
}

// FormatPercent 格式化百分比
func FormatPercent(value, total float64) string {
	if total == 0 {
		return "0.00%"
	}
	percent := (value / total) * 100
	return fmt.Sprintf("%.2f%%", percent)
}

// FormatDuration 格式化持续时间（人类可读）
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := d / time.Minute
		seconds := (d % time.Minute) / time.Second
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	hours := d / time.Hour
	minutes := (d % time.Hour) / time.Minute
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// FormatDurationShort 格式化持续时间（简短格式）
func FormatDurationShort(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// FormatBytes 格式化字节数
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatBytesShort 格式化字节数（简短格式）
func FormatBytesShort(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatTime 格式化时间
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatTimeRFC3339 格式化时间为RFC3339格式
func FormatTimeRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTimeISO8601 格式化时间为ISO8601格式
func FormatTimeISO8601(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// FormatTimeHuman 格式化时间为人类可读格式
func FormatTimeHuman(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d小时前", hours)
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	}
	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (7 * 24))
		return fmt.Sprintf("%d周前", weeks)
	}
	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (30 * 24))
		return fmt.Sprintf("%d个月前", months)
	}
	years := int(diff.Hours() / (365 * 24))
	return fmt.Sprintf("%d年前", years)
}

// FormatTimeRelative 格式化相对时间
func FormatTimeRelative(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 0 {
		diff = -diff
		if diff < time.Minute {
			return "即将"
		}
		if diff < time.Hour {
			minutes := int(diff.Minutes())
			return fmt.Sprintf("%d分钟后", minutes)
		}
		if diff < 24*time.Hour {
			hours := int(diff.Hours())
			return fmt.Sprintf("%d小时后", hours)
		}
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d天后", days)
	}

	if diff < time.Minute {
		return "刚刚"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d小时前", hours)
	}
	days := int(diff.Hours() / 24)
	return fmt.Sprintf("%d天前", days)
}

// FormatPhone 格式化电话号码
func FormatPhone(phone string) string {
	// 移除所有非数字字符
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// 根据长度格式化
	switch len(digits) {
	case 11:
		// 手机号：138 0013 8000
		return fmt.Sprintf("%s %s %s", digits[0:3], digits[3:7], digits[7:])
	case 10:
		// 固定电话：010 1234 5678
		return fmt.Sprintf("%s %s %s", digits[0:3], digits[3:7], digits[7:])
	case 8:
		// 短号：1234 5678
		return fmt.Sprintf("%s %s", digits[0:4], digits[4:])
	default:
		return phone
	}
}

// FormatIDCard 格式化身份证号
func FormatIDCard(idCard string) string {
	// 移除所有非数字和X字符
	cleaned := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == 'X' || r == 'x' {
			return r
		}
		return -1
	}, idCard)

	if len(cleaned) == 18 {
		// 18位身份证：123456 19900101 1234
		return fmt.Sprintf("%s %s %s", cleaned[0:6], cleaned[6:14], cleaned[14:])
	}
	return idCard
}

// FormatBankCard 格式化银行卡号
func FormatBankCard(card string) string {
	// 移除所有非数字字符
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, card)

	// 每4位一组
	var result strings.Builder
	for i, r := range digits {
		if i > 0 && i%4 == 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String()
}

// FormatMask 格式化掩码（隐藏部分信息）
func FormatMask(s string, start, end int, mask rune) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return s
	}

	var result strings.Builder
	result.WriteString(s[:start])
	for i := start; i < end; i++ {
		result.WriteRune(mask)
	}
	result.WriteString(s[end:])
	return result.String()
}

// FormatMaskPhone 格式化手机号（中间4位掩码）
func FormatMaskPhone(phone string) string {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	if len(digits) == 11 {
		return fmt.Sprintf("%s****%s", digits[0:3], digits[7:])
	}
	return phone
}

// FormatMaskEmail 格式化邮箱（用户名部分掩码）
func FormatMaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	username := parts[0]
	if len(username) <= 2 {
		return email
	}
	masked := fmt.Sprintf("%s***%s", username[0:1], username[len(username)-1:])
	return masked + "@" + parts[1]
}

// FormatMaskIDCard 格式化身份证号（中间掩码）
func FormatMaskIDCard(idCard string) string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == 'X' || r == 'x' {
			return r
		}
		return -1
	}, idCard)

	if len(cleaned) == 18 {
		return fmt.Sprintf("%s********%s", cleaned[0:6], cleaned[14:])
	}
	return idCard
}

// FormatMaskBankCard 格式化银行卡号（中间掩码）
func FormatMaskBankCard(card string) string {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, card)

	if len(digits) >= 8 {
		return fmt.Sprintf("%s****%s", digits[0:4], digits[len(digits)-4:])
	}
	return card
}

// FormatPlural 格式化复数形式
func FormatPlural(count int, singular, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	if plural == "" {
		plural = singular + "s"
	}
	return fmt.Sprintf("%d %s", count, plural)
}

// FormatList 格式化列表
func FormatList(items []string, separator string) string {
	return strings.Join(items, separator)
}

// FormatListWithAnd 格式化列表（最后一项用"和"连接）
func FormatListWithAnd(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + "和" + items[1]
	}
	return strings.Join(items[:len(items)-1], "、") + "和" + items[len(items)-1]
}

// FormatListWithOr 格式化列表（最后一项用"或"连接）
func FormatListWithOr(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + "或" + items[1]
	}
	return strings.Join(items[:len(items)-1], "、") + "或" + items[len(items)-1]
}

// FormatTruncate 截断字符串
func FormatTruncate(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= len(suffix) {
		return suffix
	}
	return s[:maxLen-len(suffix)] + suffix
}

// FormatPadLeft 左填充
func FormatPadLeft(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return padding + s
}

// FormatPadRight 右填充
func FormatPadRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return s + padding
}

// FormatPadCenter 居中填充
func FormatPadCenter(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padding := length - len(s)
	left := padding / 2
	right := padding - left
	return strings.Repeat(string(pad), left) + s + strings.Repeat(string(pad), right)
}

// FormatIndent 缩进
func FormatIndent(s string, indent string) string {
	lines := strings.Split(s, "\n")
	var result strings.Builder
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(indent + line)
	}
	return result.String()
}

// FormatWrap 换行
func FormatWrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && i%width == 0 {
			result.WriteString("\n")
		}
		result.WriteRune(r)
	}
	return result.String()
}

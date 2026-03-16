package time

import (
	"time"
)

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// Unix 获取当前Unix时间戳（秒）
func Unix() int64 {
	return time.Now().Unix()
}

// UnixMilli 获取当前Unix时间戳（毫秒）
func UnixMilli() int64 {
	return time.Now().UnixMilli()
}

// UnixMicro 获取当前Unix时间戳（微秒）
func UnixMicro() int64 {
	return time.Now().UnixMicro()
}

// UnixNano 获取当前Unix时间戳（纳秒）
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// Format 格式化时间
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatDefault 使用默认格式格式化时间
func FormatDefault(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// Parse 解析时间字符串
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// ParseDefault 使用默认格式解析时间字符串
func ParseDefault(value string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", value)
}

// ParseDate 解析日期字符串
func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 添加月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 添加年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// StartOfDay 获取一天的开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取一天的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek 获取一周的开始时间（周一）
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysToMonday := weekday - 1
	return StartOfDay(t.AddDate(0, 0, -daysToMonday))
}

// EndOfWeek 获取一周的结束时间（周日）
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysToSunday := 7 - weekday
	return EndOfDay(t.AddDate(0, 0, daysToSunday))
}

// StartOfMonth 获取一月的开始时间
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取一月的结束时间
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	return time.Date(year, nextMonth, 1, 0, 0, 0, 0, t.Location()).Add(-time.Nanosecond)
}

// StartOfYear 获取一年的开始时间
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取一年的结束时间
func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, t.Location())
}

// DaysBetween 计算两个时间之间的天数
func DaysBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours() / 24)
}

// HoursBetween 计算两个时间之间的小时数
func HoursBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours())
}

// MinutesBetween 计算两个时间之间的分钟数
func MinutesBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Minutes())
}

// SecondsBetween 计算两个时间之间的秒数
func SecondsBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Seconds())
}

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 判断是否是昨天
func IsYesterday(t time.Time) bool {
	return IsToday(AddDays(t, 1))
}

// IsTomorrow 判断是否是明天
func IsTomorrow(t time.Time) bool {
	return IsToday(AddDays(t, -1))
}

// IsSameDay 判断是否是同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// IsSameWeek 判断是否是同一周
func IsSameWeek(t1, t2 time.Time) bool {
	start1 := StartOfWeek(t1)
	start2 := StartOfWeek(t2)
	return IsSameDay(start1, start2)
}

// IsSameMonth 判断是否是同一月
func IsSameMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// IsSameYear 判断是否是同一年
func IsSameYear(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year()
}

// HumanizeDuration 人性化显示时长
func HumanizeDuration(d time.Duration) string {
	if d < time.Minute {
		return d.String()
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		return formatDuration(minutes, "分钟")
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		return formatDuration(hours, "小时")
	}
	days := int(d.Hours() / 24)
	return formatDuration(days, "天")
}

func formatDuration(value int, unit string) string {
	return formatInt(value) + unit
}

func formatInt(n int) string {
	if n < 0 {
		return "-" + formatInt(-n)
	}
	if n < 10 {
		return string(rune('0' + n))
	}
	return formatInt(n/10) + string(rune('0'+n%10))
}

// AddHours 添加小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// HumanizeTime 人性化显示时间
func HumanizeTime(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		return "刚刚"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return formatDuration(minutes, "分钟前")
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return formatDuration(hours, "小时前")
	}
	if IsToday(t) {
		return "今天 " + FormatTime(t)
	}
	if IsYesterday(t) {
		return "昨天 " + FormatTime(t)
	}
	if IsSameWeek(now, t) {
		weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
		return weekdays[t.Weekday()] + " " + FormatTime(t)
	}
	return FormatDefault(t)
}

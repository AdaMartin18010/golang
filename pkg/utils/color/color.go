package color

import (
	"fmt"
	"os"
)

// Color 颜色类型
type Color int

// 颜色常量
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// 颜色常量（高亮）
const (
	BrightBlack Color = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// 样式常量
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Hidden    = "\033[8m"
)

// 背景色常量
const (
	BgBlack   = 40
	BgRed     = 41
	BgGreen   = 42
	BgYellow  = 43
	BgBlue    = 44
	BgMagenta = 45
	BgCyan    = 46
	BgWhite   = 47
)

// 背景色常量（高亮）
const (
	BgBrightBlack   = 100
	BgBrightRed     = 101
	BgBrightGreen   = 102
	BgBrightYellow  = 103
	BgBrightBlue    = 104
	BgBrightMagenta = 105
	BgBrightCyan    = 106
	BgBrightWhite   = 107
)

var (
	// 是否启用颜色
	enabled = true
	// 是否检测终端支持
	autoDetect = true
)

// Enable 启用颜色
func Enable() {
	enabled = true
}

// Disable 禁用颜色
func Disable() {
	enabled = false
}

// SetEnabled 设置是否启用颜色
func SetEnabled(e bool) {
	enabled = e
}

// IsEnabled 检查是否启用颜色
func IsEnabled() bool {
	if !autoDetect {
		return enabled
	}
	// 检测终端是否支持颜色
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}
	return enabled && (term != "dumb")
}

// SetAutoDetect 设置是否自动检测终端支持
func SetAutoDetect(detect bool) {
	autoDetect = detect
}

// Colorize 为文本添加颜色
func Colorize(text string, color Color) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("\033[%dm%s%s", color, text, Reset)
}

// ColorizeWithStyle 为文本添加颜色和样式
func ColorizeWithStyle(text string, color Color, styles ...string) string {
	if !IsEnabled() {
		return text
	}
	styleStr := ""
	for _, style := range styles {
		styleStr += style
	}
	return fmt.Sprintf("%s\033[%dm%s%s", styleStr, color, text, Reset)
}

// Black 黑色文本
func Black(text string) string {
	return Colorize(text, Black)
}

// Red 红色文本
func Red(text string) string {
	return Colorize(text, Red)
}

// Green 绿色文本
func Green(text string) string {
	return Colorize(text, Green)
}

// Yellow 黄色文本
func Yellow(text string) string {
	return Colorize(text, Yellow)
}

// Blue 蓝色文本
func Blue(text string) string {
	return Colorize(text, Blue)
}

// Magenta 洋红色文本
func Magenta(text string) string {
	return Colorize(text, Magenta)
}

// Cyan 青色文本
func Cyan(text string) string {
	return Colorize(text, Cyan)
}

// White 白色文本
func White(text string) string {
	return Colorize(text, White)
}

// BrightBlack 高亮黑色文本
func BrightBlack(text string) string {
	return Colorize(text, BrightBlack)
}

// BrightRed 高亮红色文本
func BrightRed(text string) string {
	return Colorize(text, BrightRed)
}

// BrightGreen 高亮绿色文本
func BrightGreen(text string) string {
	return Colorize(text, BrightGreen)
}

// BrightYellow 高亮黄色文本
func BrightYellow(text string) string {
	return Colorize(text, BrightYellow)
}

// BrightBlue 高亮蓝色文本
func BrightBlue(text string) string {
	return Colorize(text, BrightBlue)
}

// BrightMagenta 高亮洋红色文本
func BrightMagenta(text string) string {
	return Colorize(text, BrightMagenta)
}

// BrightCyan 高亮青色文本
func BrightCyan(text string) string {
	return Colorize(text, BrightCyan)
}

// BrightWhite 高亮白色文本
func BrightWhite(text string) string {
	return Colorize(text, BrightWhite)
}

// BgBlack 黑色背景
func BgBlack(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgBlack))
}

// BgRed 红色背景
func BgRed(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgRed))
}

// BgGreen 绿色背景
func BgGreen(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgGreen))
}

// BgYellow 黄色背景
func BgYellow(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgYellow))
}

// BgBlue 蓝色背景
func BgBlue(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgBlue))
}

// BgMagenta 洋红色背景
func BgMagenta(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgMagenta))
}

// BgCyan 青色背景
func BgCyan(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgCyan))
}

// BgWhite 白色背景
func BgWhite(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgWhite))
}

// Bold 粗体文本
func Bold(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Bold, text, Reset)
}

// Dim 暗淡文本
func Dim(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Dim, text, Reset)
}

// Italic 斜体文本
func Italic(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Italic, text, Reset)
}

// Underline 下划线文本
func Underline(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Underline, text, Reset)
}

// Blink 闪烁文本
func Blink(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Blink, text, Reset)
}

// Reverse 反转文本
func Reverse(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Reverse, text, Reset)
}

// Hidden 隐藏文本
func Hidden(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", Hidden, text, Reset)
}

// Success 成功消息（绿色）
func Success(text string) string {
	return Green(text)
}

// Error 错误消息（红色）
func Error(text string) string {
	return Red(text)
}

// Warning 警告消息（黄色）
func Warning(text string) string {
	return Yellow(text)
}

// Info 信息消息（蓝色）
func Info(text string) string {
	return Blue(text)
}

// Debug 调试消息（青色）
func Debug(text string) string {
	return Cyan(text)
}

// Print 打印彩色文本
func Print(text string, color Color) {
	fmt.Print(Colorize(text, color))
}

// Println 打印彩色文本（换行）
func Println(text string, color Color) {
	fmt.Println(Colorize(text, color))
}

// Printf 格式化打印彩色文本
func Printf(format string, color Color, args ...interface{}) {
	fmt.Print(Colorize(fmt.Sprintf(format, args...), color))
}

// PrintSuccess 打印成功消息
func PrintSuccess(text string) {
	fmt.Print(Success(text))
}

// PrintError 打印错误消息
func PrintError(text string) {
	fmt.Print(Error(text))
}

// PrintWarning 打印警告消息
func PrintWarning(text string) {
	fmt.Print(Warning(text))
}

// PrintInfo 打印信息消息
func PrintInfo(text string) {
	fmt.Print(Info(text))
}

// PrintDebug 打印调试消息
func PrintDebug(text string) {
	fmt.Print(Debug(text))
}

// PrintlnSuccess 打印成功消息（换行）
func PrintlnSuccess(text string) {
	fmt.Println(Success(text))
}

// PrintlnError 打印错误消息（换行）
func PrintlnError(text string) {
	fmt.Println(Error(text))
}

// PrintlnWarning 打印警告消息（换行）
func PrintlnWarning(text string) {
	fmt.Println(Warning(text))
}

// PrintlnInfo 打印信息消息（换行）
func PrintlnInfo(text string) {
	fmt.Println(Info(text))
}

// PrintlnDebug 打印调试消息（换行）
func PrintlnDebug(text string) {
	fmt.Println(Debug(text))
}

package color

import (
	"fmt"
	"os"
)

// Color 颜色类型
type Color int

// 颜色常量
const (
	BlackColor Color = iota + 30
	RedColor
	GreenColor
	YellowColor
	BlueColor
	MagentaColor
	CyanColor
	WhiteColor
)

// 颜色常量（高亮）
const (
	BrightBlackColor Color = iota + 90
	BrightRedColor
	BrightGreenColor
	BrightYellowColor
	BrightBlueColor
	BrightMagentaColor
	BrightCyanColor
	BrightWhiteColor
)

// 样式常量
const (
	ResetCode     = "\033[0m"
	BoldCode      = "\033[1m"
	DimCode       = "\033[2m"
	ItalicCode    = "\033[3m"
	UnderlineCode = "\033[4m"
	BlinkCode     = "\033[5m"
	ReverseCode   = "\033[7m"
	HiddenCode    = "\033[8m"
)

// 背景色常量
const (
	BgBlackColor   = 40
	BgRedColor     = 41
	BgGreenColor   = 42
	BgYellowColor  = 43
	BgBlueColor    = 44
	BgMagentaColor = 45
	BgCyanColor    = 46
	BgWhiteColor   = 47
)

// 背景色常量（高亮）
const (
	BgBrightBlackColor   = 100
	BgBrightRedColor     = 101
	BgBrightGreenColor   = 102
	BgBrightYellowColor  = 103
	BgBrightBlueColor    = 104
	BgBrightMagentaColor = 105
	BgBrightCyanColor    = 106
	BgBrightWhiteColor   = 107
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
	return fmt.Sprintf("\033[%dm%s%s", color, text, ResetCode)
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
	return fmt.Sprintf("%s\033[%dm%s%s", styleStr, color, text, ResetCode)
}

// Black 黑色文本
func Black(text string) string {
	return Colorize(text, BlackColor)
}

// Red 红色文本
func Red(text string) string {
	return Colorize(text, RedColor)
}

// Green 绿色文本
func Green(text string) string {
	return Colorize(text, GreenColor)
}

// Yellow 黄色文本
func Yellow(text string) string {
	return Colorize(text, YellowColor)
}

// Blue 蓝色文本
func Blue(text string) string {
	return Colorize(text, BlueColor)
}

// Magenta 洋红色文本
func Magenta(text string) string {
	return Colorize(text, MagentaColor)
}

// Cyan 青色文本
func Cyan(text string) string {
	return Colorize(text, CyanColor)
}

// White 白色文本
func White(text string) string {
	return Colorize(text, WhiteColor)
}

// BrightBlack 高亮黑色文本
func BrightBlack(text string) string {
	return Colorize(text, BrightBlackColor)
}

// BrightRed 高亮红色文本
func BrightRed(text string) string {
	return Colorize(text, BrightRedColor)
}

// BrightGreen 高亮绿色文本
func BrightGreen(text string) string {
	return Colorize(text, BrightGreenColor)
}

// BrightYellow 高亮黄色文本
func BrightYellow(text string) string {
	return Colorize(text, BrightYellowColor)
}

// BrightBlue 高亮蓝色文本
func BrightBlue(text string) string {
	return Colorize(text, BrightBlueColor)
}

// BrightMagenta 高亮洋红色文本
func BrightMagenta(text string) string {
	return Colorize(text, BrightMagentaColor)
}

// BrightCyan 高亮青色文本
func BrightCyan(text string) string {
	return Colorize(text, BrightCyanColor)
}

// BrightWhite 高亮白色文本
func BrightWhite(text string) string {
	return Colorize(text, BrightWhiteColor)
}

// BgBlack 黑色背景
func BgBlack(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgBlackColor))
}

// BgRed 红色背景
func BgRed(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgRedColor))
}

// BgGreen 绿色背景
func BgGreen(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgGreenColor))
}

// BgYellow 黄色背景
func BgYellow(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgYellowColor))
}

// BgBlue 蓝色背景
func BgBlue(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgBlueColor))
}

// BgMagenta 洋红色背景
func BgMagenta(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgMagentaColor))
}

// BgCyan 青色背景
func BgCyan(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgCyanColor))
}

// BgWhite 白色背景
func BgWhite(text string) string {
	return ColorizeWithStyle(text, 0, fmt.Sprintf("\033[%dm", BgWhiteColor))
}

// Bold 粗体文本
func Bold(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", BoldCode, text, ResetCode)
}

// Dim 暗淡文本
func Dim(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", DimCode, text, ResetCode)
}

// Italic 斜体文本
func Italic(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", ItalicCode, text, ResetCode)
}

// Underline 下划线文本
func Underline(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", UnderlineCode, text, ResetCode)
}

// Blink 闪烁文本
func Blink(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", BlinkCode, text, ResetCode)
}

// Reverse 反转文本
func Reverse(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", ReverseCode, text, ResetCode)
}

// Hidden 隐藏文本
func Hidden(text string) string {
	if !IsEnabled() {
		return text
	}
	return fmt.Sprintf("%s%s%s", HiddenCode, text, ResetCode)
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

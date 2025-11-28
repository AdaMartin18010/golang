// Package ui 提供用户界面相关功能
package ui

import (
	"fmt"
	"os"
	"runtime"
)

// Color 终端颜色代码
type Color string

const (
	// 基本颜色
	ColorReset Color = "\033[0m"
	ColorBold  Color = "\033[1m"
	ColorDim   Color = "\033[2m"
	ColorUnder Color = "\033[4m"

	// 前景色
	ColorBlack   Color = "\033[30m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorWhite   Color = "\033[37m"

	// 明亮前景色
	ColorBrightBlack   Color = "\033[90m"
	ColorBrightRed     Color = "\033[91m"
	ColorBrightGreen   Color = "\033[92m"
	ColorBrightYellow  Color = "\033[93m"
	ColorBrightBlue    Color = "\033[94m"
	ColorBrightMagenta Color = "\033[95m"
	ColorBrightCyan    Color = "\033[96m"
	ColorBrightWhite   Color = "\033[97m"

	// 背景色
	ColorBgBlack   Color = "\033[40m"
	ColorBgRed     Color = "\033[41m"
	ColorBgGreen   Color = "\033[42m"
	ColorBgYellow  Color = "\033[43m"
	ColorBgBlue    Color = "\033[44m"
	ColorBgMagenta Color = "\033[45m"
	ColorBgCyan    Color = "\033[46m"
	ColorBgWhite   Color = "\033[47m"
)

var (
	// 是否启用颜色输出
	colorEnabled = true
)

// init 初始化颜色支持
func init() {
	// Windows 10+ 支持 ANSI 颜色
	if runtime.GOOS == "windows" {
		// 检查是否在支持颜色的终端中
		if os.Getenv("TERM") == "" && os.Getenv("ConEmuANSI") == "" {
			// 在某些 Windows 终端中可能不支持
			// 但我们默认尝试启用
		}
	}

	// 检查 NO_COLOR 环境变量
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}
}

// SetColorEnabled 设置是否启用颜色
func SetColorEnabled(enabled bool) {
	colorEnabled = enabled
}

// IsColorEnabled 返回是否启用颜色
func IsColorEnabled() bool {
	return colorEnabled
}

// Colorize 为文本添加颜色
func Colorize(text string, color Color) string {
	if !colorEnabled {
		return text
	}
	return string(color) + text + string(ColorReset)
}

// Bold 加粗文本
func Bold(text string) string {
	return Colorize(text, ColorBold)
}

// Success 成功消息（绿色）
func Success(text string) string {
	return Colorize("✅ "+text, ColorGreen)
}

// Error 错误消息（红色）
func Error(text string) string {
	return Colorize("❌ "+text, ColorRed)
}

// Warning 警告消息（黄色）
func Warning(text string) string {
	return Colorize("⚠️  "+text, ColorYellow)
}

// Info 信息消息（蓝色）
func Info(text string) string {
	return Colorize("ℹ️  "+text, ColorBlue)
}

// Debug 调试消息（暗色）
func Debug(text string) string {
	return Colorize("🐛 "+text, ColorDim)
}

// Progress 进度消息（青色）
func Progress(text string) string {
	return Colorize("🔄 "+text, ColorCyan)
}

// PrintSuccess 打印成功消息
func PrintSuccess(format string, a ...interface{}) {
	fmt.Println(Success(fmt.Sprintf(format, a...)))
}

// PrintError 打印错误消息
func PrintError(format string, a ...interface{}) {
	fmt.Println(Error(fmt.Sprintf(format, a...)))
}

// PrintWarning 打印警告消息
func PrintWarning(format string, a ...interface{}) {
	fmt.Println(Warning(fmt.Sprintf(format, a...)))
}

// PrintInfo 打印信息消息
func PrintInfo(format string, a ...interface{}) {
	fmt.Println(Info(fmt.Sprintf(format, a...)))
}

// PrintDebug 打印调试消息
func PrintDebug(format string, a ...interface{}) {
	fmt.Println(Debug(fmt.Sprintf(format, a...)))
}

// PrintProgress 打印进度消息
func PrintProgress(format string, a ...interface{}) {
	fmt.Println(Progress(fmt.Sprintf(format, a...)))
}

// Header 打印标题
func Header(text string) string {
	if !colorEnabled {
		return "\n=== " + text + " ===\n"
	}

	border := "═══════════════════════════════════════════════════════"
	return fmt.Sprintf("\n%s\n%s %s %s\n%s\n",
		Colorize(border, ColorCyan),
		Colorize("█", ColorBrightCyan),
		Colorize(text, ColorBold),
		Colorize("█", ColorBrightCyan),
		Colorize(border, ColorCyan))
}

// PrintHeader 打印标题
func PrintHeader(text string) {
	fmt.Print(Header(text))
}

// Divider 打印分隔线
func Divider() string {
	if !colorEnabled {
		return "---------------------------------------------------"
	}
	return Colorize("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━", ColorDim)
}

// PrintDivider 打印分隔线
func PrintDivider() {
	fmt.Println(Divider())
}

// Box 创建文本框
func Box(title, content string) string {
	if !colorEnabled {
		return fmt.Sprintf("\n+--- %s ---+\n%s\n+----------+\n", title, content)
	}

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n",
		Colorize("╔═══ "+title+" ═══╗", ColorCyan),
		content,
		Colorize("╚═══════════╝", ColorCyan),
		"")
}

// PrintBox 打印文本框
func PrintBox(title, content string) {
	fmt.Print(Box(title, content))
}

// Bullet 创建项目符号列表项
func Bullet(text string) string {
	if !colorEnabled {
		return "• " + text
	}
	return Colorize("●", ColorBrightBlue) + " " + text
}

// CheckMark 勾选标记
func CheckMark(text string) string {
	if !colorEnabled {
		return "[✓] " + text
	}
	return Colorize("✓", ColorGreen) + " " + text
}

// CrossMark 叉号标记
func CrossMark(text string) string {
	if !colorEnabled {
		return "[✗] " + text
	}
	return Colorize("✗", ColorRed) + " " + text
}

// Badge 创建徽章
func Badge(text string, color Color) string {
	if !colorEnabled {
		return "[" + text + "]"
	}
	return string(color) + string(ColorBgWhite) + " " + text + " " + string(ColorReset)
}

// Spinner 动画字符
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ProgressBar 创建进度条
func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	if !colorEnabled {
		return fmt.Sprintf("[%s] %d/%d (%.1f%%)", bar, current, total, percentage*100)
	}

	return fmt.Sprintf("%s %d/%d %s",
		Colorize(bar, ColorGreen),
		current,
		total,
		Colorize(fmt.Sprintf("(%.1f%%)", percentage*100), ColorDim))
}

// Table 简单表格
type Table struct {
	Headers []string
	Rows    [][]string
	Width   []int
}

// NewTable 创建新表格
func NewTable(headers ...string) *Table {
	return &Table{
		Headers: headers,
		Rows:    make([][]string, 0),
		Width:   make([]int, len(headers)),
	}
}

// AddRow 添加行
func (t *Table) AddRow(cells ...string) {
	t.Rows = append(t.Rows, cells)
}

// String 转换为字符串
func (t *Table) String() string {
	// 计算列宽
	for i, header := range t.Headers {
		if len(header) > t.Width[i] {
			t.Width[i] = len(header)
		}
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(t.Width) && len(cell) > t.Width[i] {
				t.Width[i] = len(cell)
			}
		}
	}

	result := ""

	// 打印标题
	if colorEnabled {
		result += string(ColorBold)
	}
	for i, header := range t.Headers {
		result += fmt.Sprintf("%-*s  ", t.Width[i], header)
	}
	if colorEnabled {
		result += string(ColorReset)
	}
	result += "\n"

	// 打印分隔线
	for i := range t.Headers {
		for j := 0; j < t.Width[i]; j++ {
			result += "─"
		}
		result += "  "
	}
	result += "\n"

	// 打印行
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(t.Width) {
				result += fmt.Sprintf("%-*s  ", t.Width[i], cell)
			}
		}
		result += "\n"
	}

	return result
}

// Print 打印表格
func (t *Table) Print() {
	fmt.Print(t.String())
}

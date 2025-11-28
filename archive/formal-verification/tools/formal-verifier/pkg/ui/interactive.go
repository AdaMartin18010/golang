package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt 提示用户输入
func Prompt(message string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, Colorize(defaultValue, ColorDim))
	} else {
		fmt.Printf("%s: ", message)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}

	return input
}

// Confirm 确认提示（是/否）
func Confirm(message string, defaultYes bool) bool {
	var prompt string
	if defaultYes {
		prompt = fmt.Sprintf("%s [%s/%s]", message,
			Colorize("Y", ColorGreen),
			Colorize("n", ColorDim))
	} else {
		prompt = fmt.Sprintf("%s [%s/%s]", message,
			Colorize("y", ColorDim),
			Colorize("N", ColorRed))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// Select 选择菜单
func Select(message string, options []string) int {
	PrintHeader(message)

	for i, option := range options {
		fmt.Printf("  %s %s\n",
			Colorize(fmt.Sprintf("[%d]", i+1), ColorCyan),
			option)
	}

	fmt.Println()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s (1-%d): ", Colorize("选择", ColorGreen), len(options))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var choice int
		if _, err := fmt.Sscanf(input, "%d", &choice); err == nil {
			if choice >= 1 && choice <= len(options) {
				return choice - 1
			}
		}

		PrintError("无效的选择，请输入 1-%d 之间的数字", len(options))
	}
}

// MultiSelect 多选菜单
func MultiSelect(message string, options []string) []int {
	PrintHeader(message)

	for i, option := range options {
		fmt.Printf("  %s %s\n",
			Colorize(fmt.Sprintf("[%d]", i+1), ColorCyan),
			option)
	}

	fmt.Println()
	PrintInfo("输入多个选项（用空格或逗号分隔），例如: 1 3 5")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", Colorize("选择", ColorGreen))

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 分割输入
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ' ' || r == ',' || r == ';'
	})

	selected := make([]int, 0)
	for _, part := range parts {
		var choice int
		if _, err := fmt.Sscanf(part, "%d", &choice); err == nil {
			if choice >= 1 && choice <= len(options) {
				selected = append(selected, choice-1)
			}
		}
	}

	return selected
}

// Menu 交互式菜单
type Menu struct {
	Title   string
	Options []MenuOption
}

// MenuOption 菜单选项
type MenuOption struct {
	Label       string
	Description string
	Action      func()
}

// NewMenu 创建新菜单
func NewMenu(title string) *Menu {
	return &Menu{
		Title:   title,
		Options: make([]MenuOption, 0),
	}
}

// AddOption 添加选项
func (m *Menu) AddOption(label, description string, action func()) {
	m.Options = append(m.Options, MenuOption{
		Label:       label,
		Description: description,
		Action:      action,
	})
}

// Show 显示菜单并处理用户选择
func (m *Menu) Show() {
	for {
		PrintHeader(m.Title)

		for i, option := range m.Options {
			fmt.Printf("  %s %s\n",
				Colorize(fmt.Sprintf("[%d]", i+1), ColorCyan),
				Bold(option.Label))
			if option.Description != "" {
				fmt.Printf("      %s\n", Colorize(option.Description, ColorDim))
			}
		}
		fmt.Printf("\n  %s %s\n\n",
			Colorize("[0]", ColorRed),
			"退出")

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s (0-%d): ", Colorize("选择", ColorGreen), len(m.Options))

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var choice int
		if _, err := fmt.Sscanf(input, "%d", &choice); err == nil {
			if choice == 0 {
				PrintInfo("退出菜单")
				return
			}

			if choice >= 1 && choice <= len(m.Options) {
				fmt.Println()
				m.Options[choice-1].Action()
				fmt.Println()

				if !Confirm("继续", true) {
					return
				}
				continue
			}
		}

		PrintError("无效的选择")
		fmt.Println()
	}
}

// Banner 打印应用横幅
func Banner(appName, version, description string) {
	if !colorEnabled {
		fmt.Printf("\n%s v%s\n%s\n\n", appName, version, description)
		return
	}

	banner := fmt.Sprintf(`
%s
%s  _____ __      __
%s |  ___|\ \    / /   %s
%s | |_    \ \  / /    %s
%s |  _|    \ \/ /     %s v%s
%s |_|       \__/      %s
%s
%s`,
		Colorize("╔═══════════════════════════════════════════════════════╗", ColorCyan),
		Colorize("║", ColorCyan),
		Colorize("║", ColorCyan), Colorize("Go Formal Verifier", ColorBold),
		Colorize("║", ColorCyan), Colorize("形式化验证工具", ColorDim),
		Colorize("║", ColorCyan), Colorize("FV", ColorBrightCyan), Colorize(version, ColorYellow),
		Colorize("║", ColorCyan), Colorize(description, ColorDim),
		Colorize("║", ColorCyan),
		Colorize("╚═══════════════════════════════════════════════════════╝", ColorCyan))

	fmt.Println(banner)
}

// AnimatedProgress 动画进度指示器
type AnimatedProgress struct {
	message string
	frame   int
	done    chan bool
}

// NewAnimatedProgress 创建新的动画进度指示器
func NewAnimatedProgress(message string) *AnimatedProgress {
	return &AnimatedProgress{
		message: message,
		frame:   0,
		done:    make(chan bool),
	}
}

// Start 开始动画
func (p *AnimatedProgress) Start() {
	go func() {
		for {
			select {
			case <-p.done:
				return
			default:
				fmt.Printf("\r%s %s %s",
					Colorize(SpinnerFrames[p.frame%len(SpinnerFrames)], ColorCyan),
					p.message,
					strings.Repeat(" ", 20))
				p.frame++
			}
		}
	}()
}

// Stop 停止动画
func (p *AnimatedProgress) Stop() {
	p.done <- true
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
}

// ShowProgress 显示进度（简化版）
func ShowProgress(message string, current, total int) {
	bar := ProgressBar(current, total, 30)
	fmt.Printf("\r%s %s", message, bar)
	if current == total {
		fmt.Println()
	}
}

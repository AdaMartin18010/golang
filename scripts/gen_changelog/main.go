package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// 简易变更日志追加脚本：从标准输入读取条目，写入 CHANGELOG.md 指定版本段落
func main() {
	version := os.Getenv("VERSION")
	if version == "" {
		version = time.Now().Format("v2006.01-02")
	}
	entries := readStdin()
	if len(entries) == 0 {
		fmt.Println("no entries provided")
		return
	}
	f, err := os.OpenFile("CHANGELOG.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open changelog:", err)
		return
	}
	defer f.Close()
	header := fmt.Sprintf("\n## %s (generated)\n\n", version)
	if _, err := f.WriteString(header + strings.Join(entries, "\n") + "\n"); err != nil {
		fmt.Println("write changelog:", err)
		return
	}
	fmt.Println("CHANGELOG updated for", version)
}

func readStdin() []string {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return nil
	}
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, "- "+line)
		}
	}
	return lines
}

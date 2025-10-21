package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Stats 统计信息
type Stats struct {
	FilesScanned                   int
	FilesModified                  int
	TrailingEmptyLinesInCodeBlocks int
	ExcessiveNewlines              int
	TrailingSpaces                 int
	MissingLanguageTags            int
	MermaidFormat                  int
}

var stats Stats

func main() {
	targetDir := flag.String("dir", "docs", "目标目录")
	dryRun := flag.Bool("dry-run", false, "Dry Run模式（预览）")
	flag.Parse()

	log("==================================================", "INFO")
	log("代码格式统一修复脚本", "INFO")
	log("==================================================", "INFO")
	log(fmt.Sprintf("目标目录: %s", *targetDir), "INFO")
	mode := "实际修改"
	if *dryRun {
		mode = "Dry Run (预览)"
	}
	log(fmt.Sprintf("模式: %s", mode), "INFO")
	log("==================================================", "INFO")

	// 获取所有Markdown文件
	mdFiles, err := findMarkdownFiles(*targetDir)
	if err != nil {
		log(fmt.Sprintf("Error finding markdown files: %v", err), "ERROR")
		os.Exit(1)
	}

	log(fmt.Sprintf("找到 %d 个Markdown文件", len(mdFiles)), "INFO")
	fmt.Println()

	// 处理每个文件
	for _, filepath := range mdFiles {
		processMarkdownFile(filepath, *dryRun)
	}

	// 输出统计报告
	fmt.Println()
	log("==================================================", "SUCCESS")
	log("修复完成！统计报告：", "SUCCESS")
	log("==================================================", "INFO")
	log(fmt.Sprintf("扫描文件数: %d", stats.FilesScanned), "INFO")
	log(fmt.Sprintf("修改文件数: %d", stats.FilesModified), "INFO")
	log("--------------------------------------------------", "INFO")
	log("修复问题统计：", "INFO")
	log(fmt.Sprintf("缺失语言标记: %d", stats.MissingLanguageTags), "INFO")
	log(fmt.Sprintf("代码块尾部空行: %d", stats.TrailingEmptyLinesInCodeBlocks), "INFO")
	log(fmt.Sprintf("Mermaid格式: %d", stats.MermaidFormat), "INFO")
	log(fmt.Sprintf("连续空行修复: %d", stats.ExcessiveNewlines), "INFO")
	log(fmt.Sprintf("行尾空格移除: %d", stats.TrailingSpaces), "INFO")
	log("==================================================", "INFO")

	if *dryRun {
		fmt.Println()
		log("这是Dry Run模式，未实际修改文件", "WARN")
		log("移除 -dry-run 参数以执行实际修改", "WARN")
	}
}

func log(message, level string) {
	timestamp := time.Now().Format("15:04:05")
	colors := map[string]string{
		"ERROR":   "\033[91m",
		"WARN":    "\033[93m",
		"SUCCESS": "\033[92m",
		"INFO":    "\033[97m",
	}
	reset := "\033[0m"
	color, ok := colors[level]
	if !ok {
		color = colors["INFO"]
	}
	fmt.Printf("%s[%s] %s: %s%s\n", color, timestamp, level, message, reset)
}

func findMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func processMarkdownFile(filepath string, dryRun bool) {
	stats.FilesScanned++

	// 读取文件
	content, err := os.ReadFile(filepath)
	if err != nil {
		log(fmt.Sprintf("Error reading %s: %v", filepath, err), "ERROR")
		return
	}

	if len(content) == 0 {
		return
	}

	originalContent := string(content)
	modifiedContent := originalContent
	fileModified := false

	// 1. 修复代码块格式
	modifiedContent, mod1 := fixCodeBlocks(modifiedContent)
	if mod1 {
		fileModified = true
	}

	// 2. 修复Mermaid图表格式
	modifiedContent, mod2 := fixMermaidDiagrams(modifiedContent)
	if mod2 {
		fileModified = true
	}

	// 3. 修复连续空行
	modifiedContent, mod3 := fixExcessiveNewlines(modifiedContent)
	if mod3 {
		fileModified = true
	}

	// 4. 确保文件末尾只有一个换行符
	modifiedContent = strings.TrimRight(modifiedContent, "\r\n \t") + "\n"

	// 保存文件
	if fileModified || modifiedContent != originalContent {
		if !dryRun {
			err := os.WriteFile(filepath, []byte(modifiedContent), 0644)
			if err != nil {
				log(fmt.Sprintf("Error writing %s: %v", filepath, err), "ERROR")
				return
			}
			log(fmt.Sprintf("Modified: %s", filepath), "SUCCESS")
		} else {
			log(fmt.Sprintf("Would modify: %s", filepath), "WARN")
		}
		stats.FilesModified++
	}
}

func fixCodeBlocks(content string) (string, bool) {
	modified := false
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inCodeBlock := false
	i := 0

	codeBlockStart := regexp.MustCompile(`^` + "```" + `(\w*)(.*)$`)
	goPattern := regexp.MustCompile(`^(package|func|import|type|var|const)\s`)
	bashPattern := regexp.MustCompile(`^(\$|#|cd|ls|git|npm|go\s)`)

	for i < len(lines) {
		line := lines[i]

		// 检测代码块开始
		if matches := codeBlockStart.FindStringSubmatch(line); matches != nil && !inCodeBlock {
			lang := matches[1]
			extra := strings.TrimSpace(matches[2])

			// 尝试推断缺失的语言标记
			if lang == "" && i+1 < len(lines) {
				nextLine := lines[i+1]
				if goPattern.MatchString(nextLine) {
					lang = "go"
					modified = true
					stats.MissingLanguageTags++
				} else if bashPattern.MatchString(nextLine) {
					lang = "bash"
					modified = true
					stats.MissingLanguageTags++
				}
			}

			// 统一格式
			result = append(result, "```"+lang)
			if extra != "" {
				modified = true
			}

			inCodeBlock = true
			i++
			continue
		}

		// 检测代码块结束
		if strings.TrimSpace(line) == "```" && inCodeBlock {
			// 移除代码块末尾的空行
			for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
				result = result[:len(result)-1]
				modified = true
				stats.TrailingEmptyLinesInCodeBlocks++
			}

			result = append(result, "```")
			inCodeBlock = false

			// 代码块后应该有一个空行
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				result = append(result, "")
				modified = true
			}

			i++
			continue
		}

		// 在代码块内或外，移除行尾空格
		trimmed := strings.TrimRight(line, " \t")
		if trimmed != line {
			modified = true
			stats.TrailingSpaces++
		}
		result = append(result, trimmed)
		i++
	}

	return strings.Join(result, "\n"), modified
}

func fixMermaidDiagrams(content string) (string, bool) {
	modified := false

	// 统一Mermaid开始标记
	mermaidPattern := regexp.MustCompile(`(?m)^` + "```mermaid" + `\s+.*$`)
	replacement := "```mermaid"

	newContent := mermaidPattern.ReplaceAllString(content, replacement)
	if newContent != content {
		modified = true
		// 计算修复次数
		matches := mermaidPattern.FindAllString(content, -1)
		stats.MermaidFormat += len(matches)
		content = newContent
	}

	return content, modified
}

func fixExcessiveNewlines(content string) (string, bool) {
	modified := false

	// 修复：连续3个以上空行 → 2个空行
	excessiveNewlinePattern := regexp.MustCompile(`\n\s*\n\s*\n\s*\n+`)
	replacement := "\n\n\n"

	newContent := excessiveNewlinePattern.ReplaceAllString(content, replacement)
	if newContent != content {
		modified = true
		// 计算修复次数
		matches := excessiveNewlinePattern.FindAllString(content, -1)
		stats.ExcessiveNewlines += len(matches)
		content = newContent
	}

	return content, modified
}

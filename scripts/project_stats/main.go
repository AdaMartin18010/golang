package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ProjectStats 项目统计信息
type ProjectStats struct {
	TotalFiles       int
	MarkdownFiles    int
	GoFiles          int
	TotalLines       int
	TotalWords       int
	TotalChars       int
	CodeExamples     int
	BenchmarkTests   int
	READMEFiles      int
	DocsByCategory   map[string]int
	FilesByExtension map[string]int
}

func main() {
	fmt.Println("🔍 Go 1.23+ 项目统计分析")
	fmt.Println("=" + strings.Repeat("=", 70))
	fmt.Println()

	// 获取项目根目录
	projectRoot := ".."
	if len(os.Args) > 1 {
		projectRoot = os.Args[1]
	}

	stats := &ProjectStats{
		DocsByCategory:   make(map[string]int),
		FilesByExtension: make(map[string]int),
	}

	// 遍历项目目录
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏目录和某些特殊目录
		if info.IsDir() {
			dirName := info.Name()
			if strings.HasPrefix(dirName, ".") || dirName == "node_modules" || dirName == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// 统计文件
		return processFile(path, info, stats)
	})

	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	// 打印统计结果
	printStats(stats)
}

func processFile(path string, info os.FileInfo, stats *ProjectStats) error {
	stats.TotalFiles++

	// 获取文件扩展名
	ext := filepath.Ext(path)
	stats.FilesByExtension[ext]++

	// 分类统计
	switch ext {
	case ".md":
		stats.MarkdownFiles++
		if strings.Contains(strings.ToUpper(info.Name()), "README") {
			stats.READMEFiles++
		}
		return processMarkdownFile(path, stats)
	case ".go":
		stats.GoFiles++
		return processGoFile(path, stats)
	}

	return nil
}

func processMarkdownFile(path string, stats *ProjectStats) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil // 忽略读取错误
	}

	text := string(content)
	lines := strings.Split(text, "\n")
	stats.TotalLines += len(lines)

	// 统计字数（中文按字符，英文按单词）
	words := countWords(text)
	stats.TotalWords += words
	stats.TotalChars += utf8.RuneCountInString(text)

	// 检测代码示例
	if strings.Contains(text, "```go") {
		stats.CodeExamples++
	}

	// 分类统计文档
	categorizeDoc(path, stats)

	return nil
}

func processGoFile(path string, stats *ProjectStats) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	text := string(content)
	lines := strings.Split(text, "\n")
	stats.TotalLines += len(lines)

	// 检测基准测试
	if strings.Contains(text, "func Benchmark") {
		// 计算基准测试数量
		count := strings.Count(text, "func Benchmark")
		stats.BenchmarkTests += count
	}

	return nil
}

func countWords(text string) int {
	// 简单的字数统计
	// 中文字符算1个字，英文单词算1个字
	words := 0
	inWord := false

	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			// 中文字符
			words++
		} else if r == ' ' || r == '\n' || r == '\t' {
			inWord = false
		} else if !inWord && r != ' ' && r != '\n' && r != '\t' {
			words++
			inWord = true
		}
	}

	return words
}

func categorizeDoc(path string, stats *ProjectStats) {
	pathLower := strings.ToLower(path)

	categories := map[string][]string{
		"运行时优化": {"12-Go-1.23运行时优化", "runtime", "gc", "memory"},
		"工具链增强": {"13-Go-1.23工具链增强", "toolchain", "build"},
		"并发和网络": {"14-Go-1.23并发和网络", "concurrency", "network", "http"},
		"行业应用":  {"15-Go-1.23行业应用", "industry"},
		"基础文档":  {"01-go语言基础", "basics"},
	}

	for category, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(pathLower, keyword) {
				stats.DocsByCategory[category]++
				return
			}
		}
	}

	stats.DocsByCategory["其他"]++
}

func printStats(stats *ProjectStats) {
	fmt.Println("📊 总体统计")
	fmt.Println("-" + strings.Repeat("-", 70))
	fmt.Printf("  总文件数:       %d\n", stats.TotalFiles)
	fmt.Printf("  Markdown 文件:  %d\n", stats.MarkdownFiles)
	fmt.Printf("  Go 代码文件:    %d\n", stats.GoFiles)
	fmt.Printf("  README 文件:    %d\n", stats.READMEFiles)
	fmt.Printf("  总行数:         %d\n", stats.TotalLines)
	fmt.Printf("  总字数:         %d\n", stats.TotalWords)
	fmt.Println()

	fmt.Println("💻 代码统计")
	fmt.Println("-" + strings.Repeat("-", 70))
	fmt.Printf("  代码示例:       %d\n", stats.CodeExamples)
	fmt.Printf("  基准测试:       %d\n", stats.BenchmarkTests)
	fmt.Println()

	fmt.Println("📚 文档分类")
	fmt.Println("-" + strings.Repeat("-", 70))
	for category, count := range stats.DocsByCategory {
		if count > 0 {
			fmt.Printf("  %-12s %d 个文档\n", category+":", count)
		}
	}
	fmt.Println()

	fmt.Println("📁 文件类型分布")
	fmt.Println("-" + strings.Repeat("-", 70))
	for ext, count := range stats.FilesByExtension {
		if ext != "" && count > 0 {
			fmt.Printf("  %-8s %d 个文件\n", ext+":", count)
		}
	}
	fmt.Println()

	// 计算一些有趣的指标
	fmt.Println("🎯 项目指标")
	fmt.Println("-" + strings.Repeat("-", 70))
	if stats.MarkdownFiles > 0 {
		avgWords := stats.TotalWords / stats.MarkdownFiles
		fmt.Printf("  平均文档字数:   %d 字/文档\n", avgWords)
	}
	if stats.GoFiles > 0 {
		avgLines := stats.TotalLines / stats.GoFiles
		fmt.Printf("  平均代码行数:   %d 行/文件\n", avgLines)
	}
	fmt.Println()

	fmt.Println("✨ 项目质量评估")
	fmt.Println("-" + strings.Repeat("-", 70))
	assessQuality(stats)
}

func assessQuality(stats *ProjectStats) {
	score := 0
	maxScore := 100

	// 文档数量 (25分)
	if stats.MarkdownFiles >= 20 {
		score += 25
		fmt.Println("  ✅ 文档数量充足 (+25分)")
	} else if stats.MarkdownFiles >= 10 {
		score += 15
		fmt.Println("  ⚠️  文档数量一般 (+15分)")
	} else {
		score += 5
		fmt.Println("  ❌ 文档数量不足 (+5分)")
	}

	// 代码示例 (25分)
	if stats.CodeExamples >= 50 {
		score += 25
		fmt.Println("  ✅ 代码示例丰富 (+25分)")
	} else if stats.CodeExamples >= 20 {
		score += 15
		fmt.Println("  ⚠️  代码示例一般 (+15分)")
	} else {
		score += 5
		fmt.Println("  ❌ 代码示例不足 (+5分)")
	}

	// README 文件 (20分)
	if stats.READMEFiles >= 10 {
		score += 20
		fmt.Println("  ✅ README 完善 (+20分)")
	} else if stats.READMEFiles >= 5 {
		score += 10
		fmt.Println("  ⚠️  README 一般 (+10分)")
	} else {
		score += 5
		fmt.Println("  ❌ README 不足 (+5分)")
	}

	// 基准测试 (15分)
	if stats.BenchmarkTests >= 20 {
		score += 15
		fmt.Println("  ✅ 基准测试充分 (+15分)")
	} else if stats.BenchmarkTests >= 10 {
		score += 10
		fmt.Println("  ⚠️  基准测试一般 (+10分)")
	} else {
		score += 5
		fmt.Println("  ❌ 基准测试不足 (+5分)")
	}

	// 文档质量 (15分) - 基于平均字数
	if stats.MarkdownFiles > 0 {
		avgWords := stats.TotalWords / stats.MarkdownFiles
		if avgWords >= 2000 {
			score += 15
			fmt.Println("  ✅ 文档质量优秀 (+15分)")
		} else if avgWords >= 1000 {
			score += 10
			fmt.Println("  ⚠️  文档质量一般 (+10分)")
		} else {
			score += 5
			fmt.Println("  ❌ 文档质量待提升 (+5分)")
		}
	}

	fmt.Println()
	fmt.Printf("🏆 项目总分: %d/%d\n", score, maxScore)
	fmt.Println()

	// 评级
	var rating string
	var emoji string
	switch {
	case score >= 90:
		rating = "卓越 (Excellent)"
		emoji = "🏆🏆🏆"
	case score >= 75:
		rating = "优秀 (Great)"
		emoji = "🏆🏆"
	case score >= 60:
		rating = "良好 (Good)"
		emoji = "🏆"
	case score >= 50:
		rating = "及格 (Pass)"
		emoji = "👍"
	default:
		rating = "需要改进 (Needs Improvement)"
		emoji = "💪"
	}

	fmt.Printf("📈 项目评级: %s %s\n", rating, emoji)
	fmt.Println()
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// QualityChecker 文档质量检查器
type QualityChecker struct {
	errors   []string
	warnings []string
	stats    map[string]int
}

// DocumentStats 文档统计信息
type DocumentStats struct {
	Name            string
	Lines           int
	CodeBlocks      int
	MermaidDiagrams int
	InternalLinks   int
	ExternalLinks   int
	Headers         int
}

func main() {
	checker := &QualityChecker{
		stats: make(map[string]int),
	}

	docsDir := "docs"
	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".md") && strings.HasPrefix(filepath.Base(path), "architecture_") {
			checker.checkDocument(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking docs directory: %v\n", err)
		return
	}

	checker.printReport()
}

func (qc *QualityChecker) checkDocument(filePath string) {
	fmt.Printf("检查文档: %s\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		qc.addError(fmt.Sprintf("无法打开文件 %s: %v", filePath, err))
		return
	}
	defer file.Close()

	stats := &DocumentStats{Name: filepath.Base(filePath)}
	scanner := bufio.NewScanner(file)
	lineNum := 0

	var inCodeBlock bool
	var codeBlockLang string

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		stats.Lines++

		// 检查代码块
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				// 开始代码块
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				stats.CodeBlocks++

				// 检查Go代码块是否有语言标识
				if codeBlockLang == "go" {
					qc.stats["go_blocks"]++
				} else if codeBlockLang == "yaml" {
					qc.stats["yaml_blocks"]++
				} else if codeBlockLang == "mermaid" {
					stats.MermaidDiagrams++
					qc.stats["mermaid_diagrams"]++
				} else if codeBlockLang == "" {
					qc.addWarning(fmt.Sprintf("%s:%d - 代码块缺少语言标识", filePath, lineNum))
				}
			} else {
				// 结束代码块
				inCodeBlock = false
				codeBlockLang = ""
			}
		}

		// 检查标题结构
		if strings.HasPrefix(line, "#") {
			stats.Headers++
			level := 0
			for _, char := range line {
				if char == '#' {
					level++
				} else {
					break
				}
			}

			// 检查标题格式
			if level > 6 {
				qc.addError(fmt.Sprintf("%s:%d - 标题层级过深 (>6)", filePath, lineNum))
			}
		}

		// 检查内部链接
		internalLinkRegex := regexp.MustCompile(`\]\(\.\/[^)]+\.md\)`)
		if matches := internalLinkRegex.FindAllString(line, -1); len(matches) > 0 {
			stats.InternalLinks += len(matches)
			for _, match := range matches {
				linkPath := strings.TrimSuffix(strings.TrimPrefix(match, "](./"), ")")
				fullPath := filepath.Join("docs", linkPath)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					qc.addError(fmt.Sprintf("%s:%d - 断开的内部链接: %s", filePath, lineNum, linkPath))
				}
			}
		}

		// 检查外部链接
		externalLinkRegex := regexp.MustCompile(`\]\(https?://[^)]+\)`)
		if matches := externalLinkRegex.FindAllString(line, -1); len(matches) > 0 {
			stats.ExternalLinks += len(matches)
		}

		// 检查常见问题
		qc.checkCommonIssues(filePath, lineNum, line)
	}

	if err := scanner.Err(); err != nil {
		qc.addError(fmt.Sprintf("读取文件 %s 时出错: %v", filePath, err))
	}

	// 检查文档结构完整性
	qc.checkDocumentStructure(filePath, stats)

	qc.stats["total_documents"]++
	qc.stats["total_lines"] += stats.Lines
	qc.stats["total_code_blocks"] += stats.CodeBlocks
	qc.stats["total_headers"] += stats.Headers
}

func (qc *QualityChecker) checkCommonIssues(filePath string, lineNum int, line string) {
	// 检查中英文之间是否有空格
	chineseEnglishRegex := regexp.MustCompile(`[\p{Han}][a-zA-Z]|[a-zA-Z][\p{Han}]`)
	if chineseEnglishRegex.MatchString(line) {
		// 这里可以添加更细致的空格检查
	}

	// 检查常见拼写错误
	if strings.Contains(line, "K8s") && !strings.Contains(line, "Kubernetes") {
		qc.addWarning(fmt.Sprintf("%s:%d - 建议使用完整的 'Kubernetes' 而不是 'K8s'", filePath, lineNum))
	}

	// 检查是否有TODO或FIXME
	if strings.Contains(strings.ToUpper(line), "TODO") || strings.Contains(strings.ToUpper(line), "FIXME") {
		qc.addWarning(fmt.Sprintf("%s:%d - 发现待办事项: %s", filePath, lineNum, strings.TrimSpace(line)))
	}
}

func (qc *QualityChecker) checkDocumentStructure(filePath string, stats *DocumentStats) {
	basename := filepath.Base(filePath)

	// 检查核心文档是否有足够的内容
	if strings.HasPrefix(basename, "architecture_") {
		if stats.Lines < 200 {
			qc.addWarning(fmt.Sprintf("%s - 文档内容可能不够丰富 (行数: %d)", basename, stats.Lines))
		}

		if stats.CodeBlocks < 3 {
			qc.addWarning(fmt.Sprintf("%s - 代码示例可能不够充足 (代码块: %d)", basename, stats.CodeBlocks))
		}

		if stats.Headers < 8 {
			qc.addWarning(fmt.Sprintf("%s - 文档结构可能不够完整 (标题数: %d)", basename, stats.Headers))
		}
	}
}

func (qc *QualityChecker) addError(msg string) {
	qc.errors = append(qc.errors, msg)
}

func (qc *QualityChecker) addWarning(msg string) {
	qc.warnings = append(qc.warnings, msg)
}

func (qc *QualityChecker) printReport() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 Golang架构知识库质量检查报告")
	fmt.Println(strings.Repeat("=", 80))

	// 打印统计信息
	fmt.Printf("📈 统计信息:\n")
	fmt.Printf("  - 检查文档数: %d\n", qc.stats["total_documents"])
	fmt.Printf("  - 总行数: %d\n", qc.stats["total_lines"])
	fmt.Printf("  - Go代码块: %d\n", qc.stats["go_blocks"])
	fmt.Printf("  - Mermaid图表: %d\n", qc.stats["mermaid_diagrams"])
	fmt.Printf("  - 总标题数: %d\n", qc.stats["total_headers"])
	fmt.Printf("  - 总代码块: %d\n", qc.stats["total_code_blocks"])

	// 打印错误
	if len(qc.errors) > 0 {
		fmt.Printf("\n❌ 发现 %d 个错误:\n", len(qc.errors))
		for i, err := range qc.errors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
	} else {
		fmt.Println("\n✅ 未发现严重错误!")
	}

	// 打印警告
	if len(qc.warnings) > 0 {
		fmt.Printf("\n⚠️  发现 %d 个警告:\n", len(qc.warnings))
		for i, warning := range qc.warnings {
			fmt.Printf("  %d. %s\n", i+1, warning)
		}
	} else {
		fmt.Println("\n✅ 未发现警告!")
	}

	// 质量评级
	errorCount := len(qc.errors)
	warningCount := len(qc.warnings)

	fmt.Println("\n🏆 质量评级:")
	if errorCount == 0 && warningCount <= 5 {
		fmt.Println("  🌟 优秀 (Excellent) - 文档质量达到发布标准")
	} else if errorCount == 0 && warningCount <= 15 {
		fmt.Println("  👍 良好 (Good) - 文档质量良好，建议优化警告项")
	} else if errorCount <= 3 {
		fmt.Println("  🔧 需要改进 (Needs Improvement) - 建议修复错误和警告")
	} else {
		fmt.Println("  🚨 需要重构 (Needs Refactoring) - 存在较多问题，需要系统性改进")
	}

	fmt.Println(strings.Repeat("=", 80))
}

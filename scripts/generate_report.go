package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// DocumentAnalysis 文档分析结果
type DocumentAnalysis struct {
	Name              string `json:"name"`
	FilePath          string `json:"file_path"`
	Size              int64  `json:"size"`
	Lines             int    `json:"lines"`
	GoCodeBlocks      int    `json:"go_code_blocks"`
	MermaidDiagrams   int    `json:"mermaid_diagrams"`
	Headers           int    `json:"headers"`
	InternalLinks     int    `json:"internal_links"`
	ExternalLinks     int    `json:"external_links"`
	QualityScore      int    `json:"quality_score"`
	OptimizationLevel string `json:"optimization_level"`
}

// ReportGenerator 报告生成器
type ReportGenerator struct {
	docsDir   string
	outputDir string
	analyses  []*DocumentAnalysis
}

func NewReportGenerator(docsDir, outputDir string) *ReportGenerator {
	return &ReportGenerator{
		docsDir:   docsDir,
		outputDir: outputDir,
		analyses:  make([]*DocumentAnalysis, 0),
	}
}

func (rg *ReportGenerator) AnalyzeDocuments() error {
	return filepath.Walk(rg.docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".md") && strings.HasPrefix(filepath.Base(path), "architecture_") {
			analysis, err := rg.analyzeDocument(path, info)
			if err != nil {
				fmt.Printf("Failed to analyze %s: %v\n", path, err)
				return nil
			}
			rg.analyses = append(rg.analyses, analysis)
		}
		return nil
	})
}

func (rg *ReportGenerator) analyzeDocument(filePath string, info os.FileInfo) (*DocumentAnalysis, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	analysis := &DocumentAnalysis{
		Name:     filepath.Base(filePath),
		FilePath: filePath,
		Size:     info.Size(),
	}

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	codeBlockLang := ""

	for scanner.Scan() {
		line := scanner.Text()
		analysis.Lines++

		// 检查代码块
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				if codeBlockLang == "go" {
					analysis.GoCodeBlocks++
				} else if codeBlockLang == "mermaid" {
					analysis.MermaidDiagrams++
				}
			} else {
				inCodeBlock = false
				codeBlockLang = ""
			}
		}

		// 检查标题
		if strings.HasPrefix(line, "#") {
			analysis.Headers++
		}

		// 检查内部链接
		internalLinkRegex := regexp.MustCompile(`\]\(\.\/[^)]+\.md\)`)
		analysis.InternalLinks += len(internalLinkRegex.FindAllString(line, -1))

		// 检查外部链接
		externalLinkRegex := regexp.MustCompile(`\]\(https?://[^)]+\)`)
		analysis.ExternalLinks += len(externalLinkRegex.FindAllString(line, -1))
	}

	// 计算质量评分
	analysis.QualityScore = rg.calculateQualityScore(analysis)
	analysis.OptimizationLevel = rg.determineOptimizationLevel(analysis)

	return analysis, scanner.Err()
}

func (rg *ReportGenerator) calculateQualityScore(analysis *DocumentAnalysis) int {
	score := 0

	// 文档长度评分 (0-30分)
	if analysis.Lines >= 500 {
		score += 30
	} else if analysis.Lines >= 300 {
		score += 20
	} else if analysis.Lines >= 200 {
		score += 10
	}

	// Go代码示例评分 (0-25分)
	if analysis.GoCodeBlocks >= 10 {
		score += 25
	} else if analysis.GoCodeBlocks >= 5 {
		score += 15
	} else if analysis.GoCodeBlocks >= 3 {
		score += 10
	}

	// Mermaid图表评分 (0-20分)
	if analysis.MermaidDiagrams >= 5 {
		score += 20
	} else if analysis.MermaidDiagrams >= 3 {
		score += 15
	} else if analysis.MermaidDiagrams >= 1 {
		score += 10
	}

	// 结构完整性评分 (0-15分)
	if analysis.Headers >= 15 {
		score += 15
	} else if analysis.Headers >= 10 {
		score += 10
	} else if analysis.Headers >= 8 {
		score += 5
	}

	// 链接丰富度评分 (0-10分)
	totalLinks := analysis.InternalLinks + analysis.ExternalLinks
	if totalLinks >= 10 {
		score += 10
	} else if totalLinks >= 5 {
		score += 5
	}

	return score
}

func (rg *ReportGenerator) determineOptimizationLevel(analysis *DocumentAnalysis) string {
	score := analysis.QualityScore

	if score >= 85 {
		return "优秀 (Excellent)"
	} else if score >= 70 {
		return "良好 (Good)"
	} else if score >= 50 {
		return "一般 (Average)"
	} else {
		return "需要优化 (Needs Optimization)"
	}
}

func (rg *ReportGenerator) GenerateMarkdownReport() error {
	// 按质量评分排序
	sort.Slice(rg.analyses, func(i, j int) bool {
		return rg.analyses[i].QualityScore > rg.analyses[j].QualityScore
	})

	reportPath := filepath.Join(rg.outputDir, "TECHNICAL_REPORT.md")
	file, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 生成报告内容
	fmt.Fprintf(file, "# Golang架构知识库技术报告\n\n")
	fmt.Fprintf(file, "**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 总体统计
	rg.generateOverallStats(file)

	// 文档质量排名
	rg.generateQualityRanking(file)

	// 优化建议
	rg.generateOptimizationSuggestions(file)

	fmt.Printf("✅ 技术报告已生成: %s\n", reportPath)
	return nil
}

func (rg *ReportGenerator) generateOverallStats(file *os.File) {
	totalDocs := len(rg.analyses)
	totalLines := 0
	totalGoBlocks := 0
	totalMermaidDiagrams := 0
	excellentDocs := 0
	goodDocs := 0

	for _, analysis := range rg.analyses {
		totalLines += analysis.Lines
		totalGoBlocks += analysis.GoCodeBlocks
		totalMermaidDiagrams += analysis.MermaidDiagrams

		if analysis.QualityScore >= 85 {
			excellentDocs++
		} else if analysis.QualityScore >= 70 {
			goodDocs++
		}
	}

	fmt.Fprintf(file, "## 📊 总体统计\n\n")
	fmt.Fprintf(file, "| 指标 | 数值 |\n")
	fmt.Fprintf(file, "|------|------|\n")
	fmt.Fprintf(file, "| 架构文档总数 | %d |\n", totalDocs)
	fmt.Fprintf(file, "| 总行数 | %d |\n", totalLines)
	fmt.Fprintf(file, "| Go代码示例总数 | %d |\n", totalGoBlocks)
	fmt.Fprintf(file, "| Mermaid图表总数 | %d |\n", totalMermaidDiagrams)
	fmt.Fprintf(file, "| 优秀文档数量 | %d (%.1f%%) |\n", excellentDocs, float64(excellentDocs)/float64(totalDocs)*100)
	fmt.Fprintf(file, "| 良好文档数量 | %d (%.1f%%) |\n", goodDocs, float64(goodDocs)/float64(totalDocs)*100)
	fmt.Fprintf(file, "| 平均文档长度 | %d 行 |\n", totalLines/totalDocs)
	fmt.Fprintf(file, "\n")
}

func (rg *ReportGenerator) generateQualityRanking(file *os.File) {
	fmt.Fprintf(file, "## 🏆 文档质量排名 (Top 15)\n\n")
	fmt.Fprintf(file, "| 排名 | 文档名称 | 质量评分 | 优化等级 | 行数 | Go代码 | 图表 |\n")
	fmt.Fprintf(file, "|------|----------|----------|----------|------|--------|------|\n")

	topCount := 15
	if len(rg.analyses) < topCount {
		topCount = len(rg.analyses)
	}

	for i := 0; i < topCount; i++ {
		analysis := rg.analyses[i]
		fmt.Fprintf(file, "| %d | %s | %d | %s | %d | %d | %d |\n",
			i+1,
			strings.TrimSuffix(analysis.Name, ".md"),
			analysis.QualityScore,
			analysis.OptimizationLevel,
			analysis.Lines,
			analysis.GoCodeBlocks,
			analysis.MermaidDiagrams,
		)
	}
	fmt.Fprintf(file, "\n")
}

func (rg *ReportGenerator) generateOptimizationSuggestions(file *os.File) {
	fmt.Fprintf(file, "## 💡 优化建议\n\n")

	needsOptimization := make([]*DocumentAnalysis, 0)
	for _, analysis := range rg.analyses {
		if analysis.QualityScore < 70 {
			needsOptimization = append(needsOptimization, analysis)
		}
	}

	if len(needsOptimization) == 0 {
		fmt.Fprintf(file, "🎉 所有文档质量均达到良好以上标准！\n\n")
		return
	}

	fmt.Fprintf(file, "### 需要优化的文档 (%d个)\n\n", len(needsOptimization))

	for _, analysis := range needsOptimization {
		fmt.Fprintf(file, "#### %s (评分: %d)\n\n", analysis.Name, analysis.QualityScore)

		suggestions := rg.generateSpecificSuggestions(analysis)
		for _, suggestion := range suggestions {
			fmt.Fprintf(file, "- %s\n", suggestion)
		}
		fmt.Fprintf(file, "\n")
	}
}

func (rg *ReportGenerator) generateSpecificSuggestions(analysis *DocumentAnalysis) []string {
	suggestions := make([]string, 0)

	if analysis.Lines < 300 {
		suggestions = append(suggestions, "📝 建议增加文档内容深度，当前行数较少")
	}

	if analysis.GoCodeBlocks < 5 {
		suggestions = append(suggestions, "💻 建议增加更多Go代码示例，提升实用性")
	}

	if analysis.MermaidDiagrams < 3 {
		suggestions = append(suggestions, "📊 建议添加架构图和流程图，增强可视化")
	}

	if analysis.Headers < 10 {
		suggestions = append(suggestions, "🏗️ 建议完善文档结构，增加更多章节")
	}

	if analysis.InternalLinks < 3 {
		suggestions = append(suggestions, "🔗 建议添加更多内部文档链接，增强关联性")
	}

	return suggestions
}

func (rg *ReportGenerator) GenerateJSONReport() error {
	reportPath := filepath.Join(rg.outputDir, "analysis_report.json")

	report := map[string]interface{}{
		"generated_at":    time.Now().Format(time.RFC3339),
		"total_documents": len(rg.analyses),
		"analyses":        rg.analyses,
	}

	file, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return err
	}

	fmt.Printf("✅ JSON报告已生成: %s\n", reportPath)
	return nil
}

func main() {
	docsDir := "docs"
	outputDir := "docs"

	// 确保输出目录存在
	os.MkdirAll(outputDir, 0755)

	generator := NewReportGenerator(docsDir, outputDir)

	fmt.Println("🔍 正在分析文档...")
	if err := generator.AnalyzeDocuments(); err != nil {
		fmt.Printf("❌ 文档分析失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📄 已分析 %d 个架构文档\n", len(generator.analyses))

	fmt.Println("📊 正在生成Markdown报告...")
	if err := generator.GenerateMarkdownReport(); err != nil {
		fmt.Printf("❌ Markdown报告生成失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("📋 正在生成JSON报告...")
	if err := generator.GenerateJSONReport(); err != nil {
		fmt.Printf("❌ JSON报告生成失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🎉 所有报告生成完成！")
}

package report

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/your-org/formal-verifier/pkg/project"
)

// MarkdownReport Markdown格式报告生成器
type MarkdownReport struct {
	result *project.AnalysisResult
}

// NewMarkdownReport 创建Markdown报告生成器
func NewMarkdownReport(result *project.AnalysisResult) *MarkdownReport {
	return &MarkdownReport{result: result}
}

// Generate 生成Markdown报告
func (m *MarkdownReport) Generate(output string) error {
	// 创建输出文件
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// 写入报告内容
	m.writeHeader(f)
	m.writeSummary(f)
	m.writeStats(f)
	m.writeIssues(f)
	m.writeFooter(f)

	return nil
}

// writeHeader 写入报告头部
func (m *MarkdownReport) writeHeader(f *os.File) {
	fmt.Fprintln(f, "# Go Formal Verifier - 项目分析报告")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "🔍 **Go 形式化验证工具分析报告**")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f)
}

// writeSummary 写入分析摘要
func (m *MarkdownReport) writeSummary(f *os.File) {
	fmt.Fprintln(f, "## 📊 分析摘要")
	fmt.Fprintln(f)
	fmt.Fprintln(f, m.result.Summary)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f)
}

// writeStats 写入统计信息
func (m *MarkdownReport) writeStats(f *os.File) {
	fmt.Fprintln(f, "## 📈 统计信息")
	fmt.Fprintln(f)

	// 基本统计
	fmt.Fprintln(f, "### 基本统计")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "- **文件数**: %d\n", m.result.Stats.TotalFiles)
	fmt.Fprintf(f, "- **代码行数**: %d\n", m.result.Stats.TotalLines)
	fmt.Fprintf(f, "- **总问题数**: %d\n", m.result.Stats.TotalIssues)
	fmt.Fprintf(f, "- **质量评分**: %d/100 ", m.result.Stats.QualityScore)
	fmt.Fprintln(f, m.scoreEmoji(m.result.Stats.QualityScore))
	fmt.Fprintln(f)

	// 按严重程度分类
	fmt.Fprintln(f, "### 按严重程度分类")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "- ❌ **错误**: %d\n", m.result.Stats.ErrorCount)
	fmt.Fprintf(f, "- ⚠️ **警告**: %d\n", m.result.Stats.WarningCount)
	fmt.Fprintf(f, "- ℹ️ **提示**: %d\n", m.result.Stats.InfoCount)
	fmt.Fprintln(f)

	// 按类别分类
	fmt.Fprintln(f, "### 按类别分类")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "- ⚡ **并发问题**: %d\n", m.result.Stats.ConcurrencyIssues)
	fmt.Fprintf(f, "- 🔤 **类型问题**: %d\n", m.result.Stats.TypeIssues)
	fmt.Fprintf(f, "- 📊 **数据流问题**: %d\n", m.result.Stats.DataFlowIssues)
	fmt.Fprintf(f, "- ⚙️ **优化建议**: %d\n", m.result.Stats.OptimizationIssues)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f)
}

// writeIssues 写入问题详情
func (m *MarkdownReport) writeIssues(f *os.File) {
	if m.result.Stats.TotalIssues == 0 {
		fmt.Fprintln(f, "## ✅ 代码质量优秀")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "没有发现任何问题！")
		fmt.Fprintln(f)
		return
	}

	fmt.Fprintln(f, "## 🔍 问题详情")
	fmt.Fprintln(f)

	// 错误
	if m.result.Stats.ErrorCount > 0 {
		fmt.Fprintln(f, "### ❌ 错误")
		fmt.Fprintln(f)
		errors := m.result.GetIssuesBySeverity("error")
		for i, issue := range errors {
			m.writeIssue(f, issue, i+1)
		}
		fmt.Fprintln(f)
	}

	// 警告
	if m.result.Stats.WarningCount > 0 {
		fmt.Fprintln(f, "### ⚠️ 警告")
		fmt.Fprintln(f)
		warnings := m.result.GetIssuesBySeverity("warning")
		maxDisplay := 20
		for i, issue := range warnings {
			if i >= maxDisplay {
				fmt.Fprintf(f, "*... 还有 %d 个警告*\n", len(warnings)-maxDisplay)
				fmt.Fprintln(f)
				break
			}
			m.writeIssue(f, issue, i+1)
		}
		fmt.Fprintln(f)
	}

	// 提示
	if m.result.Stats.InfoCount > 0 && m.result.Stats.InfoCount <= 10 {
		fmt.Fprintln(f, "### ℹ️ 提示")
		fmt.Fprintln(f)
		infos := m.result.GetIssuesBySeverity("info")
		for i, issue := range infos {
			m.writeIssue(f, issue, i+1)
		}
		fmt.Fprintln(f)
	} else if m.result.Stats.InfoCount > 10 {
		fmt.Fprintf(f, "### ℹ️ 提示: %d 个\n", m.result.Stats.InfoCount)
		fmt.Fprintln(f)
		fmt.Fprintln(f, "*详情请查看JSON报告*")
		fmt.Fprintln(f)
	}

	fmt.Fprintln(f, "---")
	fmt.Fprintln(f)
}

// writeIssue 写入单个问题
func (m *MarkdownReport) writeIssue(f *os.File, issue project.Issue, index int) {
	fmt.Fprintf(f, "#### %d. [%s] %s\n",
		index,
		issue.Category,
		filepath.Base(issue.File))
	fmt.Fprintln(f)
	fmt.Fprintf(f, "**位置**: `%s:%d:%d`\n", issue.File, issue.Line, issue.Column)
	fmt.Fprintln(f)
	fmt.Fprintf(f, "**问题**: %s\n", issue.Message)
	fmt.Fprintln(f)

	if issue.Suggestion != "" {
		fmt.Fprintf(f, "💡 **建议**: %s\n", issue.Suggestion)
		fmt.Fprintln(f)
	}
}

// writeFooter 写入报告尾部
func (m *MarkdownReport) writeFooter(f *os.File) {
	fmt.Fprintln(f, "## 📚 关于")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "**Go Formal Verifier** - 基于 Go 1.25.3 形式化理论体系")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "### 理论基础")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "- 文档02: CSP并发模型与形式化证明")
	fmt.Fprintln(f, "- 文档03: Go类型系统形式化定义")
	fmt.Fprintln(f, "- 文档13: Go控制流形式化完整分析")
	fmt.Fprintln(f, "- 文档15: Go编译器优化形式化证明")
	fmt.Fprintln(f, "- 文档16: Go并发模式完整形式化分析")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "### 文档位置")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "`docs/01-语言基础/00-Go-1.25.3形式化理论体系/`")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "### 链接")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "- [GitHub](https://github.com/your-org/formal-verifier)")
	fmt.Fprintln(f, "- [文档](https://github.com/your-org/formal-verifier/docs)")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "---")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "*生成时间: 由 Go Formal Verifier 自动生成*")
}

// scoreEmoji 返回质量评分对应的表情
func (m *MarkdownReport) scoreEmoji(score int) string {
	if score >= 90 {
		return "✅"
	} else if score >= 70 {
		return "✓"
	} else if score >= 50 {
		return "⚠️"
	}
	return "❌"
}

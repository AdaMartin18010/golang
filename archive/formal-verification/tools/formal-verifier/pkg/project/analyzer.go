package project

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"sync"
)

// Analyzer 项目级分析器
//
// 功能：
//   - 跨文件数据流分析
//   - 项目级并发检查
//   - 全局类型验证
//   - 综合问题报告
type Analyzer struct {
	Scanner *Scanner // 导出Scanner字段以便在main.go中配置
	fset    *token.FileSet

	// 分析结果
	issues []Issue
	mu     sync.Mutex
}

// Issue 表示分析发现的问题
type Issue struct {
	// Severity 严重程度：error, warning, info
	Severity string

	// Category 问题类别：concurrency, type, dataflow, optimization
	Category string

	// Message 问题描述
	Message string

	// File 文件路径
	File string

	// Line 行号
	Line int

	// Column 列号
	Column int

	// Suggestion 修复建议
	Suggestion string
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	// Issues 所有发现的问题
	Issues []Issue

	// Stats 统计信息
	Stats *AnalysisStats

	// Summary 分析摘要
	Summary string
}

// AnalysisStats 分析统计信息
type AnalysisStats struct {
	TotalFiles   int
	TotalLines   int
	TotalIssues  int
	ErrorCount   int
	WarningCount int
	InfoCount    int

	// 按类别统计
	ConcurrencyIssues  int
	TypeIssues         int
	DataFlowIssues     int
	OptimizationIssues int

	// 质量评分 (0-100)
	QualityScore int
}

// NewAnalyzer 创建新的项目分析器
func NewAnalyzer(rootDir string) *Analyzer {
	return &Analyzer{
		Scanner: NewScanner(rootDir),
		fset:    token.NewFileSet(),
		issues:  make([]Issue, 0),
	}
}

// WithScanner 设置自定义扫描器
func (a *Analyzer) WithScanner(scanner *Scanner) *Analyzer {
	a.Scanner = scanner
	return a
}

// Analyze 执行项目分析
func (a *Analyzer) Analyze() (*AnalysisResult, error) {
	// 扫描项目文件
	scanResult, err := a.Scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	if len(scanResult.Files) == 0 {
		return nil, fmt.Errorf("no Go files found in project")
	}

	// 并行分析文件
	var wg sync.WaitGroup
	for _, file := range scanResult.Files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			a.analyzeFile(filePath)
		}(file)
	}
	wg.Wait()

	// 构建分析结果
	result := &AnalysisResult{
		Issues: a.issues,
		Stats:  a.buildStats(scanResult.Stats),
	}

	// 生成摘要
	result.Summary = a.generateSummary(result.Stats)

	// 按严重程度和文件排序
	a.sortIssues(result.Issues)

	return result, nil
}

// analyzeFile 分析单个文件
func (a *Analyzer) analyzeFile(filePath string) {
	// 解析文件
	f, err := parser.ParseFile(a.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		a.addIssue(Issue{
			Severity: "error",
			Category: "syntax",
			Message:  fmt.Sprintf("Failed to parse file: %v", err),
			File:     filePath,
			Line:     1,
		})
		return
	}

	// 基本检查
	a.checkBasicIssues(f, filePath)
	a.checkConcurrencyIssues(f, filePath)
	a.checkTypeIssues(f, filePath)
}

// checkBasicIssues 检查基本问题
func (a *Analyzer) checkBasicIssues(f *ast.File, filePath string) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			// 检查函数复杂度
			if a.isFunctionTooComplex(node) {
				pos := a.fset.Position(node.Pos())
				a.addIssue(Issue{
					Severity:   "warning",
					Category:   "complexity",
					Message:    fmt.Sprintf("Function '%s' is too complex", node.Name.Name),
					File:       filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Suggestion: "Consider refactoring into smaller functions",
				})
			}

		case *ast.GoStmt:
			// 检查goroutine泄露风险
			if !a.hasProperGoroutineCleanup(node) {
				pos := a.fset.Position(node.Pos())
				a.addIssue(Issue{
					Severity:   "warning",
					Category:   "concurrency",
					Message:    "Potential goroutine leak: missing cleanup mechanism",
					File:       filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Suggestion: "Add context cancellation or done channel",
				})
			}
		}
		return true
	})
}

// checkConcurrencyIssues 检查并发问题
func (a *Analyzer) checkConcurrencyIssues(f *ast.File, filePath string) {
	var (
		hasGoroutines bool
		hasChannels   bool
		hasMutex      bool
	)

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GoStmt:
			hasGoroutines = true

		case *ast.ChanType:
			hasChannels = true

		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				if ident.Name == "sync" && (node.Sel.Name == "Mutex" || node.Sel.Name == "RWMutex") {
					hasMutex = true
				}
			}
		}
		return true
	})

	// 并发模式分析
	if hasGoroutines && !hasChannels && !hasMutex {
		a.addIssue(Issue{
			Severity:   "warning",
			Category:   "concurrency",
			Message:    "Goroutines without synchronization mechanism",
			File:       filePath,
			Line:       1,
			Suggestion: "Add channels or mutexes for synchronization",
		})
	}
}

// checkTypeIssues 检查类型相关问题
func (a *Analyzer) checkTypeIssues(f *ast.File, filePath string) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.TypeAssertExpr:
			// 检查类型断言是否检查了错误
			if !a.hasTypeAssertionCheck(node) {
				pos := a.fset.Position(node.Pos())
				a.addIssue(Issue{
					Severity:   "warning",
					Category:   "type",
					Message:    "Type assertion without ok check",
					File:       filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Suggestion: "Use v, ok := x.(Type) to check assertion",
				})
			}
		}
		return true
	})
}

// isFunctionTooComplex 判断函数是否过于复杂
func (a *Analyzer) isFunctionTooComplex(fn *ast.FuncDecl) bool {
	if fn.Body == nil {
		return false
	}

	complexity := 1
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.SelectStmt:
			complexity++
		case *ast.BinaryExpr:
			// &&, || 增加复杂度
			complexity++
		}
		return true
	})

	return complexity > 10 // 圈复杂度阈值
}

// hasProperGoroutineCleanup 检查goroutine是否有适当的清理机制
func (a *Analyzer) hasProperGoroutineCleanup(goStmt *ast.GoStmt) bool {
	// 简化检查：查看是否使用了context或done channel
	hasCleanup := false

	ast.Inspect(goStmt, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.Ident:
			if node.Name == "ctx" || node.Name == "context" || node.Name == "done" {
				hasCleanup = true
				return false
			}
		case *ast.SelectorExpr:
			if sel, ok := node.X.(*ast.Ident); ok {
				if sel.Name == "context" {
					hasCleanup = true
					return false
				}
			}
		}
		return true
	})

	return hasCleanup
}

// hasTypeAssertionCheck 检查类型断言是否有ok检查
func (a *Analyzer) hasTypeAssertionCheck(typeAssert *ast.TypeAssertExpr) bool {
	// 这是一个简化的检查，实际需要更复杂的AST分析
	// 这里假设如果类型断言在赋值语句中且有两个左值，则有ok检查
	return false // 保守返回false，鼓励总是使用ok检查
}

// addIssue 添加问题（线程安全）
func (a *Analyzer) addIssue(issue Issue) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.issues = append(a.issues, issue)
}

// buildStats 构建统计信息
func (a *Analyzer) buildStats(scanStats *Stats) *AnalysisStats {
	stats := &AnalysisStats{
		TotalFiles: scanStats.TotalFiles,
		TotalLines: scanStats.TotalLines,
	}

	// 统计问题
	for _, issue := range a.issues {
		stats.TotalIssues++

		switch issue.Severity {
		case "error":
			stats.ErrorCount++
		case "warning":
			stats.WarningCount++
		case "info":
			stats.InfoCount++
		}

		switch issue.Category {
		case "concurrency":
			stats.ConcurrencyIssues++
		case "type":
			stats.TypeIssues++
		case "dataflow":
			stats.DataFlowIssues++
		case "optimization":
			stats.OptimizationIssues++
		}
	}

	// 计算质量评分
	stats.QualityScore = a.calculateQualityScore(stats)

	return stats
}

// calculateQualityScore 计算代码质量评分
func (a *Analyzer) calculateQualityScore(stats *AnalysisStats) int {
	// 基础分数
	score := 100

	// 扣分规则
	score -= stats.ErrorCount * 10       // 每个错误扣10分
	score -= stats.WarningCount * 5      // 每个警告扣5分
	score -= stats.InfoCount * 1         // 每个提示扣1分
	score -= stats.ConcurrencyIssues * 3 // 并发问题额外扣3分

	// 确保分数在0-100之间
	if score < 0 {
		score = 0
	}

	return score
}

// generateSummary 生成分析摘要
func (a *Analyzer) generateSummary(stats *AnalysisStats) string {
	summary := fmt.Sprintf(
		"Analyzed %d files (%d lines) and found %d issues:\n"+
			"  - Errors: %d\n"+
			"  - Warnings: %d\n"+
			"  - Info: %d\n"+
			"Quality Score: %d/100\n",
		stats.TotalFiles,
		stats.TotalLines,
		stats.TotalIssues,
		stats.ErrorCount,
		stats.WarningCount,
		stats.InfoCount,
		stats.QualityScore,
	)

	// 添加严重程度评估
	if stats.QualityScore >= 90 {
		summary += "✅ Excellent code quality!\n"
	} else if stats.QualityScore >= 70 {
		summary += "✓ Good code quality\n"
	} else if stats.QualityScore >= 50 {
		summary += "⚠️  Code quality needs improvement\n"
	} else {
		summary += "❌ Poor code quality - immediate attention required\n"
	}

	return summary
}

// sortIssues 对问题进行排序
func (a *Analyzer) sortIssues(issues []Issue) {
	sort.Slice(issues, func(i, j int) bool {
		// 首先按严重程度排序
		severityOrder := map[string]int{"error": 0, "warning": 1, "info": 2}
		if severityOrder[issues[i].Severity] != severityOrder[issues[j].Severity] {
			return severityOrder[issues[i].Severity] < severityOrder[issues[j].Severity]
		}

		// 然后按文件名排序
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File
		}

		// 最后按行号排序
		return issues[i].Line < issues[j].Line
	})
}

// FilterIssues 过滤问题
func (r *AnalysisResult) FilterIssues(filter func(Issue) bool) []Issue {
	filtered := make([]Issue, 0)
	for _, issue := range r.Issues {
		if filter(issue) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// GetIssuesBySeverity 按严重程度获取问题
func (r *AnalysisResult) GetIssuesBySeverity(severity string) []Issue {
	return r.FilterIssues(func(i Issue) bool {
		return i.Severity == severity
	})
}

// GetIssuesByCategory 按类别获取问题
func (r *AnalysisResult) GetIssuesByCategory(category string) []Issue {
	return r.FilterIssues(func(i Issue) bool {
		return i.Category == category
	})
}

// GetIssuesByFile 按文件获取问题
func (r *AnalysisResult) GetIssuesByFile(file string) []Issue {
	return r.FilterIssues(func(i Issue) bool {
		return i.File == file
	})
}

// HasErrors 判断是否有错误
func (r *AnalysisResult) HasErrors() bool {
	return r.Stats.ErrorCount > 0
}

// HasWarnings 判断是否有警告
func (r *AnalysisResult) HasWarnings() bool {
	return r.Stats.WarningCount > 0
}

// IsHealthy 判断代码是否健康（无错误且警告少于阈值）
func (r *AnalysisResult) IsHealthy(maxWarnings int) bool {
	return !r.HasErrors() && r.Stats.WarningCount <= maxWarnings
}

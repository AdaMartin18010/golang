package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer(".")
	if analyzer == nil {
		t.Fatal("NewAnalyzer returned nil")
	}

	if analyzer.scanner == nil {
		t.Error("Expected scanner to be initialized")
	}

	if analyzer.fset == nil {
		t.Error("Expected fset to be initialized")
	}
}

func TestAnalyze(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

func complexFunction(x int) int {
	if x > 0 {
		if x > 10 {
			if x > 20 {
				if x > 30 {
					if x > 40 {
						if x > 50 {
							return x * 2
						}
					}
				}
			}
		}
	}
	return x
}
`

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	analyzer := NewAnalyzer(tmpDir)
	result, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Stats == nil {
		t.Fatal("Expected stats, got nil")
	}

	if result.Stats.TotalFiles != 1 {
		t.Errorf("Expected 1 file, got %d", result.Stats.TotalFiles)
	}

	// 应该检测到复杂函数
	hasComplexityIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "complexity" {
			hasComplexityIssue = true
			break
		}
	}

	if !hasComplexityIssue {
		t.Error("Expected to find complexity issue")
	}
}

func TestAnalyzeWithGoroutines(t *testing.T) {
	tmpDir := t.TempDir()

	// 代码包含goroutine但没有适当的清理机制
	testCode := `package main

func main() {
	go func() {
		// This goroutine has no cleanup mechanism
		for {
			// do something
		}
	}()
}
`

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	analyzer := NewAnalyzer(tmpDir)
	result, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 应该检测到goroutine泄露风险
	hasGoroutineIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "concurrency" {
			hasGoroutineIssue = true
			break
		}
	}

	if !hasGoroutineIssue {
		t.Error("Expected to find goroutine leak issue")
	}
}

func TestAnalyzeWithTypeAssertion(t *testing.T) {
	tmpDir := t.TempDir()

	testCode := `package main

func main() {
	var x interface{} = 42
	_ = x.(int) // Type assertion without ok check
}
`

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	analyzer := NewAnalyzer(tmpDir)
	result, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 应该检测到类型断言问题
	hasTypeIssue := false
	for _, issue := range result.Issues {
		if issue.Category == "type" {
			hasTypeIssue = true
			break
		}
	}

	if !hasTypeIssue {
		t.Error("Expected to find type assertion issue")
	}
}

func TestAnalyzeEmptyProject(t *testing.T) {
	tmpDir := t.TempDir()

	analyzer := NewAnalyzer(tmpDir)
	_, err := analyzer.Analyze()
	if err == nil {
		t.Error("Expected error for empty project")
	}
}

func TestAnalyzeInvalidSyntax(t *testing.T) {
	tmpDir := t.TempDir()

	invalidCode := `package main

func main( {
	// Invalid syntax - missing closing parenthesis
}
`

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(invalidCode), 0644); err != nil {
		t.Fatal(err)
	}

	analyzer := NewAnalyzer(tmpDir)
	result, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 应该有语法错误
	hasSyntaxError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && issue.Category == "syntax" {
			hasSyntaxError = true
			break
		}
	}

	if !hasSyntaxError {
		t.Error("Expected to find syntax error")
	}
}

func TestFilterIssues(t *testing.T) {
	result := &AnalysisResult{
		Issues: []Issue{
			{Severity: "error", Category: "type", File: "a.go"},
			{Severity: "warning", Category: "concurrency", File: "b.go"},
			{Severity: "info", Category: "optimization", File: "c.go"},
		},
		Stats: &AnalysisStats{},
	}

	// 测试按严重程度过滤
	errors := result.GetIssuesBySeverity("error")
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	// 测试按类别过滤
	concurrency := result.GetIssuesByCategory("concurrency")
	if len(concurrency) != 1 {
		t.Errorf("Expected 1 concurrency issue, got %d", len(concurrency))
	}

	// 测试按文件过滤
	fileIssues := result.GetIssuesByFile("a.go")
	if len(fileIssues) != 1 {
		t.Errorf("Expected 1 issue in a.go, got %d", len(fileIssues))
	}
}

func TestHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		result   *AnalysisResult
		hasError bool
	}{
		{
			name: "with errors",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 1},
			},
			hasError: true,
		},
		{
			name: "without errors",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 0, WarningCount: 5},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasErrors(); got != tt.hasError {
				t.Errorf("HasErrors() = %v, want %v", got, tt.hasError)
			}
		})
	}
}

func TestIsHealthy(t *testing.T) {
	tests := []struct {
		name        string
		result      *AnalysisResult
		maxWarnings int
		want        bool
	}{
		{
			name: "healthy with no issues",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 0, WarningCount: 0},
			},
			maxWarnings: 5,
			want:        true,
		},
		{
			name: "healthy with few warnings",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 0, WarningCount: 3},
			},
			maxWarnings: 5,
			want:        true,
		},
		{
			name: "unhealthy with errors",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 1, WarningCount: 0},
			},
			maxWarnings: 5,
			want:        false,
		},
		{
			name: "unhealthy with too many warnings",
			result: &AnalysisResult{
				Stats: &AnalysisStats{ErrorCount: 0, WarningCount: 10},
			},
			maxWarnings: 5,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsHealthy(tt.maxWarnings); got != tt.want {
				t.Errorf("IsHealthy(%d) = %v, want %v", tt.maxWarnings, got, tt.want)
			}
		})
	}
}

func TestCalculateQualityScore(t *testing.T) {
	analyzer := NewAnalyzer(".")

	tests := []struct {
		name      string
		stats     *AnalysisStats
		wantRange [2]int // min, max
	}{
		{
			name: "perfect score",
			stats: &AnalysisStats{
				ErrorCount:   0,
				WarningCount: 0,
				InfoCount:    0,
			},
			wantRange: [2]int{100, 100},
		},
		{
			name: "with errors",
			stats: &AnalysisStats{
				ErrorCount:   2,
				WarningCount: 0,
				InfoCount:    0,
			},
			wantRange: [2]int{80, 80},
		},
		{
			name: "with warnings",
			stats: &AnalysisStats{
				ErrorCount:   0,
				WarningCount: 5,
				InfoCount:    0,
			},
			wantRange: [2]int{75, 75},
		},
		{
			name: "poor quality",
			stats: &AnalysisStats{
				ErrorCount:        5,
				WarningCount:      10,
				InfoCount:         10,
				ConcurrencyIssues: 5,
			},
			wantRange: [2]int{0, 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := analyzer.calculateQualityScore(tt.stats)
			if score < tt.wantRange[0] || score > tt.wantRange[1] {
				t.Errorf("calculateQualityScore() = %d, want range [%d, %d]",
					score, tt.wantRange[0], tt.wantRange[1])
			}
		})
	}
}

func TestGenerateSummary(t *testing.T) {
	analyzer := NewAnalyzer(".")

	stats := &AnalysisStats{
		TotalFiles:   10,
		TotalLines:   1000,
		TotalIssues:  15,
		ErrorCount:   2,
		WarningCount: 8,
		InfoCount:    5,
		QualityScore: 75,
	}

	summary := analyzer.generateSummary(stats)

	if summary == "" {
		t.Error("Expected non-empty summary")
	}

	// 检查摘要是否包含关键信息
	expectedStrings := []string{
		"10 files",
		"1000 lines",
		"15 issues",
		"Errors: 2",
		"Warnings: 8",
		"Info: 5",
		"Quality Score: 75/100",
	}

	for _, expected := range expectedStrings {
		if !contains(summary, expected) {
			t.Errorf("Summary missing expected string: %s", expected)
		}
	}
}

func TestSortIssues(t *testing.T) {
	analyzer := NewAnalyzer(".")

	issues := []Issue{
		{Severity: "info", File: "b.go", Line: 10},
		{Severity: "error", File: "a.go", Line: 5},
		{Severity: "warning", File: "a.go", Line: 10},
		{Severity: "error", File: "a.go", Line: 1},
	}

	analyzer.sortIssues(issues)

	// 验证排序顺序：错误优先，然后按文件名，最后按行号
	if issues[0].Severity != "error" || issues[0].Line != 1 {
		t.Error("First issue should be error at line 1")
	}

	if issues[1].Severity != "error" || issues[1].Line != 5 {
		t.Error("Second issue should be error at line 5")
	}

	if issues[2].Severity != "warning" {
		t.Error("Third issue should be warning")
	}

	if issues[3].Severity != "info" {
		t.Error("Fourth issue should be info")
	}
}

func BenchmarkAnalyze(b *testing.B) {
	// 创建临时测试目录
	tmpDir := b.TempDir()

	testCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

func helper() {
	// some code
}
`

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(testCode), 0644); err != nil {
		b.Fatal(err)
	}

	analyzer := NewAnalyzer(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.Analyze()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

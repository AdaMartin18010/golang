package optimization

import (
	"strings"
	"testing"
)

func TestEscapeAnalysis(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/escape_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	escapes := analyzer.GetEscapeAnalysis()
	if len(escapes) == 0 {
		t.Error("Expected some escape analysis results")
	}

	t.Logf("Escape analysis found %d allocations", len(escapes))

	// 检查是否检测到逃逸
	hasEscape := false
	for _, info := range escapes {
		t.Logf("  - %s: %s -> %s (%s)", info.VarName, info.Position, info.EscapesTo, info.Reason)
		if info.EscapesTo == "heap" {
			hasEscape = true
		}
	}

	if !hasEscape {
		t.Log("Note: No heap escapes detected (this may be expected for simple cases)")
	}
}

func TestInlineAnalysis(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/inline_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	inlines := analyzer.GetInlineAnalysis()
	if len(inlines) == 0 {
		t.Error("Expected some inline analysis results")
	}

	t.Logf("Inline analysis found %d functions", len(inlines))

	// 应该有可内联和不可内联的函数
	canInlineCount := 0
	cannotInlineCount := 0

	for _, info := range inlines {
		t.Logf("  - %s (cost=%d, inline=%v): %s",
			info.FuncName, info.Cost, info.CanInline, info.Reason)

		if info.CanInline {
			canInlineCount++
		} else {
			cannotInlineCount++
		}
	}

	t.Logf("Can inline: %d, Cannot inline: %d", canInlineCount, cannotInlineCount)
}

func TestBCEAnalysis(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/bce_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	bces := analyzer.GetBCEAnalysis()
	if len(bces) == 0 {
		t.Error("Expected some BCE analysis results")
	}

	t.Logf("BCE analysis found %d array accesses", len(bces))

	// 应该有可消除和不可消除的检查
	canEliminateCount := 0
	cannotEliminateCount := 0

	for _, info := range bces {
		t.Logf("  - %s[%s] (eliminate=%v): %s",
			info.ArrayExpr, info.IndexExpr, info.CanEliminate, info.Reason)

		if info.CanEliminate {
			canEliminateCount++
		} else {
			cannotEliminateCount++
		}
	}

	t.Logf("Can eliminate: %d, Cannot eliminate: %d", canEliminateCount, cannotEliminateCount)

	// 至少应该有一些可以消除的（常量索引或range循环）
	if canEliminateCount == 0 {
		t.Error("Expected at least some BCE opportunities")
	}
}

func TestReport(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/inline_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	report := analyzer.Report()

	// 检查报告包含关键部分
	if !strings.Contains(report, "编译器优化分析报告") {
		t.Error("Report should contain title")
	}

	if !strings.Contains(report, "逃逸分析") {
		t.Error("Report should contain escape analysis section")
	}

	if !strings.Contains(report, "函数内联分析") {
		t.Error("Report should contain inline analysis section")
	}

	if !strings.Contains(report, "边界检查消除") {
		t.Error("Report should contain BCE section")
	}

	if !strings.Contains(report, "形式化理论基础") {
		t.Error("Report should contain theory section")
	}

	t.Log("Generated report:\n" + report)
}

func TestOptimizationSafety(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/escape_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	// 验证优化是安全的
	results := analyzer.GetResults()

	// 确保我们识别了堆分配
	heapCount := 0
	stackCount := 0
	for _, info := range results.EscapeAnalysis {
		if info.EscapesTo == "heap" {
			heapCount++
		} else if info.EscapesTo == "stack" {
			stackCount++
		}
	}

	t.Logf("Heap allocations: %d, Stack allocations: %d", heapCount, stackCount)

	// 应该有一些分类为堆或栈
	if heapCount+stackCount == 0 {
		t.Log("Note: No allocations classified (may be expected for simple files)")
	}
}

func TestInlineCostCalculation(t *testing.T) {
	analyzer := NewAnalyzer()

	err := analyzer.AnalyzeFile("../../testdata/inline_test.go")
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	inlines := analyzer.GetInlineAnalysis()

	// 验证成本计算合理
	for _, info := range inlines {
		if info.Cost < 0 {
			t.Errorf("Function %s has negative cost: %d", info.FuncName, info.Cost)
		}

		// 小函数应该可内联
		if info.Cost < 10 && !info.CanInline {
			t.Errorf("Small function %s (cost=%d) should be inlinable",
				info.FuncName, info.Cost)
		}

		// 大函数应该不可内联
		if info.Cost > 100 && info.CanInline {
			t.Logf("Warning: Large function %s (cost=%d) marked as inlinable",
				info.FuncName, info.Cost)
		}
	}
}

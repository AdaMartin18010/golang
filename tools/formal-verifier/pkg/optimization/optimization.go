// Package optimization implements compiler optimization analysis
// based on formal verification theory.
//
// 理论基础: 文档15 Go编译器优化形式化证明
package optimization

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// =============================================================================
// 数据结构
// =============================================================================

// OptimizerAnalyzer 执行编译器优化分析
type OptimizerAnalyzer struct {
	fset    *token.FileSet
	results *AnalysisResults
}

// AnalysisResults 包含所有优化分析结果
type AnalysisResults struct {
	EscapeAnalysis   []*EscapeInfo
	InlineAnalysis   []*InlineInfo
	BCEAnalysis      []*BCEInfo
	TotalFunctions   int
	TotalAllocations int
}

// EscapeInfo 记录逃逸分析结果
type EscapeInfo struct {
	Position    token.Position
	VarName     string
	EscapesTo   string // "heap", "stack", "argument"
	Reason      string
	CanOptimize bool
}

// InlineInfo 记录内联分析结果
type InlineInfo struct {
	Position  token.Position
	FuncName  string
	Cost      int
	CanInline bool
	Reason    string
}

// BCEInfo 记录边界检查消除分析结果
type BCEInfo struct {
	Position     token.Position
	ArrayExpr    string
	IndexExpr    string
	CanEliminate bool
	Reason       string
}

// =============================================================================
// 构造函数
// =============================================================================

// NewAnalyzer 创建新的优化分析器
func NewAnalyzer() *OptimizerAnalyzer {
	return &OptimizerAnalyzer{
		fset: token.NewFileSet(),
		results: &AnalysisResults{
			EscapeAnalysis: make([]*EscapeInfo, 0),
			InlineAnalysis: make([]*InlineInfo, 0),
			BCEAnalysis:    make([]*BCEInfo, 0),
		},
	}
}

// =============================================================================
// 主分析函数
// =============================================================================

// AnalyzeFile 分析指定的Go源文件
func (oa *OptimizerAnalyzer) AnalyzeFile(filename string) error {
	// 直接解析单个文件
	file, err := parser.ParseFile(oa.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// 执行三个优化分析
	oa.analyzeEscape(file)
	oa.analyzeInline(file)
	oa.analyzeBCE(file)

	return nil
}

// =============================================================================
// 逃逸分析 (Escape Analysis)
// =============================================================================

// analyzeEscape 执行逃逸分析
//
// 形式化定义：
//
//	逃逸分析判定对象的分配位置(栈或堆)
//	Escape(obj) = {
//	  "heap"  if obj escapes function scope
//	  "stack" if obj lifetime ⊆ function lifetime
//	}
//
// 理论基础：
//
//	obj escapes ⟺
//	  ∃ reference to obj that outlives function OR
//	  obj stored in heap location
func (oa *OptimizerAnalyzer) analyzeEscape(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			oa.results.TotalFunctions++
			oa.analyzeEscapeInFunc(node)
		}
		return true
	})
}

func (oa *OptimizerAnalyzer) analyzeEscapeInFunc(fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		// 分析变量声明
		case *ast.AssignStmt:
			for i, rhs := range node.Rhs {
				oa.analyzeEscapeExpr(rhs, node.Lhs[i])
			}

		// 分析return语句 - 返回的局部变量逃逸
		case *ast.ReturnStmt:
			for _, result := range node.Results {
				if ident, ok := result.(*ast.Ident); ok {
					pos := oa.fset.Position(ident.Pos())
					oa.results.EscapeAnalysis = append(oa.results.EscapeAnalysis, &EscapeInfo{
						Position:    pos,
						VarName:     ident.Name,
						EscapesTo:   "heap",
						Reason:      "returned from function",
						CanOptimize: false,
					})
				}
			}

		// 分析闭包 - 闭包捕获的变量逃逸
		case *ast.FuncLit:
			// 检查闭包内引用的外部变量
			capturedVars := make(map[string]bool)
			ast.Inspect(node.Body, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if ident.Obj != nil && ident.Obj.Kind == ast.Var {
						capturedVars[ident.Name] = true
					}
				}
				return true
			})

			for varName := range capturedVars {
				pos := oa.fset.Position(node.Pos())
				oa.results.EscapeAnalysis = append(oa.results.EscapeAnalysis, &EscapeInfo{
					Position:    pos,
					VarName:     varName,
					EscapesTo:   "heap",
					Reason:      "captured by closure",
					CanOptimize: false,
				})
			}
		}
		return true
	})
}

func (oa *OptimizerAnalyzer) analyzeEscapeExpr(expr ast.Expr, lhs ast.Expr) {
	switch e := expr.(type) {
	// new() 和 make() 分配
	case *ast.CallExpr:
		if ident, ok := e.Fun.(*ast.Ident); ok {
			if ident.Name == "new" {
				oa.results.TotalAllocations++

				var varName string
				if lhsIdent, ok := lhs.(*ast.Ident); ok {
					varName = lhsIdent.Name
				} else {
					varName = "temp"
				}

				pos := oa.fset.Position(e.Pos())

				// 简单启发式: 如果直接赋值给局部变量且不再使用，可能可以栈分配
				canOptimize := lhs != nil && isLocalVar(lhs)
				escapesTo := "heap"
				reason := "heap allocation"

				if canOptimize {
					escapesTo = "stack"
					reason = "local variable, may optimize to stack"
				}

				oa.results.EscapeAnalysis = append(oa.results.EscapeAnalysis, &EscapeInfo{
					Position:    pos,
					VarName:     varName,
					EscapesTo:   escapesTo,
					Reason:      reason,
					CanOptimize: canOptimize,
				})
			} else if ident.Name == "make" {
				oa.results.TotalAllocations++

				var varName string
				if lhsIdent, ok := lhs.(*ast.Ident); ok {
					varName = lhsIdent.Name
				} else {
					varName = "slice/map/chan"
				}

				pos := oa.fset.Position(e.Pos())
				oa.results.EscapeAnalysis = append(oa.results.EscapeAnalysis, &EscapeInfo{
					Position:    pos,
					VarName:     varName,
					EscapesTo:   "heap",
					Reason:      "make allocation",
					CanOptimize: false,
				})
			}
		}

	// 复合字面量 (struct, slice, map)
	case *ast.CompositeLit:
		var varName string
		if lhsIdent, ok := lhs.(*ast.Ident); ok {
			varName = lhsIdent.Name
		} else {
			varName = "literal"
		}

		pos := oa.fset.Position(e.Pos())
		canOptimize := lhs != nil && isLocalVar(lhs)
		escapesTo := "heap"
		reason := "composite literal"

		if canOptimize {
			escapesTo = "stack"
			reason = "local composite literal, may optimize to stack"
		}

		oa.results.EscapeAnalysis = append(oa.results.EscapeAnalysis, &EscapeInfo{
			Position:    pos,
			VarName:     varName,
			EscapesTo:   escapesTo,
			Reason:      reason,
			CanOptimize: canOptimize,
		})
	}
}

func isLocalVar(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Obj != nil && ident.Obj.Kind == ast.Var
	}
	return false
}

// =============================================================================
// 函数内联分析 (Function Inlining Analysis)
// =============================================================================

// analyzeInline 执行函数内联分析
//
// 形式化定义：
//
//	InlineCost(f) = Σ(instruction weights)
//	CanInline(f) ⟺
//	  InlineCost(f) < threshold ∧
//	  ¬isRecursive(f) ∧
//	  ¬hasComplexControl(f)
//
// Go编译器内联阈值：
//   - 小函数: cost < 80
//   - 叶子函数: 优先级高
//   - 调用开销: 20
//   - 循环惩罚: 30
func (oa *OptimizerAnalyzer) analyzeInline(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			oa.analyzeInlineFunc(fn)
		}
		return true
	})
}

func (oa *OptimizerAnalyzer) analyzeInlineFunc(fn *ast.FuncDecl) {
	if fn.Body == nil {
		return
	}

	funcName := fn.Name.Name
	pos := oa.fset.Position(fn.Pos())

	// 计算内联成本
	cost := oa.computeInlineCost(fn.Body)

	// 内联决策
	const (
		InlineMaxBudget = 80 // 指令预算
		InlineSmallFunc = 10 // 小函数阈值
	)

	canInline := false
	reason := ""

	// 检查递归
	if oa.isRecursive(fn) {
		reason = "recursive function"
	} else if oa.hasComplexControl(fn.Body) {
		reason = "complex control flow"
	} else if cost > InlineMaxBudget {
		reason = fmt.Sprintf("cost too high (%d > %d)", cost, InlineMaxBudget)
	} else if cost < InlineSmallFunc {
		canInline = true
		reason = "small function, always inline"
	} else if oa.isLeafFunction(fn.Body) {
		canInline = true
		reason = "leaf function"
	} else {
		canInline = true
		reason = fmt.Sprintf("cost acceptable (%d < %d)", cost, InlineMaxBudget)
	}

	oa.results.InlineAnalysis = append(oa.results.InlineAnalysis, &InlineInfo{
		Position:  pos,
		FuncName:  funcName,
		Cost:      cost,
		CanInline: canInline,
		Reason:    reason,
	})
}

func (oa *OptimizerAnalyzer) computeInlineCost(body *ast.BlockStmt) int {
	cost := 0

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt, *ast.BinaryExpr:
			cost += 1 // 基本指令

		case *ast.CallExpr:
			cost += 20 // 调用开销

		case *ast.ForStmt, *ast.RangeStmt:
			cost += 30 // 循环惩罚

		case *ast.IfStmt:
			cost += 2 // 分支开销

		case *ast.SwitchStmt:
			cost += 5 // switch开销

		case *ast.DeferStmt:
			cost += 10 // defer开销

		case *ast.GoStmt:
			cost += 15 // goroutine开销

		case *ast.FuncLit:
			// 闭包很昂贵
			cost += 50

		case *ast.ReturnStmt:
			if len(node.Results) > 0 {
				cost += 1
			}
		}
		return true
	})

	return cost
}

func (oa *OptimizerAnalyzer) isRecursive(fn *ast.FuncDecl) bool {
	funcName := fn.Name.Name
	isRecursive := false

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == funcName {
					isRecursive = true
					return false
				}
			}
		}
		return true
	})

	return isRecursive
}

func (oa *OptimizerAnalyzer) hasComplexControl(body *ast.BlockStmt) bool {
	hasComplex := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.LabeledStmt, *ast.BranchStmt:
			// goto, break label等
			hasComplex = true
			return false
		}
		return true
	})

	return hasComplex
}

func (oa *OptimizerAnalyzer) isLeafFunction(body *ast.BlockStmt) bool {
	hasCalls := false

	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.CallExpr); ok {
			hasCalls = true
			return false
		}
		return true
	})

	return !hasCalls
}

// =============================================================================
// 边界检查消除 (Bounds Check Elimination, BCE)
// =============================================================================

// analyzeBCE 执行边界检查消除分析
//
// 形式化定义：
//
//	CanEliminate(a[i]) ⟺
//	  编译器能证明: 0 ≤ i < len(a)
//
// 方法：
//  1. 常量索引
//  2. range循环
//  3. 条件保护
//  4. 重复访问
func (oa *OptimizerAnalyzer) analyzeBCE(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Body != nil {
				oa.analyzeBCEInFunc(node.Body)
			}
		}
		return true
	})
}

func (oa *OptimizerAnalyzer) analyzeBCEInFunc(body *ast.BlockStmt) {
	// 跟踪已检查的数组访问
	checkedAccesses := make(map[string]bool)

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		// 分析数组/切片索引
		case *ast.IndexExpr:
			oa.analyzeBCEIndexExpr(node, checkedAccesses)

		// 分析range循环
		case *ast.RangeStmt:
			oa.analyzeBCERange(node)
		}
		return true
	})
}

func (oa *OptimizerAnalyzer) analyzeBCEIndexExpr(node *ast.IndexExpr, checkedAccesses map[string]bool) {
	pos := oa.fset.Position(node.Pos())

	arrayExpr := exprToString(node.X)
	indexExpr := exprToString(node.Index)

	accessKey := arrayExpr + "[" + indexExpr + "]"

	canEliminate := false
	reason := "unknown"

	// 1. 常量索引
	if lit, ok := node.Index.(*ast.BasicLit); ok && lit.Kind == token.INT {
		canEliminate = true
		reason = "constant index"
	}

	// 2. 重复访问
	if checkedAccesses[accessKey] {
		canEliminate = true
		reason = "already checked"
	}

	// 3. 简单的循环变量（启发式）
	if ident, ok := node.Index.(*ast.Ident); ok {
		// 如果索引是一个标识符，可能是循环变量
		// 这里简化处理，实际需要更复杂的控制流分析
		if strings.Contains(ident.Name, "i") || strings.Contains(ident.Name, "idx") {
			canEliminate = true
			reason = "likely loop induction variable"
		}
	}

	if !canEliminate {
		reason = "cannot prove bounds safety"
	}

	oa.results.BCEAnalysis = append(oa.results.BCEAnalysis, &BCEInfo{
		Position:     pos,
		ArrayExpr:    arrayExpr,
		IndexExpr:    indexExpr,
		CanEliminate: canEliminate,
		Reason:       reason,
	})

	// 记录已检查的访问
	checkedAccesses[accessKey] = true
}

func (oa *OptimizerAnalyzer) analyzeBCERange(node *ast.RangeStmt) {
	// range循环中的索引访问可以消除边界检查
	if node.Body == nil {
		return
	}

	// 获取循环变量名
	var indexVar string
	if key, ok := node.Key.(*ast.Ident); ok {
		indexVar = key.Name
	}

	if indexVar == "" {
		return
	}

	rangeExpr := exprToString(node.X)

	// 在循环体中查找使用该索引的数组访问
	ast.Inspect(node.Body, func(n ast.Node) bool {
		if indexExpr, ok := n.(*ast.IndexExpr); ok {
			if ident, ok := indexExpr.Index.(*ast.Ident); ok {
				if ident.Name == indexVar {
					arrayExpr := exprToString(indexExpr.X)
					if arrayExpr == rangeExpr {
						pos := oa.fset.Position(indexExpr.Pos())
						oa.results.BCEAnalysis = append(oa.results.BCEAnalysis, &BCEInfo{
							Position:     pos,
							ArrayExpr:    arrayExpr,
							IndexExpr:    indexVar,
							CanEliminate: true,
							Reason:       "range loop index",
						})
					}
				}
			}
		}
		return true
	})
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.BasicLit:
		return e.Value
	default:
		return "?"
	}
}

// =============================================================================
// 报告生成
// =============================================================================

// Report 生成完整的优化分析报告
func (oa *OptimizerAnalyzer) Report() string {
	var sb strings.Builder

	sb.WriteString("==========================================================================\n")
	sb.WriteString("               编译器优化分析报告 (Compiler Optimization Analysis)\n")
	sb.WriteString("==========================================================================\n\n")

	// 统计信息
	sb.WriteString("📊 统计信息:\n")
	sb.WriteString(fmt.Sprintf("   - 总函数数: %d\n", oa.results.TotalFunctions))
	sb.WriteString(fmt.Sprintf("   - 总分配数: %d\n", oa.results.TotalAllocations))
	sb.WriteString(fmt.Sprintf("   - 逃逸分析: %d 个对象\n", len(oa.results.EscapeAnalysis)))
	sb.WriteString(fmt.Sprintf("   - 内联分析: %d 个函数\n", len(oa.results.InlineAnalysis)))
	sb.WriteString(fmt.Sprintf("   - BCE分析: %d 个数组访问\n\n", len(oa.results.BCEAnalysis)))

	// 逃逸分析结果
	sb.WriteString("🔍 1. 逃逸分析 (Escape Analysis):\n\n")
	sb.WriteString("   理论: obj escapes ⟺ ∃ reference outliving function\n\n")

	if len(oa.results.EscapeAnalysis) == 0 {
		sb.WriteString("   ✅ 无逃逸对象\n\n")
	} else {
		stackCount := 0
		heapCount := 0
		for _, info := range oa.results.EscapeAnalysis {
			if info.EscapesTo == "stack" {
				stackCount++
			} else {
				heapCount++
			}
		}

		sb.WriteString(fmt.Sprintf("   📊 栈分配: %d 个 | 堆分配: %d 个\n\n", stackCount, heapCount))

		for i, info := range oa.results.EscapeAnalysis {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("   ... and %d more\n", len(oa.results.EscapeAnalysis)-10))
				break
			}

			icon := "⚠️ "
			if info.CanOptimize {
				icon = "✅"
			}

			sb.WriteString(fmt.Sprintf("   %s %s:%d:%d: %s -> %s (%s)\n",
				icon,
				info.Position.Filename,
				info.Position.Line,
				info.Position.Column,
				info.VarName,
				info.EscapesTo,
				info.Reason))
		}
		sb.WriteString("\n")
	}

	// 内联分析结果
	sb.WriteString("🔍 2. 函数内联分析 (Function Inlining Analysis):\n\n")
	sb.WriteString("   理论: CanInline(f) ⟺ cost < 80 ∧ ¬recursive ∧ ¬complex\n\n")

	if len(oa.results.InlineAnalysis) == 0 {
		sb.WriteString("   ℹ️  无函数分析\n\n")
	} else {
		canInlineCount := 0
		for _, info := range oa.results.InlineAnalysis {
			if info.CanInline {
				canInlineCount++
			}
		}

		sb.WriteString(fmt.Sprintf("   📊 可内联: %d 个 | 不可内联: %d 个\n\n",
			canInlineCount, len(oa.results.InlineAnalysis)-canInlineCount))

		for i, info := range oa.results.InlineAnalysis {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("   ... and %d more\n", len(oa.results.InlineAnalysis)-10))
				break
			}

			icon := "❌"
			if info.CanInline {
				icon = "✅"
			}

			sb.WriteString(fmt.Sprintf("   %s %s (cost: %d): %s\n",
				icon,
				info.FuncName,
				info.Cost,
				info.Reason))
		}
		sb.WriteString("\n")
	}

	// BCE分析结果
	sb.WriteString("🔍 3. 边界检查消除 (Bounds Check Elimination):\n\n")
	sb.WriteString("   理论: CanEliminate(a[i]) ⟺ provable: 0 ≤ i < len(a)\n\n")

	if len(oa.results.BCEAnalysis) == 0 {
		sb.WriteString("   ℹ️  无数组访问\n\n")
	} else {
		eliminateCount := 0
		for _, info := range oa.results.BCEAnalysis {
			if info.CanEliminate {
				eliminateCount++
			}
		}

		sb.WriteString(fmt.Sprintf("   📊 可消除: %d 个 | 不可消除: %d 个\n\n",
			eliminateCount, len(oa.results.BCEAnalysis)-eliminateCount))

		for i, info := range oa.results.BCEAnalysis {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("   ... and %d more\n", len(oa.results.BCEAnalysis)-10))
				break
			}

			icon := "❌"
			if info.CanEliminate {
				icon = "✅"
			}

			sb.WriteString(fmt.Sprintf("   %s %s:%d: %s[%s] (%s)\n",
				icon,
				info.Position.Filename,
				info.Position.Line,
				info.ArrayExpr,
				info.IndexExpr,
				info.Reason))
		}
		sb.WriteString("\n")
	}

	// 形式化理论基础
	sb.WriteString("📐 形式化理论基础:\n")
	sb.WriteString("   - 逃逸分析: obj escapes ⟺ ∃ ref outliving function\n")
	sb.WriteString("   - 内联分析: InlineCost < threshold ∧ ¬recursive\n")
	sb.WriteString("   - BCE: provable(0 ≤ i < len(a)) ⟹ eliminate check\n\n")

	sb.WriteString("==========================================================================\n")

	return sb.String()
}

// GetResults 返回分析结果
func (oa *OptimizerAnalyzer) GetResults() *AnalysisResults {
	return oa.results
}

// GetEscapeAnalysis 返回逃逸分析结果
func (oa *OptimizerAnalyzer) GetEscapeAnalysis() []*EscapeInfo {
	return oa.results.EscapeAnalysis
}

// GetInlineAnalysis 返回内联分析结果
func (oa *OptimizerAnalyzer) GetInlineAnalysis() []*InlineInfo {
	return oa.results.InlineAnalysis
}

// GetBCEAnalysis 返回BCE分析结果
func (oa *OptimizerAnalyzer) GetBCEAnalysis() []*BCEInfo {
	return oa.results.BCEAnalysis
}

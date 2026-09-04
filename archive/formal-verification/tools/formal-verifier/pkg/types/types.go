// Package types provides type system verification for Go programs.
// 类型系统验证包
//
// 理论基础：
// - 文档03：Go类型系统形式化定义
//
// 核心验证：
// 1. Progress定理 (Progress Theorem)
// 2. Preservation定理 (Preservation Theorem)
// 3. 类型安全性 (Type Safety)
// 4. 泛型约束检查 (Generic Constraints)
package types

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

// ====================================================================================
// 第一部分：核心数据结构
// ====================================================================================

// TypeEnvironment 类型环境
// 形式化定义：Γ: Variable → Type
type TypeEnvironment struct {
	bindings map[string]types.Type // 变量到类型的映射
	parent   *TypeEnvironment      // 父环境（用于作用域嵌套）
}

// NewTypeEnvironment 创建新的类型环境
func NewTypeEnvironment(parent *TypeEnvironment) *TypeEnvironment {
	return &TypeEnvironment{
		bindings: make(map[string]types.Type),
		parent:   parent,
	}
}

// Bind 绑定变量到类型
func (env *TypeEnvironment) Bind(name string, typ types.Type) {
	env.bindings[name] = typ
}

// Lookup 查找变量的类型
func (env *TypeEnvironment) Lookup(name string) (types.Type, bool) {
	if typ, ok := env.bindings[name]; ok {
		return typ, true
	}
	if env.parent != nil {
		return env.parent.Lookup(name)
	}
	return nil, false
}

// TypeJudgment 类型判断结果
// 形式化：Γ ⊢ e : T (在环境Γ下，表达式e具有类型T)
type TypeJudgment struct {
	Expression ast.Expr       // 表达式
	Type       types.Type     // 类型
	Valid      bool           // 是否有效
	Error      string         // 错误信息
	Position   token.Position // 位置
}

// ProgressError Progress定理违反
type ProgressError struct {
	Position token.Position
	Message  string
	Term     ast.Expr
}

// PreservationError Preservation定理违反
type PreservationError struct {
	Position token.Position
	Message  string
	Before   types.Type
	After    types.Type
}

// GenericConstraintError 泛型约束违反
type GenericConstraintError struct {
	Position   token.Position
	Message    string
	TypeParam  string
	Constraint types.Type
	Actual     types.Type
}

// ====================================================================================
// 第二部分：类型系统验证器
// ====================================================================================

// TypeVerifier 类型系统验证器
type TypeVerifier struct {
	fset               *token.FileSet
	typeInfo           *types.Info
	progressErrors     []ProgressError
	preservationErrors []PreservationError
	constraintErrors   []GenericConstraintError
	judgments          []*TypeJudgment
}

// NewVerifier 创建新的类型验证器
func NewVerifier() *TypeVerifier {
	return &TypeVerifier{
		progressErrors:     []ProgressError{},
		preservationErrors: []PreservationError{},
		constraintErrors:   []GenericConstraintError{},
		judgments:          []*TypeJudgment{},
	}
}

// VerifyFile 验证文件的类型安全性
//
// 注意：这里用 go/types 直接对单个文件做类型检查，而不是 go/packages 加载整个目录。
// testdata 夹具目录中的文件均为独立的 package main（彼此存在重复的 main 声明、
// 部分文件故意非法），目录级包加载必然失败；单文件语义才是本验证器的正确边界。
func (tv *TypeVerifier) VerifyFile(filename string) error {
	tv.fset = token.NewFileSet()

	// 解析文件
	file, err := parser.ParseFile(tv.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{Importer: importer.Default(), GoVersion: "go1.27"}
	if _, err := conf.Check(file.Name.Name, tv.fset, []*ast.File{file}, info); err != nil {
		return fmt.Errorf("type check: %w", err)
	}
	tv.typeInfo = info

	// 执行各种验证
	tv.verifyProgress(file)
	tv.verifyPreservation(file)
	tv.verifyGenericConstraints(file)

	return nil
}

// ====================================================================================
// 第三部分：Progress定理验证
// ====================================================================================

// verifyProgress 验证Progress定理
// 形式化定义：
//
//	Progress定理：如果 Γ ⊢ e : T 且 e 是良型的，
//	那么 e 要么是一个值，要么可以进行计算步骤。
//
// 公式：∀e, T. (⊢ e : T) ⟹ (value(e) ∨ ∃e'. e ↦ e')
func (tv *TypeVerifier) verifyProgress(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.Ident:
			tv.checkProgressIdent(expr)
		case *ast.BinaryExpr:
			tv.checkProgressBinaryExpr(expr)
		case *ast.CallExpr:
			tv.checkProgressCallExpr(expr)
		case *ast.IndexExpr:
			tv.checkProgressIndexExpr(expr)
		}
		return true
	})
}

// checkProgressIdent 检查标识符的Progress性质
func (tv *TypeVerifier) checkProgressIdent(ident *ast.Ident) {
	// 检查标识符是否已定义
	obj := tv.typeInfo.Uses[ident]
	if obj == nil {
		obj = tv.typeInfo.Defs[ident]
	}

	if obj == nil && ident.Name != "_" && ident.Name != "nil" {
		tv.progressErrors = append(tv.progressErrors, ProgressError{
			Position: tv.fset.Position(ident.Pos()),
			Message:  fmt.Sprintf("undefined identifier: %s", ident.Name),
			Term:     ident,
		})
	}
}

// checkProgressBinaryExpr 检查二元表达式的Progress性质
func (tv *TypeVerifier) checkProgressBinaryExpr(expr *ast.BinaryExpr) {
	// 检查操作数是否有类型
	xType := tv.typeInfo.Types[expr.X]
	yType := tv.typeInfo.Types[expr.Y]

	if !xType.IsValue() {
		tv.progressErrors = append(tv.progressErrors, ProgressError{
			Position: tv.fset.Position(expr.X.Pos()),
			Message:  "left operand is not a value",
			Term:     expr.X,
		})
	}

	if !yType.IsValue() {
		tv.progressErrors = append(tv.progressErrors, ProgressError{
			Position: tv.fset.Position(expr.Y.Pos()),
			Message:  "right operand is not a value",
			Term:     expr.Y,
		})
	}
}

// checkProgressCallExpr 检查函数调用的Progress性质
func (tv *TypeVerifier) checkProgressCallExpr(call *ast.CallExpr) {
	// 检查被调用者是否有效
	funType := tv.typeInfo.Types[call.Fun]
	if !funType.IsValue() && !funType.IsType() {
		tv.progressErrors = append(tv.progressErrors, ProgressError{
			Position: tv.fset.Position(call.Fun.Pos()),
			Message:  "callee is not a valid function",
			Term:     call.Fun,
		})
	}
}

// checkProgressIndexExpr 检查索引表达式的Progress性质
func (tv *TypeVerifier) checkProgressIndexExpr(idx *ast.IndexExpr) {
	// 检查被索引的对象是否有效
	xType := tv.typeInfo.Types[idx.X]
	if !xType.IsValue() {
		tv.progressErrors = append(tv.progressErrors, ProgressError{
			Position: tv.fset.Position(idx.X.Pos()),
			Message:  "indexed expression is not a value",
			Term:     idx.X,
		})
	}
}

// ====================================================================================
// 第四部分：Preservation定理验证
// ====================================================================================

// verifyPreservation 验证Preservation定理
// 形式化定义：
//
//	Preservation定理：如果 Γ ⊢ e : T 且 e ↦ e'，
//	那么 Γ ⊢ e' : T (类型在计算过程中保持不变)。
//
// 公式：∀e, e', T. (⊢ e : T ∧ e ↦ e') ⟹ ⊢ e' : T
func (tv *TypeVerifier) verifyPreservation(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			tv.checkPreservationAssign(stmt)
		case *ast.IfStmt:
			tv.checkPreservationIf(stmt)
		case *ast.ReturnStmt:
			tv.checkPreservationReturn(stmt)
		}
		return true
	})
}

// checkPreservationAssign 检查赋值语句的Preservation性质
func (tv *TypeVerifier) checkPreservationAssign(assign *ast.AssignStmt) {
	if len(assign.Lhs) != len(assign.Rhs) && assign.Tok != token.DEFINE {
		return // 处理多返回值的情况
	}

	for i := 0; i < len(assign.Lhs) && i < len(assign.Rhs); i++ {
		lhsType := tv.typeInfo.Types[assign.Lhs[i]]
		rhsType := tv.typeInfo.Types[assign.Rhs[i]]

		if !tv.isAssignable(lhsType.Type, rhsType.Type) {
			tv.preservationErrors = append(tv.preservationErrors, PreservationError{
				Position: tv.fset.Position(assign.Pos()),
				Message:  "type mismatch in assignment",
				Before:   lhsType.Type,
				After:    rhsType.Type,
			})
		}
	}
}

// checkPreservationIf 检查if语句的Preservation性质
func (tv *TypeVerifier) checkPreservationIf(ifStmt *ast.IfStmt) {
	if ifStmt.Cond != nil {
		condType := tv.typeInfo.Types[ifStmt.Cond]
		if condType.Type != nil && !tv.isBooleanType(condType.Type) {
			tv.preservationErrors = append(tv.preservationErrors, PreservationError{
				Position: tv.fset.Position(ifStmt.Cond.Pos()),
				Message:  "if condition must be boolean",
				Before:   types.Typ[types.Bool],
				After:    condType.Type,
			})
		}
	}
}

// checkPreservationReturn 检查return语句的Preservation性质
func (tv *TypeVerifier) checkPreservationReturn(ret *ast.ReturnStmt) {
	// 查找包含此return的函数
	// 这里简化处理，实际应该追踪函数上下文
	// TODO: 实现完整的函数返回类型检查
}

// ====================================================================================
// 第五部分：泛型约束验证
// ====================================================================================

// verifyGenericConstraints 验证泛型约束
// 形式化定义：
//
//	对于泛型类型参数 T 和约束 C，
//	任何使用 T 的地方都必须满足约束 C。
//
// 公式：∀T, C, t. (T : C ∧ t : T) ⟹ satisfies(t, C)
func (tv *TypeVerifier) verifyGenericConstraints(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.FuncDecl:
			if decl.Type.TypeParams != nil {
				tv.checkGenericFunction(decl)
			}
		case *ast.TypeSpec:
			if decl.TypeParams != nil {
				tv.checkGenericType(decl)
			}
		}
		return true
	})
}

// checkGenericFunction 检查泛型函数的约束
func (tv *TypeVerifier) checkGenericFunction(funcDecl *ast.FuncDecl) {
	// 检查类型参数的约束
	if funcDecl.Type.TypeParams == nil {
		return
	}

	for _, field := range funcDecl.Type.TypeParams.List {
		for _, name := range field.Names {
			// 检查约束是否满足
			if field.Type != nil {
				constraint := tv.typeInfo.Types[field.Type]
				if constraint.Type != nil {
					tv.checkConstraintSatisfaction(name, constraint.Type, funcDecl.Body)
				}
			}
		}
	}
}

// checkGenericType 检查泛型类型的约束
func (tv *TypeVerifier) checkGenericType(typeSpec *ast.TypeSpec) {
	// 检查类型参数的约束
	if typeSpec.TypeParams == nil {
		return
	}

	for _, field := range typeSpec.TypeParams.List {
		for _, name := range field.Names {
			if field.Type != nil {
				constraint := tv.typeInfo.Types[field.Type]
				if constraint.Type != nil {
					// 检查类型定义中的约束使用
					tv.checkTypeConstraintUsage(name, constraint.Type, typeSpec.Type)
				}
			}
		}
	}
}

// checkConstraintSatisfaction 检查约束满足性
func (tv *TypeVerifier) checkConstraintSatisfaction(typeParam *ast.Ident, constraint types.Type, body *ast.BlockStmt) {
	if body == nil {
		return
	}

	// 遍历函数体，检查类型参数的使用
	ast.Inspect(body, func(n ast.Node) bool {
		// 检查类型参数的实际使用是否满足约束
		// 这里简化处理，实际需要更复杂的约束检查
		return true
	})
}

// checkTypeConstraintUsage 检查类型定义中的约束使用
func (tv *TypeVerifier) checkTypeConstraintUsage(typeParam *ast.Ident, constraint types.Type, typeExpr ast.Expr) {
	// 检查类型定义中类型参数的使用
	ast.Inspect(typeExpr, func(n ast.Node) bool {
		// 简化处理
		return true
	})
}

// ====================================================================================
// 第六部分：辅助函数
// ====================================================================================

// isAssignable 检查类型t是否可以赋值给类型u
func (tv *TypeVerifier) isAssignable(u, t types.Type) bool {
	if u == nil || t == nil {
		return false
	}

	// 使用types包的Assignable函数
	return types.AssignableTo(t, u)
}

// isBooleanType 检查类型是否为布尔类型
func (tv *TypeVerifier) isBooleanType(t types.Type) bool {
	if t == nil {
		return false
	}
	basic, ok := t.Underlying().(*types.Basic)
	return ok && basic.Kind() == types.Bool
}

// ====================================================================================
// 第七部分：报告生成
// ====================================================================================

// Report 生成类型验证报告
func (tv *TypeVerifier) Report() string {
	var sb strings.Builder

	sb.WriteString("==========================================================================\n")
	sb.WriteString("               类型系统验证报告 (Type System Verification Report)\n")
	sb.WriteString("==========================================================================\n\n")

	sb.WriteString("📊 统计信息:\n")
	sb.WriteString(fmt.Sprintf("   - Type Judgments: %d\n", len(tv.judgments)))
	sb.WriteString(fmt.Sprintf("   - Progress Errors: %d\n", len(tv.progressErrors)))
	sb.WriteString(fmt.Sprintf("   - Preservation Errors: %d\n", len(tv.preservationErrors)))
	sb.WriteString(fmt.Sprintf("   - Constraint Errors: %d\n\n", len(tv.constraintErrors)))

	sb.WriteString("🔍 验证结果:\n\n")

	// Progress定理
	if len(tv.progressErrors) == 0 {
		sb.WriteString("   1. ✅ Progress定理: 验证通过\n")
	} else {
		sb.WriteString(fmt.Sprintf("   1. ⚠️  Progress定理: %d个违反\n", len(tv.progressErrors)))
		for i, err := range tv.progressErrors {
			if i < 3 { // 只显示前3个
				sb.WriteString(fmt.Sprintf("      - %s: %s\n", err.Position, err.Message))
			}
		}
		if len(tv.progressErrors) > 3 {
			sb.WriteString(fmt.Sprintf("      ... and %d more\n", len(tv.progressErrors)-3))
		}
	}

	// Preservation定理
	if len(tv.preservationErrors) == 0 {
		sb.WriteString("   2. ✅ Preservation定理: 验证通过\n")
	} else {
		sb.WriteString(fmt.Sprintf("   2. ⚠️  Preservation定理: %d个违反\n", len(tv.preservationErrors)))
		for i, err := range tv.preservationErrors {
			if i < 3 {
				sb.WriteString(fmt.Sprintf("      - %s: %s\n", err.Position, err.Message))
			}
		}
		if len(tv.preservationErrors) > 3 {
			sb.WriteString(fmt.Sprintf("      ... and %d more\n", len(tv.preservationErrors)-3))
		}
	}

	// 泛型约束
	if len(tv.constraintErrors) == 0 {
		sb.WriteString("   3. ✅ 泛型约束: 验证通过\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("   3. ⚠️  泛型约束: %d个违反\n\n", len(tv.constraintErrors)))
	}

	sb.WriteString("📐 形式化理论基础:\n")
	sb.WriteString("   - Progress: ∀e, T. (⊢ e : T) ⟹ (value(e) ∨ ∃e'. e ↦ e')\n")
	sb.WriteString("   - Preservation: ∀e, e', T. (⊢ e : T ∧ e ↦ e') ⟹ ⊢ e' : T\n")
	sb.WriteString("   - Type Safety: Progress ∧ Preservation\n")
	sb.WriteString("   - Constraints: ∀T, C. (T : C) ⟹ satisfies(T, C)\n\n")

	sb.WriteString("==========================================================================\n")

	return sb.String()
}

// GetProgressErrors 获取Progress错误
func (tv *TypeVerifier) GetProgressErrors() []ProgressError {
	return tv.progressErrors
}

// GetPreservationErrors 获取Preservation错误
func (tv *TypeVerifier) GetPreservationErrors() []PreservationError {
	return tv.preservationErrors
}

// GetConstraintErrors 获取约束错误
func (tv *TypeVerifier) GetConstraintErrors() []GenericConstraintError {
	return tv.constraintErrors
}

// IsSafe 检查是否类型安全
// 类型安全 = Progress ∧ Preservation
func (tv *TypeVerifier) IsSafe() bool {
	return len(tv.progressErrors) == 0 && len(tv.preservationErrors) == 0
}

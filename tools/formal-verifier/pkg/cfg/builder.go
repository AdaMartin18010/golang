// Package cfg implements Control Flow Graph construction for Go programs
// Based on Document 13: Go Control Flow Formalization
package cfg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// Node represents a basic block in the CFG
type Node struct {
	ID           int
	Label        string
	Stmts        []ast.Stmt
	Successors   []*Node
	Predecessors []*Node
}

// CFG represents a Control Flow Graph
// Formal Definition (from Document 13):
//
//	CFG = (N, E, entry, exit)
//	- N: set of nodes (basic blocks)
//	- E: set of edges (control flow transitions)
//	- entry: entry node
//	- exit: exit node
type CFG struct {
	Nodes       []*Node
	Entry       *Node
	Exit        *Node
	FSet        *token.FileSet
	nodeCounter int
}

// Builder constructs CFG from Go source code
type Builder struct {
	fset         *token.FileSet
	cfg          *CFG
	currentBlock *Node
}

// NewBuilder creates a new CFG builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// BuildFromFile constructs a CFG from a Go source file
func (b *Builder) BuildFromFile(filename string) (*CFG, error) {
	// Parse the source file
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	file, err := parser.ParseFile(b.fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Initialize CFG
	b.cfg = &CFG{
		Nodes:       make([]*Node, 0),
		FSet:        b.fset,
		nodeCounter: 0,
	}

	// Create entry and exit nodes
	b.cfg.Entry = b.newNode("Entry")
	b.cfg.Exit = b.newNode("Exit")

	b.currentBlock = b.cfg.Entry

	// Build CFG for each function in the file
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			b.buildFunction(funcDecl)
		}
	}

	return b.cfg, nil
}

// newNode creates a new CFG node
func (b *Builder) newNode(label string) *Node {
	node := &Node{
		ID:           b.cfg.nodeCounter,
		Label:        label,
		Stmts:        make([]ast.Stmt, 0),
		Successors:   make([]*Node, 0),
		Predecessors: make([]*Node, 0),
	}
	b.cfg.nodeCounter++
	b.cfg.Nodes = append(b.cfg.Nodes, node)
	return node
}

// addEdge adds a control flow edge from 'from' to 'to'
func (b *Builder) addEdge(from, to *Node) {
	from.Successors = append(from.Successors, to)
	to.Predecessors = append(to.Predecessors, from)
}

// buildFunction builds CFG for a function declaration
func (b *Builder) buildFunction(fn *ast.FuncDecl) {
	if fn.Body == nil {
		return
	}

	// Create function entry block
	funcEntry := b.newNode(fmt.Sprintf("func_%s_entry", fn.Name.Name))
	b.addEdge(b.cfg.Entry, funcEntry)

	b.currentBlock = funcEntry

	// Build CFG for function body
	exitBlock := b.buildStmtList(fn.Body.List)

	// Connect to exit
	if exitBlock != nil {
		b.addEdge(exitBlock, b.cfg.Exit)
	}
}

// buildStmtList builds CFG for a list of statements
func (b *Builder) buildStmtList(stmts []ast.Stmt) *Node {
	if len(stmts) == 0 {
		return b.currentBlock
	}

	for _, stmt := range stmts {
		b.currentBlock = b.buildStmt(stmt)
		if b.currentBlock == nil {
			return nil
		}
	}

	return b.currentBlock
}

// buildStmt builds CFG for a single statement
func (b *Builder) buildStmt(stmt ast.Stmt) *Node {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return b.buildExprStmt(s)

	case *ast.AssignStmt:
		return b.buildAssignStmt(s)

	case *ast.IfStmt:
		return b.buildIfStmt(s)

	case *ast.ForStmt:
		return b.buildForStmt(s)

	case *ast.RangeStmt:
		return b.buildRangeStmt(s)

	case *ast.SwitchStmt:
		return b.buildSwitchStmt(s)

	case *ast.ReturnStmt:
		return b.buildReturnStmt(s)

	case *ast.GoStmt:
		return b.buildGoStmt(s)

	case *ast.DeferStmt:
		return b.buildDeferStmt(s)

	case *ast.SelectStmt:
		return b.buildSelectStmt(s)

	case *ast.BlockStmt:
		return b.buildStmtList(s.List)

	default:
		// For other statement types, add to current block
		b.currentBlock.Stmts = append(b.currentBlock.Stmts, stmt)
		return b.currentBlock
	}
}

// buildExprStmt builds CFG for an expression statement
func (b *Builder) buildExprStmt(stmt *ast.ExprStmt) *Node {
	b.currentBlock.Stmts = append(b.currentBlock.Stmts, stmt)
	return b.currentBlock
}

// buildAssignStmt builds CFG for an assignment statement
func (b *Builder) buildAssignStmt(stmt *ast.AssignStmt) *Node {
	b.currentBlock.Stmts = append(b.currentBlock.Stmts, stmt)
	return b.currentBlock
}

// buildIfStmt builds CFG for an if statement
// Formal: if condition { then_branch } else { else_branch }
//
//	Creates diamond-shaped CFG structure
func (b *Builder) buildIfStmt(stmt *ast.IfStmt) *Node {
	// Init statement (if any)
	if stmt.Init != nil {
		b.currentBlock = b.buildStmt(stmt.Init)
	}

	// Condition block
	condBlock := b.newNode("if_cond")
	condBlock.Stmts = []ast.Stmt{&ast.ExprStmt{X: stmt.Cond}}
	b.addEdge(b.currentBlock, condBlock)

	// Then branch
	thenBlock := b.newNode("if_then")
	b.addEdge(condBlock, thenBlock)
	b.currentBlock = thenBlock
	thenExit := b.buildStmtList(stmt.Body.List)

	// Merge block (where control flow converges)
	mergeBlock := b.newNode("if_merge")

	// Else branch (if exists)
	var elseExit *Node
	if stmt.Else != nil {
		elseBlock := b.newNode("if_else")
		b.addEdge(condBlock, elseBlock)
		b.currentBlock = elseBlock
		elseExit = b.buildStmt(stmt.Else)
		if elseExit != nil {
			b.addEdge(elseExit, mergeBlock)
		}
	} else {
		// No else branch: condition -> merge
		b.addEdge(condBlock, mergeBlock)
	}

	// Then branch -> merge
	if thenExit != nil {
		b.addEdge(thenExit, mergeBlock)
	}

	return mergeBlock
}

// buildForStmt builds CFG for a for loop
// Formal: for init; cond; post { body }
//
//	Creates loop structure with back edge
func (b *Builder) buildForStmt(stmt *ast.ForStmt) *Node {
	// Init statement (if any)
	if stmt.Init != nil {
		b.currentBlock = b.buildStmt(stmt.Init)
	}

	// Loop header (condition)
	headerBlock := b.newNode("for_header")
	b.addEdge(b.currentBlock, headerBlock)
	if stmt.Cond != nil {
		headerBlock.Stmts = []ast.Stmt{&ast.ExprStmt{X: stmt.Cond}}
	}

	// Loop body
	bodyBlock := b.newNode("for_body")
	b.addEdge(headerBlock, bodyBlock)
	b.currentBlock = bodyBlock
	bodyExit := b.buildStmtList(stmt.Body.List)

	// Post statement
	var postBlock *Node
	if stmt.Post != nil {
		postBlock = b.newNode("for_post")
		if bodyExit != nil {
			b.addEdge(bodyExit, postBlock)
		}
		postBlock.Stmts = []ast.Stmt{stmt.Post}
		// Back edge: post -> header
		b.addEdge(postBlock, headerBlock)
	} else {
		// Back edge: body -> header
		if bodyExit != nil {
			b.addEdge(bodyExit, headerBlock)
		}
	}

	// Exit block (loop exit)
	exitBlock := b.newNode("for_exit")
	b.addEdge(headerBlock, exitBlock)

	return exitBlock
}

// buildRangeStmt builds CFG for a range loop
func (b *Builder) buildRangeStmt(stmt *ast.RangeStmt) *Node {
	// Range header
	headerBlock := b.newNode("range_header")
	b.addEdge(b.currentBlock, headerBlock)
	headerBlock.Stmts = []ast.Stmt{stmt}

	// Range body
	bodyBlock := b.newNode("range_body")
	b.addEdge(headerBlock, bodyBlock)
	b.currentBlock = bodyBlock
	bodyExit := b.buildStmtList(stmt.Body.List)

	// Back edge: body -> header
	if bodyExit != nil {
		b.addEdge(bodyExit, headerBlock)
	}

	// Exit block
	exitBlock := b.newNode("range_exit")
	b.addEdge(headerBlock, exitBlock)

	return exitBlock
}

// buildSwitchStmt builds CFG for a switch statement
func (b *Builder) buildSwitchStmt(stmt *ast.SwitchStmt) *Node {
	// Init statement (if any)
	if stmt.Init != nil {
		b.currentBlock = b.buildStmt(stmt.Init)
	}

	// Switch header
	switchBlock := b.newNode("switch_header")
	b.addEdge(b.currentBlock, switchBlock)
	if stmt.Tag != nil {
		switchBlock.Stmts = []ast.Stmt{&ast.ExprStmt{X: stmt.Tag}}
	}

	// Merge block
	mergeBlock := b.newNode("switch_merge")

	// Process each case
	for i, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			continue
		}

		caseBlock := b.newNode(fmt.Sprintf("case_%d", i))
		b.addEdge(switchBlock, caseBlock)
		b.currentBlock = caseBlock
		caseExit := b.buildStmtList(caseClause.Body)

		// Connect to merge (unless there's a fallthrough)
		if caseExit != nil {
			b.addEdge(caseExit, mergeBlock)
		}
	}

	return mergeBlock
}

// buildReturnStmt builds CFG for a return statement
func (b *Builder) buildReturnStmt(stmt *ast.ReturnStmt) *Node {
	returnBlock := b.newNode("return")
	returnBlock.Stmts = []ast.Stmt{stmt}
	b.addEdge(b.currentBlock, returnBlock)
	b.addEdge(returnBlock, b.cfg.Exit)
	return nil // No successor (return terminates control flow)
}

// buildGoStmt builds CFG for a go statement (goroutine spawn)
func (b *Builder) buildGoStmt(stmt *ast.GoStmt) *Node {
	goBlock := b.newNode("go_spawn")
	goBlock.Stmts = []ast.Stmt{stmt}
	b.addEdge(b.currentBlock, goBlock)

	// TODO: Model concurrent execution (for now, just continue)
	nextBlock := b.newNode("after_go")
	b.addEdge(goBlock, nextBlock)

	return nextBlock
}

// buildDeferStmt builds CFG for a defer statement
func (b *Builder) buildDeferStmt(stmt *ast.DeferStmt) *Node {
	deferBlock := b.newNode("defer")
	deferBlock.Stmts = []ast.Stmt{stmt}
	b.addEdge(b.currentBlock, deferBlock)

	// Defer doesn't change control flow immediately
	nextBlock := b.newNode("after_defer")
	b.addEdge(deferBlock, nextBlock)

	return nextBlock
}

// buildSelectStmt builds CFG for a select statement
func (b *Builder) buildSelectStmt(stmt *ast.SelectStmt) *Node {
	// Select header
	selectBlock := b.newNode("select_header")
	b.addEdge(b.currentBlock, selectBlock)

	// Merge block
	mergeBlock := b.newNode("select_merge")

	// Process each case
	for i, clause := range stmt.Body.List {
		commClause, ok := clause.(*ast.CommClause)
		if !ok {
			continue
		}

		caseBlock := b.newNode(fmt.Sprintf("select_case_%d", i))
		b.addEdge(selectBlock, caseBlock)
		b.currentBlock = caseBlock
		caseExit := b.buildStmtList(commClause.Body)

		if caseExit != nil {
			b.addEdge(caseExit, mergeBlock)
		}
	}

	return mergeBlock
}

// Print prints the CFG in a human-readable format
func (cfg *CFG) Print() {
	fmt.Println("=== Control Flow Graph ===")
	fmt.Printf("Total Nodes: %d\n", len(cfg.Nodes))
	fmt.Printf("Entry: %s (ID: %d)\n", cfg.Entry.Label, cfg.Entry.ID)
	fmt.Printf("Exit: %s (ID: %d)\n", cfg.Exit.Label, cfg.Exit.ID)
	fmt.Println()

	for _, node := range cfg.Nodes {
		fmt.Printf("Node %d: %s\n", node.ID, node.Label)
		fmt.Printf("  Statements: %d\n", len(node.Stmts))

		if len(node.Successors) > 0 {
			fmt.Print("  Successors: ")
			for i, succ := range node.Successors {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", succ.ID)
			}
			fmt.Println()
		}

		if len(node.Predecessors) > 0 {
			fmt.Print("  Predecessors: ")
			for i, pred := range node.Predecessors {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", pred.ID)
			}
			fmt.Println()
		}

		fmt.Println()
	}
}

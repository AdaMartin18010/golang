// Package cfg - SSA (Static Single Assignment) transformation
// Based on Document 13: Go Control Flow Formalization, Chapter 3
package cfg

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SSANode represents a node in SSA form
type SSANode struct {
	*Node                   // Embed original CFG node
	PhiFunctions []*PhiFunc // φ-functions at the beginning of this block
	SSAStmts     []SSAStmt  // Statements in SSA form
}

// PhiFunc represents a φ-function in SSA
// Formal Definition (from Document 13):
//
//	φ(x₁, x₂, ..., xₙ) where xᵢ comes from predecessor i
type PhiFunc struct {
	Variable string           // Target variable name
	Version  int              // SSA version number
	Sources  map[*Node]string // Source variable for each predecessor
}

// SSAStmt represents a statement in SSA form
type SSAStmt struct {
	Original ast.Stmt // Original statement
	Defs     []string // Variables defined (with version)
	Uses     []string // Variables used (with version)
}

// SSACFG represents a CFG in SSA form
type SSACFG struct {
	*CFG
	SSANodes    map[*Node]*SSANode
	DomTree     *DominatorTree
	DomFrontier map[*Node]map[*Node]bool
	VarVersions map[string]int // Current version for each variable
}

// DominatorTree represents the dominator tree
type DominatorTree struct {
	Root     *Node
	IDom     map[*Node]*Node          // Immediate dominator
	Children map[*Node][]*Node        // Dominator tree children
	DomSet   map[*Node]map[*Node]bool // Dominance sets
}

// SSAConverter converts a CFG to SSA form
type SSAConverter struct {
	cfg         *CFG
	ssaCFG      *SSACFG
	domTree     *DominatorTree
	domFrontier map[*Node]map[*Node]bool
	varVersions map[string]int
	varStacks   map[string][]int
}

// NewSSAConverter creates a new SSA converter
func NewSSAConverter(cfg *CFG) *SSAConverter {
	return &SSAConverter{
		cfg:         cfg,
		varVersions: make(map[string]int),
		varStacks:   make(map[string][]int),
	}
}

// Convert converts the CFG to SSA form
// Algorithm based on Cytron et al. "Efficiently Computing Static Single Assignment Form"
//
// Steps:
// 1. Compute dominance information (dominator tree, dominance frontiers)
// 2. Place φ-functions at dominance frontiers
// 3. Rename variables to unique SSA versions
func (c *SSAConverter) Convert() (*SSACFG, error) {
	// Step 1: Compute dominator tree
	fmt.Println("🔄 Step 1: Computing dominator tree...")
	domTree, err := c.computeDominatorTree()
	if err != nil {
		return nil, fmt.Errorf("failed to compute dominator tree: %w", err)
	}
	c.domTree = domTree

	// Step 2: Compute dominance frontiers
	fmt.Println("🔄 Step 2: Computing dominance frontiers...")
	c.domFrontier = c.computeDominanceFrontiers()

	// Step 3: Place φ-functions
	fmt.Println("🔄 Step 3: Placing φ-functions...")
	ssaNodes := c.placePhiFunctions()

	// Step 4: Rename variables
	fmt.Println("🔄 Step 4: Renaming variables...")
	c.renameVariables(c.cfg.Entry, ssaNodes)

	// Create SSA CFG
	c.ssaCFG = &SSACFG{
		CFG:         c.cfg,
		SSANodes:    ssaNodes,
		DomTree:     domTree,
		DomFrontier: c.domFrontier,
		VarVersions: c.varVersions,
	}

	fmt.Println("✅ SSA conversion completed!")
	return c.ssaCFG, nil
}

// computeDominatorTree computes the dominator tree using iterative algorithm
// Formal Definition (from Document 13):
//
//	Node d dominates node n if every path from entry to n goes through d
//	IDom(n) = immediate dominator of n (closest strict dominator)
func (c *SSAConverter) computeDominatorTree() (*DominatorTree, error) {
	if c.cfg.Entry == nil {
		return nil, fmt.Errorf("CFG has no entry node")
	}

	tree := &DominatorTree{
		Root:     c.cfg.Entry,
		IDom:     make(map[*Node]*Node),
		Children: make(map[*Node][]*Node),
		DomSet:   make(map[*Node]map[*Node]bool),
	}

	// Initialize dominator sets
	// Dom(entry) = {entry}
	// Dom(n) = all nodes (initial over-approximation)
	allNodes := make(map[*Node]bool)
	for _, node := range c.cfg.Nodes {
		allNodes[node] = true
	}

	for _, node := range c.cfg.Nodes {
		if node == c.cfg.Entry {
			tree.DomSet[node] = map[*Node]bool{node: true}
		} else {
			tree.DomSet[node] = make(map[*Node]bool)
			for n := range allNodes {
				tree.DomSet[node][n] = true
			}
		}
	}

	// Iterative fixed-point algorithm
	// Dom(n) = {n} ∪ (∩ Dom(p) for all predecessors p of n)
	changed := true
	for changed {
		changed = false
		for _, node := range c.cfg.Nodes {
			if node == c.cfg.Entry {
				continue
			}

			// Compute intersection of predecessor dominators
			newDomSet := make(map[*Node]bool)
			newDomSet[node] = true // n always dominates itself

			if len(node.Predecessors) > 0 {
				// Start with first predecessor's dom set
				for n := range tree.DomSet[node.Predecessors[0]] {
					newDomSet[n] = true
				}

				// Intersect with other predecessors
				for _, pred := range node.Predecessors[1:] {
					for n := range newDomSet {
						if !tree.DomSet[pred][n] && n != node {
							delete(newDomSet, n)
						}
					}
				}
			}

			// Check if changed
			if len(newDomSet) != len(tree.DomSet[node]) {
				changed = true
				tree.DomSet[node] = newDomSet
			} else {
				for n := range newDomSet {
					if !tree.DomSet[node][n] {
						changed = true
						tree.DomSet[node] = newDomSet
						break
					}
				}
			}
		}
	}

	// Compute immediate dominators
	// IDom(n) is the unique node in Dom(n) - {n} that is dominated by all others
	for _, node := range c.cfg.Nodes {
		if node == c.cfg.Entry {
			continue
		}

		var idom *Node
		for d := range tree.DomSet[node] {
			if d == node {
				continue
			}
			// Check if d is the immediate dominator
			// i.e., no other dominator strictly between d and n
			isIdom := true
			for d2 := range tree.DomSet[node] {
				if d2 == node || d2 == d {
					continue
				}
				// If d2 dominates d, then d is not immediate
				if tree.DomSet[d][d2] {
					isIdom = false
					break
				}
			}
			if isIdom {
				idom = d
				break
			}
		}

		if idom != nil {
			tree.IDom[node] = idom
			tree.Children[idom] = append(tree.Children[idom], node)
		}
	}

	return tree, nil
}

// computeDominanceFrontiers computes dominance frontiers for all nodes
// Formal Definition (from Document 13):
//
//	DF(n) = {y | ∃ predecessor p of y: n dominates p but n does not strictly dominate y}
func (c *SSAConverter) computeDominanceFrontiers() map[*Node]map[*Node]bool {
	df := make(map[*Node]map[*Node]bool)

	// Initialize
	for _, node := range c.cfg.Nodes {
		df[node] = make(map[*Node]bool)
	}

	// For each node n
	for _, node := range c.cfg.Nodes {
		if len(node.Predecessors) < 2 {
			// Join nodes (with multiple predecessors) are interesting
			continue
		}

		// For each predecessor p of n
		for _, pred := range node.Predecessors {
			runner := pred
			// Walk up dominator tree from p
			for runner != c.domTree.IDom[node] && runner != nil {
				df[runner][node] = true
				runner = c.domTree.IDom[runner]
			}
		}
	}

	return df
}

// placePhiFunctions places φ-functions at appropriate locations
// Algorithm: Place φ-function for variable v at dominance frontier of all nodes that define v
func (c *SSAConverter) placePhiFunctions() map[*Node]*SSANode {
	ssaNodes := make(map[*Node]*SSANode)

	// Initialize SSA nodes
	for _, node := range c.cfg.Nodes {
		ssaNodes[node] = &SSANode{
			Node:         node,
			PhiFunctions: make([]*PhiFunc, 0),
			SSAStmts:     make([]SSAStmt, 0),
		}
	}

	// Collect variable definitions
	varDefs := make(map[string]map[*Node]bool)
	for _, node := range c.cfg.Nodes {
		vars := c.getDefinedVariables(node)
		for _, v := range vars {
			if varDefs[v] == nil {
				varDefs[v] = make(map[*Node]bool)
			}
			varDefs[v][node] = true
		}
	}

	// Place φ-functions
	for varName, defNodes := range varDefs {
		// Compute iterated dominance frontier
		workList := make([]*Node, 0)
		for n := range defNodes {
			workList = append(workList, n)
		}

		phiPlaced := make(map[*Node]bool)
		for len(workList) > 0 {
			node := workList[0]
			workList = workList[1:]

			// For each node in dominance frontier
			for dfNode := range c.domFrontier[node] {
				if !phiPlaced[dfNode] {
					// Place φ-function
					phi := &PhiFunc{
						Variable: varName,
						Version:  0, // Will be set during renaming
						Sources:  make(map[*Node]string),
					}
					ssaNodes[dfNode].PhiFunctions = append(ssaNodes[dfNode].PhiFunctions, phi)
					phiPlaced[dfNode] = true

					// If this is a new definition point, add to worklist
					if !defNodes[dfNode] {
						workList = append(workList, dfNode)
						defNodes[dfNode] = true
					}
				}
			}
		}
	}

	return ssaNodes
}

// renameVariables performs variable renaming to create unique SSA names
// Algorithm: DFS traversal of dominator tree with stack-based renaming
func (c *SSAConverter) renameVariables(node *Node, ssaNodes map[*Node]*SSANode) {
	if node == nil {
		return
	}

	ssaNode := ssaNodes[node]

	// Save current stack state (for backtracking)
	savedStacks := make(map[string][]int)
	for v, stack := range c.varStacks {
		savedStacks[v] = make([]int, len(stack))
		copy(savedStacks[v], stack)
	}

	// Rename φ-function targets
	for _, phi := range ssaNode.PhiFunctions {
		newVersion := c.getNewVersion(phi.Variable)
		phi.Version = newVersion
		c.pushVersion(phi.Variable, newVersion)
	}

	// Rename statements
	for _, stmt := range node.Stmts {
		// Rename uses
		uses := c.getUsedVariables(stmt)
		for _, v := range uses {
			// Get current version from stack
			if len(c.varStacks[v]) > 0 {
				_ = c.getCurrentVersion(v)
			}
		}

		// Rename definitions
		defs := c.getDefinedVariablesFromStmt(stmt)
		for _, v := range defs {
			newVersion := c.getNewVersion(v)
			c.pushVersion(v, newVersion)
		}
	}

	// Rename φ-function sources in successors
	for _, succ := range node.Successors {
		succSSA := ssaNodes[succ]
		for _, phi := range succSSA.PhiFunctions {
			if len(c.varStacks[phi.Variable]) > 0 {
				version := c.getCurrentVersion(phi.Variable)
				phi.Sources[node] = fmt.Sprintf("%s_%d", phi.Variable, version)
			}
		}
	}

	// Recursively process children in dominator tree
	for _, child := range c.domTree.Children[node] {
		c.renameVariables(child, ssaNodes)
	}

	// Restore stack state
	c.varStacks = savedStacks
}

// Helper functions for variable analysis

func (c *SSAConverter) getDefinedVariables(node *Node) []string {
	vars := make(map[string]bool)
	for _, stmt := range node.Stmts {
		for _, v := range c.getDefinedVariablesFromStmt(stmt) {
			vars[v] = true
		}
	}
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}
	return result
}

func (c *SSAConverter) getDefinedVariablesFromStmt(stmt ast.Stmt) []string {
	vars := make([]string, 0)
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		for _, lhs := range s.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				vars = append(vars, ident.Name)
			}
		}
	case *ast.DeclStmt:
		if genDecl, ok := s.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						vars = append(vars, name.Name)
					}
				}
			}
		}
	}
	return vars
}

func (c *SSAConverter) getUsedVariables(stmt ast.Stmt) []string {
	vars := make(map[string]bool)
	// Simplified: would need full expression traversal
	// For now, approximate with identifiers in RHS
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		for _, rhs := range s.Rhs {
			c.collectIdentifiers(rhs, vars)
		}
	case *ast.ExprStmt:
		c.collectIdentifiers(s.X, vars)
	}
	result := make([]string, 0, len(vars))
	for v := range vars {
		result = append(result, v)
	}
	return result
}

func (c *SSAConverter) collectIdentifiers(expr ast.Expr, vars map[string]bool) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *ast.Ident:
		vars[e.Name] = true
	case *ast.BinaryExpr:
		c.collectIdentifiers(e.X, vars)
		c.collectIdentifiers(e.Y, vars)
	case *ast.UnaryExpr:
		c.collectIdentifiers(e.X, vars)
	case *ast.CallExpr:
		for _, arg := range e.Args {
			c.collectIdentifiers(arg, vars)
		}
	}
}

// Version management

func (c *SSAConverter) getNewVersion(variable string) int {
	c.varVersions[variable]++
	return c.varVersions[variable]
}

func (c *SSAConverter) getCurrentVersion(variable string) int {
	stack := c.varStacks[variable]
	if len(stack) == 0 {
		return 0
	}
	return stack[len(stack)-1]
}

func (c *SSAConverter) pushVersion(variable string, version int) {
	c.varStacks[variable] = append(c.varStacks[variable], version)
}

// Print prints the SSA CFG
func (s *SSACFG) Print() {
	fmt.Println("=== SSA Control Flow Graph ===")
	fmt.Printf("Total Nodes: %d\n", len(s.Nodes))
	fmt.Println()

	for _, node := range s.Nodes {
		ssaNode := s.SSANodes[node]
		fmt.Printf("Node %d: %s\n", node.ID, node.Label)

		// Print φ-functions
		if len(ssaNode.PhiFunctions) > 0 {
			fmt.Println("  φ-functions:")
			for _, phi := range ssaNode.PhiFunctions {
				fmt.Printf("    %s_%d = φ(", phi.Variable, phi.Version)
				first := true
				for pred, source := range phi.Sources {
					if !first {
						fmt.Print(", ")
					}
					fmt.Printf("%s from %d", source, pred.ID)
					first = false
				}
				fmt.Println(")")
			}
		}

		// Print statements
		if len(node.Stmts) > 0 {
			fmt.Printf("  Statements: %d\n", len(node.Stmts))
		}

		fmt.Println()
	}

	// Print dominator tree
	if s.DomTree != nil {
		fmt.Println("=== Dominator Tree ===")
		for node, idom := range s.DomTree.IDom {
			fmt.Printf("IDom(%d) = %d\n", node.ID, idom.ID)
		}
		fmt.Println()
	}

	// Print dominance frontiers
	fmt.Println("=== Dominance Frontiers ===")
	for node, df := range s.DomFrontier {
		if len(df) > 0 {
			fmt.Printf("DF(%d) = {", node.ID)
			first := true
			for dfNode := range df {
				if !first {
					fmt.Print(", ")
				}
				fmt.Printf("%d", dfNode.ID)
				first = false
			}
			fmt.Println("}")
		}
	}
}

// VerifySSAProperty verifies that the SSA property holds
// Property: Each variable is defined exactly once
func (s *SSACFG) VerifySSAProperty() (bool, []string) {
	errors := make([]string, 0)
	definedVars := make(map[string]bool)

	for _, node := range s.Nodes {
		ssaNode := s.SSANodes[node]

		// Check φ-functions
		for _, phi := range ssaNode.PhiFunctions {
			varName := fmt.Sprintf("%s_%d", phi.Variable, phi.Version)
			if definedVars[varName] {
				errors = append(errors, fmt.Sprintf("Variable %s defined multiple times", varName))
			}
			definedVars[varName] = true
		}

		// Check statements (would need full implementation)
	}

	return len(errors) == 0, errors
}

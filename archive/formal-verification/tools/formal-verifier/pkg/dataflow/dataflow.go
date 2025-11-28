// Package dataflow implements data flow analysis algorithms
// Based on Document 13: Go Control Flow Formalization, Chapter 4
package dataflow

import (
	"fmt"

	"github.com/your-org/formal-verifier/pkg/cfg"
)

// Direction represents the direction of data flow analysis
type Direction int

const (
	Forward  Direction = iota // Forward analysis (e.g., reaching definitions)
	Backward                  // Backward analysis (e.g., liveness)
)

// Framework represents a data flow analysis framework
type Framework struct {
	Direction Direction
	Meet      MeetOperator
	Transfer  TransferFunction
	Boundary  interface{} // Boundary condition
	Initial   interface{} // Initial value for other nodes
}

// MeetOperator represents the meet operation (∪ or ∩)
type MeetOperator func(values []interface{}) interface{}

// TransferFunction represents the transfer function for a node
type TransferFunction func(node *cfg.Node, input interface{}) interface{}

// Analysis represents a data flow analysis instance
type Analysis struct {
	CFG       *cfg.CFG
	Framework *Framework
	In        map[*cfg.Node]interface{} // IN[n] sets
	Out       map[*cfg.Node]interface{} // OUT[n] sets
}

// NewAnalysis creates a new data flow analysis
func NewAnalysis(cfgGraph *cfg.CFG, framework *Framework) *Analysis {
	return &Analysis{
		CFG:       cfgGraph,
		Framework: framework,
		In:        make(map[*cfg.Node]interface{}),
		Out:       make(map[*cfg.Node]interface{}),
	}
}

// Run executes the data flow analysis
// Generic iterative algorithm (Kildall's algorithm)
func (a *Analysis) Run() {
	// Initialize
	for _, node := range a.CFG.Nodes {
		if node == a.CFG.Entry && a.Framework.Direction == Forward {
			a.In[node] = a.Framework.Boundary
			a.Out[node] = a.Framework.Transfer(node, a.In[node])
		} else if node == a.CFG.Exit && a.Framework.Direction == Backward {
			a.Out[node] = a.Framework.Boundary
			a.In[node] = a.Framework.Transfer(node, a.Out[node])
		} else {
			a.In[node] = a.Framework.Initial
			a.Out[node] = a.Framework.Initial
		}
	}

	// Iterate until fixpoint
	changed := true
	iterations := 0
	maxIterations := 1000 // Prevent infinite loops

	for changed && iterations < maxIterations {
		changed = false
		iterations++

		for _, node := range a.CFG.Nodes {
			// Skip boundary node
			if (node == a.CFG.Entry && a.Framework.Direction == Forward) ||
				(node == a.CFG.Exit && a.Framework.Direction == Backward) {
				continue
			}

			// Compute IN[n] or OUT[n] using meet operator
			var newValue interface{}
			if a.Framework.Direction == Forward {
				// IN[n] = meet(OUT[p] for all predecessors p)
				predValues := make([]interface{}, len(node.Predecessors))
				for i, pred := range node.Predecessors {
					predValues[i] = a.Out[pred]
				}
				newValue = a.Framework.Meet(predValues)

				if !equal(newValue, a.In[node]) {
					changed = true
					a.In[node] = newValue
				}

				// OUT[n] = Transfer(n, IN[n])
				newOut := a.Framework.Transfer(node, a.In[node])
				if !equal(newOut, a.Out[node]) {
					changed = true
					a.Out[node] = newOut
				}
			} else { // Backward
				// OUT[n] = meet(IN[s] for all successors s)
				succValues := make([]interface{}, len(node.Successors))
				for i, succ := range node.Successors {
					succValues[i] = a.In[succ]
				}
				newValue = a.Framework.Meet(succValues)

				if !equal(newValue, a.Out[node]) {
					changed = true
					a.Out[node] = newValue
				}

				// IN[n] = Transfer(n, OUT[n])
				newIn := a.Framework.Transfer(node, a.Out[node])
				if !equal(newIn, a.In[node]) {
					changed = true
					a.In[node] = newIn
				}
			}
		}
	}

	if iterations >= maxIterations {
		fmt.Printf("⚠️  Warning: Data flow analysis did not converge after %d iterations\n", iterations)
	} else {
		fmt.Printf("✅ Data flow analysis converged after %d iterations\n", iterations)
	}
}

// Print prints the analysis results
func (a *Analysis) Print() {
	fmt.Println("=== Data Flow Analysis Results ===")
	for _, node := range a.CFG.Nodes {
		fmt.Printf("Node %d: %s\n", node.ID, node.Label)
		fmt.Printf("  IN:  %v\n", a.In[node])
		fmt.Printf("  OUT: %v\n", a.Out[node])
		fmt.Println()
	}
}

// Helper function to compare values (simplified)
func equal(a, b interface{}) bool {
	// Simplified comparison - would need type-specific logic
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// ===== Specific Data Flow Analyses =====

// VariableSet represents a set of variables
type VariableSet map[string]bool

// NewVariableSet creates a new variable set
func NewVariableSet() VariableSet {
	return make(VariableSet)
}

// Add adds a variable to the set
func (s VariableSet) Add(variable string) {
	s[variable] = true
}

// Remove removes a variable from the set
func (s VariableSet) Remove(variable string) {
	delete(s, variable)
}

// Union computes the union of two sets
func (s VariableSet) Union(other VariableSet) VariableSet {
	result := NewVariableSet()
	for v := range s {
		result.Add(v)
	}
	for v := range other {
		result.Add(v)
	}
	return result
}

// Intersection computes the intersection of two sets
func (s VariableSet) Intersection(other VariableSet) VariableSet {
	result := NewVariableSet()
	for v := range s {
		if other[v] {
			result.Add(v)
		}
	}
	return result
}

// Difference computes the set difference s - other
func (s VariableSet) Difference(other VariableSet) VariableSet {
	result := NewVariableSet()
	for v := range s {
		if !other[v] {
			result.Add(v)
		}
	}
	return result
}

// Copy creates a copy of the set
func (s VariableSet) Copy() VariableSet {
	result := NewVariableSet()
	for v := range s {
		result.Add(v)
	}
	return result
}

// ===== Liveness Analysis =====

// LivenessAnalysis performs liveness analysis
// Formal Definition (from Document 13):
//
//	OUT[n] = ⋃(s∈succ(n)) IN[s]
//	IN[n] = use[n] ∪ (OUT[n] - def[n])
type LivenessAnalysis struct {
	*Analysis
	Use map[*cfg.Node]VariableSet // Variables used in node
	Def map[*cfg.Node]VariableSet // Variables defined in node
}

// NewLivenessAnalysis creates a new liveness analysis
func NewLivenessAnalysis(cfgGraph *cfg.CFG) *LivenessAnalysis {
	la := &LivenessAnalysis{
		Use: make(map[*cfg.Node]VariableSet),
		Def: make(map[*cfg.Node]VariableSet),
	}

	// Compute use and def sets for each node
	for _, node := range cfgGraph.Nodes {
		la.Use[node] = NewVariableSet()
		la.Def[node] = NewVariableSet()
		// TODO: Analyze statements to compute use/def
		// For now, placeholder
	}

	// Create framework
	framework := &Framework{
		Direction: Backward,
		Boundary:  NewVariableSet(),
		Initial:   NewVariableSet(),
		Meet: func(values []interface{}) interface{} {
			// Union meet operator
			result := NewVariableSet()
			for _, v := range values {
				if set, ok := v.(VariableSet); ok {
					result = result.Union(set)
				}
			}
			return result
		},
		Transfer: func(node *cfg.Node, input interface{}) interface{} {
			// IN[n] = use[n] ∪ (OUT[n] - def[n])
			out := input.(VariableSet)
			result := la.Use[node].Union(out.Difference(la.Def[node]))
			return result
		},
	}

	la.Analysis = NewAnalysis(cfgGraph, framework)
	return la
}

// ===== Reaching Definitions Analysis =====

// Definition represents a variable definition
type Definition struct {
	Variable string
	NodeID   int
}

// DefinitionSet represents a set of definitions
type DefinitionSet map[Definition]bool

// NewDefinitionSet creates a new definition set
func NewDefinitionSet() DefinitionSet {
	return make(DefinitionSet)
}

// Add adds a definition to the set
func (s DefinitionSet) Add(def Definition) {
	s[def] = true
}

// Union computes the union of two sets
func (s DefinitionSet) Union(other DefinitionSet) DefinitionSet {
	result := NewDefinitionSet()
	for d := range s {
		result.Add(d)
	}
	for d := range other {
		result.Add(d)
	}
	return result
}

// Difference computes the set difference s - other
func (s DefinitionSet) Difference(other DefinitionSet) DefinitionSet {
	result := NewDefinitionSet()
	for d := range s {
		if !other[d] {
			result.Add(d)
		}
	}
	return result
}

// ReachingDefinitionsAnalysis performs reaching definitions analysis
// Formal Definition (from Document 13):
//
//	OUT[n] = gen[n] ∪ (IN[n] - kill[n])
//	IN[n] = ⋃(p∈pred(n)) OUT[p]
type ReachingDefinitionsAnalysis struct {
	*Analysis
	Gen  map[*cfg.Node]DefinitionSet // Definitions generated in node
	Kill map[*cfg.Node]DefinitionSet // Definitions killed in node
}

// NewReachingDefinitionsAnalysis creates a new reaching definitions analysis
func NewReachingDefinitionsAnalysis(cfgGraph *cfg.CFG) *ReachingDefinitionsAnalysis {
	rda := &ReachingDefinitionsAnalysis{
		Gen:  make(map[*cfg.Node]DefinitionSet),
		Kill: make(map[*cfg.Node]DefinitionSet),
	}

	// Compute gen and kill sets for each node
	for _, node := range cfgGraph.Nodes {
		rda.Gen[node] = NewDefinitionSet()
		rda.Kill[node] = NewDefinitionSet()
		// TODO: Analyze statements to compute gen/kill
		// For now, placeholder
	}

	// Create framework
	framework := &Framework{
		Direction: Forward,
		Boundary:  NewDefinitionSet(),
		Initial:   NewDefinitionSet(),
		Meet: func(values []interface{}) interface{} {
			// Union meet operator
			result := NewDefinitionSet()
			for _, v := range values {
				if set, ok := v.(DefinitionSet); ok {
					result = result.Union(set)
				}
			}
			return result
		},
		Transfer: func(node *cfg.Node, input interface{}) interface{} {
			// OUT[n] = gen[n] ∪ (IN[n] - kill[n])
			in := input.(DefinitionSet)
			result := rda.Gen[node].Union(in.Difference(rda.Kill[node]))
			return result
		},
	}

	rda.Analysis = NewAnalysis(cfgGraph, framework)
	return rda
}

// ===== Available Expressions Analysis =====

// Expression represents an expression
type Expression struct {
	Op   string
	Args []string
}

// ExpressionSet represents a set of expressions
type ExpressionSet map[string]bool

// NewExpressionSet creates a new expression set
func NewExpressionSet() ExpressionSet {
	return make(ExpressionSet)
}

// Add adds an expression to the set
func (s ExpressionSet) Add(expr string) {
	s[expr] = true
}

// Intersection computes the intersection of two sets
func (s ExpressionSet) Intersection(other ExpressionSet) ExpressionSet {
	result := NewExpressionSet()
	for e := range s {
		if other[e] {
			result.Add(e)
		}
	}
	return result
}

// Union computes the union of two sets
func (s ExpressionSet) Union(other ExpressionSet) ExpressionSet {
	result := NewExpressionSet()
	for e := range s {
		result.Add(e)
	}
	for e := range other {
		result.Add(e)
	}
	return result
}

// Difference computes the set difference s - other
func (s ExpressionSet) Difference(other ExpressionSet) ExpressionSet {
	result := NewExpressionSet()
	for e := range s {
		if !other[e] {
			result.Add(e)
		}
	}
	return result
}

// AvailableExpressionsAnalysis performs available expressions analysis
// Formal Definition (from Document 13):
//
//	OUT[n] = gen[n] ∪ (IN[n] - kill[n])
//	IN[n] = ⋂(p∈pred(n)) OUT[p]
type AvailableExpressionsAnalysis struct {
	*Analysis
	Gen  map[*cfg.Node]ExpressionSet // Expressions generated in node
	Kill map[*cfg.Node]ExpressionSet // Expressions killed in node
}

// NewAvailableExpressionsAnalysis creates a new available expressions analysis
func NewAvailableExpressionsAnalysis(cfgGraph *cfg.CFG) *AvailableExpressionsAnalysis {
	aea := &AvailableExpressionsAnalysis{
		Gen:  make(map[*cfg.Node]ExpressionSet),
		Kill: make(map[*cfg.Node]ExpressionSet),
	}

	// Compute gen and kill sets for each node
	for _, node := range cfgGraph.Nodes {
		aea.Gen[node] = NewExpressionSet()
		aea.Kill[node] = NewExpressionSet()
		// TODO: Analyze statements to compute gen/kill
		// For now, placeholder
	}

	// Create universal set (all expressions)
	universal := NewExpressionSet()

	// Create framework
	framework := &Framework{
		Direction: Forward,
		Boundary:  NewExpressionSet(),
		Initial:   universal, // Start with all expressions available
		Meet: func(values []interface{}) interface{} {
			// Intersection meet operator
			if len(values) == 0 {
				return universal
			}
			result := values[0].(ExpressionSet)
			for i := 1; i < len(values); i++ {
				if set, ok := values[i].(ExpressionSet); ok {
					result = result.Intersection(set)
				}
			}
			return result
		},
		Transfer: func(node *cfg.Node, input interface{}) interface{} {
			// OUT[n] = gen[n] ∪ (IN[n] - kill[n])
			in := input.(ExpressionSet)
			result := aea.Gen[node].Union(in.Difference(aea.Kill[node]))
			return result
		},
	}

	aea.Analysis = NewAnalysis(cfgGraph, framework)
	return aea
}

// RunAll runs all common data flow analyses
func RunAll(cfgGraph *cfg.CFG) {
	fmt.Println("=== Running Data Flow Analyses ===")
	fmt.Println()

	// Liveness Analysis
	fmt.Println("📊 Liveness Analysis (Backward)")
	liveness := NewLivenessAnalysis(cfgGraph)
	liveness.Run()
	fmt.Println()

	// Reaching Definitions
	fmt.Println("📊 Reaching Definitions Analysis (Forward)")
	reaching := NewReachingDefinitionsAnalysis(cfgGraph)
	reaching.Run()
	fmt.Println()

	// Available Expressions
	fmt.Println("📊 Available Expressions Analysis (Forward)")
	available := NewAvailableExpressionsAnalysis(cfgGraph)
	available.Run()
	fmt.Println()

	fmt.Println("✅ All data flow analyses completed!")
}

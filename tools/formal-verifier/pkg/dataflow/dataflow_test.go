package dataflow

import (
	"testing"

	"github.com/your-org/formal-verifier/pkg/cfg"
)

func TestLivenessAnalysis(t *testing.T) {
	// Build CFG
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Run liveness analysis
	liveness := NewLivenessAnalysis(cfgGraph)
	liveness.Run()

	// Verify results
	if len(liveness.In) == 0 {
		t.Error("Liveness analysis produced no results")
	}

	t.Logf("Liveness analysis completed successfully")
	t.Logf("Analyzed %d nodes", len(cfgGraph.Nodes))
}

func TestReachingDefinitionsAnalysis(t *testing.T) {
	// Build CFG
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Run reaching definitions analysis
	reaching := NewReachingDefinitionsAnalysis(cfgGraph)
	reaching.Run()

	// Verify results
	if len(reaching.In) == 0 {
		t.Error("Reaching definitions analysis produced no results")
	}

	t.Logf("Reaching definitions analysis completed successfully")
	t.Logf("Analyzed %d nodes", len(cfgGraph.Nodes))
}

func TestAvailableExpressionsAnalysis(t *testing.T) {
	// Build CFG
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Run available expressions analysis
	available := NewAvailableExpressionsAnalysis(cfgGraph)
	available.Run()

	// Verify results
	if len(available.In) == 0 {
		t.Error("Available expressions analysis produced no results")
	}

	t.Logf("Available expressions analysis completed successfully")
	t.Logf("Analyzed %d nodes", len(cfgGraph.Nodes))
}

func TestRunAll(t *testing.T) {
	// Build CFG
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Run all analyses
	RunAll(cfgGraph)

	t.Log("All data flow analyses completed successfully")
}

func TestVariableSet(t *testing.T) {
	s1 := NewVariableSet()
	s1.Add("x")
	s1.Add("y")

	s2 := NewVariableSet()
	s2.Add("y")
	s2.Add("z")

	// Test union
	union := s1.Union(s2)
	if !union["x"] || !union["y"] || !union["z"] {
		t.Error("Union failed")
	}

	// Test intersection
	inter := s1.Intersection(s2)
	if !inter["y"] || len(inter) != 1 {
		t.Error("Intersection failed")
	}

	// Test difference
	diff := s1.Difference(s2)
	if !diff["x"] || len(diff) != 1 {
		t.Error("Difference failed")
	}

	t.Log("Variable set operations work correctly")
}

func BenchmarkLivenessAnalysis(b *testing.B) {
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		b.Fatalf("Failed to build CFG: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		liveness := NewLivenessAnalysis(cfgGraph)
		liveness.Run()
	}
}

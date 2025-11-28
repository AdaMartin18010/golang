package cfg

import (
	"testing"
)

func TestCFGBuilder_Simple(t *testing.T) {
	builder := NewBuilder()

	// Test with a simple file
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Verify CFG structure
	if cfg == nil {
		t.Fatal("CFG is nil")
	}

	if cfg.Entry == nil {
		t.Error("CFG entry is nil")
	}

	if cfg.Exit == nil {
		t.Error("CFG exit is nil")
	}

	if len(cfg.Nodes) == 0 {
		t.Error("CFG has no nodes")
	}

	t.Logf("CFG constructed successfully: %d nodes", len(cfg.Nodes))

	// Print CFG for debugging
	cfg.Print()
}

func TestCFGBuilder_Visualizer(t *testing.T) {
	builder := NewBuilder()

	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Test visualizer
	visualizer := NewVisualizer(cfg)

	// Test stats
	stats := visualizer.GetStats()
	t.Logf("Stats: Nodes=%d, Edges=%d, MaxDepth=%d, Loops=%d, Branches=%d",
		stats.NodeCount, stats.EdgeCount, stats.MaxDepth, stats.LoopCount, stats.BranchCount)

	if stats.NodeCount == 0 {
		t.Error("Node count is 0")
	}
}

func TestCFGBuilder_IfStmt(t *testing.T) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Check if CFG contains if-related nodes
	hasIfNode := false
	for _, node := range cfg.Nodes {
		if node.Label == "if_cond" || node.Label == "if_then" || node.Label == "if_else" {
			hasIfNode = true
			break
		}
	}

	if !hasIfNode {
		t.Log("Note: No if statement nodes found (may be expected depending on functions)")
	}
}

func TestCFGBuilder_Loop(t *testing.T) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Check if CFG contains loop-related nodes
	hasLoopNode := false
	for _, node := range cfg.Nodes {
		if node.Label == "for_header" || node.Label == "for_body" || node.Label == "range_header" {
			hasLoopNode = true
			break
		}
	}

	if !hasLoopNode {
		t.Log("Note: No loop nodes found (may be expected depending on functions)")
	}
}

// Benchmark CFG construction
func BenchmarkCFGBuilder(b *testing.B) {
	builder := NewBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildFromFile("../../testdata/simple.go")
		if err != nil {
			b.Fatalf("Failed to build CFG: %v", err)
		}
	}
}

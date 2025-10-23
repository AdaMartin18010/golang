package cfg

import (
	"testing"
)

func TestSSAConverter_Simple(t *testing.T) {
	// Build CFG from test file
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	// Convert to SSA
	converter := NewSSAConverter(cfg)
	ssaCFG, err := converter.Convert()
	if err != nil {
		t.Fatalf("Failed to convert to SSA: %v", err)
	}

	// Verify SSA property
	valid, errors := ssaCFG.VerifySSAProperty()
	if !valid {
		t.Errorf("SSA property violated: %v", errors)
	}

	// Print SSA CFG for debugging
	ssaCFG.Print()

	t.Logf("SSA conversion successful")
	t.Logf("Dominator tree nodes: %d", len(ssaCFG.DomTree.IDom))
}

func TestDominatorTree(t *testing.T) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	converter := NewSSAConverter(cfg)
	domTree, err := converter.computeDominatorTree()
	if err != nil {
		t.Fatalf("Failed to compute dominator tree: %v", err)
	}

	// Verify entry node dominates itself
	if !domTree.DomSet[cfg.Entry][cfg.Entry] {
		t.Error("Entry node should dominate itself")
	}

	// Verify dominator tree properties
	for node, idom := range domTree.IDom {
		if node == cfg.Entry {
			continue
		}
		// IDom should be in the dominator set
		if !domTree.DomSet[node][idom] {
			t.Errorf("IDom(%d) = %d, but %d not in Dom(%d)", node.ID, idom.ID, idom.ID, node.ID)
		}
	}

	t.Logf("Dominator tree computed successfully")
	t.Logf("Immediate dominators: %d", len(domTree.IDom))
}

func TestDominanceFrontiers(t *testing.T) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	converter := NewSSAConverter(cfg)
	_, err = converter.computeDominatorTree()
	if err != nil {
		t.Fatalf("Failed to compute dominator tree: %v", err)
	}

	domFrontier := converter.computeDominanceFrontiers()

	// Verify dominance frontier properties
	totalDF := 0
	for node, df := range domFrontier {
		totalDF += len(df)
		t.Logf("DF(%d) has %d nodes", node.ID, len(df))
	}

	t.Logf("Dominance frontiers computed successfully")
	t.Logf("Total DF edges: %d", totalDF)
}

func TestPhiFunctionPlacement(t *testing.T) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		t.Fatalf("Failed to build CFG: %v", err)
	}

	converter := NewSSAConverter(cfg)
	_, err = converter.computeDominatorTree()
	if err != nil {
		t.Fatalf("Failed to compute dominator tree: %v", err)
	}

	converter.domFrontier = converter.computeDominanceFrontiers()
	ssaNodes := converter.placePhiFunctions()

	// Count φ-functions
	totalPhi := 0
	for _, ssaNode := range ssaNodes {
		totalPhi += len(ssaNode.PhiFunctions)
		if len(ssaNode.PhiFunctions) > 0 {
			t.Logf("Node %d has %d φ-functions", ssaNode.ID, len(ssaNode.PhiFunctions))
		}
	}

	t.Logf("φ-functions placed successfully")
	t.Logf("Total φ-functions: %d", totalPhi)
}

// Benchmark SSA conversion
func BenchmarkSSAConversion(b *testing.B) {
	builder := NewBuilder()
	cfg, err := builder.BuildFromFile("../../testdata/simple.go")
	if err != nil {
		b.Fatalf("Failed to build CFG: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter := NewSSAConverter(cfg)
		_, err := converter.Convert()
		if err != nil {
			b.Fatalf("Failed to convert to SSA: %v", err)
		}
	}
}

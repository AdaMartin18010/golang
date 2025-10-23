// Package cfg - Visualizer for Control Flow Graphs
// Generates DOT format output for Graphviz visualization
package cfg

import (
	"fmt"
	"go/ast"
	"go/printer"
	"os"
	"strings"
)

// Visualizer generates visual representations of CFGs
type Visualizer struct {
	cfg *CFG
}

// NewVisualizer creates a new CFG visualizer
func NewVisualizer(cfg *CFG) *Visualizer {
	return &Visualizer{cfg: cfg}
}

// ExportDOT exports the CFG to DOT format (Graphviz)
// DOT is a graph description language for visualization
func (v *Visualizer) ExportDOT(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create DOT file: %w", err)
	}
	defer file.Close()

	// Write DOT header
	fmt.Fprintln(file, "digraph CFG {")
	fmt.Fprintln(file, "  rankdir=TB;")
	fmt.Fprintln(file, "  node [shape=box, style=rounded];")
	fmt.Fprintln(file, "  edge [fontsize=10];")
	fmt.Fprintln(file, "")

	// Special styling for entry and exit nodes
	fmt.Fprintf(file, "  node%d [label=\"%s\", shape=circle, style=filled, fillcolor=lightgreen];\n",
		v.cfg.Entry.ID, v.cfg.Entry.Label)
	fmt.Fprintf(file, "  node%d [label=\"%s\", shape=circle, style=filled, fillcolor=lightcoral];\n",
		v.cfg.Exit.ID, v.cfg.Exit.Label)
	fmt.Fprintln(file, "")

	// Write nodes
	for _, node := range v.cfg.Nodes {
		if node == v.cfg.Entry || node == v.cfg.Exit {
			continue // Already written with special styling
		}

		label := v.formatNodeLabel(node)
		fillColor := v.getNodeColor(node)

		fmt.Fprintf(file, "  node%d [label=\"%s\", fillcolor=%s, style=filled];\n",
			node.ID, label, fillColor)
	}
	fmt.Fprintln(file, "")

	// Write edges
	for _, node := range v.cfg.Nodes {
		for _, succ := range node.Successors {
			edgeLabel := v.getEdgeLabel(node, succ)
			if edgeLabel != "" {
				fmt.Fprintf(file, "  node%d -> node%d [label=\"%s\"];\n",
					node.ID, succ.ID, edgeLabel)
			} else {
				fmt.Fprintf(file, "  node%d -> node%d;\n",
					node.ID, succ.ID)
			}
		}
	}

	// Write DOT footer
	fmt.Fprintln(file, "}")

	return nil
}

// formatNodeLabel creates a formatted label for a node
func (v *Visualizer) formatNodeLabel(node *Node) string {
	var builder strings.Builder

	// Node label
	builder.WriteString(node.Label)

	// Add statements if any
	if len(node.Stmts) > 0 {
		builder.WriteString("\\n")
		builder.WriteString(strings.Repeat("-", 20))
		builder.WriteString("\\n")

		for i, stmt := range node.Stmts {
			if i > 0 {
				builder.WriteString("\\n")
			}
			stmtStr := v.formatStatement(stmt)
			// Limit statement length for readability
			if len(stmtStr) > 40 {
				stmtStr = stmtStr[:37] + "..."
			}
			builder.WriteString(stmtStr)
		}
	}

	return builder.String()
}

// formatStatement formats an AST statement for display
func (v *Visualizer) formatStatement(stmt ast.Stmt) string {
	var builder strings.Builder
	printer.Fprint(&builder, v.cfg.FSet, stmt)

	// Escape special characters for DOT format
	result := builder.String()
	result = strings.ReplaceAll(result, "\"", "\\\"")
	result = strings.ReplaceAll(result, "\n", "\\n")
	result = strings.ReplaceAll(result, "\t", " ")

	return result
}

// getNodeColor returns the color for a node based on its type
func (v *Visualizer) getNodeColor(node *Node) string {
	label := node.Label

	switch {
	case strings.HasPrefix(label, "if_"):
		if strings.Contains(label, "cond") {
			return "lightyellow"
		} else if strings.Contains(label, "then") {
			return "lightblue"
		} else if strings.Contains(label, "else") {
			return "lightcyan"
		}
		return "white"

	case strings.HasPrefix(label, "for_"), strings.HasPrefix(label, "range_"):
		if strings.Contains(label, "header") {
			return "lightpink"
		} else if strings.Contains(label, "body") {
			return "lavender"
		}
		return "white"

	case strings.HasPrefix(label, "switch_"), strings.HasPrefix(label, "select_"):
		if strings.Contains(label, "header") {
			return "lightgoldenrodyellow"
		} else if strings.Contains(label, "case") {
			return "lightsteelblue"
		}
		return "white"

	case strings.Contains(label, "return"):
		return "lightsalmon"

	case strings.Contains(label, "go_"):
		return "lightgreen"

	case strings.Contains(label, "defer"):
		return "peachpuff"

	default:
		return "white"
	}
}

// getEdgeLabel returns a label for an edge (for conditional branches)
func (v *Visualizer) getEdgeLabel(from, to *Node) string {
	// Check if this is a conditional branch
	if strings.Contains(from.Label, "if_cond") {
		// Try to determine if this is the true or false branch
		if strings.Contains(to.Label, "then") {
			return "true"
		} else if strings.Contains(to.Label, "else") || strings.Contains(to.Label, "merge") {
			return "false"
		}
	}

	// Check if this is a loop back edge
	if strings.Contains(from.Label, "body") && strings.Contains(to.Label, "header") {
		return "back"
	}
	if strings.Contains(from.Label, "post") && strings.Contains(to.Label, "header") {
		return "back"
	}

	return ""
}

// ExportJSON exports the CFG to JSON format
func (v *Visualizer) ExportJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "{")
	fmt.Fprintf(file, "  \"nodes\": [\n")

	// Write nodes
	for i, node := range v.cfg.Nodes {
		if i > 0 {
			fmt.Fprintln(file, ",")
		}
		fmt.Fprintf(file, "    {\n")
		fmt.Fprintf(file, "      \"id\": %d,\n", node.ID)
		fmt.Fprintf(file, "      \"label\": \"%s\",\n", node.Label)
		fmt.Fprintf(file, "      \"stmt_count\": %d", len(node.Stmts))
		fmt.Fprintf(file, "\n    }")
	}

	fmt.Fprintln(file, "\n  ],")
	fmt.Fprintf(file, "  \"edges\": [\n")

	// Write edges
	firstEdge := true
	for _, node := range v.cfg.Nodes {
		for _, succ := range node.Successors {
			if !firstEdge {
				fmt.Fprintln(file, ",")
			}
			firstEdge = false
			fmt.Fprintf(file, "    {\"from\": %d, \"to\": %d}", node.ID, succ.ID)
		}
	}

	fmt.Fprintln(file, "\n  ],")
	fmt.Fprintf(file, "  \"entry\": %d,\n", v.cfg.Entry.ID)
	fmt.Fprintf(file, "  \"exit\": %d\n", v.cfg.Exit.ID)
	fmt.Fprintln(file, "}")

	return nil
}

// ExportHTML exports the CFG to an interactive HTML visualization
func (v *Visualizer) ExportHTML(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	// Write HTML header
	fmt.Fprintln(file, `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>CFG Visualization</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
        }
        #graph {
            border: 1px solid #ccc;
            background-color: #f9f9f9;
        }
        .node {
            cursor: pointer;
        }
        .node circle {
            fill: #69b3a2;
            stroke: #333;
            stroke-width: 2px;
        }
        .node.entry circle {
            fill: #90EE90;
        }
        .node.exit circle {
            fill: #FFB6C6;
        }
        .node text {
            font-size: 12px;
            text-anchor: middle;
        }
        .link {
            fill: none;
            stroke: #999;
            stroke-width: 2px;
            marker-end: url(#arrowhead);
        }
        .link.back-edge {
            stroke: #ff6b6b;
            stroke-dasharray: 5,5;
        }
        h1 {
            color: #333;
        }
        .stats {
            margin: 20px 0;
            padding: 10px;
            background-color: #e8f4f8;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <h1>Control Flow Graph Visualization</h1>
    <div class="stats">
        <strong>Statistics:</strong> 
        Total Nodes: <span id="nodeCount"></span> | 
        Total Edges: <span id="edgeCount"></span> | 
        Entry: <span id="entryNode"></span> | 
        Exit: <span id="exitNode"></span>
    </div>
    <svg id="graph" width="1200" height="800"></svg>
    <script>`)

	// Write graph data
	fmt.Fprintln(file, "const graphData = {")
	fmt.Fprintln(file, "  nodes: [")

	for i, node := range v.cfg.Nodes {
		if i > 0 {
			fmt.Fprintln(file, ",")
		}
		nodeType := "normal"
		if node == v.cfg.Entry {
			nodeType = "entry"
		} else if node == v.cfg.Exit {
			nodeType = "exit"
		}
		fmt.Fprintf(file, "    {id: %d, label: \"%s\", type: \"%s\"}",
			node.ID, node.Label, nodeType)
	}

	fmt.Fprintln(file, "\n  ],")
	fmt.Fprintln(file, "  links: [")

	firstEdge := true
	for _, node := range v.cfg.Nodes {
		for _, succ := range node.Successors {
			if !firstEdge {
				fmt.Fprintln(file, ",")
			}
			firstEdge = false
			edgeLabel := v.getEdgeLabel(node, succ)
			isBackEdge := edgeLabel == "back"
			fmt.Fprintf(file, "    {source: %d, target: %d, label: \"%s\", backEdge: %t}",
				node.ID, succ.ID, edgeLabel, isBackEdge)
		}
	}

	fmt.Fprintln(file, "\n  ]")
	fmt.Fprintln(file, "};")

	// Write D3.js visualization code
	fmt.Fprintln(file, `
// Update statistics
document.getElementById('nodeCount').textContent = graphData.nodes.length;
document.getElementById('edgeCount').textContent = graphData.links.length;
document.getElementById('entryNode').textContent = graphData.nodes.find(n => n.type === 'entry').label;
document.getElementById('exitNode').textContent = graphData.nodes.find(n => n.type === 'exit').label;

const svg = d3.select("#graph");
const width = +svg.attr("width");
const height = +svg.attr("height");

// Define arrowhead marker
svg.append("defs").append("marker")
    .attr("id", "arrowhead")
    .attr("viewBox", "-0 -5 10 10")
    .attr("refX", 20)
    .attr("refY", 0)
    .attr("orient", "auto")
    .attr("markerWidth", 8)
    .attr("markerHeight", 8)
    .append("svg:path")
    .attr("d", "M 0,-5 L 10 ,0 L 0,5")
    .attr("fill", "#999");

// Create force simulation
const simulation = d3.forceSimulation(graphData.nodes)
    .force("link", d3.forceLink(graphData.links).id(d => d.id).distance(100))
    .force("charge", d3.forceManyBody().strength(-500))
    .force("center", d3.forceCenter(width / 2, height / 2));

// Draw links
const link = svg.append("g")
    .selectAll("line")
    .data(graphData.links)
    .enter().append("line")
    .attr("class", d => d.backEdge ? "link back-edge" : "link");

// Draw nodes
const node = svg.append("g")
    .selectAll("g")
    .data(graphData.nodes)
    .enter().append("g")
    .attr("class", d => "node " + d.type)
    .call(d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended));

node.append("circle")
    .attr("r", 20);

node.append("text")
    .attr("dy", 4)
    .text(d => d.label.length > 10 ? d.label.substring(0, 10) + "..." : d.label)
    .style("fill", "#333");

// Add title for full label on hover
node.append("title")
    .text(d => d.label);

// Update positions on simulation tick
simulation.on("tick", () => {
    link
        .attr("x1", d => d.source.x)
        .attr("y1", d => d.source.y)
        .attr("x2", d => d.target.x)
        .attr("y2", d => d.target.y);

    node
        .attr("transform", d => "translate(" + d.x + "," + d.y + ")");
});

// Drag functions
function dragstarted(event, d) {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
}

function dragged(event, d) {
    d.fx = event.x;
    d.fy = event.y;
}

function dragended(event, d) {
    if (!event.active) simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
}
    </script>
</body>
</html>`)

	return nil
}

// Stats returns statistics about the CFG
type Stats struct {
	NodeCount   int
	EdgeCount   int
	MaxDepth    int
	LoopCount   int
	BranchCount int
}

// GetStats calculates statistics for the CFG
func (v *Visualizer) GetStats() *Stats {
	stats := &Stats{
		NodeCount: len(v.cfg.Nodes),
	}

	// Count edges
	for _, node := range v.cfg.Nodes {
		stats.EdgeCount += len(node.Successors)
	}

	// Count loops (nodes with back edges)
	for _, node := range v.cfg.Nodes {
		if strings.Contains(node.Label, "for") || strings.Contains(node.Label, "range") {
			stats.LoopCount++
		}
	}

	// Count branches (nodes with multiple successors)
	for _, node := range v.cfg.Nodes {
		if len(node.Successors) > 1 {
			stats.BranchCount++
		}
	}

	// Calculate max depth (BFS from entry)
	stats.MaxDepth = v.calculateMaxDepth()

	return stats
}

// calculateMaxDepth calculates the maximum depth of the CFG using BFS
func (v *Visualizer) calculateMaxDepth() int {
	if v.cfg.Entry == nil {
		return 0
	}

	visited := make(map[*Node]bool)
	depth := make(map[*Node]int)
	queue := []*Node{v.cfg.Entry}
	visited[v.cfg.Entry] = true
	depth[v.cfg.Entry] = 0
	maxDepth := 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		currentDepth := depth[node]
		if currentDepth > maxDepth {
			maxDepth = currentDepth
		}

		for _, succ := range node.Successors {
			if !visited[succ] {
				visited[succ] = true
				depth[succ] = currentDepth + 1
				queue = append(queue, succ)
			}
		}
	}

	return maxDepth
}

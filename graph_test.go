package main

import (
	"testing"
)

// TestGraph tests graph
func TestGraph(t *testing.T) {
	g := newGraph()
	g.addNode(newGraphNode(0))
	g.addEdge(newGraphEdge(0, 4, 2.9))
	g.addEdge(newGraphEdge(0, 5, 1.0))
	g.addNode(newGraphNode(1))
	g.addEdge(newGraphEdge(1, 2, 3.1))
	g.show()
}

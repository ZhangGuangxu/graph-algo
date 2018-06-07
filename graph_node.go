package main

const (
	invalidNodeIndex = -1
)

func isValidNodeIndex(idx int) bool {
	return idx >= 0
}

type graphNode struct {
	Index int
}

func newGraphNode(idx int) graphNode {
	return graphNode{Index: idx}
}

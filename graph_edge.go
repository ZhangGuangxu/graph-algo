package main

type graphEdge struct {
	From int
	To   int
	Cost float32
}

func newGraphEdgeDefault() graphEdge {
	return graphEdge{From: invalidNodeIndex, To: invalidNodeIndex, Cost: 1.0}
}

func newGraphEdge(f int, t int, c float32) graphEdge {
	return graphEdge{From: f, To: t, Cost: c}
}

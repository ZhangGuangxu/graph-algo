package main

import (
	"errors"
	"fmt"
	"github.com/ZhangGuangxu/circularqueue"
	"github.com/ZhangGuangxu/stack"
)

var errInvalidNodeIndex = errors.New("invalid node index")
var errPathNotFound = errors.New("path not found")
var errEdgeTypeWrong = errors.New("edge type wrong")

type graphEdges []graphEdge

type graph struct {
	nodes     []graphNode
	edges     []graphEdges
	nextIndex int
}

func newGraph() *graph {
	return &graph{}
}

func (g *graph) addNode(n graphNode) {
	if n.Index != g.nextIndex {
		return
	}

	g.nodes = append(g.nodes, n)
	g.nextIndex++

	g.edges = append(g.edges, graphEdges{})
}

func (g *graph) addEdge(e graphEdge) {
	if !isValidNodeIndex(e.From) || !isValidNodeIndex(e.To) {
		return
	}
	if e.From >= len(g.nodes) {
		return
	}

	g.edges[e.From] = append(g.edges[e.From], e)
}

func (g *graph) show() {
	for i, n := range g.nodes {
		if !isValidNodeIndex(n.Index) {
			continue
		}
		fmt.Printf("%d-> ", n.Index)
		for _, edge := range g.edges[i] {
			//fmt.Printf("f:%d, t:%d, c:%v; ", edge.From, edge.To, edge.Cost)
			fmt.Printf("%v; ", edge)
		}
		fmt.Println()
	}
	fmt.Println()
}

// dft is deep first traverse.
func (g *graph) dft() {

}

func reversePath(path []graphEdge) []graphEdge {
	length := len(path)
	half := length / 2
	for i := 0; i < half; i++ {
		j := length - (i + 1)
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// dfs is deep first search.
func (g *graph) dfs(begin graphNode, end graphNode) ([]graphEdge, error) {
	b := begin.Index
	if b < 0 || b >= len(g.nodes) {
		return nil, errInvalidNodeIndex
	}

	e := end.Index
	if e < 0 || e >= len(g.nodes) {
		return nil, errInvalidNodeIndex
	}

	if b == e {
		return []graphEdge{}, nil
	}

	s := stack.NewStack()
	for _, tmp := range g.edges[b] {
		s.Push(tmp)
	}
	record := make(map[int]int) // 用于记录曾经加入过栈的边，key是To, value是From
	record[b] = b               // 起始点比较特殊

	for !s.IsEmpty() {
		tmpEdge, err := s.Pop()
		if err != nil {
			return nil, err
		}
		edge, ok := tmpEdge.(graphEdge)
		if !ok {
			return nil, errEdgeTypeWrong
		}

		if edge.To == e {
			var path []graphEdge
			for {
				path = append(path, edge)
				if edge.From == b {
					return reversePath(path), nil
				}

				edge = graphEdge{From: record[edge.From], To: edge.From}
			}
		}

		if _, ok := record[edge.To]; ok {
			continue
		}
		for _, tmp := range g.edges[edge.To] {
			s.Push(tmp)
		}
		record[edge.To] = edge.From
	}

	return nil, errPathNotFound
}

func (g *graph) bfs(begin graphNode, end graphNode) ([]graphEdge, error) {
	b := begin.Index
	if b < 0 || b >= len(g.nodes) {
		return nil, errInvalidNodeIndex
	}

	e := end.Index
	if e < 0 || e >= len(g.nodes) {
		return nil, errInvalidNodeIndex
	}

	if b == e {
		return []graphEdge{}, nil
	}

	q := circularqueue.NewCircularQueue()
	for _, tmp := range g.edges[b] {
		q.Push(tmp)
	}
	record := make(map[int]int)
	record[b] = b

	for !q.IsEmpty() {
		tmp, err := q.Pop()
		if err != nil {
			return nil, err
		}
		edge, ok := tmp.(graphEdge)
		if !ok {
			return nil, errEdgeTypeWrong
		}

		if edge.To == e {
			var path []graphEdge
			for {
				path = append(path, edge)
				if edge.From == b {
					return reversePath(path), nil
				}

				edge = graphEdge{From: record[edge.From], To: edge.From}
			}
		}

		if _, ok := record[edge.To]; ok {
			continue
		}
		for _, tmp := range g.edges[edge.To] {
			q.Push(tmp)
		}
		record[edge.To] = edge.From
	}

	return nil, errPathNotFound
}

package main

import (
	"fmt"
	"github.com/ZhangGuangxu/stack"
	"sort"
)

var pn = fmt.Println
var pf = fmt.Printf

func main() {
	pn("directional")
	{
		md := mapData{}
		err := md.load("bin/a.map")
		if err != nil {
			err = md.load("./a.map")
			if err != nil {
				err = md.load("../../bin/a.map")
				if err != nil {
					pf("load a.map got error %v\n", err)
					return
				}
			}
		}
		md.show()
		pn()

		var allIndex sort.IntSlice
		for k := range md.edgesMap {
			allIndex = append(allIndex, k)
		}
		sort.Sort(allIndex)

		g := newGraph()
		for _, from := range allIndex {
			g.addNode(newGraphNode(from))
			edges := md.edgesMap[from]
			for to, edge := range edges {
				g.addEdge(newGraphEdge(from, to, edge.Cost))
			}
		}
		g.show()

		path, err := g.dfs(newGraphNode(0), newGraphNode(1))
		if err != nil {
			pn(err)
			return
		}
		pf("dfs %d->%d    %v\n", 0, 1, path)

		path, err = g.bfs(newGraphNode(0), newGraphNode(2))
		if err != nil {
			pn(err)
			return
		}
		pf("bfs %d->%d    %v\n", 0, 2, path)
	}

	{
		g := newGraph()
		g.addNode(newGraphNode(0))
		g.addEdge(newGraphEdge(0, 4, 2.9))
		g.addEdge(newGraphEdge(0, 5, 1.0))
		g.addNode(newGraphNode(1))
		g.addEdge(newGraphEdge(1, 2, 3.1))
		g.show()
	}

	{
		s := stack.NewStack()
		s.Push(1)
		s.Pop()
	}

	pn("non-directional Bidirectional BFS")
	{
		md := mapData{}
		err := md.load("bin/a2.map")
		if err != nil {
			err = md.load("./a2.map")
			if err != nil {
				err = md.load("../../bin/a2.map")
				if err != nil {
					pf("load a.map got error %v\n", err)
					return
				}
			}
		}
		md.show()
		pn()

		var allIndex sort.IntSlice
		for k := range md.edgesMap {
			allIndex = append(allIndex, k)
		}
		sort.Sort(allIndex)

		m := make(map[int]map[int]float32)

		// non-directional
		g := newGraph()
		for _, from := range allIndex {
			g.addNode(newGraphNode(from))

			edges := md.edgesMap[from]
			for to, edge := range edges {
				g.addEdge(newGraphEdge(from, to, edge.Cost))
				v, ok := m[to]
				if ok {
					v[from] = edge.Cost
				} else {
					m[to] = make(map[int]float32)
					m[to][from] = edge.Cost
				}
			}

			for to, cost := range m[from] {
				g.addEdge(newGraphEdge(from, to, cost))
				delete(m[from], to)
			}
		}
		g.show()

		bs := newBiBFS(g, newGraphNode(0), newGraphNode(5))
		pn(bs.search())
		pn(bs.index)
	}

	pn("directional")
	{
		md := mapData{}
		err := md.load("bin/a.map")
		if err != nil {
			err = md.load("./a.map")
			if err != nil {
				err = md.load("../../bin/a.map")
				if err != nil {
					pf("load a.map got error %v\n", err)
					return
				}
			}
		}
		md.show()
		pn()

		var allIndex sort.IntSlice
		for k := range md.edgesMap {
			allIndex = append(allIndex, k)
		}
		sort.Sort(allIndex)

		g := newGraph()
		for _, from := range allIndex {
			g.addNode(newGraphNode(from))
			edges := md.edgesMap[from]
			for to, edge := range edges {
				g.addEdge(newGraphEdge(from, to, edge.Cost))
			}
		}
		g.show()

		d := NewDijkstra(g, 4, 2)
		d.Search()
		pn(d.PathToTarget())
	}

	pn("directional astar")
	{
		md := mapData{}
		err := md.load("bin/a.map")
		if err != nil {
			err = md.load("./a.map")
			if err != nil {
				err = md.load("../../bin/a.map")
				if err != nil {
					pf("load a.map got error %v\n", err)
					return
				}
			}
		}
		md.show()
		pn()

		var allIndex sort.IntSlice
		for k := range md.edgesMap {
			allIndex = append(allIndex, k)
		}
		sort.Sort(allIndex)

		g := newGraph()
		for _, from := range allIndex {
			g.addNode(newGraphNode(from))
			edges := md.edgesMap[from]
			for to, edge := range edges {
				g.addEdge(newGraphEdge(from, to, edge.Cost))
			}
		}
		g.show()

		d := NewAstar(g, 4, 2)
		d.Search()
		pn(d.PathToTarget())
	}
}

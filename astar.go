package main

// Astar algorithm
type Astar struct {
	graph  *graph
	source int
	target int
	hFn    func(nd1, nd2 int) float32 // heuristic

	frontier map[int]graphEdge // search frontier
	// fcost = gcost + hcost (heuristic)
	fcost map[int]float32   // cost to target
	gcost map[int]float32   // cost to some node
	spt   map[int]graphEdge // shortest path tree

	err error
}

// NewAstar returns a instance of Dijkstra.
func NewAstar(g *graph, s, t int) *Astar {
	return &Astar{
		graph:    g,
		source:   s,
		target:   t,
		hFn:      func(nd1, nd2 int) float32 { return 0 },
		frontier: make(map[int]graphEdge),
		fcost:    make(map[int]float32),
		gcost:    make(map[int]float32),
		spt:      make(map[int]graphEdge),
	}
}

// NewAstarWithH returns a instance of Dijkstra.
func NewAstarWithH(g *graph, s, t int, h func(nd1, nd2 int) float32) *Astar {
	return &Astar{
		graph:    g,
		source:   s,
		target:   t,
		hFn:      h,
		frontier: make(map[int]graphEdge),
		fcost:    make(map[int]float32),
		gcost:    make(map[int]float32),
		spt:      make(map[int]graphEdge),
	}
}

// Search trys to find the shortest path from source to target.
// source is a node index, same as target.
func (d *Astar) Search() {
	d.gcost[d.source] = 0
	d.fcost[d.source] = 0
	pq := NewIndexedPriorityQueueMin(d.fcost)
	pq.Insert(d.source)

	for !pq.IsEmpty() {
		i, err := pq.Pop()
		if err != nil {
			d.err = err
			return
		}

		edge, ok := d.frontier[i]
		if ok {
			d.spt[i] = edge
			i = edge.To
		}

		if i == d.target {
			return
		}

		for _, e := range d.graph.edges[i] {
			t := e.To
			g := d.gcost[i] + e.Cost
			f := g + d.hFn(t, d.target)
			if _, ok := d.frontier[t]; !ok {
				d.frontier[t] = e
				d.gcost[t] = g
				d.fcost[t] = f
				pq.Insert(t)
			} else if g < d.gcost[t] {
				if _, ok := d.spt[t]; !ok {
					d.frontier[t] = e
					d.gcost[t] = g
					d.fcost[t] = f
					pq.ChangePriority(t)
				}
			}
		}
	}
}

// PathToTarget returns shortest path from source to target.
func (d *Astar) PathToTarget() ([]graphEdge, error) {
	if d.err != nil {
		return []graphEdge{}, d.err
	}

	var path []graphEdge
	idx := d.target
	for {
		if idx == d.source {
			break
		}
		e, ok := d.spt[idx]
		if !ok {
			break
		}
		path = append(path, e)
		idx = e.From
	}

	return reversePath(path), nil
}

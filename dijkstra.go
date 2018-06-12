package main

// Dijkstra algorithm
type Dijkstra struct {
	graph  *graph
	source int
	target int

	frontier map[int]graphEdge // search frontier
	cost     map[int]float32   // cost to some node
	spt      map[int]graphEdge // shortest path tree

	err error
}

// NewDijkstra returns a instance of Dijkstra.
func NewDijkstra(g *graph, s, t int) *Dijkstra {
	return &Dijkstra{
		graph:    g,
		source:   s,
		target:   t,
		frontier: make(map[int]graphEdge),
		cost:     make(map[int]float32),
		spt:      make(map[int]graphEdge),
	}
}

// Search trys to find the shortest path from source to target.
// source is a node index, same as target.
func (d *Dijkstra) Search() {
	d.cost[d.source] = 0
	pq := NewIndexedPriorityQueueMin(d.cost)
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
			newCost := d.cost[i] + e.Cost
			t := e.To
			if _, ok := d.frontier[t]; !ok {
				d.frontier[t] = e
				d.cost[t] = newCost
				pq.Insert(t)
			} else if newCost < d.cost[t] {
				if _, ok := d.spt[t]; !ok {
					d.frontier[t] = e
					d.cost[t] = newCost
					pq.ChangePriority(t)
				}
			}
		}
	}
}

// PathToTarget returns shortest path from source to target.
func (d *Dijkstra) PathToTarget() ([]graphEdge, error) {
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

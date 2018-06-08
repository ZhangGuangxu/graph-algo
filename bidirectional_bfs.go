package main

import (
	"github.com/ZhangGuangxu/circularqueue"
	"log"
	"sync"
)

type biBFS struct {
	graph  *graph
	source graphNode
	target graphNode

	wg *sync.WaitGroup
	s  *bfs
	rs *rbfs

	stop  chan bool
	index int
}

func newBiBFS(graph *graph, source graphNode, target graphNode) *biBFS {
	s := &biBFS{
		graph:  graph,
		source: source,
		target: target,
		wg:     &sync.WaitGroup{},
		stop:   make(chan bool),
		index:  invalidNodeIndex,
	}
	s.s = newBFS(s)
	s.rs = newRBFS(s)
	return s
}

func (s *biBFS) search() ([]graphEdge, error) {
	b := s.source.Index
	e := s.target.Index
	if b == e {
		return []graphEdge{}, nil
	}

	s.s.start()
	s.rs.start()
	s.wg.Wait()

	if s.index == invalidNodeIndex {
		path, err := s.s.result()
		if err != nil {
			log.Println(err)

			path, err = s.rs.result()
			if err != nil {
				log.Println(err)
				return nil, err
			}

			return path, nil
		}

		return path, nil
	}

	path1, err := s.s.result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	path2, err := s.rs.result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("%v, %v\n", path1, path2)
	len1 := len(path1)
	path := make([]graphEdge, len1+len(path2))
	copy(path, path1)
	copy(path[len1:], path2)
	return path, nil
}

func (s *biBFS) checkJoin(idx int) bool {
	if s.rs.checkJoin(idx) {
		s.index = idx
		close(s.stop)
		return true
	}
	return false
}

func (s *biBFS) shouldStop() bool {
	select {
	case _, ok := <-s.stop:
		if !ok {
			return true
		}
	default:
	}
	return false
}

type bfs struct {
	parent *biBFS
	record map[int]int
	err    error
}

func newBFS(parent *biBFS) *bfs {
	return &bfs{
		parent: parent,
		record: make(map[int]int),
	}
}

func (s *bfs) start() {
	s.parent.wg.Add(1)
	go s.run()
}

func (s *bfs) run() {
	defer s.parent.wg.Done()

	g := s.parent.graph

	b := s.parent.source.Index
	if b < 0 || b >= len(g.nodes) {
		s.err = errInvalidNodeIndex
		return
	}

	e := s.parent.target.Index
	if e < 0 || e >= len(g.nodes) {
		s.err = errInvalidNodeIndex
		return
	}

	q := circularqueue.NewCircularQueue()
	for _, tmp := range g.edges[b] {
		q.Push(tmp)
	}
	s.record[b] = b

	for !q.IsEmpty() {
		tmp, err := q.Pop()
		if err != nil {
			s.err = err
			return
		}
		edge, ok := tmp.(graphEdge)
		if !ok {
			s.err = errEdgeTypeWrong
			return
		}

		if s.parent.checkJoin(edge.To) {
			s.record[edge.To] = edge.From
			return
		}

		if edge.To == e {
			s.record[edge.To] = edge.From
			return
		}

		if _, ok := s.record[edge.To]; ok {
			continue
		}
		for _, tmp := range g.edges[edge.To] {
			q.Push(tmp)
		}
		s.record[edge.To] = edge.From
	}

	s.err = errPathNotFound
}

func (s *bfs) result() ([]graphEdge, error) {
	if s.err != nil {
		return nil, s.err
	}

	idx := s.parent.index
	if idx == invalidNodeIndex {
		idx = s.parent.target.Index
	}
	edge := graphEdge{From: s.record[idx], To: idx}
	b := s.parent.source.Index
	var path []graphEdge

	for {
		path = append(path, edge)
		if edge.From == b {
			return reversePath(path), nil
		}

		edge = graphEdge{From: s.record[edge.From], To: edge.From}
	}
}

// The leading r represents reverse, which means search begins with target node.
type rbfs struct {
	parent *biBFS

	mx     *sync.Mutex
	record map[int]int

	err error
}

func newRBFS(parent *biBFS) *rbfs {
	return &rbfs{
		parent: parent,
		mx:     &sync.Mutex{},
		record: make(map[int]int),
	}
}

func (s *rbfs) start() {
	s.parent.wg.Add(1)
	go s.run()
}

func (s *rbfs) run() {
	defer s.parent.wg.Done()

	g := s.parent.graph

	b := s.parent.target.Index
	if b < 0 || b >= len(g.nodes) {
		s.err = errInvalidNodeIndex
		return
	}

	e := s.parent.source.Index
	if e < 0 || e >= len(g.nodes) {
		s.err = errInvalidNodeIndex
		return
	}

	q := circularqueue.NewCircularQueue()
	for _, tmp := range g.edges[b] {
		q.Push(tmp)
	}
	s.addRecord(b, b)

	for !q.IsEmpty() {
		if s.parent.shouldStop() {
			return
		}

		tmp, err := q.Pop()
		if err != nil {
			s.err = err
			return
		}
		edge, ok := tmp.(graphEdge)
		if !ok {
			s.err = errEdgeTypeWrong
			return
		}

		if edge.To == e {
			s.addRecord(edge.To, edge.From)
			return
		}

		if s.hasRecord(edge.To) {
			continue
		}
		for _, tmp := range g.edges[edge.To] {
			q.Push(tmp)
		}
		s.addRecord(edge.To, edge.From)
	}

	s.err = errPathNotFound
}

func (s *rbfs) checkJoin(idx int) bool {
	s.mx.Lock()
	_, ok := s.record[idx]
	s.mx.Unlock()
	return ok
}

func (s *rbfs) hasRecord(idx int) bool {
	_, ok := s.record[idx]
	return ok
}

func (s *rbfs) addRecord(to, from int) {
	s.mx.Lock()
	s.record[to] = from
	s.mx.Unlock()
}

func (s *rbfs) result() ([]graphEdge, error) {
	if s.err != nil {
		return nil, s.err
	}

	idx := s.parent.index
	if idx == invalidNodeIndex {
		idx = s.parent.source.Index
	}
	edge := graphEdge{From: idx, To: s.record[idx]}
	e := s.parent.target.Index
	var path []graphEdge

	for {
		path = append(path, edge)
		if edge.To == e {
			return path, nil
		}

		edge = graphEdge{From: edge.To, To: s.record[edge.To]}
	}
}

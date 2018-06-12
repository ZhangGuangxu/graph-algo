package main

import (
	"errors"
	"fmt"
)

// ErrEmptyHeap tells us the heap is empty.
var ErrEmptyHeap = errors.New("empty heap")

// ErrCostNotExist tells us the cost to some node is no filled yet.
var ErrCostNotExist = errors.New("cost not exist")

// ErrNodeIndexNotExist tells us some node index does not exist.
var ErrNodeIndexNotExist = errors.New("node index not exist")

// ErrInvalidItemIndex tells us some item index is invalid.
var ErrInvalidItemIndex = errors.New("invalid item index")

const (
	invalidTail = -1
)

// IndexedPriorityQueueMin is a min-heap.
// data-item is index to graphNode.
// compare policy is based on cost.
type IndexedPriorityQueueMin struct {
	cost                 map[int]float32 // key is index to graphNode, value is cost
	nodeIndexToItemIndex map[int]int     // key is index to graphNode, value is index to data item
	way                  int
	data                 []int
	tail                 int // the index of the last item in data
}

// NewIndexedPriorityQueueMin returns an instance of IndexedPriorityQueueMin.
func NewIndexedPriorityQueueMin(cost map[int]float32) *IndexedPriorityQueueMin {
	return NewIndexedPriorityQueueMinWithNWayAndSize(cost, 2, 1)
}

// NewIndexedPriorityQueueMinWithNWay returns an instance of IndexedPriorityQueueMin.
func NewIndexedPriorityQueueMinWithNWay(cost map[int]float32, nWay int) *IndexedPriorityQueueMin {
	return NewIndexedPriorityQueueMinWithNWayAndSize(cost, nWay, 1)
}

// NewIndexedPriorityQueueMinWithSize returns an instance of IndexedPriorityQueueMin.
func NewIndexedPriorityQueueMinWithSize(cost map[int]float32, s int) *IndexedPriorityQueueMin {
	return NewIndexedPriorityQueueMinWithNWayAndSize(cost, 2, s)
}

// NewIndexedPriorityQueueMinWithNWayAndSize returns an instance of IndexedPriorityQueueMin with init-size.
func NewIndexedPriorityQueueMinWithNWayAndSize(cost map[int]float32, nWay int, s int) *IndexedPriorityQueueMin {
	return &IndexedPriorityQueueMin{
		cost:                 cost,
		nodeIndexToItemIndex: make(map[int]int, s),
		way:                  nWay,
		data:                 make([]int, s),
		tail:                 invalidTail,
	}
}

// nodeIndexA is smaller than nodeIndexB
func (h *IndexedPriorityQueueMin) compare(nodeIndexA, nodeIndexB int) bool {
	costA, ok := h.cost[nodeIndexA]
	if !ok {
		panic(ErrCostNotExist)
	}
	costB, ok := h.cost[nodeIndexB]
	if !ok {
		panic(ErrCostNotExist)
	}
	return costA > costB
}

// IsEmpty returns true when heap is empty.
func (h *IndexedPriorityQueueMin) IsEmpty() bool {
	return h.tail == invalidTail
}

// Insert inserts an item into heap.
func (h *IndexedPriorityQueueMin) Insert(x int) {
	if h.tail+1 >= len(h.data) {
		h.makeSpace()
	}

	h.tail++
	h.data[h.tail] = x
	h.nodeIndexToItemIndex[x] = h.tail
	h.siftUp(h.tail)
}

func (h *IndexedPriorityQueueMin) makeSpace() {
	d := make([]int, len(h.data)*2+1)
	copy(d, h.data)
	h.data = d
}

// ChangePriority changes the priority of nodeIndex.
func (h *IndexedPriorityQueueMin) ChangePriority(nodeIndex int) {
	itemIndex, ok := h.nodeIndexToItemIndex[nodeIndex]
	if !ok {
		panic(ErrNodeIndexNotExist)
	}

	if itemIndex < 0 || itemIndex > h.tail {
		panic(ErrInvalidItemIndex)
	}

	if !h.siftUp(itemIndex) {
		h.siftDown(itemIndex)
	}
}

// Pop removes the root node from the heap and returns that node.
// Pop returns error when heap is empty. So you had better make sure
// the heap is not empty before you invode Pop on it.
func (h *IndexedPriorityQueueMin) Pop() (int, error) {
	if h.IsEmpty() {
		return 0, ErrEmptyHeap
	}

	v := h.data[0]
	h.data[0] = h.data[h.tail]
	delete(h.nodeIndexToItemIndex, v)
	h.nodeIndexToItemIndex[h.data[h.tail]] = 0
	h.tail--
	h.siftDown(0)
	return v, nil
}

func (h *IndexedPriorityQueueMin) siftUp(begin int) (swap bool) {
	idx := begin
	parentIdx := 0

	for {
		if idx == 0 {
			return
		}

		r := idx % h.way
		if r == 0 {
			parentIdx = (idx - h.way) / h.way
		} else {
			parentIdx = (idx - r) / h.way
		}
		if h.compare(h.data[parentIdx], h.data[idx]) {
			h.nodeIndexToItemIndex[h.data[parentIdx]] = idx
			h.nodeIndexToItemIndex[h.data[idx]] = parentIdx
			h.data[parentIdx], h.data[idx] = h.data[idx], h.data[parentIdx]
			idx = parentIdx
			swap = true
		} else {
			return
		}
	}
}

func (h *IndexedPriorityQueueMin) siftDown(begin int) (swap bool) {
	if h.tail <= 0 {
		return
	}

	idx := begin

	for {
		idx, swap = h.compareWithChildren(idx)
		if !swap {
			return
		}
	}
}

func (h *IndexedPriorityQueueMin) compareWithChildren(idx int) (newIdx int, swap bool) {
	newIdx = idx
	max := h.data[idx]

	for a := 1; a <= h.way; a++ {
		i := idx*h.way + a
		if i > h.tail {
			break
		}
		if h.compare(max, h.data[i]) {
			newIdx = i
			max = h.data[i]
		}
	}

	swap = newIdx != idx
	if swap {
		h.nodeIndexToItemIndex[h.data[idx]] = newIdx
		h.nodeIndexToItemIndex[h.data[newIdx]] = idx
		h.data[idx], h.data[newIdx] = h.data[newIdx], h.data[idx]
	}
	return
}

func (h *IndexedPriorityQueueMin) show() {
	for i := 0; i <= h.tail; i++ {
		fmt.Printf("%d ", h.data[i])
	}
	fmt.Println()
}

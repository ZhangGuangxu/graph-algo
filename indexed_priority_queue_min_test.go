package main

import (
	"testing"
)

func TestMaxIntBinaryHeap(t *testing.T) {
	{
		cost := make(map[int]float32)
		cost[1] = 1.9
		cost[2] = 5.0
		cost[3] = 4.1
		cost[4] = 0
		cost[5] = 3.0

		a := NewIndexedPriorityQueueMin(cost)
		if !a.IsEmpty() {
			t.Error("a is not empty")
		}

		h := NewIndexedPriorityQueueMin(cost)
		h.Insert(4)
		if h.IsEmpty() {
			t.Error("h should not be empty after Insert")
		}

		i, _ := h.Pop()
		if i != 4 {
			t.Errorf("h.Pop() got %d, want %d", i, 4)
		}
		h.Insert(1)
		h.Insert(5)

		i, _ = h.Pop()
		if i != 1 {
			t.Errorf("h.Pop got %d, want %d", i, 1)
		}
		h.Insert(2)

		i, _ = h.Pop()
		if i != 5 {
			t.Errorf("h.Pop() got %d, want %d", i, 5)
		}
		h.Insert(3)

		i, _ = h.Pop()
		if i != 3 {
			t.Errorf("h.Pop() got %d, want %d", i, 3)
		}

		i, _ = h.Pop()
		if i != 2 {
			t.Errorf("h.Pop() got %d, want %d", i, 2)
		}

		if !h.IsEmpty() {
			t.Error("h should be empty")
		}
	}

	{
		cost := make(map[int]float32)
		//cost[1] = 1.9
		//cost[2] = 4.2
		//cost[3] = 4.1
		cost[4] = 0
		//cost[5] = 3.0

		h := NewIndexedPriorityQueueMin(cost)
		h.Insert(4)

		i, _ := h.Pop()
		if i != 4 {
			t.Errorf("h.Pop() got %d, want %d", i, 4)
		}
		cost[1] = 1.9
		h.Insert(1)
		cost[5] = 3.0
		h.Insert(5)

		i, _ = h.Pop()
		if i != 1 {
			t.Errorf("h.Pop got %d, want %d", i, 1)
		}
		cost[2] = 5.0
		h.Insert(2)

		i, _ = h.Pop()
		if i != 5 {
			t.Errorf("h.Pop() got %d, want %d", i, 5)
		}
		cost[3] = 4.1
		h.Insert(3)

		i, _ = h.Pop()
		if i != 3 {
			t.Errorf("h.Pop() got %d, want %d", i, 3)
		}
		cost[2] = 4.2
		h.ChangePriority(2)

		i, _ = h.Pop()
		if i != 2 {
			t.Errorf("h.Pop() got %d, want %d", i, 2)
		}

		if !h.IsEmpty() {
			t.Error("h should be empty")
		}
	}
}

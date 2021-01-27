package heap

import (
	"testing"
)

type blockData struct {
	key      string
	priority int
}

// This data is static, and never mutated by our heap
var testBlockData []blockData = []blockData{
	{"a", 0},
	{"b", 1},
	{"c", 2},
	{"d", 3},
	{"e", 4},
	{"f", 5},
	{"g", 6},
	{"h", 7},
	{"i", 8},
	{"j", 9},
}

type blockDataPtr int

func (p blockDataPtr) Order() int {
	return testBlockData[p].priority
}

var blockDataHeap []Orderable

func resetBlockDataHeap() {
	blockDataHeap = make([]Orderable, len(testBlockData))
	for i := range testBlockData {
		blockDataHeap[i] = blockDataPtr(i)
	}
}

func TestEnclosed(t *testing.T) {

	t.Run("Heapify", func(t *testing.T) {
		resetBlockDataHeap()
		subject, _ := Heapify(blockDataHeap, -1)
		assertHeapOrdering(subject, t)
	})

	t.Run("Pop", func(t *testing.T) {
		resetBlockDataHeap()
		subject, _ := Heapify(blockDataHeap, -1)

		prev := -1
		for subject.Size() > 0 {
			nextItem, _ := subject.Pop()
			ptr := nextItem.(blockDataPtr)
			priority := testBlockData[ptr].priority
			if priority < prev {
				t.Logf(
					"Popped item of priority %d before item of %d, violating heap ordering invariant",
					priority,
					prev,
				)
				prev = priority
			}
		}
	})

}

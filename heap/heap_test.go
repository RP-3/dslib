package heap

import (
	"math"
	"math/rand"
	"testing"
)

// test implementation
type intHeap []int

func (h intHeap) Less(i, j int) bool { return h[i] <= h[j] }
func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h intHeap) Peak() Any          { return h[0] }
func (h *intHeap) Push(v Any)        { (*h) = append(*h, v.(int)) }
func (h *intHeap) Pop() Any {
	result := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]
	return result
}

func newIntHeap(storage intHeap, maxSize int) *Heap {
	return NewHeap(&storage, maxSize)
}

func TestEmpty(t *testing.T) {
	subject := newIntHeap([]int{}, -1)
	assertInt(subject.storage.Len(), 0, t)
}

func TestUnbounded(t *testing.T) {
	var unboundedHeap = func() *Heap {
		return newIntHeap([]int{}, -1)
	}

	t.Run("Capacity", func(t *testing.T) {
		subject := unboundedHeap()
		assertInt(subject.Capacity(), maxInt, t)
	})

	t.Run("Push", func(t *testing.T) {
		t.Run("when the heap is empty", func(t *testing.T) {
			subject, item := unboundedHeap(), 1
			subject.Push(item)

			// increases in size
			assertInt(subject.Size(), 1, t)

			// places new item at the head
			obj, ok := subject.Peak()
			assertBool(ok, true, t)
			assertBool(obj.(int) == item, true, t)
		})

		t.Run("when the heap has a lower-priority item at the head", func(t *testing.T) {
			subject := unboundedHeap()
			a, b := 1, 2
			subject.Push(a)
			subject.Push(b)

			// increases in size
			assertInt(subject.Size(), 2, t)

			// places higher-priority item at tail
			item, ok := subject.Peak()
			assertBool(ok, true, t)
			assertInt(item.(int), 1, t)
		})

		t.Run("when the heap has a higher-priority item at the head", func(t *testing.T) {
			subject := unboundedHeap()
			a, b := 1, 2
			subject.Push(b)
			subject.Push(a)

			// Increases in size
			assertInt(subject.Size(), 2, t)

			// does not replace the head item
			item, ok := subject.Peak()
			assertBool(ok, true, t)
			assertInt(item.(int), 1, t)
		})

		t.Run("when the newest item requires just one swap", func(t *testing.T) {
			subject := unboundedHeap()
			subject.Push(4)
			subject.Push(5)
			subject.Push(8)
			subject.Push(6)
			subject.Push(9)
			subject.Push(9)
			subject.Push(7)
			assertHeapOrdering(subject, t)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Run("when the heap is empty", func(t *testing.T) {
			subject := unboundedHeap()

			// returns nil
			_, exists := subject.Pop()
			assertBool(exists, false, t)
		})

		t.Run("when the heap has a single item", func(t *testing.T) {
			subject, item := unboundedHeap(), 1
			subject.Push(item)

			// returns the correct item
			obj, ok := subject.Pop()
			assertBool(ok, true, t)
			assertBool(obj.(int) == item, true, t)
		})

		t.Run("when the heap contains both higher and lower priority items", func(t *testing.T) {
			subject := unboundedHeap()
			subject.Push(0)
			subject.Push(5)
			subject.Push(1)
			subject.Push(4)
			subject.Push(3)

			// should contain all items as expected
			assertInt(subject.Size(), 5, t)

			subject.Push(2) // should sift to the middle
			assertInt(subject.Size(), 6, t)

			// sorts items by their given order
			lastVal := math.MinInt64
			for subject.Size() > 0 {
				assertHeapOrdering(subject, t)
				top, ok := subject.Pop()
				assertBool(ok, true, t)
				assertBool(top.(int) > lastVal, true, t)
				lastVal = top.(int)
			}
		})
	})
}

func TestFixedSize(t *testing.T) {
	heapSize := 5
	var fixedHeap = func() *Heap {
		return NewHeap(&intHeap{}, heapSize)
	}

	t.Run("Capacity", func(t *testing.T) {
		subject := fixedHeap()
		assertInt(subject.Capacity(), heapSize, t)
	})

	t.Run("Push", func(t *testing.T) {
		t.Run("when <= size items are inserted", func(t *testing.T) {
			subject := fixedHeap()
			subject.Push(1)
			subject.Push(5)
			subject.Push(2)
			subject.Push(4)
			subject.Push(3)

			assertInt(subject.Size(), heapSize, t)
		})

		t.Run("when > size items are inserted", func(t *testing.T) {
			subject := fixedHeap()
			subject.Push(0)
			subject.Push(5)
			subject.Push(1)
			subject.Push(4)
			subject.Push(3)

			item, overFlowed := subject.Push(2)

			// it does not exceed max size
			assertInt(subject.Size(), heapSize, t)

			// it retains the min items
			sortedContents := make([]int, 0, 5)
			for subject.Size() > 0 {
				assertHeapOrdering(subject, t)
				item, ok := subject.Pop()
				assertBool(ok, true, t)
				sortedContents = append(sortedContents, item.(int))
			}
			assertSlice(sortedContents, []int{1, 2, 3, 4, 5}, t) // zero is missing

			assertBool(overFlowed, true, t)
			assertInt(item.(int), 0, t)
		})
	})
}

func TestRobustness(t *testing.T) {
	heapSize := -1 // unbounded
	testSize := 200
	popPercent := 25

	t.Run("heap ordering robustness", func(t *testing.T) {
		subject := NewHeap(&intHeap{}, heapSize)
		for i := 0; i < testSize; i++ {
			if rand.Intn(100) > popPercent {
				item := rand.Int()
				subject.Push(item)
			} else {
				subject.Pop()
			}
			assertHeapOrdering(subject, t)
		}
	})
}

func TestHeapify(t *testing.T) {
	t.Run("when the provided slice is empty", func(t *testing.T) {
		subject, discarded := Heapify(&intHeap{}, -1)
		assertHeapOrdering(subject, t)
		assertInt(len(discarded), 0, t) // nothing discarded
	})

	t.Run("when the provided heap has items within it", func(t *testing.T) {
		nums := intHeap{1, 9, 2, 8, 3, 7, 4, 6, 5, 4, 6, 3, 7, 2, 8, 1, 9}
		subject, discarded := Heapify(&nums, -1)
		assertInt(subject.Capacity(), maxInt, t)

		// generates a valid heap out of the given slice
		assertHeapOrdering(subject, t)

		assertInt(len(discarded), 0, t) // nothing discarded
	})

	t.Run("when the provided heap is larger than the specified size", func(t *testing.T) {
		nums := intHeap{1, 9, 2, 8, 3, 7, 4} // seven numbers

		subject, discarded := Heapify(&nums, 5)
		assertHeapOrdering(subject, t) // valid
		assertInt(subject.Capacity(), 5, t)

		// should remove 1 and 2 (the smallest two)
		sortedContents := make([]int, 0, 5)
		for subject.Size() > 0 {
			assertHeapOrdering(subject, t)
			item, ok := subject.Pop()
			assertBool(ok, true, t)
			sortedContents = append(sortedContents, item.(int))
		}
		assertSlice(sortedContents, []int{3, 4, 7, 8, 9}, t) // zero is missing

		assertInt(discarded[0].(int), 1, t)
		assertInt(discarded[1].(int), 2, t)
	})
}

// test helpers
func assertInt(a int, b int, t *testing.T) {
	if a != b {
		t.Errorf("Expected %d but got %d\n", a, b)
	}
}

func assertBool(a bool, b bool, t *testing.T) {
	if a != b {
		t.Errorf("Expected %t but got %t\n", a, b)
	}
}

func assertNil(a Any, t *testing.T) {
	if a != nil {
		t.Errorf("Expected nil but it wasn't\n")
	}
}

func assertSlice(a, b []int, t *testing.T) {
	if len(a) != len(b) {
		t.Error("Slice lengths are not equal")
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Elements at %d differ. %d vs %d\n", i, a[i], b[i])
		}
	}
}

func assertHeapOrdering(heap *Heap, t *testing.T) {
	storageLen := heap.storage.Len()
	for i := 0; i < storageLen/2; i++ {
		left, right := i*2+1, i*2+2
		if left < storageLen {
			assertBool(heap.storage.Less(i, left), true, t)
		}
		if right < storageLen {
			assertBool(heap.storage.Less(i, right), true, t)
		}
	}
}

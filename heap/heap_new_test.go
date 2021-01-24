package heap

import (
	"math"
	"math/rand"
	"testing"
)

func TestEmpty(t *testing.T) {
	subject := NewHeap(-1)
	_, exists := subject.Peak()
	assertBool(exists, false, t)
	assertInt(subject.Size(), 0, t)
}

func TestUnbounded(t *testing.T) {
	var unboundedHeap = func() *Heap {
		return NewHeap(-1)
	}

	t.Run("Capacity", func(t *testing.T) {
		subject := unboundedHeap()
		assertInt(subject.Capacity(), maxInt, t)
	})

	t.Run("Push", func(t *testing.T) {
		t.Run("when the heap is empty", func(t *testing.T) {
			subject, item := unboundedHeap(), testItem{1}
			subject.Push(item)

			// increases in size
			assertInt(subject.Size(), 1, t)

			// places new item at the head
			obj, ok := subject.Peak()
			assertBool(ok, true, t)
			assertBool(equal(obj, item), true, t)
		})

		t.Run("when the heap has a lower-priority item at the head", func(t *testing.T) {
			subject := unboundedHeap()
			a, b := testItem{1}, testItem{2}
			subject.Push(a)
			subject.Push(b)

			// increases in size
			assertInt(subject.Size(), 2, t)

			// places higher-priority item at tail
			item, ok := subject.Peak()
			assertBool(ok, true, t)
			assertInt(item.Order(), 1, t)
		})

		t.Run("when the heap has a higher-priority item at the head", func(t *testing.T) {
			subject := unboundedHeap()
			a, b := testItem{1}, testItem{2}
			subject.Push(b)
			subject.Push(a)

			// Increases in size
			assertInt(subject.Size(), 2, t)

			// does not replace the head item
			item, ok := subject.Peak()
			assertBool(ok, true, t)
			assertInt(item.Order(), 1, t)
		})

		t.Run("when the newest item requires just one swap", func(t *testing.T) {
			subject := unboundedHeap()
			subject.Push(testItem{4})
			subject.Push(testItem{5})
			subject.Push(testItem{8})
			subject.Push(testItem{6})
			subject.Push(testItem{9})
			subject.Push(testItem{9})
			subject.Push(testItem{7})
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
			subject, item := unboundedHeap(), testItem{1}
			subject.Push(item)

			// returns the correct item
			obj, ok := subject.Pop()
			assertBool(ok, true, t)
			assertBool(equal(obj, item), true, t)
		})

		t.Run("when the heap contains both higher and lower priority items", func(t *testing.T) {
			subject := unboundedHeap()
			subject.Push(testItem{key: 0})
			subject.Push(testItem{key: 5})
			subject.Push(testItem{key: 1})
			subject.Push(testItem{key: 4})
			subject.Push(testItem{key: 3})

			// should contain all items as expected
			assertInt(subject.Size(), 5, t)

			subject.Push(testItem{key: 2}) // should sift to the middle
			assertInt(subject.Size(), 6, t)

			// sorts items by their given order
			lastVal := math.MinInt64
			for subject.Size() > 0 {
				assertHeapOrdering(subject, t)
				top, ok := subject.Pop()
				assertBool(ok, true, t)
				assertBool(top.Order() > lastVal, true, t)
				lastVal = top.Order()
			}
		})
	})
}

func TestFixedSize(t *testing.T) {
	heapSize := 5
	var fixedHeap = func() *Heap {
		return NewHeap(heapSize)
	}

	t.Run("Capacity", func(t *testing.T) {
		subject := fixedHeap()
		assertInt(subject.Capacity(), heapSize, t)
	})

	t.Run("Push", func(t *testing.T) {
		t.Run("when <= size items are inserted", func(t *testing.T) {
			subject := fixedHeap()
			subject.Push(testItem{key: 1})
			subject.Push(testItem{key: 5})
			subject.Push(testItem{key: 2})
			subject.Push(testItem{key: 4})
			subject.Push(testItem{key: 3})

			assertInt(subject.Size(), heapSize, t)
		})

		t.Run("when > size items are inserted", func(t *testing.T) {
			subject := fixedHeap()
			subject.Push(testItem{key: 0})
			subject.Push(testItem{key: 5})
			subject.Push(testItem{key: 1})
			subject.Push(testItem{key: 4})
			subject.Push(testItem{key: 3})

			item, overFlowed := subject.Push(testItem{key: 2})

			// it does not exceed max size
			assertInt(subject.Size(), heapSize, t)

			// it retains the min items
			sortedContents := make([]int, 0, 5)
			for subject.Size() > 0 {
				assertHeapOrdering(subject, t)
				item, ok := subject.Pop()
				assertBool(ok, true, t)
				sortedContents = append(sortedContents, item.Order())
			}
			assertSlice(sortedContents, []int{1, 2, 3, 4, 5}, t) // zero is missing

			assertBool(overFlowed, true, t)
			assertInt(item.Order(), 0, t)
		})
	})
}

func TestRobustness(t *testing.T) {
	heapSize := -1 // unbounded
	testSize := 200
	popPercent := 25

	t.Run("heap ordering robustness", func(t *testing.T) {
		subject := NewHeap(heapSize)
		for i := 0; i < testSize; i++ {
			if rand.Intn(100) > popPercent {
				item := testItem{key: rand.Int()}
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
		subject, discarded := Heapify(make([]Orderable, 0), -1)
		assertHeapOrdering(subject, t)
		assertInt(len(discarded), 0, t) // nothing discarded
	})

	t.Run("when the provided heap has items within it", func(t *testing.T) {
		nums := []Orderable{
			testItem{key: 1},
			testItem{key: 9},
			testItem{key: 2},
			testItem{key: 8},
			testItem{key: 3},
			testItem{key: 7},
			testItem{key: 4},
			testItem{key: 6},
			testItem{key: 5},
			testItem{key: 4},
			testItem{key: 6},
			testItem{key: 3},
			testItem{key: 7},
			testItem{key: 2},
			testItem{key: 8},
			testItem{key: 1},
			testItem{key: 9},
		}
		subject, discarded := Heapify(nums, -1)
		assertInt(subject.Capacity(), maxInt, t)

		// generates a valid heap out of the given slice
		assertHeapOrdering(subject, t)

		assertInt(len(discarded), 0, t) // nothing discarded
	})

	t.Run("when the provided heap is larger than the specified size", func(t *testing.T) {
		nums := []Orderable{ // seven numbers
			testItem{key: 1},
			testItem{key: 9},
			testItem{key: 2},
			testItem{key: 8},
			testItem{key: 3},
			testItem{key: 7},
			testItem{key: 4},
		}

		subject, discarded := Heapify(nums, 5)
		assertHeapOrdering(subject, t) // valid
		assertInt(subject.Capacity(), 5, t)

		// should remove 1 and 2 (the smallest two)
		sortedContents := make([]int, 0, 5)
		for subject.Size() > 0 {
			assertHeapOrdering(subject, t)
			item, ok := subject.Pop()
			assertBool(ok, true, t)
			sortedContents = append(sortedContents, item.Order())
		}
		assertSlice(sortedContents, []int{3, 4, 7, 8, 9}, t) // zero is missing

		assertInt(discarded[0].Order(), 1, t)
		assertInt(discarded[1].Order(), 2, t)
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

func assertNil(a *Orderable, t *testing.T) {
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

type testItem struct {
	key int
}

func (t testItem) Order() int {
	return t.key
}

func equal(a Orderable, b testItem) bool {
	obj, coerced := a.(testItem)
	if !coerced {
		return false
	}
	return obj.key == b.key
}

func assertHeapOrdering(heap *Heap, t *testing.T) {
	storage := heap.storage
	for i, item := range storage {
		left, right := i*2+1, i*2+2
		if left < len(storage) {
			assertBool(storage[left].Order() >= item.Order(), true, t)
		}
		if right < len(storage) {
			assertBool(storage[right].Order() >= item.Order(), true, t)
		}
	}
}

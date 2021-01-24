package heap

import (
	"math"
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoHeap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoHeap Suite")
}

var _ = Describe("Heap", func() {
	var subject *Heap

	Describe("empty state inspection", func() {
		BeforeEach(func() {
			subject = NewHeap(-1)
		})

		Describe("peak", func() {
			It("returns nil", func() {
				_, exists := subject.Peak()
				Expect(exists).To(Equal(false))
			})
		})

		Describe("Size", func() {
			It("returns 0", func() {
				Expect(subject.Size()).To(Equal(0))
			})
		})
	})

	Context("when size is unbounded", func() {

		Describe("Push", func() {
			BeforeEach(func() {
				subject = NewHeap(-1)
			})

			Context("when the heap is empty", func() {
				item := testItem{1}
				BeforeEach(func() {
					Expect(subject.Size()).To(Equal(0))
					subject.Push(item)
				})

				It("increases in size", func() {
					Expect(subject.Size()).To(Equal(1))
				})

				It("places the new item at the head", func() {
					obj, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(equal(obj, item)).To(Equal(true))
				})
			})

			Context("when the heap has a lower-priority item at the head", func() {
				a, b := testItem{1}, testItem{2}
				BeforeEach(func() {
					subject.Push(a)
					subject.Push(b)
				})

				It("Increases in size", func() {
					Expect(subject.Size()).To(Equal(2))
				})

				It("does not replace the head item", func() {
					item, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(item.Order()).To(Equal(1))
				})
			})

			Context("when the heap has a higher-priority item at the head", func() {
				a, b := testItem{1}, testItem{2}
				BeforeEach(func() {
					subject.Push(b)
					subject.Push(a)
				})

				It("Increases in size", func() {
					Expect(subject.Size()).To(Equal(2))
				})

				It("does not replace the head item", func() {
					item, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(item.Order()).To(Equal(1))
				})
			})

			Context("when the newest item requires just one swap", func() {
				BeforeEach(func() {
					subject.Push(testItem{4})
					subject.Push(testItem{5})
					subject.Push(testItem{8})
					subject.Push(testItem{6})
					subject.Push(testItem{9})
					subject.Push(testItem{9})
					subject.Push(testItem{7})
				})

				It("does not violate the heap ordering property", func() {
					assertHeapOrdering(subject)
				})
			})
		})

		Describe("Pop", func() {
			BeforeEach(func() {
				subject = NewHeap(-1)
			})

			Context("when the heap is empty", func() {
				It("returns nil", func() {
					_, exists := subject.Pop()
					Expect(exists).To(Equal(false))
				})
			})

			Context("when the heap has a single item", func() {
				item := testItem{1}
				BeforeEach(func() {
					subject.Push(item)
				})

				It("returns that item", func() {
					obj, ok := subject.Pop()
					Expect(ok).To(Equal(true))
					Expect(equal(obj, item)).To(Equal(true))
				})
			})

			Context("when the heap contains both higher and lower priority items", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 0})
					subject.Push(testItem{key: 5})
					subject.Push(testItem{key: 1})
					subject.Push(testItem{key: 4})
					subject.Push(testItem{key: 3})
					Expect(subject.Size()).To(Equal(5))
					subject.Push(testItem{key: 2}) // should sift to the middle
					Expect(subject.Size()).To(Equal(6))
				})

				It("sorts items by their given order", func() {
					lastVal := math.MinInt64
					for subject.Size() > 0 {
						assertHeapOrdering(subject)
						top, ok := subject.Pop()
						Expect(ok).To(Equal(true))
						Expect(top.Order() > lastVal).To(Equal(true))
						lastVal = top.Order()
					}
				})
			})
		})
	})

	Context("when a size is specified", func() {
		heapSize := 5

		Describe("Push", func() {
			BeforeEach(func() {
				subject = NewHeap(heapSize)
			})

			Context("when <= size items are inserted", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 1})
					subject.Push(testItem{key: 5})
					subject.Push(testItem{key: 2})
					subject.Push(testItem{key: 4})
					subject.Push(testItem{key: 3})
				})

				It("allows all items to exist inside", func() {
					Expect(subject.Size()).To(Equal(heapSize))
				})
			})

			Context("when additional items are inserted", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 0})
					subject.Push(testItem{key: 5})
					subject.Push(testItem{key: 1})
					subject.Push(testItem{key: 4})
					subject.Push(testItem{key: 3})
					Expect(subject.Size()).To(Equal(heapSize))
				})

				It("does not exceed maximum size", func() {
					subject.Push(testItem{key: 2})
					Expect(subject.Size()).To(Equal(heapSize))
				})

				It("retains the lower-priority items", func() {
					subject.Push(testItem{key: 2})
					sortedContents := make([]int, 0, 5)
					for subject.Size() > 0 {
						assertHeapOrdering(subject)
						item, ok := subject.Pop()
						Expect(ok).To(Equal(true))
						sortedContents = append(sortedContents, item.Order())
					}
					Expect(sortedContents).To(BeEquivalentTo([]int{1, 2, 3, 4, 5})) // zero is missing
				})

				It("ejects the highest-priority item", func() {
					item, overflowed := subject.Push(testItem{key: 2})
					Expect(overflowed).To(Equal(true))
					Expect(item.Order()).To(Equal(0))
				})
			})
		})
	})

	Describe("robustness", func() {

		heapSize := -1 // unbounded
		testSize := 200
		popPercent := 25

		Describe("heap ordering", func() {
			BeforeEach(func() {
				subject = NewHeap(heapSize)
			})

			It("never violates the heap ordering property", func() {
				for i := 0; i < testSize; i++ {
					if rand.Intn(100) > popPercent {
						item := testItem{key: rand.Int()}
						subject.Push(item)
					} else {
						subject.Pop()
					}
					assertHeapOrdering(subject)

				}
			})
		})
	})

	Describe("heapify", func() {
		Context("when the provided slice is empty", func() {
			BeforeEach(func() {
				subject = Heapify(make([]Orderable, 0), -1)
			})

			It("generates a valid (albeit empty) heap out of the given slice", func() {
				assertHeapOrdering(subject)
			})
		})

		Context("when the provided heap has items within it", func() {
			BeforeEach(func() {
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
				subject = Heapify(nums, -1)
			})

			It("generates a valid heap out of the given slice", func() {
				assertHeapOrdering(subject)
			})
		})
	})
})

// testing helpers
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

func assertHeapOrdering(heap *Heap) {
	storage := heap.UnsafeStorage()
	for i, item := range storage {
		left, right := i*2+1, i*2+2
		if left < len(storage) {
			Expect(storage[left].Order() >= item.Order()).To(Equal(true))
		}
		if right < len(storage) {
			Expect(storage[right].Order() >= item.Order()).To(Equal(true))
		}
	}
}

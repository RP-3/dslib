package heap

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

// Orderable defines the properties that any item must have
// to be heap-ordered.
type Orderable interface {
	// Order dictates the internal ordering of the items in the heap. Heap is
	// min-ordered, so lowest-order items have the highest priority
	Order() int
}

// Heap is a priority queue (min-heap), which maintains the lowest-ordered
// values at the front. Items contained within the heap must conform to the
// Orderable interface, which just means they must have an integer Order
// which can be used to compare items. Heaps can be pushed to or popped
// from. Popping always yields the lowest-orederd item in the heap, and
// pushing adds an item to the heap. If the heap is at capacity, pushing
// ejects and returns the lowest-ordered item inside the heap. For a general
// description see https://en.wikipedia.org/wiki/Heap_(data_structure).
// For a detailed explanation see https://bradfieldcs.com/algos/trees/priority-queues-with-binary-heaps/
type Heap struct {
	storage []Orderable
	maxSize int
}

// NewHeap returns a Heap of the specified size. If size <= 0
// heap size is unbounded.
func NewHeap(maxSize int) *Heap {
	if maxSize <= 0 {
		return &Heap{maxSize: maxInt}
	}
	return &Heap{maxSize: maxSize}
}

// Heapify returns a Heap of the specified size using the given source
// slice as its backing storage, and heap-sorts it in <= O(n) time. If
// the given heap is larger than the specified size the second return
// value contains the lowest-ordered values in the heap, which have
// been discarded
func Heapify(source []Orderable, maxSize int) (*Heap, []Orderable) {
	result := &Heap{storage: source, maxSize: maxSize}
	result.heapify()

	discarded := []Orderable{}

	if maxSize > 0 {
		for result.Size() > maxSize {
			excessVal, _ := result.Pop()
			discarded = append(discarded, excessVal)
		}
	}

	return result, discarded
}

// Push adds an item to the heap.
// The second return val, if true, indicates that the heap is at its
// maximum capacity the highest priority item was popped and returned
// to you as the first return val
func (h *Heap) Push(val Orderable) (Orderable, bool) {
	h.storage = append(h.storage, val)
	h.percolateUp(len(h.storage) - 1)
	if len(h.storage) > h.maxSize {
		return h.Pop()
	}
	return nil, false
}

// UnsafeStorage yields a shallow copy of the underlying storage of the heap.
// The behaviour following mutation of the copy or its pointers is undefined
func (h *Heap) UnsafeStorage() []Orderable {
	result := make([]Orderable, 0, len(h.storage))
	copy(result, h.storage)
	return result
}

// Pop removes the highest priority item from the heap.
// The second return val, if false, indicates that the heap is empty
// and that a nil value was returned to you as the first return val
func (h *Heap) Pop() (Orderable, bool) {
	switch len(h.storage) {
	case 0:
		return nil, false
	case 1:
		return h.removeLast(), true
	default:
		result := h.storage[0]
		h.storage[0] = h.removeLast()
		h.percolateDown(0)
		return result, true
	}
}

// Peak returns the highest priority item in the heap without
// dequeuing it.
// The second return val, if false, indicates that the heap is empty
// and that a nil value was returned to you as the first return val
func (h *Heap) Peak() (Orderable, bool) {
	if len(h.storage) > 0 {
		return h.storage[0], true
	}
	return nil, false
}

// Size returns the number of items in the Heap
func (h *Heap) Size() int {
	return len(h.storage)
}

func (h *Heap) removeLast() Orderable {
	result := h.storage[len(h.storage)-1]
	h.storage = h.storage[:len(h.storage)-1]
	return result
}

func (h *Heap) percolateUp(i int) {
	parentIndex := h.parentIndex(i)
	for parentIndex >= 0 && parentIndex < i && !h.inOrder(parentIndex, i) {
		h.storage[parentIndex], h.storage[i] = h.storage[i], h.storage[parentIndex]
		i = parentIndex
		parentIndex = h.parentIndex(i)
	}
}

func (h *Heap) percolateDown(i int) {
	childIndex := h.highestPriorityChildIndex(i)
	for childIndex > -1 && !h.inOrder(i, childIndex) {
		h.storage[i], h.storage[childIndex] = h.storage[childIndex], h.storage[i]
		i = childIndex
		childIndex = h.highestPriorityChildIndex(i)
	}
}

// Returns the highest priority child index.
// If there are no children, returns -1
func (h *Heap) highestPriorityChildIndex(parentIndex int) int {
	left, right := h.leftChildIndex(parentIndex), h.rightChildIndex(parentIndex)
	switch {
	case left >= len(h.storage):
		return -1 // no children
	case right >= len(h.storage):
		return left // no right child
	// both children exist
	case h.storage[left].Order() <= h.storage[right].Order():
		return left // left child greater or equal priority
	default:
		return right // right child greater priority
	}
}

func (h *Heap) inOrder(parentIndex, childIndex int) bool {
	return h.storage[parentIndex].Order() < h.storage[childIndex].Order()
}

func (h *Heap) parentIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

func (h *Heap) leftChildIndex(parentIndex int) int {
	return parentIndex*2 + 1
}

func (h *Heap) rightChildIndex(parentIndex int) int {
	return parentIndex*2 + 2
}

func (h *Heap) heapify() {
	if len(h.storage) == 0 {
		return
	}
	parentIndex := (len(h.storage) - 1) / 2 // skip the bottom row
	for parentIndex >= 0 {
		h.percolateDown(parentIndex)
		parentIndex--
	}
}

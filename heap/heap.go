// Package heap contains handy methods for instantiating and using heaps.
// Heaps can be pushed to or popped from. Popping always yields the
// lowest-ordered item in the heap, and pushing adds an item to the heap.
// Unlike in https://pkg.go.dev/container/heap@go1.15.6, if the the heap is at
// capacity, pushing ejects and returns the lowest-ordered item.
// For a general description see https://en.wikipedia.org/wiki/Heap_(data_structure).
// For a detailed explanation see https://bradfieldcs.com/algos/trees/priority-queues-with-binary-heaps/
package heap

import "sort"

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)

// Heap is an instance of a heap structure
type Heap struct {
	storage Interface
	maxSize int
}

// The Interface type describes the requirements for a type using the routines
// in this package
type Interface interface {
	sort.Interface
	Push(x Any)
	Pop() Any
	Peak() Any
}

// Any is just an alias for the empty interface `interface{}`
type Any interface{}

// NewHeap returns a Heap of the specified size. If size <= 0 heap size is
// unbounded.
func NewHeap(data Interface, maxSize int) *Heap {
	if maxSize <= 0 {
		maxSize = maxInt
	}
	return &Heap{storage: data, maxSize: maxSize}
}

// Heapify returns a Heap of the specified size using the given source slice as
// its backing storage, and heap-sorts it in O(n) time. If the given heap is
// larger than the specified size the second return value contains the
// lowest-ordered values in the heap, which have been discarded
func Heapify(source Interface, maxSize int) (*Heap, []Any) {
	if maxSize <= 0 {
		maxSize = maxInt
	}
	result := Heap{storage: source, maxSize: maxSize}
	result.heapify()

	discarded := make([]Any, 0)

	if maxSize > 0 {
		for result.storage.Len() > maxSize {
			excessVal, _ := result.Pop()
			discarded = append(discarded, excessVal)
		}
	}

	return &result, discarded
}

// Push adds an item to the heap in O(log(n)) time. The second return val, if
// true, indicates that the heap is at its maximum capacity the highest
// priority item was popped and returned to you as the first return val
func (h *Heap) Push(val Any) (Any, bool) {
	h.storage.Push(val)
	h.percolateUp(h.storage.Len() - 1)
	if h.storage.Len() > h.maxSize {
		return h.Pop()
	}
	return nil, false
}

// Pop removes the highest priority item from the heap in O(log(n)) time. The
// second return val, if false, indicates that the heap is empty and that a nil
// value was returned to you as the first return val
func (h *Heap) Pop() (Any, bool) {
	switch h.storage.Len() {
	case 0:
		return nil, false
	case 1:
		return h.storage.Pop(), true
	default:
		h.storage.Swap(0, h.storage.Len()-1)
		result := h.storage.Pop()
		h.percolateDown(0)
		return result, true
	}
}

// Capacity returns the maximum size of the heap. O(1).
func (h *Heap) Capacity() int {
	return h.maxSize
}

// Size returns the number of items in the heap using the `Len` method of
// the underlying `Interface.Len()`.
func (h *Heap) Size() int {
	return h.storage.Len()
}

// Peak returns the highest priority item from the heap in O(1) time. without
// removing it. second return val, if false, indicates that the heap is empty
// and that a nil value was returned to you as the first return val
func (h *Heap) Peak() (Any, bool) {
	if h.storage.Len() > 0 {
		return h.storage.Peak(), true
	}
	return nil, false
}

/*
 * Private methods
 */

func (h *Heap) percolateUp(i int) {
	parentIndex := h.parentIndex(i)
	for parentIndex >= 0 && parentIndex < i && !h.storage.Less(parentIndex, i) {
		h.storage.Swap(parentIndex, i)
		i = parentIndex
		parentIndex = h.parentIndex(i)
	}
}

func (h *Heap) percolateDown(i int) {
	childIndex := h.highestPriorityChildIndex(i)
	for childIndex > -1 && !h.storage.Less(i, childIndex) {
		h.storage.Swap(i, childIndex)
		i = childIndex
		childIndex = h.highestPriorityChildIndex(i)
	}
}

// Returns the highest priority child index. If there are no children, returns -1
func (h *Heap) highestPriorityChildIndex(parentIndex int) int {
	left, right := h.leftChildIndex(parentIndex), h.rightChildIndex(parentIndex)
	switch {
	case left >= h.storage.Len():
		return -1 // no children
	case right >= h.storage.Len():
		return left // no right child
	// both children exist
	case h.storage.Less(left, right):
		return left // left child greater or equal priority
	default:
		return right // right child greater priority
	}
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
	if h.storage.Len() == 0 {
		return
	}
	parentIndex := (h.storage.Len() - 1) / 2 // skip the bottom row
	for parentIndex >= 0 {
		h.percolateDown(parentIndex)
		parentIndex--
	}
}

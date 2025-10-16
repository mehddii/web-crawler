package main

type Queue[T comparable] interface {
	Front() (T, bool)
	Back() (T, bool)
	Enqueue(values ...T)
	Dequeue() (T, bool)
	Size() int
}

type ArrayQueue[T comparable] struct {
	array    []T
	len, cap int
}

func NewQueue[T comparable]() *ArrayQueue[T] {
	return &ArrayQueue[T]{
		array: make([]T, 0),
		len:   0,
		cap:   0,
	}
}

func (q *ArrayQueue[T]) Size() int {
	return q.len
}

func (q *ArrayQueue[T]) Front() (T, bool) {
	if q.len == 0 {
		var front T
		return front, false
	}

	return q.array[0], true
}

func (q *ArrayQueue[T]) Back() (T, bool) {
	if q.len == 0 {
		var back T
		return back, false
	}

	ind := q.len - 1
	return q.array[ind], true
}

// Adds one or more elements to the end of the queue.
// If the current capacity isn't enough, it is doubled.
func (q *ArrayQueue[T]) Enqueue(values ...T) {
	if len(values) == 0 {
		return
	}

	if q.cap <= q.len+len(values) {
		newCapacity := (q.len + len(values)) * 2
		extendedArray := make([]T, q.len, newCapacity)
		copy(extendedArray, q.array)
		q.array = extendedArray
	}

	q.array = append(q.array, values...)
	q.len = len(q.array)
	q.cap = cap(q.array)
}

// Removes and returns the front element from the queue.
// Whenever the length falls below 25% of the capacity,
// the queue shrinks to half of its current capacity.
func (q *ArrayQueue[T]) Dequeue() (T, bool) {
	var (
		value T
		ok    bool
	)

	if q.len == 0 {
		return value, ok
	}

	value = q.array[0]
	q.array = q.array[1:]
	ok = true
	q.len = len(q.array)

	if q.len < q.cap/4 {
		shrinkedArray := make([]T, q.len, q.cap/2)
		copy(shrinkedArray, q.array)
		q.array = shrinkedArray
		q.cap = cap(q.array)
	}

	return value, ok
}

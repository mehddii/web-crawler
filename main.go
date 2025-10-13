package main

import (
	"fmt"
	"slices"
	"sync"
)

type Queue[T comparable] interface {
	Front() (T, bool)
	Back() (T, bool)
	Enqueue(value T)
	Dequeue() (T, bool)
	Size() int
}

type Fetcher interface {
	Fetch(uri string) ([]byte, error)
}

type Parser interface {
	Parse(html []byte) (text string, links []string, err error)
}

type Persister interface {
	Save(uri string, text string, metadata string) error
}

type Worker interface {
	Work(frontier Queue[string]) error
}

type Master struct {
	frontier  Queue[string]
	fetcher   Fetcher
	parser    Parser
	persister Persister
	mu        sync.Mutex
	wg        sync.WaitGroup
}

type ArrayQueue[T comparable] struct {
	array []T
}

func NewQueue[T comparable]() *ArrayQueue[T] {
	return &ArrayQueue[T]{
		array: make([]T, 0),
	}
}

func (q *ArrayQueue[T]) Size() int {
	return len(q.array)
}

func (q *ArrayQueue[T]) Front() (T, bool) {
	var front T

	if len(q.array) == 0 {
		return front, false
	}

	return q.array[0], true
}

func (q *ArrayQueue[T]) Back() (T, bool) {
	var back T

	if len(q.array) == 0 {
		return back, false
	}

	ind := len(q.array) - 1
	return q.array[ind], true
}

func (q *ArrayQueue[T]) Enqueue(value T) {
	capacity := cap(q.array)
	length := len(q.array)

	if capacity == length {
		newCapacity := max(capacity*2, 1)
		extendedArray := make([]T, length, newCapacity)
		copy(extendedArray, q.array)
		q.array = extendedArray
	}

	q.array = append(q.array, value)
}

func (q *ArrayQueue[T]) Dequeue() (T, bool) {
	capacity := cap(q.array)
	length := len(q.array)

	var value T
	var ok bool = false

	if length > 0 {
		value = q.array[0]
		q.array = slices.Delete(q.array, 0, 1)
		ok = true
		length = len(q.array)

		if length < capacity/4 {
			shrinkedArray := make([]T, length, capacity/2)
			copy(shrinkedArray, q.array)
			q.array = shrinkedArray
		}

	}

	return value, ok
}

func main() {
	var q Queue[int] = NewQueue[int]()
	for i := range 10 {
		q.Enqueue(i)
		v, ok := q.Back()

		if ok {
			fmt.Printf("Back of the queue %v, its size %v\n", v, q.Size())
		}
	}

	for i := range 10 {
		v, ok := q.Front()
		if ok {
			fmt.Printf("At iteration %v Front is %v\n", i+1, v)
			q.Dequeue()
		}
	}

	fmt.Println("Queue final state: ", q)
}

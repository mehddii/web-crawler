package main

import (
	"fmt"
)

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

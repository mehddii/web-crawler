package main

import (
	"sync"
)

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

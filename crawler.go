package main

import (
	"log"
	stdurl "net/url"
	"strings"
	"sync"
)

// The goal here is to have one data structure that stitch
// together all the components the crawler workers will need
// to perform there task.
// Its also its responsability to handle orchestrating the workers.
type Master struct {
	frontier  Queue[string]
	fetcher   Fetcher
	parser    Parser
	persister Persister
	visited   map[string]struct{}
	ch        chan string
	mu        sync.Mutex
	wg        sync.WaitGroup
}

// Dependencies are injected and not created to change implementation easily in the future
func NewMaster(frontier Queue[string], fetcher Fetcher, parser Parser, persister Persister) *Master {
	return &Master{
		frontier:  frontier,
		fetcher:   fetcher,
		parser:    parser,
		persister: persister,
		visited:   make(map[string]struct{}),
		ch:        make(chan string),
		mu:        sync.Mutex{},
		wg:        sync.WaitGroup{},
	}
}

// The method fires some number of goroutines, and
// dones't consider for how much urls exist in the queue
// in order to have a simple retry mechanism and not to
// worry about when to stop crawling.
func (master *Master) FireWorkers(limit int) {
	defer close(master.ch)

	log.Println("Starting to crawl with", limit, "workers")
	master.wg.Add(limit)
	for i := 1; i <= limit; {
		master.mu.Lock()
		url, ok := master.frontier.Dequeue()
		if !ok {
			master.mu.Unlock()
			continue
		}
		master.mu.Unlock()

		master.mu.Lock()
		_, ok = master.visited[url]
		if ok {
			master.mu.Unlock()
			continue
		}
		master.mu.Unlock()

		go func(master *Master, url string) {
			defer master.wg.Done()
			crawler := Crawler{
				master: master,
			}

			url, err := crawler.Work()

			if err != nil {
				if strings.TrimSpace(url) == "" {
					return
				}

				master.mu.Lock()
				master.frontier.Enqueue(url)
				master.mu.Unlock()
			}
		}(master, url)

		master.ch <- url
		i++
	}

	master.wg.Wait()
}

type Worker interface {
	Work() error
}

type Crawler struct {
	master *Master
}

// For simplicity reasons the work is done sequentially
// for each worker.
func (crawler *Crawler) Work() (string, error) {
	master := crawler.master
	url := <-master.ch

	log.Println("Worker", &crawler, "started processing", url)
	html, err := master.fetcher.Fetch(url)
	if err != nil {
		return "", err
	}
	log.Println(url, "was fetched successfully")

	text, links, err := master.parser.Parse(html)
	if err != nil {
		return "", err
	}
	log.Println(url, "was parsed successfully")
	master.mu.Lock()
	master.frontier.Enqueue(links...)
	master.mu.Unlock()
	log.Println(len(links), "new urls were added to the queue")

	parsedUrl, err := stdurl.Parse(url)
	if err != nil {
		return "", err
	}

	err = master.persister.Save(url, text, parsedUrl.Host)
	if err != nil {
		return "", err
	}

	master.mu.Lock()
	master.visited[url] = struct{}{}
	master.mu.Unlock()

	return url, nil
}

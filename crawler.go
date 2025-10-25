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
func (master *Master) fireWorkers(limit int) {
	log.Println("Starting to crawl with", limit, "workers")

	frontierLength := master.frontier.Size()
	visitedLength := len(master.visited)

	for frontierLength > 0 && visitedLength <= limit {
		log.Println("___________________________________________________")
		log.Println("Frontier Length", frontierLength)
		log.Println("Visited Length", visitedLength)
		log.Println("___________________________________________________")

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


		master.wg.Add(1)
		go func() {
			defer master.wg.Done()
			crawler := Crawler{
				master: master,
			}

			url, err := crawler.Work()
			if err != nil {
				master.mu.Lock()
				_, ok := master.visited[url]
				master.mu.Unlock()

				if strings.TrimSpace(url) == "" || ok {
					return
				}

				master.mu.Lock()
				master.frontier.Enqueue(url)
				master.mu.Unlock()

			}
		}()

		// Profilling
		// f, err := os.Create("crawler.pprof")
		// if err != nil {
		// 	log.Panicln("Error occured while creating the file", err)
		// }
		// defer f.Close()
		// log.Println("Profile file was created successfully")
		// // pprof.WriteHeapProfile(f)
		// pprof.Lookup("goroutine").WriteTo(f, 0)

		frontierLength = master.frontier.Size()
		visitedLength = len(master.visited)
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

	/*
		if !ok {
			return "", errors.New("The queue is Empty! no work to be done.")
		}

		_, ok = master.visited[url]
			master.mu.Unlock()
			if ok {
				log.Println(ok)
				return "", nil
			}
	*/

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
	log.Println("___________________________________________________")
	log.Println(master.visited)
	log.Println(master.frontier)
	log.Println("___________________________________________________")
	master.mu.Unlock()

	return url, nil
}

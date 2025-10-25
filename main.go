package main

import (
	"log"
)

func main() {
	persister, err := NewSqlitePersister("./data/webpages.db", "webpages")
	if err != nil {
		log.Fatalln("Failed to establish a connection to the database", err)
	}

	parser := &HTMLParser{}
	fetcher := &SimpleFetcher{}
	frontier := NewQueue[string]()
	frontier.Enqueue("https://go.dev/doc/")

	crawler := NewMaster(frontier, fetcher, parser, persister)
	crawler.FireWorkers(10)
}

package main

import (
	"context"
	"log"
)

func main() {
	persister, err := Init(context.Background(), "file:./data/webpages.db")
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

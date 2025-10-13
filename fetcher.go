package main

type Fetcher interface {
	Fetch(uri string) ([]byte, error)
}

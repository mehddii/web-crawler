package main

import (
	"io"
	"net/http"
)

type Fetcher interface {
	Fetch(uri string) ([]byte, error)
}

type SimpleFetcher struct{}

func (fetcher *SimpleFetcher) Fetch(uri string) ([]byte, error) {
	resp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

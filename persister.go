package main

type Persister interface {
	Save(uri string, text string, metadata string) error
}

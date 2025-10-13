package main

type Parser interface {
	Parse(html []byte) (text string, links []string, err error)
}

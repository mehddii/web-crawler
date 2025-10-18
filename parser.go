package main

import (
	"bytes"
	"strings"

	stdhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Parser interface {
	Parse(html []byte) (text string, links []string, err error)
}

type HTMLParser struct{}

func (parser *HTMLParser) Parse(html []byte) (text string, links []string, err error) {
	r := bytes.NewReader(html)
	doc, err := stdhtml.Parse(r)
	if err != nil {
		return "", nil, err
	}

	links = make([]string, 0)
	for node := range doc.Descendants() {
		// Extracts http/https links from <a> tags.
		if node.Type == stdhtml.ElementNode && node.DataAtom == atom.A {
			for _, attr := range node.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
					links = append(links, attr.Val)
				}
			}
		}
	}

	return extractText(doc), links, nil
}

// Traverses recursively the html tree and extracts text.
// All the <script>, <style> and <iframe> nodes are skipped
// to avoid noise.
func extractText(node *stdhtml.Node) string {
	var (
		builder strings.Builder
		f       func(node *stdhtml.Node)
	)

	f = func(node *stdhtml.Node) {
		if node.Type == stdhtml.ElementNode &&
			(node.Data == "script" || node.Data == "style" || node.Data == "iframe") {
			return
		}

		if node.Type == stdhtml.TextNode && !strings.HasPrefix(node.Data, "<iframe") {
			data := strings.TrimSpace(node.Data)
			if data != "" {
				_, err := builder.WriteString(data)
				builder.WriteString(" ")
				if err != nil {
					return
				}
			}
		}

		for n := node.FirstChild; n != nil; n = n.NextSibling {
			f(n)
		}
	}
	f(node)

	return strings.TrimSpace(builder.String())
}

package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

// fetchBody wraps a simple GET request returning the body.
func fetchBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

// attr returns the value of a named attribute from an HTML node.
func attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

package utils

import (
	"golang.org/x/net/html"
)

// Attr returns the value of a named attribute from an HTML node.
func Attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

package parser

import (
	"golang.org/x/net/html"
)

func GetNodeById(node *html.Node, id string) *html.Node {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == id {
				return node
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := GetNodeById(child, id); found != nil {
			return found
		}
	}
	return nil
}

func GetNodeByClassName(node *html.Node, className string) *html.Node {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == className {
				return node
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := GetNodeByClassName(child, className); found != nil {
			return found
		}
	}
	return nil
}

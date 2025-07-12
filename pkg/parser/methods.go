package parser

import (
	"golang.org/x/net/html"
)

const (
	SelectorIdPrefix    = '#'
	SelectorClassPrefix = '.'
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

func GetNodeByTagName(node *html.Node, tagName string) *html.Node {
	if node.Type == html.ElementNode && node.Data == tagName {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := GetNodeByTagName(child, tagName); found != nil {
			return found
		}
	}
	return nil
}

func QuerySelector(node *html.Node, selector string) *html.Node {
	if len(selector) == 0 {
		return nil
	}

	if rune(selector[0]) == SelectorIdPrefix {
		id := selector[1:]
		return GetNodeById(node, id)
	}
	if rune(selector[0]) == SelectorClassPrefix {
		className := selector[1:]
		return GetNodeByClassName(node, className)
	}

	tagName := selector
	if tagName == "*" {
		return node
	}

	return GetNodeByTagName(node, tagName)
}

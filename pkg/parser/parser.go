package parser

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseString(htmlString string) (*html.Node, error) {
	if htmlString == "<nil>" {
		return nil, fmt.Errorf("received <nil> string, cannot parse HTML")
	}

	context := &html.Node{
		Type:      html.ElementNode,
		Data:      "div",
		Namespace: "",
		DataAtom:  atom.Div,
	}

	doc, err := html.ParseFragment(strings.NewReader(htmlString), context)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML string: %w", err)
	}

	if len(doc) != 1 {
		return nil, fmt.Errorf("expected a single root node, got %d nodes", len(doc))
	}

	return doc[0], nil
}

func ParseFile(filePath string) (*html.Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	if strings.HasSuffix(filePath, "index.html") {
		index, err := html.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse index HTML file: %w", err)
		}

		return index, nil
	}

	context := &html.Node{
		Type:      html.ElementNode,
		Data:      "div",
		Namespace: "",
		DataAtom:  atom.Div,
	}

	nodes, err := html.ParseFragment(file, context)

	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML fragment: %w", err)
	}

	if len(nodes) != 1 {
		return nil, fmt.Errorf("expected a single root node, got %d nodes", len(nodes))
	}

	return nodes[0], nil
}

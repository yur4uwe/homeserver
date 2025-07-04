package parser

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func MergePage(newPageContent *html.Node) (*html.Node, error) {
	file_name := "../static/index.html"
	index_file, err := ParseFile(file_name)
	if err != nil {
		return nil, fmt.Errorf("failed to parse index file: %w", err)
	}

	return InsertNodeIntoPage(index_file, newPageContent, "main-content")
}

func InsertNodeIntoPage(file_node, node_to_insert *html.Node, selector string) (*html.Node, error) {
	var selectorNode *html.Node
	if selector == "" {
		selectorNode = file_node
	} else {
		selectorNode = GetNodeById(file_node, selector)
		if selectorNode == nil {
			return nil, fmt.Errorf("node with selector '%s' not found", selector)
		}
	}

	selectorNode.AppendChild(node_to_insert)

	return file_node, nil
}

func HTMLToByte(node *html.Node) ([]byte, error) {
	var buf strings.Builder
	if err := html.Render(&buf, node); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return []byte(buf.String()), nil
}

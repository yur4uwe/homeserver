package parser

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func MergePage(newPageContent string) ([]byte, error) {
	return InsertNodeIntoPage("../static/index.html", newPageContent, "main-content")
}

func InsertNodeIntoPage(file_name string, node_string string, selector string) ([]byte, error) {
	file, err := os.Open(file_name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", file_name, err)
	}
	defer file.Close()

	htmlFile, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var selectorNode *html.Node
	if selector == "" {
		selectorNode = htmlFile
	} else {
		selectorNode = getNodeById(htmlFile, selector)
		if selectorNode == nil {
			return nil, fmt.Errorf("node with selector '%s' not found", selector)
		}
	}

	node, err := ParseString(node_string)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node: %w", err)
	}

	selectorNode.AppendChild(node)

	return HTMLToByte(htmlFile)
}

func HTMLToByte(node *html.Node) ([]byte, error) {
	var buf strings.Builder
	if err := html.Render(&buf, node); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return []byte(buf.String()), nil
}

func ParseString(htmlString string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML string: %w", err)
	}
	return doc, nil
}

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

func HTMLToString(node *html.Node) (string, error) {
	data, err := HTMLToByte(node)
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to string: %w", err)
	}

	return string(data), nil
}

func LoadIcon(icon_name string) string {
	icon_path := fmt.Sprintf("../static/assets/%s.svg", icon_name)
	icon_node, err := ParseFile(icon_path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load icon %s: %v\n", icon_name, err))
	}

	icon_data, err := HTMLToString(icon_node)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert icon %s to string: %v\n", icon_name, err))
	}

	return icon_data
}

const (
	AlertSuccess = "success"
	AlertError   = "error"
	AlertWarning = "warning"
	AlertInfo    = "info"
)

func NewAlert(message string, status string) *html.Node {
	alertNode := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "class", Val: fmt.Sprintf("alert alert-%s", status)},
			{Key: "hx-get", Val: "/nothing"},
			{Key: "hx-trigger", Val: "load delay:3s"},
			{Key: "hx-swap", Val: "delete"},
		},
	}

	alertText := &html.Node{
		Type: html.TextNode,
		Data: message,
	}

	alertNode.AppendChild(alertText)
	return alertNode
}

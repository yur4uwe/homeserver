package parser

import (
	"fmt"
	"homeserver/pkg/templates"
	"html/template"
	"os"
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

	icon_data, err := os.ReadFile(icon_path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load icon %s: %v\n", icon_name, err))
	}

	return string(icon_data)
}

func NewAlert(title, iconName, message, status string) string {

	alert := templates.ModalContent{
		Title:   title,
		Icon:    LoadIcon(iconName),
		Type:    status,
		Message: message,
	}

	tmpl := template.Must(template.ParseFiles("../static/templates/modal.html"))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, alert); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute template: %v\n", err)
		return ""
	}

	return buf.String()
}

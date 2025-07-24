package parser

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Must(node *html.Node, err error) *html.Node {
	if err != nil {
		panic(fmt.Sprintf("Failed to parse HTML: %v\n", err))
	}
	return node
}

func MustAsString(node *html.Node, err error) string {
	if err != nil {
		panic(fmt.Sprintf("Failed to convert HTML to string: %v\n", err))
	}

	data, err := HTMLToString(node)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert HTML to string: %v\n", err))
	}

	return data
}

func ParseString(htmlString string) (*html.Node, error) {
	if htmlString == "<nil>" {
		return nil, fmt.Errorf("received <nil> string, cannot parse HTML")
	}

	htmlString = strings.TrimSpace(htmlString)
	if htmlString == "" {
		return nil, fmt.Errorf("received empty string, cannot parse HTML")
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

func ParseTemplate(template_filepath string, template_data interface{}) (*html.Node, error) {
	tmpl, err := ParseTemplateToString(template_filepath, template_data)
	if err != nil {
		return nil, err
	}
	return ParseString(tmpl)
}

func ParseTemplateToString(template_filepath string, template_data interface{}) (string, error) {
	tmpl := template.Must(template.ParseFiles(template_filepath))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, template_data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

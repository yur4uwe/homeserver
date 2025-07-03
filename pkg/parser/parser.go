package parser

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func MergePage(newPageContent string) ([]byte, error) {
	file, err := os.Open("../static/index.html")
	if err != nil {
		return nil, err
	}

	htmlFile, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	content_node := getNodeById(htmlFile, "main-content")
	if content_node == nil {
		return nil, fmt.Errorf("content node with id 'main-content' not found")
	}

	new_content, err := html.Parse(strings.NewReader(newPageContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse new page content: %w", err)
	}

	content_node.AppendChild(new_content)

	var buf strings.Builder
	if err := html.Render(&buf, htmlFile); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return []byte(buf.String()), nil
}

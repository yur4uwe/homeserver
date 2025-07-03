package parser

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func ParseString(htmlString string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML string: %w", err)
	}
	return doc, nil
}

func ParseFile(filePath string) (*html.Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	return html.Parse(file)
}

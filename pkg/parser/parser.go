package parser

import (
	"fmt"
	"html/template"
	"strings"
)

func ParseTemplateToString(template_filepath string, template_data interface{}) (string, error) {
	tmpl := template.Must(template.ParseFiles(template_filepath))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, template_data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

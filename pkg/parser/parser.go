package parser

import (
	"fmt"
	"strings"
)

func ParseTemplateToString(template_filepath string, template_data interface{}) (string, error) {
	tmpl := templateMust(templateParseFiles(template_filepath))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, template_data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

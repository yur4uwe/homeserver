package parser

import (
	"fmt"
	"homeserver/pkg/templates"
	"html/template"
	"os"
	"strings"
)

func LoadIcon(icon_name string) string {
	icon_path := fmt.Sprintf("../static/assets/%s.svg", icon_name)

	icon_data, err := os.ReadFile(icon_path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load icon %s: %v\n", icon_name, err))
	}

	return string(icon_data)
}

func NewAlert(title, iconName, message, status string) string {

	alert := &templates.ModalContent{
		Title:   title,
		Icon:    template.HTML("<icon name=\"" + iconName + "\"></icon>"),
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

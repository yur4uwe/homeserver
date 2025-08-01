package parser

import (
	"fmt"
	"homeserver/pkg/templates"
	"html/template"
	"os"
)

var (
	osReadFile = os.ReadFile         // For testing purposes
	Must       = template.Must       // For testing purposes
	ParseFiles = template.ParseFiles // For testing purposes
)

func LoadIcon(icon_name string) string {
	icon_path := fmt.Sprintf("../static/assets/%s.svg", icon_name)

	icon_data, err := osReadFile(icon_path)
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

	str, err := ParseTemplateToString("../static/templates/modal.html", alert)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse alert template: %v\n", err))
	}

	return str
}

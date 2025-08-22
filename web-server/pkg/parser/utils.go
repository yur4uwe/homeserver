package parser

import (
	"fmt"
	"homeserver/constants"
	"homeserver/web-server/pkg/templates"
	"html/template"
	"os"
)

var (
	OsReadFile = os.ReadFile         // For testing purposes
	Must       = template.Must       // For testing purposes
	ParseFiles = template.ParseFiles // For testing purposes
)

func LoadIcon(icon_name string) string {
	icon_path := fmt.Sprintf(constants.Rootd+"static/assets/%s.svg", icon_name)

	icon_data, err := OsReadFile(icon_path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load icon %s: %v\n", icon_name, err))
	}

	return string(icon_data)
}

func NewAlert(title, iconName, message, status string) string {

	alert := &templates.ModalContent{
		Title:   title,
		Icon:    template.HTML(LoadIcon(iconName)),
		Type:    status,
		Message: message,
	}

	str, err := ParseTemplateToString(constants.Rootd+"static/templates/modal.html", alert)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse alert template: %v\n", err))
	}

	return str
}

package templates

import "html/template"

type ErrorData struct {
	Title      string
	Message    string
	Additional template.HTML
}

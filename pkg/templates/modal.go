package templates

import "html/template"

type ModalContent struct {
	Title   string
	Icon    template.HTML
	Type    string // e.g., "error", "info", "warning"
	Message string
}

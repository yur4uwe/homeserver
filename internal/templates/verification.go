package templates

import "html/template"

type VerificationPageData struct {
	Success bool
	Error   string
	Causes  []template.HTML
	Icon    template.HTML
}

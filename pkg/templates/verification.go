package templates

import "html/template"

type VerificationPageData struct {
	Success bool
	Error   string
	Icon    template.HTML
}

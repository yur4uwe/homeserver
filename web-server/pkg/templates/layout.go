package templates

import "html/template"

type LayoutPageData struct {
	Header       template.HTML
	Content      template.HTML
	ModalContent template.HTML
	Footer       template.HTML
}

package templates

type ModalContent struct {
	Title    string
	IconName string
	Type     string // e.g., "error", "info", "warning"
	Message  string
}

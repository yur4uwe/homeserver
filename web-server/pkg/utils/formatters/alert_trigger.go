package formatters

import (
	"fmt"
	"homeserver/constants"
)

func AlertEventBuilder(title, message, status string) string {
	icon := "check"
	switch status {
	case constants.Success:
		icon = "check"
	case constants.Error:
		icon = "x"
	default:
		panic(fmt.Sprintf("Unknown status: %s", status))
	}

	return fmt.Sprintf(`{"raise-alert": {"title": "%s", "value": "%s", "status": "%s", "icon": "%s"}}`, title, message, status, icon)
}

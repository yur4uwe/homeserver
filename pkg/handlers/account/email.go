package account

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas"
	"net/http"
)

func HandleEmailChange(r *http.Request, connection interfaces.DatabaseInterface) error {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("missing request context")
	}

	reqCont.LogHistory.Log("[Handler] Email change requested")

	new_email := r.FormValue("email")
	if new_email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	err := connection.ModifyUserField(r.Context(), reqCont.Requester, schemas.Email, new_email)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

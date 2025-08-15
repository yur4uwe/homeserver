package account

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas"
	"net/http"
)

func HandlePasswordChange(r *http.Request, connection interfaces.DatabaseInterface) error {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("missing request context")
	}

	reqCont.LogHistory.Log("[Handler] Password change requested")

	new_password := r.FormValue("password")
	if new_password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err := connection.ModifyUserField(r.Context(), reqCont.Requester, schemas.Password, new_password)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

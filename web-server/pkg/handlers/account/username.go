package account

import (
	"fmt"
	"homeserver/db/repositories"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas"
	"homeserver/pkg/utils"
	"net/http"
)

var (
	ErrorResponse = utils.ErrorResponse
)

func HandleUsernameChange(r *http.Request, repos *repositories.Repositories) error {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("missing request context")
	}

	reqCont.LogHistory.Log("[Handler] Username change requested")

	new_username := r.FormValue("username")
	if new_username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	err := connection.ModifyUserField(r.Context(), reqCont.Requester, schemas.Username, new_username)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

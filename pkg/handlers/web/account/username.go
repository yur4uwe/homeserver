package account

import (
	"fmt"
	"homeserver/constants/uconsts"
	"homeserver/db/repositories"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
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

	fvp := map[string]any{
		uconsts.Username: new_username,
	}

	err := repos.UserRepo.Update(r.Context(), reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

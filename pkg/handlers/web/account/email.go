package account

import (
	"fmt"
	"homeserver/constants/uconsts"
	"homeserver/db/repositories"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
)

func HandleEmailChange(r *http.Request, repos *repositories.Repositories) error {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("missing request context")
	}

	reqCont.LogHistory.Log("[Handler] Email change requested")

	new_email := r.FormValue("email")
	if new_email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	fvp := map[string]any{
		uconsts.Email: new_email,
	}

	err := repos.UserRepo.Update(r.Context(), reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

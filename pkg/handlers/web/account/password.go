package account

import (
	"fmt"
	"homeserver/constants/uconsts"
	"homeserver/db/repositories"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
)

func HandlePasswordChange(r *http.Request, repos *repositories.Repositories) error {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("missing request context")
	}

	reqCont.LogHistory.Log("[Handler] Password change requested")

	new_password := r.FormValue("password")
	if new_password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	fvp := map[string]any{
		uconsts.Password: new_password,
	}

	err := repos.UserRepo.Update(r.Context(), reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

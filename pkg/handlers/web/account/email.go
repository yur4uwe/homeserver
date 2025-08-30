package account

import (
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
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

	fvp := uschema.New().WithEmail(new_email)

	_, err := repos.UserRepo.Update(r.Context(), reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

package account

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	middleware "homeserver/pkg/middleware/web"
)

func HandlePasswordChange(
	ctx context.Context,
	repos *repositories.Repositories,
	reqCont *middleware.AdditionalRequestContext,
	password string,
	confirm_password string,
) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if password != confirm_password {
		return fmt.Errorf("passwords do not match")
	}

	fvp := uschema.New().WithPassword(password)

	_, err := repos.UserRepo.Update(ctx, reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

package account

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	middleware "homeserver/pkg/middleware/web"
)

func HandleUsernameChange(
	ctx context.Context,
	repos *repositories.Repositories,
	reqCont *middleware.AdditionalRequestContext,
	uname string,
) error {
	if uname == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if uname == reqCont.Requester.Username {
		return fmt.Errorf("new username is the same as the current username")
	}

	fvp := uschema.New().WithUsername(uname)

	_, err := repos.UserRepo.Update(ctx, reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update username")
	}

	return nil
}

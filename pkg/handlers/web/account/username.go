package account

import (
	"context"
	"fmt"
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
)

func HandleUsernameChange(
	ctx context.Context,
	repos *db.Database,
	reqCont *middleware.AdditionalRequestContext,
	uname string,
) error {
	if uname == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if uname == reqCont.Requester.Username {
		return fmt.Errorf("new username is the same as the current username")
	}

	updates := &db.User{Username: uname}

	_, err := repos.Users.Update(ctx, reqCont.Requester, updates)
	if err != nil {
		return fmt.Errorf("failed to update username")
	}

	return nil
}

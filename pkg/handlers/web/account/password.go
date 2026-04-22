package account

import (
	"context"
	"fmt"
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
)

func HandlePasswordChange(
	ctx context.Context,
	repos *db.Database,
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

	updates := &db.User{Password: password}

	_, err := repos.Users.Update(ctx, reqCont.Requester, updates)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

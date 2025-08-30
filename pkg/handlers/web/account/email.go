package account

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/pkg/crypto"
	"homeserver/pkg/handlers/web/auth"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils/formatters"
	"strings"
	"time"
)

func HandleEmailChange(
	ctx context.Context,
	repos *repositories.Repositories,
	reqCont *middleware.AdditionalRequestContext,
	email string,
) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if email == reqCont.Requester.Email {
		return fmt.Errorf("new email is the same as the current email")
	}

	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}

	fvp := uschema.New(true).WithEmail(email).WithEmailVerified(false)

	_, err := repos.UserRepo.Update(ctx, reqCont.Requester, fvp)
	if err != nil {
		return fmt.Errorf("failed to update email")
	}

	tkn, err := crypto.NewToken(email, time.Now().Add(24*time.Hour))
	if err != nil {
		panic(err)
	}

	if err := repos.TokenRepo.CreateToken(ctx, reqCont.Requester.ID, tkn); err != nil {
		panic(err)
	}

	go auth.SendVerificationEmail(formatters.RegisterForm{
		Username: reqCont.Requester.Username,
		Email:    email,
	}, tkn)

	return nil
}

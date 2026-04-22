package account

import (
	"context"
	"fmt"
	"homeserver/db"
	"homeserver/pkg/crypto"
	"homeserver/pkg/handlers/web/auth"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils"
	"strings"
	"time"
)

func HandleEmailChange(
	ctx context.Context,
	repos *db.Database,
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

	updates := &db.User{Email: email, EmailVerified: false}

	_, err := repos.Users.Update(ctx, reqCont.Requester, updates)
	if err != nil {
		return fmt.Errorf("failed to update email")
	}

	tkn, err := crypto.NewToken(email, time.Now().Add(24*time.Hour))
	if err != nil {
		panic(err)
	}

	if err := repos.Tokens.CreateToken(ctx, reqCont.Requester.ID, tkn); err != nil {
		panic(err)
	}

	go auth.SendVerificationEmail(utils.RegisterForm{
		Username: reqCont.Requester.Username,
		Email:    email,
	}, tkn)

	return nil
}

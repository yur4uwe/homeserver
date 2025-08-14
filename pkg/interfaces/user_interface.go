package interfaces

import (
	"context"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/user_schema"
)

type UserDBInterface interface {
	CreateUserEntry(ctx context.Context, user *user_schema.User) (*user_schema.User, error)
	ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error)
	CheckUserFolder(ctx context.Context, user *user_schema.User) error
	QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (*user_schema.User, bool, error)
	QueryUsername(ctx context.Context, username string) (*user_schema.User, error)
	QueryEmail(ctx context.Context, email string) (*user_schema.User, error)
	QueryPassword(ctx context.Context, password string) (*user_schema.User, error)
	ModifyUserField(ctx context.Context, user *user_schema.User, field string, value interface{}) error
}

package interfaces

import (
	"context"
	"homeserver/pkg/schemas/user_schema"
)

type TokenInterface interface {
	CheckToken(ctx context.Context, token_str string) (int, error)
	CreateToken(ctx context.Context, userID int, token_str string) error
	QueryToken(ctx context.Context, token_str string) (*user_schema.User, error)
	VerifyToken(ctx context.Context, token_str string) error
	DeleteToken(ctx context.Context, token_str string) error
}

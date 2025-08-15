package db

import (
	"context"
	"homeserver/pkg/db/token"
	"homeserver/pkg/schemas/user_schema"
)

func (db *Database) CheckToken(ctx context.Context, token_str string) (int, error) {
	return token.CheckToken(ctx, db, token_str)
}

func (db *Database) CreateToken(ctx context.Context, userID int, token_str string) error {
	return token.CreateToken(ctx, db, userID, token_str)
}

func (db *Database) QueryToken(ctx context.Context, token_str string) (*user_schema.User, error) {
	return token.QueryToken(ctx, db, token_str)
}

func (db *Database) VerifyToken(ctx context.Context, token_str string) error {
	return token.VerifyToken(ctx, db, token_str)
}

func (db *Database) DeleteToken(ctx context.Context, token_str string) error {
	return token.DeleteToken(ctx, db, token_str)
}

package db

import (
	"context"
	"homeserver/pkg/db/user"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/user_schema"
)

func (db *Database) CreateUserEntry(ctx context.Context, user_data *user_schema.User) (*user_schema.User, error) {
	return user.CreateUserEntry(ctx, db, user_data)
}

func (db *Database) ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error) {
	return user.ExtractUserEntries(ctx, db)
}

func (db *Database) CheckUserFolder(ctx context.Context, user_data *user_schema.User) error {
	return user.CheckUserFolder(ctx, db, user_data)
}

func (db *Database) QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (int, bool, error) {
	return user.QueryUserOnLogin(ctx, db, form)
}

func (db *Database) QueryUsername(ctx context.Context, username string) (int, error) {
	return user.QueryUsername(ctx, db, username)
}

func (db *Database) QueryEmail(ctx context.Context, email string) (int, error) {
	return user.QueryEmail(ctx, db, email)
}

func (db *Database) QueryPassword(ctx context.Context, password string) (int, error) {
	return user.QueryPassword(ctx, db, password)
}

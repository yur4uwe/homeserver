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

func (db *Database) DeleteUser(ctx context.Context, user_data *user_schema.User) error {
	return user.DeleteUser(ctx, db, user_data)
}

func (db *Database) ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error) {
	return user.ExtractUserEntries(ctx, db)
}

func (db *Database) CheckUserFolder(ctx context.Context, user_data *user_schema.User) error {
	return user.CheckUserFolder(ctx, db, user_data)
}

func (db *Database) QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (*user_schema.User, bool, error) {
	return user.QueryUserOnLogin(ctx, db, form)
}

func (db *Database) QueryUsername(ctx context.Context, username string) (*user_schema.User, error) {
	return user.QueryUsername(ctx, db, username)
}

func (db *Database) QueryEmail(ctx context.Context, email string) (*user_schema.User, error) {
	return user.QueryEmail(ctx, db, email)
}

func (db *Database) QueryPassword(ctx context.Context, password string) (*user_schema.User, error) {
	return user.QueryPassword(ctx, db, password)
}

func (db *Database) ModifyUserField(ctx context.Context, user_data *user_schema.User, field string, value interface{}) error {
	return user.ModifyUserField(ctx, db, user_data, field, value)
}

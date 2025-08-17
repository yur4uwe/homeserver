package database

import (
	"context"
	"fmt"
	"homeserver/db"
)

func ResetDatabase(connection *db.Database) error {
	ctx := context.Background()

	_, err := connection.Pool.Exec(ctx, `TRUNCATE TABLE users RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating users table: %v", err)
	}

	_, err = connection.Pool.Exec(ctx, `TRUNCATE TABLE user_folders RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating user_folders table: %v", err)
	}

	_, err = connection.Pool.Exec(ctx, `TRUNCATE TABLE user_files RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating user_files table: %v", err)
	}

	return nil
}

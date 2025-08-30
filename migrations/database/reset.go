package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ResetDatabase(connection *pgxpool.Pool) error {
	ctx := context.Background()

	_, err := connection.Exec(ctx, `TRUNCATE TABLE users RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating users table: %v", err)
	}

	_, err = connection.Exec(ctx, `TRUNCATE TABLE user_folders RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating user_folders table: %v", err)
	}

	_, err = connection.Exec(ctx, `TRUNCATE TABLE user_files RESTART IDENTITY CASCADE;`)
	if err != nil {
		return fmt.Errorf("error truncating user_files table: %v", err)
	}

	return nil
}

package database

import (
	"context"
	"fmt"
	"homeserver/db"
)

func SeparateExtensionsFromNames(connection *db.Database) error {
	ctx := context.Background()

	tag, err := connection.Pool.Exec(ctx, `SELECT
  		filename,
  		reverse(split_part(reverse(filename), '.', 1)) AS extension,
  		left(filename, length(filename) - length(reverse(split_part(reverse(filename), '.', 1))) - 1) AS base_filename
		FROM files;`)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

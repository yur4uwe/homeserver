package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateFolderStructure(connection *pgxpool.Pool) error {
	ctx := context.Background()

	rows, err := connection.Query(ctx, `
		SELECT
    		u.username,
			u.id AS owner_id,
    		f.name AS root_folder_name
		FROM
    		users u
		LEFT JOIN
    		user_folders f
    	ON f.owner_id = u.id AND f.parent IS NULL
	`)

	if err != nil {
		return fmt.Errorf("error querying folder structure: %v", err)
	}

	for rows.Next() {
		var username, rootFolderName string
		var ownerID int
		if err := rows.Scan(&username, &ownerID, &rootFolderName); err != nil {
			return fmt.Errorf("error scanning folder row: %v", err)
		}
		if rootFolderName != "" {
			tag, err := connection.Exec(ctx, `
				INSERT INTO user_folders (name, path, parent, owner_id, created_at)
				VALUES ($1, $2, $3, $4, now())
			`, rootFolderName, "", nil, ownerID)
			if err != nil {
				return fmt.Errorf("error creating root folder: %v", err)
			}

			if tag.RowsAffected() == 0 {
				return fmt.Errorf("no rows affected when creating root folder for user %s", username)
			}
		}
	}

	return nil
}

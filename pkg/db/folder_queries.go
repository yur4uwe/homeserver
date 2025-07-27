package db

import (
	"context"
	"fmt"
)

func (db *Database) CreateFolderEntry(ctx context.Context, folder *UserFolders) error {
	_, err := db.pool.Exec(
		ctx,
		`INSERT INTO user_folders (name, parent, owner_id, created_at) VALUES ($1, $2, $3, NOW())`,
		folder.Name,
		folder.Parent,
		folder.OwnerID,
	)
	if err != nil {
		return fmt.Errorf("error creating folder entry: %v", err)
	}
	return nil
}

func (db *Database) ExtractUserFolders(ctx context.Context, userID int) ([]*UserFolders, error) {
	rows, err := db.pool.Query(
		ctx,
		`SELECT id, name, parent, owner_id, created_at FROM user_folders WHERE owner_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user folders: %v", err)
	}
	defer rows.Close()

	var entries []*UserFolders
	for rows.Next() {
		entry := &UserFolders{}
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Parent, &entry.OwnerID, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning user folder entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

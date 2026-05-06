package db

import (
	"context"
	"fmt"
	"homeserver/constants"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UserFolder entity representing the user_folders table
type UserFolder struct {
	ID            int       `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Directory     string    `db:"path" json:"path"` // Path starts from user root folder, e.g. "username/documents"
	Parent        *int      `db:"parent" json:"parent,omitempty"`
	OwnerID       int       `db:"owner_id" json:"owner_id"`
	TrustedSource []byte    `db:"trusted_source" json:"trusted_source"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// FolderRepo defines the data access layer for Folders
type FolderRepo interface {
	Create(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error)
	Get(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error)
	GetByOwner(ctx context.Context, ownerID int) ([]UserFolder, error)
	GetChildren(ctx context.Context, path, name string, owner *User) ([]UserFolder, error)
	Update(ctx context.Context, folder *UserFolder, updates *UserFolder, owner *User) (*UserFolder, error)
	Delete(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error)
}

type pgFolderRepo struct {
	db QueryExecutor
}

func NewFolderRepo(db QueryExecutor) FolderRepo {
	return &pgFolderRepo{db: db}
}

// Logic & Helpers

func (f *UserFolder) Validate(ownerName string) error {
	if f.Name == "" {
		return fmt.Errorf("folder name cannot be empty")
	}
	if f.OwnerID <= 0 {
		return fmt.Errorf("invalid owner ID")
	}
	if f.Parent == nil && f.Name != ownerName {
		return fmt.Errorf("root folder must have the same name as the owner")
	}
	if (f.Directory == "" && f.Parent != nil) || (f.Directory != "" && f.Parent == nil) {
		// Exception for root folder where Path is empty and Parent is nil
		if !(f.Directory == "" && f.Parent == nil && f.Name == ownerName) {
			return fmt.Errorf("folder path and parent folder must be set together for non-root folders")
		}
	}
	return nil
}

func (f *UserFolder) RelativePath() string {
	if f.Directory != "" {
		return fmt.Sprintf("%s/%s", f.Directory, f.Name)
	}
	return f.Name
}

func (f *UserFolder) AbsolutePath() string {
	storagePath := os.Getenv(constants.StoragePath)
	return filepath.Join(storagePath, f.RelativePath())
}

func (f *UserFolder) Identifier(ownerName string, startParam int) (string, []any) {
	if f.ID != 0 {
		return fmt.Sprintf("id = $%d", startParam), []any{f.ID}
	}
	if f.Parent == nil && f.Name == ownerName && f.OwnerID != 0 {
		return fmt.Sprintf("name = $%d AND owner_id = $%d", startParam, startParam+1), []any{f.Name, f.OwnerID}
	}
	if f.Directory != "" && f.Name != "" && f.OwnerID != 0 {
		return fmt.Sprintf("path = $%d AND name = $%d AND owner_id = $%d",
				startParam, startParam+1, startParam+2),
			[]any{f.Directory, f.Name, f.OwnerID}
	}
	return "", nil
}

// Implementation

func (r *pgFolderRepo) Create(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error) {
	if err := folder.Validate(owner.Username); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_folders (name, path, parent, owner_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name, path, owner_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, path, parent, owner_id, trusted_source, created_at`

	var result UserFolder
	err := r.db.QueryRow(ctx, query, folder.Name, folder.Directory, folder.Parent, folder.OwnerID).
		Scan(&result.ID, &result.Name, &result.Directory, &result.Parent, &result.OwnerID, &result.TrustedSource, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return &result, nil
}

func (r *pgFolderRepo) Get(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error) {
	uname := ""
	if owner != nil {
		uname = owner.Username
	}

	where, args := folder.Identifier(uname, 1)
	if where == "" {
		return nil, fmt.Errorf("folder is not identifiable")
	}

	query := fmt.Sprintf(`
		SELECT id, name, path, parent, owner_id, trusted_source, created_at 
		FROM user_folders WHERE %s`, where)

	var result UserFolder
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Name, &result.Directory, &result.Parent, &result.OwnerID, &result.TrustedSource, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return &result, nil
}

func (r *pgFolderRepo) GetByOwner(ctx context.Context, ownerID int) ([]UserFolder, error) {
	query := `SELECT id, name, path, parent, owner_id, trusted_source, created_at FROM user_folders WHERE owner_id = $1`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []UserFolder
	for rows.Next() {
		var f UserFolder
		if err := rows.Scan(&f.ID, &f.Name, &f.Directory, &f.Parent, &f.OwnerID, &f.TrustedSource, &f.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *pgFolderRepo) GetChildren(ctx context.Context, path, name string, owner *User) ([]UserFolder, error) {
	query := `
		SELECT id, name, path, parent, owner_id, trusted_source, created_at 
		FROM user_folders 
		WHERE owner_id = $1 AND parent = (
			SELECT id FROM user_folders WHERE path = $2 AND name = $3 AND owner_id = $1
		)`

	rows, err := r.db.Query(ctx, query, owner.ID, path, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []UserFolder
	for rows.Next() {
		var f UserFolder
		if err := rows.Scan(&f.ID, &f.Name, &f.Directory, &f.Parent, &f.OwnerID, &f.TrustedSource, &f.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *pgFolderRepo) Update(ctx context.Context, folder *UserFolder, updates *UserFolder, owner *User) (*UserFolder, error) {
	where, whereArgs := folder.Identifier(owner.Username, 1)
	if where == "" {
		return nil, fmt.Errorf("invalid folder data")
	}

	var setParts []string
	var args []any
	argID := 1

	if updates.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argID))
		args = append(args, updates.Name)
		argID++
	}
	if updates.Directory != "" {
		setParts = append(setParts, fmt.Sprintf("path = $%d", argID))
		args = append(args, updates.Directory)
		argID++
	}
	if updates.Parent != nil {
		setParts = append(setParts, fmt.Sprintf("parent = $%d", argID))
		args = append(args, updates.Parent)
		argID++
	}
	if len(updates.TrustedSource) > 0 {
		setParts = append(setParts, fmt.Sprintf("trusted_source = $%d", argID))
		args = append(args, updates.TrustedSource)
		argID++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Re-index WHERE
	where, _ = folder.Identifier(owner.Username, argID)
	args = append(args, whereArgs...)

	query := fmt.Sprintf(`
		UPDATE user_folders SET %s WHERE %s 
		RETURNING id, name, path, parent, owner_id, trusted_source, created_at`,
		strings.Join(setParts, ", "), where)

	var result UserFolder
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Name, &result.Directory, &result.Parent, &result.OwnerID, &result.TrustedSource, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return &result, nil
}

func (r *pgFolderRepo) Delete(ctx context.Context, folder *UserFolder, owner *User) (*UserFolder, error) {
	where, args := folder.Identifier(owner.Username, 1)
	if where == "" {
		return nil, fmt.Errorf("invalid folder data")
	}

	query := fmt.Sprintf(`
		DELETE FROM user_folders WHERE %s 
		RETURNING id, name, path, parent, owner_id, trusted_source, created_at`, where)

	var result UserFolder
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Name, &result.Directory, &result.Parent, &result.OwnerID, &result.TrustedSource, &result.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to delete folder: %w", err)
	}

	return &result, nil
}

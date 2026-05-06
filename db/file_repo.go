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

// UserFile entity representing the user_files table
type UserFile struct {
	ID            int       `db:"id" json:"id"`
	Filename      string    `db:"filename" json:"filename"`
	Extension     string    `db:"extension" json:"extension"`
	Path          string    `db:"path" json:"path"` // Path starts from user root folder, e.g. "username/documents"
	Directory     int       `db:"directory_id" json:"directory_id"`
	OwnerID       int       `db:"owner_id" json:"owner_id"`
	Size          int64     `db:"size" json:"size"`
	ContentType   string    `db:"content_type" json:"content_type"`
	TrustedSource []byte    `db:"trusted_source" json:"trusted_source"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	Metadata      *string   `db:"metadata" json:"metadata"`
}

// FileRepo defines the data access layer for Files
type FileRepo interface {
	Create(ctx context.Context, file *UserFile, owner *User) (*UserFile, error)
	Get(ctx context.Context, file *UserFile, owner *User) (*UserFile, error)
	GetByUser(ctx context.Context, owner *User) ([]UserFile, error)
	GetByDirectory(ctx context.Context, directoryPath string, owner *User) ([]UserFile, error)
	FuzzySearch(ctx context.Context, search_query string, owner *User) ([]UserFile, error)
	Update(ctx context.Context, file *UserFile, updates *UserFile, owner *User) (*UserFile, error)
	Delete(ctx context.Context, file *UserFile, owner *User) (*UserFile, error)
}

type pgFileRepo struct {
	db QueryExecutor
}

func NewFileRepo(db QueryExecutor) FileRepo {
	return &pgFileRepo{db: db}
}

// Logic & Helpers

func (f *UserFile) Validate(ownerName string) error {
	if f.Filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if f.Extension == "" || strings.HasPrefix(f.Extension, ".") {
		return fmt.Errorf("invalid extension")
	}
	if f.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if strings.HasPrefix(f.Path, "/") || strings.HasSuffix(f.Path, "/") {
		return fmt.Errorf("path cannot start or end with a slash")
	}
	if ownerName != "" && !strings.HasPrefix(f.Path, ownerName) {
		return fmt.Errorf("path must start with owner name")
	}
	if f.OwnerID <= 0 {
		return fmt.Errorf("invalid owner ID")
	}
	return nil
}

func (f *UserFile) RelativePath() string {
	ext := f.Extension
	if !strings.HasPrefix(ext, ".") && ext != "" {
		ext = "." + ext
	}
	if f.Path != "" {
		return fmt.Sprintf("%s/%s%s", f.Path, f.Filename, ext)
	}
	return fmt.Sprintf("%s%s", f.Filename, ext)
}

func (f *UserFile) AbsolutePath() string {
	storagePath := os.Getenv(constants.StoragePath)
	return filepath.Join(storagePath, f.RelativePath())
}

func (f *UserFile) Identifier(startParam int) (string, []any) {
	if f.ID != 0 {
		return fmt.Sprintf("id = $%d", startParam), []any{f.ID}
	}
	if f.Path != "" && f.Filename != "" && f.Extension != "" && f.OwnerID != 0 {
		return fmt.Sprintf("path = $%d AND filename = $%d AND extension = $%d AND owner_id = $%d",
				startParam, startParam+1, startParam+2, startParam+3),
			[]any{f.Path, f.Filename, f.Extension, f.OwnerID}
	}
	return "", nil
}

func PreventDuplicateFileName(file *UserFile) (*UserFile, error) {
	newFile := *file
	counter := 1
	for {
		path := newFile.AbsolutePath()
		if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
			return &newFile, nil
		}
		newFile.Filename = fmt.Sprintf("%s_%d", file.Filename, counter)
		counter++
	}
}

// Implementation

func (r *pgFileRepo) Create(ctx context.Context, file *UserFile, owner *User) (*UserFile, error) {
	if err := file.Validate(owner.Username); err != nil {
		return nil, err
	}

	file, err := PreventDuplicateFileName(file)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_files (filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(ctx, query,
		file.Filename, file.Extension, file.Path, file.Directory, file.OwnerID,
		file.Size, file.ContentType, file.Metadata, file.TrustedSource,
	).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create file entry: %w", err)
	}

	return file, nil
}

func (r *pgFileRepo) Get(ctx context.Context, file *UserFile, owner *User) (*UserFile, error) {
	where, args := file.Identifier(1)
	if where == "" {
		// Try minimal check if owner is provided
		if owner != nil && file.Filename != "" && file.Path != "" {
			file.OwnerID = owner.ID
			where, args = file.Identifier(1)
		}
	}

	if where == "" {
		return nil, fmt.Errorf("missing file identifiers")
	}

	query := fmt.Sprintf(`
		SELECT id, filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source, created_at, updated_at
		FROM user_files WHERE %s`, where)

	var result UserFile
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&result.ID, &result.Filename, &result.Extension, &result.Path, &result.Directory,
		&result.OwnerID, &result.Size, &result.ContentType, &result.Metadata, &result.TrustedSource,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &result, nil
}

func (r *pgFileRepo) GetByUser(ctx context.Context, owner *User) ([]UserFile, error) {
	if owner.ID == 0 {
		return nil, fmt.Errorf("invalid owner ID")
	}

	query := `
		SELECT id, filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source, created_at, updated_at
		FROM user_files WHERE owner_id = $1`

	rows, err := r.db.Query(ctx, query, owner.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []UserFile
	for rows.Next() {
		var f UserFile
		err := rows.Scan(
			&f.ID, &f.Filename, &f.Extension, &f.Path, &f.Directory,
			&f.OwnerID, &f.Size, &f.ContentType, &f.Metadata, &f.TrustedSource,
			&f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *pgFileRepo) GetByDirectory(ctx context.Context, directoryPath string, owner *User) ([]UserFile, error) {
	query := `
		SELECT id, filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source, created_at, updated_at
		FROM user_files WHERE path = $1 AND owner_id = $2`

	rows, err := r.db.Query(ctx, query, directoryPath, owner.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []UserFile
	for rows.Next() {
		var f UserFile
		err := rows.Scan(
			&f.ID, &f.Filename, &f.Extension, &f.Path, &f.Directory,
			&f.OwnerID, &f.Size, &f.ContentType, &f.Metadata, &f.TrustedSource,
			&f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *pgFileRepo) FuzzySearch(ctx context.Context, search_query string, owner *User) ([]UserFile, error) {
	query := `
		SELECT f.id, f.filename, f.extension, f.path, f.directory_id, f.owner_id, f.size, f.content_type, f.metadata, f.trusted_source, f.created_at, f.updated_at
		FROM user_files f
		WHERE f.owner_id = $2 AND f.filename % $1
		ORDER BY similarity(f.filename, $1) DESC
		LIMIT 20`

	rows, err := r.db.Query(ctx, query, search_query, owner.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []UserFile
	for rows.Next() {
		var f UserFile
		err := rows.Scan(
			&f.ID, &f.Filename, &f.Extension, &f.Path, &f.Directory,
			&f.OwnerID, &f.Size, &f.ContentType, &f.Metadata, &f.TrustedSource,
			&f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *pgFileRepo) Update(ctx context.Context, file *UserFile, updates *UserFile, owner *User) (*UserFile, error) {
	where, whereArgs := file.Identifier(1)
	if where == "" {
		return nil, fmt.Errorf("no identifier provided for update")
	}

	var setParts []string
	var args []any
	argID := 1

	if updates.Filename != "" {
		setParts = append(setParts, fmt.Sprintf("filename = $%d", argID))
		args = append(args, updates.Filename)
		argID++
	}
	if updates.Extension != "" {
		setParts = append(setParts, fmt.Sprintf("extension = $%d", argID))
		args = append(args, updates.Extension)
		argID++
	}
	if updates.Path != "" {
		setParts = append(setParts, fmt.Sprintf("path = $%d", argID))
		args = append(args, updates.Path)
		argID++
	}
	if updates.Directory != 0 {
		setParts = append(setParts, fmt.Sprintf("directory_id = $%d", argID))
		args = append(args, updates.Directory)
		argID++
	}
	if updates.Size != 0 {
		setParts = append(setParts, fmt.Sprintf("size = $%d", argID))
		args = append(args, updates.Size)
		argID++
	}
	if updates.ContentType != "" {
		setParts = append(setParts, fmt.Sprintf("content_type = $%d", argID))
		args = append(args, updates.ContentType)
		argID++
	}
	if updates.Metadata != nil {
		setParts = append(setParts, fmt.Sprintf("metadata = $%d", argID))
		args = append(args, updates.Metadata)
		argID++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	// Re-index WHERE
	where, _ = file.Identifier(argID)
	args = append(args, whereArgs...)

	query := fmt.Sprintf(`
		UPDATE user_files SET %s WHERE %s
		RETURNING id, filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source, created_at, updated_at`,
		strings.Join(setParts, ", "), where)

	var result UserFile
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&result.ID, &result.Filename, &result.Extension, &result.Path, &result.Directory,
		&result.OwnerID, &result.Size, &result.ContentType, &result.Metadata, &result.TrustedSource,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update file entry: %w", err)
	}

	return &result, nil
}

func (r *pgFileRepo) Delete(ctx context.Context, file *UserFile, owner *User) (*UserFile, error) {
	where, args := file.Identifier(1)
	if where == "" {
		return nil, fmt.Errorf("no identifier provided for delete")
	}

	query := fmt.Sprintf(`
		DELETE FROM user_files WHERE %s
		RETURNING id, filename, extension, path, directory_id, owner_id, size, content_type, metadata, trusted_source, created_at, updated_at`, where)

	var result UserFile
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&result.ID, &result.Filename, &result.Extension, &result.Path, &result.Directory,
		&result.OwnerID, &result.Size, &result.ContentType, &result.Metadata, &result.TrustedSource,
		&result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to delete file entry: %w", err)
	}

	return &result, nil
}

package db

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// User entity representing the users table
type User struct {
	ID            int       `db:"id" json:"id"`
	Username      string    `db:"username" json:"username"`
	Password      string    `db:"password" json:"password"`
	Email         string    `db:"email" json:"email"`
	AuthToken     *string   `db:"auth_token" json:"auth_token,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
}

// UserRepo defines the data access layer for Users
type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	Query(ctx context.Context, user *User) (*User, error)
	QueryAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User, updates *User) (*User, error)
	Delete(ctx context.Context, user *User) (*User, error)
}

type pgUserRepo struct {
	db QueryExecutor
}

func NewUserRepo(db QueryExecutor) UserRepo {
	return &pgUserRepo{db: db}
}

// Validation Logic

func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		return fmt.Errorf("invalid or empty email")
	}
	return nil
}

// Identifier returns the WHERE clause and arguments based on the most specific available field
func (u *User) Identifier(startParam int) (string, []any) {
	if u.ID != 0 {
		return fmt.Sprintf("id = $%d", startParam), []any{u.ID}
	}
	if u.Username != "" {
		return fmt.Sprintf("username = $%d", startParam), []any{u.Username}
	}
	if u.Email != "" {
		return fmt.Sprintf("email = $%d", startParam), []any{u.Email}
	}
	return "", nil
}

// Implementation

func (r *pgUserRepo) Create(ctx context.Context, user *User) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (username, password, email) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, email_verified`

	err := r.db.QueryRow(ctx, query, user.Username, user.Password, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.EmailVerified)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *pgUserRepo) Query(ctx context.Context, user *User) (*User, error) {
	where, args := user.Identifier(1)
	if where == "" {
		return nil, fmt.Errorf("no identifier provided for query")
	}

	query := fmt.Sprintf(`
		SELECT id, username, password, email, auth_token, created_at, email_verified 
		FROM users WHERE %s`, where)

	var result User
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Username, &result.Password, &result.Email, &result.AuthToken, &result.CreatedAt, &result.EmailVerified)

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &result, nil
}

func (r *pgUserRepo) QueryAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, username, password, email, auth_token, created_at, email_verified FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.AuthToken, &u.CreatedAt, &u.EmailVerified); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *pgUserRepo) Update(ctx context.Context, user *User, updates *User) (*User, error) {
	where, whereArgs := user.Identifier(1) // This is just for the WHERE clause
	if where == "" {
		return nil, fmt.Errorf("no identifier provided for update")
	}

	// Build dynamic SET clause
	var setParts []string
	var args []any
	argID := 1

	if updates.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argID))
		args = append(args, updates.Username)
		argID++
	}
	if updates.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argID))
		args = append(args, updates.Email)
		argID++
	}
	if updates.Password != "" {
		setParts = append(setParts, fmt.Sprintf("password = $%d", argID))
		args = append(args, updates.Password)
		argID++
	}
	if updates.AuthToken != nil {
		setParts = append(setParts, fmt.Sprintf("auth_token = $%d", argID))
		args = append(args, updates.AuthToken)
		argID++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Re-index the WHERE clause
	where, _ = user.Identifier(argID)
	args = append(args, whereArgs...)

	query := fmt.Sprintf(`
		UPDATE users SET %s WHERE %s 
		RETURNING id, username, password, email, auth_token, created_at, email_verified`,
		strings.Join(setParts, ", "), where)

	var result User
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Username, &result.Password, &result.Email, &result.AuthToken, &result.CreatedAt, &result.EmailVerified)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &result, nil
}

func (r *pgUserRepo) Delete(ctx context.Context, user *User) (*User, error) {
	where, args := user.Identifier(1)
	if where == "" {
		return nil, fmt.Errorf("no identifier provided for delete")
	}

	query := fmt.Sprintf(`
		DELETE FROM users WHERE %s 
		RETURNING id, username, password, email, auth_token, created_at, email_verified`, where)

	var result User
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&result.ID, &result.Username, &result.Password, &result.Email, &result.AuthToken, &result.CreatedAt, &result.EmailVerified)

	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &result, nil
}

package globals

import (
	"net/http"
)

type DBTable struct {
	Name        string
	SQLCreation string
}

const (
	EnvPort          = "PORT"
	EnvDbURL         = "DB_URL"
	EnvStoragePath   = "STORAGE_PATH"
	EnvEmailPassword = "EMAIL_PASSWORD"
	EnvEmailHost     = "EMAIL_HOST"
	EnvEmailPort     = "EMAIL_PORT"
	EnvEmailUser     = "EMAIL_USER"
)

var (
	DBUsersTable = &DBTable{
		Name: "users",
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS users (
    			id SERIAL PRIMARY KEY,
    			username TEXT NOT NULL,
    			password TEXT NOT NULL,
    			email TEXT NOT NULL,
    			auth_token TEXT,
    			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    			email_verified BOOLEAN NOT NULL DEFAULT false
			);`,
	}
	DBFilesTable = &DBTable{
		Name: "user_files",
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS user_files (
    			id SERIAL PRIMARY KEY,
    			filename TEXT NOT NULL,
    			path TEXT NOT NULL,
    			owner_id INTEGER NOT NULL,
    			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    			metadata JSONB,
    			size BIGINT NOT NULL DEFAULT 0,
    			content_type TEXT NOT NULL DEFAULT '',
    			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    			CONSTRAINT user_files_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id)
			);`,
	}
)

type LoginForm struct {
	Identificator string `json:"identificator"`
	Password      string `json:"password"`
}

func LoginFormFromRequest(r *http.Request) LoginForm {
	var form LoginForm

	form.Identificator = r.FormValue("identificator")
	form.Password = r.FormValue("password")
	return form
}

type RegisterForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func RegisterFormFromRequest(r *http.Request) RegisterForm {
	var form RegisterForm

	register_fields := []string{"username", "email", "password", "confirm-password"}
	for _, field := range register_fields {
		switch field {
		case "username":
			form.Username = r.FormValue(field)
		case "email":
			form.Email = r.FormValue(field)
		case "password":
			form.Password = r.FormValue(field)
		case "confirm-password":
			form.Confirm = r.FormValue(field)
		}
	}

	return form
}

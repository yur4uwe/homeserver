package globals

type DBTable struct {
	Name        string
	Fields      []string
	SQLCreation string
}

var (
	DBUsersTable = &DBTable{
		Name: "users",
		Fields: []string{
			"id",
			"username",
			"password",
			"email",
			"auth_token",
			"created_at",
			"email_verified",
		},
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				auth_token TEXT UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				email_verified BOOLEAN NOT NULL DEFAULT false
			);
		`,
	}
	DBFilesTable = &DBTable{
		Name: "user_files",
		Fields: []string{
			"id",
			"filename",
			"path",
			"size",
			"content_type",
			"owner_id",
			"created_at",
			"updated_at",
			"metadata",
		},
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
			);
		`,
	}
	DBFoldersTable = &DBTable{
		Name: "user_folders",
		Fields: []string{
			"id",
			"name",
			"parent",
			"owner_id",
			"created_at",
		},
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS user_folders (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				parent INTEGER REFERENCES user_folders(id),
				owner_id INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT user_folders_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id)
			);
		`,
	}
)

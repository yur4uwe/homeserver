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
			"directory",
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
				directory_id INTEGER NOT NULL REFERENCES user_folders(id),
				owner_id INTEGER NOT NULL,
				size BIGINT NOT NULL DEFAULT 0,
				content_type TEXT NOT NULL DEFAULT '',
				metadata JSONB,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			CREATE INDEX IF NOT EXISTS user_files_filename_trgm_idx
        	ON user_files
        	USING gin (filename gin_trgm_ops);
		`,
	}
	DBFoldersTable = &DBTable{
		Name: "user_folders",
		Fields: []string{
			"id",
			"name",
			"path",
			"parent",
			"owner_id",
			"created_at",
		},
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS user_folders (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				path TEXT NOT NULL,
				parent INTEGER REFERENCES user_folders(id),
				owner_id INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT user_folders_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id)
			);
		`,
	}
)

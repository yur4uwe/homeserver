package db

type DBTable struct {
	Name        string
	Fields      []string
	SQLCreation string
}

var (
	DBUsersTable = &DBTable{
		Name: "users",
		Fields: []string{"id", "username", "password", "email", "auth_token", "created_at", "email_verified"},
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
		Fields: []string{"id", "filename", "extension", "path", "directory_id", "owner_id", "size", "content_type", "trusted_source", "metadata", "created_at", "updated_at"},
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS user_files (
				id SERIAL PRIMARY KEY,
				filename TEXT NOT NULL,
				path TEXT NOT NULL,
				extension TEXT NOT NULL,
				directory_id INTEGER NOT NULL REFERENCES user_folders(id),
				owner_id INTEGER NOT NULL,
				size BIGINT NOT NULL DEFAULT 0,
				content_type TEXT NOT NULL DEFAULT '',
				metadata JSONB,
				trusted_source BYTEA UNIQUE,
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
		Fields: []string{"id", "name", "path", "parent", "owner_id", "trusted_source", "created_at"},
		SQLCreation: `
			CREATE TABLE IF NOT EXISTS user_folders (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				path TEXT NOT NULL,
				parent INTEGER REFERENCES user_folders(id),
				owner_id INTEGER NOT NULL,
				trusted_source BYTEA UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT user_folders_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id),
				CONSTRAINT user_folders_unique UNIQUE (name, path, owner_id)
			);
		`,
	}
)

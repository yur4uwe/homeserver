package dbconsts

import (
	"homeserver/constants/fconsts"
	"homeserver/constants/folconsts"
	"homeserver/constants/uconsts"
)

type DBTable struct {
	Name        string
	Fields      []string
	SQLCreation string
}

var (
	DBUsersTable = &DBTable{
		Name: "users",
		Fields: []string{
			uconsts.ID,
			uconsts.Username,
			uconsts.Password,
			uconsts.Email,
			uconsts.AuthToken,
			uconsts.CreatedAt,
			uconsts.EmailVerified,
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
			fconsts.ID,
			fconsts.Filename,
			fconsts.Extension,
			fconsts.Path,
			fconsts.Directory,
			fconsts.Size,
			fconsts.ContentType,
			fconsts.OwnerID,
			fconsts.CreatedAt,
			fconsts.UpdatedAt,
			fconsts.Metadata,
		},
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
			folconsts.ID,
			folconsts.FolderName,
			folconsts.Path,
			folconsts.Parent,
			folconsts.OwnerID,
			folconsts.CreatedAt,
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

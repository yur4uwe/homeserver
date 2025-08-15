package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/fconsts"
	"homeserver/constants/uconsts"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (repo *FileRepo) DeleteFileByID(ctx context.Context, fileID string, owner *uschema.User) error {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return err
	}

	tag, err := repo.db.Pool.Exec(ctx, "DELETE FROM files WHERE id = $1 AND owner_id = $2", fileID, owner.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete file")
	}

	return nil
}

func (repo *FileRepo) DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) error {
	if err := owner.ShouldHaveFields(uconsts.Username, uconsts.ID); err != nil {
		return err
	}

	if err := file.ShouldHaveFields(owner.Username)(fconsts.Filename, fconsts.Path); err != nil {
		return err
	}

	tag, err := repo.db.Pool.Exec(ctx, "DELETE FROM files WHERE filename = $1 AND path = $2 AND owner_id = $3", file.Filename, file.Path, owner.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete file")
	}

	return nil
}

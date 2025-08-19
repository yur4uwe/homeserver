package foltests

import (
	"homeserver/constants/folconsts"
	folderrepo "homeserver/db/repositories/folder_repo"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderUpdateFields_AllFields(t *testing.T) {
	now := time.Now()
	parent := 42
	folder := &folschema.UserFolder{
		ID:        1,
		Name:      "docs",
		Path:      "/home/user/docs",
		Parent:    &parent,
		OwnerID:   7,
		CreatedAt: now,
	}

	fields := []string{
		folconsts.ID,
		folconsts.FolderName,
		folconsts.Path,
		folconsts.Parent,
		folconsts.OwnerID,
		folconsts.CreatedAt,
	}
	args := folderrepo.OrderUpdateFields(fields, folder)

	assert.Equal(t, []any{1, "docs", "/home/user/docs", &parent, 7, now}, args)
}

func TestOrderUpdateFields_SubsetFields(t *testing.T) {
	parent := 42

	folder := &folschema.UserFolder{
		Name:   "music",
		Path:   "/home/user/music",
		Parent: &parent,
	}

	fields := []string{folconsts.Path, folconsts.FolderName}
	args := folderrepo.OrderUpdateFields(fields, folder)

	assert.Equal(t, []any{"/home/user/music", "music"}, args)
}

func TestOrderUpdateFields_EmptyFields(t *testing.T) {
	parent := 0
	folder := &folschema.UserFolder{
		Name:   "empty",
		Path:   "/empty",
		Parent: &parent,
	}

	fields := []string{}
	args := folderrepo.OrderUpdateFields(fields, folder)

	assert.Empty(t, args)
}

func TestOrderUpdateFields_UnknownField(t *testing.T) {
	parent := 1
	folder := &folschema.UserFolder{
		Name:   "misc",
		Path:   "/misc",
		Parent: &parent,
	}

	fields := []string{"not_a_field", folconsts.FolderName}
	args := folderrepo.OrderUpdateFields(fields, folder)

	assert.Equal(t, []any{nil, "misc"}, args)
}

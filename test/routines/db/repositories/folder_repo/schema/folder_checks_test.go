package folschema_test

import (
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderChecksSuite(t *testing.T) {
	t.Run("TestCheckFolderID", func(t *testing.T) {
		folder := &folschema.UserFolder{ID: 1}
		err := folschema.CheckFolderID(folder)
		assert.NoError(t, err, "CheckFolderID should not return an error for valid ID")

		folder = &folschema.UserFolder{ID: 0}
		err = folschema.CheckFolderID(folder)
		assert.Error(t, err, "CheckFolderID should return an error for invalid ID")
		assert.EqualError(t, err, "ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderName", func(t *testing.T) {
		folder := &folschema.UserFolder{Name: "documents"}
		err := folschema.CheckFolderName(folder)
		assert.NoError(t, err, "CheckFolderName should not return an error for valid name")

		folder = &folschema.UserFolder{Name: ""}
		err = folschema.CheckFolderName(folder)
		assert.Error(t, err, "CheckFolderName should return an error for empty name")
		assert.EqualError(t, err, "name cannot be empty", "Error message should match")
	})

	t.Run("TestCheckFolderPath", func(t *testing.T) {
		checkPath := folschema.CheckFolderPath("john_doe")

		folder := &folschema.UserFolder{Path: "john_doe/documents"}
		err := checkPath(folder)
		assert.NoError(t, err, "CheckFolderPath should not return an error for valid path")

		folder = &folschema.UserFolder{Path: ""}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for empty path")
		assert.EqualError(t, err, "path cannot be empty", "Error message should match")

		folder = &folschema.UserFolder{Path: "/invalid_path"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path starting with a slash")
		assert.EqualError(t, err, "path cannot start with a slash", "Error message should match")

		folder = &folschema.UserFolder{Path: "invalid_path/"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path ending with a slash")
		assert.EqualError(t, err, "path cannot end with a slash", "Error message should match")

		folder = &folschema.UserFolder{Path: "other_user/documents"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, "path must start with the owner's name", "Error message should match")
	})

	t.Run("TestCheckParent", func(t *testing.T) {
		parentID := 42
		folder := &folschema.UserFolder{Parent: &parentID}
		err := folschema.CheckParent(folder)
		assert.NoError(t, err, "CheckParent should not return an error for valid parent ID")

		parentID = 0
		folder = &folschema.UserFolder{Parent: &parentID}
		err = folschema.CheckParent(folder)
		assert.Error(t, err, "CheckParent should return an error for invalid parent ID")
		assert.EqualError(t, err, "parent ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderOwnerID", func(t *testing.T) {
		folder := &folschema.UserFolder{OwnerID: 100}
		err := folschema.CheckFolderOwnerID(folder)
		assert.NoError(t, err, "CheckFolderOwnerID should not return an error for valid owner ID")

		folder = &folschema.UserFolder{OwnerID: 0}
		err = folschema.CheckFolderOwnerID(folder)
		assert.Error(t, err, "CheckFolderOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, "owner ID must be greater than 0", "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		mockInvalidParentID := -1

		folder := &folschema.UserFolder{
			ID:      0,
			Name:    "",
			Path:    "/invalid_path/",
			Parent:  &mockInvalidParentID,
			OwnerID: 0,
		}

		err := folder.ShouldHaveFields("john_doe")(folconsts.ID, folconsts.FolderName, folconsts.Path, folconsts.FolderParent, folconsts.OwnerID)
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid fields")
		errStr := err.Error()
		assert.Contains(t, errStr, "id:ID must be greater than 0", "Error message should contain ID error")
		assert.Contains(t, errStr, "name:name cannot be empty", "Error message should contain name error")
		assert.Contains(t, errStr, "path:path cannot start with a slash", "Error message should contain path error")
		assert.Contains(t, errStr, "parent:parent ID must be greater than 0", "Error message should contain parent error")
		assert.Contains(t, errStr, "owner_id:owner ID must be greater than 0", "Error message should contain owner ID error")

		mockParentID := 1

		folder = &folschema.UserFolder{
			ID:      1,
			Name:    "documents",
			Path:    "john_doe/documents",
			Parent:  &mockParentID,
			OwnerID: 100,
		}

		err = folder.ShouldHaveFields("john_doe")(folconsts.ID, folconsts.FolderName, folconsts.Path, folconsts.FolderParent, folconsts.OwnerID, folconsts.FolderParent)
		assert.NoError(t, err, "ShouldHaveFields should not return an error for valid fields")

		err = folder.ShouldHaveFields("john_doe")()
		assert.Error(t, err, "ShouldHaveFields should return an error when no fields are specified")
		assert.EqualError(t, err, "no fields specified for validation", "Error message should match")

		err = folder.ShouldHaveFields("john_doe")("InvalidField")
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid field name")
		assert.Contains(t, err.Error(), "field 'InvalidField' is not a valid folder field", "Error message should match")
	})
}

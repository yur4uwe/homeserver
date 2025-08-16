package folschema

import (
	"homeserver/constants/folconsts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderChecksSuite(t *testing.T) {
	t.Run("TestCheckFolderID", func(t *testing.T) {
		folder := &UserFolder{ID: 1}
		err := CheckFolderID(folder)
		assert.NoError(t, err, "CheckFolderID should not return an error for valid ID")

		folder = &UserFolder{ID: 0}
		err = CheckFolderID(folder)
		assert.Error(t, err, "CheckFolderID should return an error for invalid ID")
		assert.EqualError(t, err, "ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderName", func(t *testing.T) {
		folder := &UserFolder{Name: "documents"}
		err := CheckFolderName(folder)
		assert.NoError(t, err, "CheckFolderName should not return an error for valid name")

		folder = &UserFolder{Name: ""}
		err = CheckFolderName(folder)
		assert.Error(t, err, "CheckFolderName should return an error for empty name")
		assert.EqualError(t, err, "name cannot be empty", "Error message should match")
	})

	t.Run("TestCheckFolderPath", func(t *testing.T) {
		checkPath := CheckFolderPath("john_doe")

		folder := &UserFolder{Path: "john_doe/documents"}
		err := checkPath(folder)
		assert.NoError(t, err, "CheckFolderPath should not return an error for valid path")

		folder = &UserFolder{Path: ""}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for empty path")
		assert.EqualError(t, err, "path cannot be empty", "Error message should match")

		folder = &UserFolder{Path: "/invalid_path"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path starting with a slash")
		assert.EqualError(t, err, "path cannot start with a slash", "Error message should match")

		folder = &UserFolder{Path: "invalid_path/"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path ending with a slash")
		assert.EqualError(t, err, "path cannot end with a slash", "Error message should match")

		folder = &UserFolder{Path: "other_user/documents"}
		err = checkPath(folder)
		assert.Error(t, err, "CheckFolderPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, "path must start with the owner's name", "Error message should match")
	})

	t.Run("TestCheckParent", func(t *testing.T) {
		parentID := 42
		folder := &UserFolder{Parent: &parentID}
		err := CheckParent(folder)
		assert.NoError(t, err, "CheckParent should not return an error for valid parent ID")

		parentID = 0
		folder = &UserFolder{Parent: &parentID}
		err = CheckParent(folder)
		assert.Error(t, err, "CheckParent should return an error for invalid parent ID")
		assert.EqualError(t, err, "parent ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderOwnerID", func(t *testing.T) {
		folder := &UserFolder{OwnerID: 100}
		err := CheckFolderOwnerID(folder)
		assert.NoError(t, err, "CheckFolderOwnerID should not return an error for valid owner ID")

		folder = &UserFolder{OwnerID: 0}
		err = CheckFolderOwnerID(folder)
		assert.Error(t, err, "CheckFolderOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, "owner ID must be greater than 0", "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		mockInvalidParentID := -1

		folder := &UserFolder{
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

		folder = &UserFolder{
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

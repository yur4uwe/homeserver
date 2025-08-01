package folder_schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderChecksSuite(t *testing.T) {
	t.Run("TestCheckFolderID", func(t *testing.T) {
		folder := &UserFolder{ID: 1}
		result, err := CheckFolderID(folder, nil)
		assert.NoError(t, err, "CheckFolderID should not return an error for valid ID")
		assert.Equal(t, folder, result, "CheckFolderID should return the same folder object")

		folder = &UserFolder{ID: 0}
		_, err = CheckFolderID(folder, nil)
		assert.Error(t, err, "CheckFolderID should return an error for invalid ID")
		assert.EqualError(t, err, "folder ID is invalid: ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderName", func(t *testing.T) {
		folder := &UserFolder{Name: "documents"}
		result, err := CheckFolderName(folder, nil)
		assert.NoError(t, err, "CheckFolderName should not return an error for valid name")
		assert.Equal(t, folder, result, "CheckFolderName should return the same folder object")

		folder = &UserFolder{Name: ""}
		_, err = CheckFolderName(folder, nil)
		assert.Error(t, err, "CheckFolderName should return an error for empty name")
		assert.EqualError(t, err, "folder name is invalid: name cannot be empty", "Error message should match")
	})

	t.Run("TestCheckFolderPath", func(t *testing.T) {
		checkPath := CheckFolderPath("john_doe")

		folder := &UserFolder{Path: "john_doe/documents"}
		result, err := checkPath(folder, nil)
		assert.NoError(t, err, "CheckFolderPath should not return an error for valid path")
		assert.Equal(t, folder, result, "CheckFolderPath should return the same folder object")

		folder = &UserFolder{Path: ""}
		_, err = checkPath(folder, nil)
		assert.Error(t, err, "CheckFolderPath should return an error for empty path")
		assert.EqualError(t, err, "path is invalid: path cannot be empty", "Error message should match")

		folder = &UserFolder{Path: "/invalid_path"}
		_, err = checkPath(folder, nil)
		assert.Error(t, err, "CheckFolderPath should return an error for path starting with a slash")
		assert.EqualError(t, err, "path is invalid: path cannot start with a slash", "Error message should match")

		folder = &UserFolder{Path: "invalid_path/"}
		_, err = checkPath(folder, nil)
		assert.Error(t, err, "CheckFolderPath should return an error for path ending with a slash")
		assert.EqualError(t, err, "path is invalid: path cannot end with a slash", "Error message should match")

		folder = &UserFolder{Path: "other_user/documents"}
		_, err = checkPath(folder, nil)
		assert.Error(t, err, "CheckFolderPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, "path is invalid: path must start with the owner's name 'john_doe'", "Error message should match")
	})

	t.Run("TestCheckParent", func(t *testing.T) {
		parentID := 42
		folder := &UserFolder{Parent: &parentID}
		result, err := CheckParent(folder, nil)
		assert.NoError(t, err, "CheckParent should not return an error for valid parent ID")
		assert.Equal(t, folder, result, "CheckParent should return the same folder object")

		parentID = 0
		folder = &UserFolder{Parent: &parentID}
		_, err = CheckParent(folder, nil)
		assert.Error(t, err, "CheckParent should return an error for invalid parent ID")
		assert.EqualError(t, err, "parent is invalid: parent ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFolderOwnerID", func(t *testing.T) {
		folder := &UserFolder{OwnerID: 100}
		result, err := CheckFolderOwnerID(folder, nil)
		assert.NoError(t, err, "CheckFolderOwnerID should not return an error for valid owner ID")
		assert.Equal(t, folder, result, "CheckFolderOwnerID should return the same folder object")

		folder = &UserFolder{OwnerID: 0}
		_, err = CheckFolderOwnerID(folder, nil)
		assert.Error(t, err, "CheckFolderOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, "owner ID is invalid: owner ID must be greater than 0", "Error message should match")
	})
}

package db

import (
	"context"
	"fmt"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMethods(t *testing.T) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	defer func() {
		for _, table := range setup.RequiredDBTables {
			err := connection.DropTable(context.Background(), table.Name)
			if err != nil {
				t.Errorf("Failed to drop table %s: %v", table.Name, err)
			}
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	err := setup.DatabaseSetup(connection)
	assert.NoError(t, err, "Failed to check database tables")

	mockUser := user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("securepassword")

	mockIncompleteUser := user_schema.NewUserBuilder().
		WithUsername("incompleteuser").
		WithEmail("incompleteuser@example.com")

	t.Run("CreateUserEntry happy path", func(t *testing.T) {
		con_copy := connection.Copy()
		end_transaction := con_copy.UseTransaction()
		defer end_transaction()

		inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
		assert.NoError(t, err, "Failed to create user")
		assert.NotNil(t, inserted_user, "Inserted user should not be nil")
		assert.Equal(t, mockUser.Username, inserted_user.Username, "Usernames should match")
		assert.Equal(t, mockUser.Email, inserted_user.Email, "Emails should match")
		assert.NotEmpty(t, inserted_user.ID, "User ID should not be empty")
		assert.NotEmpty(t, inserted_user.CreatedAt, "User creation time should not be empty")
	})

	t.Run("CreateUserEntry with missing fields", func(t *testing.T) {
		con_copy := connection.Copy()
		end_transaction := con_copy.UseTransaction()
		defer end_transaction()

		_, err := con_copy.CreateUserEntry(context.Background(), mockIncompleteUser)
		assert.Error(t, err, "Creating user with missing fields should return an error")
		assert.Contains(t, err.Error(), "invalid user data", "Error message should indicate invalid user data")
	})

	t.Run("DeleteUser happy path", func(t *testing.T) {
		con_copy := connection.Copy()
		end_transaction := con_copy.UseTransaction()
		defer end_transaction()

		inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
		assert.NoError(t, err, "Failed to create user")
		assert.NotNil(t, inserted_user, "Inserted user should not be nil")

		err = con_copy.DeleteUser(context.Background(), inserted_user)
		assert.NoError(t, err, "Failed to delete user")
	})

	t.Run("DeleteUser with missing ID", func(t *testing.T) {
		con_copy := connection.Copy()
		end_transaction := con_copy.UseTransaction()
		defer end_transaction()

		err := con_copy.DeleteUser(context.Background(), mockIncompleteUser)
		assert.Error(t, err, "Deleting user with missing ID should return an error")
		assert.Contains(t, err.Error(), "invalid user data", "Error message should indicate invalid user data")
	})

	t.Run("ExtractUserEntries", func(t *testing.T) {
		t.Skip("Skipping ExtractUserEntries test as it requires a populated database")
		// Implement the test for extracting user entries
	})

	t.Run("CheckUserFolder", func(t *testing.T) {
		t.Skip("Skipping CheckUserFolder test as it requires a valid user folder")
		// Implement the test for checking user folder
	})

	t.Run("QueryUserOnLogin", func(t *testing.T) {
		t.Skip("Skipping QueryUserOnLogin test as it requires a valid login form")
		// Implement the test for querying user on login
	})

	t.Run("QueryUsername", func(t *testing.T) {
		t.Skip("Skipping QueryUsername test as it requires a valid username")
		// Implement the test for querying username
	})

	t.Run("QueryEmail", func(t *testing.T) {
		t.Skip("Skipping QueryEmail test as it requires a valid email")
		// Implement the test for querying email
	})

	t.Run("QueryPassword", func(t *testing.T) {
		t.Skip("Skipping QueryPassword test as it requires a valid password")
		// Implement the test for querying password
	})
}

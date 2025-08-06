package db

import (
	"context"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupUserFileTestEnvironment(t *testing.T) (*Database, cleanupFunc) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	err := setup.DatabaseSetup(connection)
	assert.NoError(t, err, "Failed to check database tables")

	mockUser = user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("securepassword")

	mockSecondUser = user_schema.NewUserBuilder().
		WithUsername("seconduser").
		WithEmail("seconduser@example.com").
		WithPassword("securepassword")

	mockIncompleteUser = user_schema.NewUserBuilder().
		WithUsername("incompleteuser").
		WithEmail("incompleteuser@example.com")

	con_copy := connection.Copy()
	end_transaction := con_copy.UseTransaction()
	assert.NotNil(t, end_transaction)
	return con_copy, end_transaction
}

func cleanupUserTestEnvironment(t *testing.T, cleanupFunc cleanupFunc) {
	cleanupFunc()
	for _, table := range setup.RequiredDBTables {
		err := connection.DropTable(context.Background(), table.Name)
		if err != nil {
			t.Errorf("Failed to drop table %s: %v", table.Name, err)
		}
	}
}

func TestCreateUserEntryHappyPath(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.Equal(t, mockUser.Username, inserted_user.Username, "Usernames should match")
	assert.Equal(t, mockUser.Email, inserted_user.Email, "Emails should match")
	assert.NotEmpty(t, inserted_user.ID, "User ID should not be empty")
	assert.NotEmpty(t, inserted_user.CreatedAt, "User creation time should not be empty")
}

func TestCreateUserEntryWithMissingFields(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	_, err := con_copy.CreateUserEntry(context.Background(), mockIncompleteUser)
	assert.Error(t, err, "Creating user with missing fields should return an error")
	assert.Contains(t, err.Error(), "invalid user data", "Error message should indicate invalid user data")
}

func TestDeleteUserHappyPath(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")

	err = con_copy.DeleteUser(context.Background(), inserted_user)
	assert.NoError(t, err, "Failed to delete user")
}

func TestDeleteUserWithMissingID(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	err := con_copy.DeleteUser(context.Background(), mockIncompleteUser)
	assert.Error(t, err, "Deleting user with missing ID should return an error")
	assert.Contains(t, err.Error(), "invalid user data", "Error message should indicate invalid user data")
}

func TestExtractUserEntries(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user1, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user1, "User should not be nil")

	user2, err := con_copy.CreateUserEntry(context.Background(), mockSecondUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user2, "User should not be nil")

	users, err := con_copy.ExtractUserEntries(context.Background())
	assert.NoError(t, err, "Failed to extract user entries")
	assert.NotEmpty(t, users, "Users should not be empty")
	assert.Len(t, users, 2, "There should be two users in the database")
	assert.Contains(t, users, user1, "Extracted users should contain the first user")
	assert.Contains(t, users, user2, "Extracted users should contain the second user")
}

func TestCheckUserFolder(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user, err := con_copy.CreateUserEntry(context.Background(), mockSecondUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user, "User should not be nil")

	err = con_copy.CheckUserFolder(context.Background(), user)
	assert.NoError(t, err, "Failed to check user folder")

	exists, err := con_copy.FolderExists(context.Background(), user.Username, user.ID, "")
	assert.NoError(t, err, "Failed to create root folder")
	assert.True(t, exists, "User folder should exist")
}

func TestQueryUserOnLogin(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user, "User should not be nil")

	loginForm := globals.LoginForm{
		Identificator: mockUser.Username,
		Password:      mockUser.Password,
	}

	retrieved_user, email_verified, err := con_copy.QueryUserOnLogin(context.Background(), loginForm)
	assert.NoError(t, err, "Failed to query user on login")
	assert.False(t, email_verified, "User email should not be verified")
	assert.Equal(t, user.ID, retrieved_user.ID, "User ID should match the created user ID")
}

func TestQueryUsername(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user, "User should not be nil")

	retrieved_user, err := con_copy.QueryUsername(context.Background(), mockUser.Username)
	assert.NoError(t, err, "Failed to query username")
	assert.Equal(t, user.ID, retrieved_user.ID, "Username ID should match the created user ID")
}

func TestQueryEmail(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user, "User should not be nil")

	retrieved_user, err := con_copy.QueryEmail(context.Background(), mockUser.Email)
	assert.NoError(t, err, "Failed to query email")
	assert.Equal(t, user.ID, retrieved_user.ID, "Email ID should match the created user ID")
}

func TestQueryPassword(t *testing.T) {
	con_copy, end_transaction := setupUserFileTestEnvironment(t)
	defer cleanupUserTestEnvironment(t, end_transaction)

	user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, user, "User should not be nil")

	retrieved_user, err := con_copy.QueryPassword(context.Background(), mockUser.Password)
	assert.NoError(t, err, "Failed to query password")
	assert.Equal(t, user.ID, retrieved_user.ID, "Password ID should match the created user ID")
}

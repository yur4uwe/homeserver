package db

import (
	"context"
	"homeserver/pkg/db/token"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"
	"homeserver/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTokenTestEnvironment(t *testing.T) (*Database, cleanupFunc) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	mockOS := new(mocks.MockOS)
	mockOS.On("Mkdir", mock.Anything, mock.Anything).Return(nil)

	cleanupMkdir := test.MockFunction(&token.Mkdir, mockOS.Mkdir, func() {})

	err := setup.DatabaseSetup(connection)
	assert.NoError(t, err, "Failed to check database tables")

	mockUser = user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("securepassword")

	con_copy := connection.Copy()
	end_transaction := con_copy.UseTransaction()
	assert.NotNil(t, end_transaction)

	cleanup := func() {
		cleanupMkdir()
		end_transaction()
	}

	return con_copy, cleanup
}

func cleanupTokenTestEnvironment(t *testing.T, cleanupFunc cleanupFunc) {
	cleanupFunc()
	for _, table := range setup.RequiredDBTables {
		err := connection.DropTable(context.Background(), table.Name)
		if err != nil {
			t.Errorf("Failed to drop table %s: %v", table.Name, err)
		}
	}
}

func TestCreateToken(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "User ID should not be zero")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	retrieved_user, err := con_copy.QueryUsername(context.Background(), inserted_user.Username)
	assert.NoError(t, err, "Failed to retrieve user by username")
	assert.NotNil(t, retrieved_user, "Retrieved user should not be nil")
	assert.Equal(t, inserted_user.ID, retrieved_user.ID, "User IDs should match")
	assert.Equal(t, "testtoken123", *retrieved_user.AuthToken, "Auth token should match")
}

func TestQueryToken(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "User ID should not be zero")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	retrieved_user, err := con_copy.QueryToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to retrieve user by token")
	assert.NotNil(t, retrieved_user, "Retrieved user should not be nil")
	assert.Equal(t, inserted_user.ID, retrieved_user.ID, "User IDs should match")
	assert.Equal(t, "testtoken123", *retrieved_user.AuthToken, "Auth token should match")
}

func TestVerifyToken(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "User ID should not be zero")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	err = con_copy.VerifyToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to verify token")

	retrieved_user, err := con_copy.QueryToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to retrieve user by token")
	assert.NotNil(t, retrieved_user, "Retrieved user should not be nil")
	assert.Equal(t, inserted_user.ID, retrieved_user.ID, "User IDs should match")
	assert.True(t, retrieved_user.EmailVerified, "Email should be verified after token verification")
}

func TestCheckToken(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "User ID should not be zero")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	err = con_copy.VerifyToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to verify token")

	retrieved_user_id, err := con_copy.CheckToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to retrieve user by token")
	assert.NotEqual(t, 0, retrieved_user_id, "Retrieved user ID should not be zero")
}

func TestDeleteToken(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "User ID should not be zero")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	err = con_copy.DeleteToken(context.Background(), "testtoken123")
	assert.NoError(t, err, "Failed to delete token")

	retrieved_user, err := con_copy.QueryToken(context.Background(), "testtoken123")
	assert.Error(t, err, "Querying deleted token should return an error")
	assert.Nil(t, retrieved_user, "Retrieved user should be nil after token deletion")
}

func TestCreateTokenUnhappyPath(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	err := con_copy.CreateToken(context.Background(), 0, "testtoken123")
	assert.Error(t, err, "Creating token with invalid user ID should return an error")
	assert.Contains(t, err.Error(), token.ErrInvalidUserID, "Error message should indicate token creation failure")

	err = con_copy.CreateToken(context.Background(), 1, "")
	assert.Error(t, err, "Creating token with empty token should return an error")
	assert.Contains(t, err.Error(), token.ErrEmptyToken, "Error message should indicate token creation failure")

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")

	err = con_copy.CreateToken(context.Background(), 3, "testtoken123")
	assert.Error(t, err, "Creating token with empty token should return an error")
	assert.Contains(t, err.Error(), "no user found with ID 3", "Error message should indicate token creation failure")
}

func TestQueryTokenUnhappyPath(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	retrieved_user, err := con_copy.QueryToken(context.Background(), "")
	assert.Error(t, err, "Querying token with empty token should return an error")
	assert.Nil(t, retrieved_user, "Retrieved user should be nil for empty token")
	assert.Contains(t, err.Error(), token.ErrEmptyToken, "Error message should indicate empty token")

	retrieved_user, err = con_copy.QueryToken(context.Background(), "nonexistenttoken")
	assert.Error(t, err, "Querying token with non-existent token should return an error")
	assert.Nil(t, retrieved_user, "Retrieved user should be nil for non-existent token")
	assert.Contains(t, err.Error(), "error querying token", "Error message should indicate token query failure")
}

func TestVerifyTokenUnhappyPath(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	err := con_copy.VerifyToken(context.Background(), "")
	assert.Error(t, err, "Verifying token with empty token should return an error")
	assert.Contains(t, err.Error(), token.ErrEmptyToken, "Error message should indicate token verification failure")

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create user")
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")

	err = con_copy.CreateToken(context.Background(), inserted_user.ID, "testtoken123")
	assert.NoError(t, err, "Failed to create token")

	err = con_copy.VerifyToken(context.Background(), "nonexistenttoken")
	assert.Error(t, err, "Verifying token with non-existent token should return an error")
	assert.Contains(t, err.Error(), "error verifying token", "Error message should indicate token verification failure")

	err = con_copy.pool.QueryRow(
		context.Background(),
		`UPDATE users SET email_verified = true WHERE auth_token = $1 RETURNING username`,
		"testtoken123",
	).Scan(&inserted_user.Username)
	assert.NoError(t, err, "Failed to update email verification status")
	assert.NotEmpty(t, inserted_user.Username, "Username should not be empty after verification")

	err = con_copy.VerifyToken(context.Background(), "testtoken123")
	assert.Error(t, err, "Verifying token should fail if email is already verified")
	assert.Contains(t, err.Error(), "error verifying token", "Error message should indicate token verification failure")
}

func TestCheckTokenUnhappyPath(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	retrieved_user_id, err := con_copy.CheckToken(context.Background(), "")
	assert.Error(t, err, "Checking token with empty token should return an error")
	assert.Equal(t, 0, retrieved_user_id, "Retrieved user ID should be zero for empty token")
	assert.Contains(t, err.Error(), token.ErrEmptyToken, "Error message should indicate token check failure")

	retrieved_user_id, err = con_copy.CheckToken(context.Background(), "nonexistenttoken")
	assert.Error(t, err, "Checking token with non-existent token should return an error")
	assert.Equal(t, 0, retrieved_user_id, "Retrieved user ID should be zero for non-existent token")
	assert.Contains(t, err.Error(), "error checking token", "Error message should indicate token check failure")
}

func TestDeleteTokenUnhappyPath(t *testing.T) {
	con_copy, end_transaction := setupTokenTestEnvironment(t)
	defer cleanupTokenTestEnvironment(t, end_transaction)

	err := con_copy.DeleteToken(context.Background(), "")
	assert.Error(t, err, "Deleting token with empty token should return an error")
	assert.Contains(t, err.Error(), token.ErrEmptyToken, "Error message should indicate token deletion failure")

	err = con_copy.DeleteToken(context.Background(), "nonexistenttoken")
	assert.Error(t, err, "Deleting non-existent token should not return an error")
	assert.Contains(t, err.Error(), "no token found with value nonexistenttoken", "Error message should indicate token deletion failure")
}

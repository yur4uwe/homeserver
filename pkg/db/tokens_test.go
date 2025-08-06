package db

import (
	"context"
	"homeserver/pkg/db/token"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTokenTestEnvironment(t *testing.T) (*Database, cleanupFunc) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	mockOS := new(MockOS)
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

package userrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"homeserver/db"
	userrepo "homeserver/db/repositories/user_repo"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeUserForUpdate() *uschema.User {
	return &uschema.User{
		ID:       10,
		Username: "john_doe",
		Password: "securepassword",
		Email:    "john@example.com",
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	fieldsToUpdate := map[string]any{"username": "john_updated", "email": "john_updated@example.com"}

	expected := &uschema.User{
		ID:            10,
		Username:      "john_updated",
		Password:      "securepassword",
		Email:         "john_updated@example.com",
		AuthToken:     nil,
		CreatedAt:     time.Now(),
		EmailVerified: false,
	}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{
			expected.ID,
			expected.Username,
			"",
			expected.Email,
			expected.AuthToken,
			expected.CreatedAt,
			expected.EmailVerified,
		}
		mocks.RunScan(args[0].([]any), values)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.Update(context.Background(), user, fieldsToUpdate)
	assert.NoError(t, err)
	assert.Equal(t, expected.Username, user.Username)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.EmailVerified, user.EmailVerified)
}

func TestUpdateUser_NoFields(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	err := repo.Update(context.Background(), user, map[string]any{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
}

func TestUpdateUser_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	fieldsToUpdate := map[string]any{"username": "john_updated", "email": "john_updated@example.com"}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.Update(context.Background(), user, fieldsToUpdate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestUpdateUser_InvalidUserData(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	user.ID = 0 // Should trigger validation error

	err := repo.Update(context.Background(), user, map[string]any{"username": "john_updated"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user data")
}

func TestUpdateUser_UsernameAlreadyExists(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	existingUser := makeUserForUpdate()
	existingUser.Username = "existing_user"

	fieldsToUpdate := map[string]any{"username": "john_updated"}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(errors.New("username already exists"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.Update(context.Background(), user, fieldsToUpdate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
}

func TestUpdateUser_EmailAlreadyExists(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	existingUser := makeUserForUpdate()
	existingUser.Email = "existing_email@example.com"

	fieldsToUpdate := map[string]any{"email": "john_updated@example.com"}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(errors.New("email already exists"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.Update(context.Background(), user, fieldsToUpdate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestUpdateUser_InvalidField(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForUpdate()
	fieldsToUpdate := map[string]any{"invalid_field": "value"}

	err := repo.Update(context.Background(), user, fieldsToUpdate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field")
}

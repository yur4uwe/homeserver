package userrepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	userrepo "homeserver/db/repositories/user_repo"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeUser() *uschema.User {
	return &uschema.User{
		Username: "john_doe",
		Password: "securepassword",
		Email:    "john@example.com",
	}
}

func TestCreateUser_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUser()
	expected := &uschema.User{
		ID:       10,
		Username: user.Username,
		Email:    user.Email,
	}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Username, expected.Email}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})

	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	inserted, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, inserted.ID)
	assert.Equal(t, expected.Username, inserted.Username)
	assert.Equal(t, expected.Email, inserted.Email)
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUser()
	user.Username = "" // Should trigger validation error

	inserted, err := repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "invalid user data")
}

func TestCreateUser_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUser()

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	inserted, err := repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "failed to insert user")
}

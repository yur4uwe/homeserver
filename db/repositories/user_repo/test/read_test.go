package userrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"homeserver/db"
	userrepo "homeserver/db/repositories/user_repo"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeUserForRead() *uschema.User {
	return &uschema.User{
		ID:       10,
		Username: "john_doe",
		Email:    "john@example.com",
	}
}

func TestQueryUser_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForRead()
	expected := &uschema.User{
		ID:            10,
		Username:      "john_doe",
		Email:         "john@example.com",
		AuthToken:     nil,
		CreatedAt:     time.Now(),
		EmailVerified: true,
	}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{
			expected.ID,
			expected.Username,
			expected.Email,
			expected.AuthToken,
			expected.CreatedAt,
			expected.EmailVerified,
		}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})

	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	result, err := repo.Query(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Username, result.Username)
	assert.Equal(t, expected.Email, result.Email)
	assert.Equal(t, expected.EmailVerified, result.EmailVerified)
}

func TestQueryUser_NoValidFields(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := &uschema.User{}
	result, err := repo.Query(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no valid fields to query")
}

func TestQueryUser_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForRead()

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	result, err := repo.Query(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error querying user")
}

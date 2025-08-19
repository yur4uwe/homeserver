package userrepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	userrepo "homeserver/db/repositories/user_repo"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/jackc/pgx/v5/pgconn"
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

var fieldDescs = []pgconn.FieldDescription{
	{Name: "id"},
	{Name: "username"},
	{Name: "email"},
	{Name: "created_at"},
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

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		mocks.RunScan(ptrs, []any{expected.ID, expected.Username, expected.Email})
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	inserted, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, inserted.ID)
	assert.Equal(t, expected.Username, inserted.Username)
	assert.Equal(t, expected.Email, inserted.Email)

	mockRows.AssertExpectations(t)
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

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("failed to insert user"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	inserted, err := repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "failed to insert user")
}

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

func makeUserForDelete() *uschema.User {
	return &uschema.User{
		ID:       10,
		Username: "john_doe",
		Password: "securepassword",
		Email:    "john@example.com",
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForDelete()
	tag := pgconn.NewCommandTag("DELETE 1")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(tag, nil)

	err := repo.Delete(context.Background(), user)
	assert.NoError(t, err)
}

func TestDeleteUser_InvalidUserData(t *testing.T) {
	mockDB := &db.Database{}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForDelete()
	user.ID = 0 // Should trigger validation error

	err := repo.Delete(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user data")
}

func TestDeleteUser_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForDelete()
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(pgconn.CommandTag{}, errors.New("db error"))

	err := repo.Delete(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting user entry")
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := userrepo.NewUserRepo(mockDB)

	user := makeUserForDelete()
	tag := pgconn.NewCommandTag("DELETE 0")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(tag, nil)

	err := repo.Delete(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

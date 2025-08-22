package tokenrepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	tokenrepo "homeserver/db/repositories/token_repo"
	"homeserver/test/helpers/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockTag := pgconn.NewCommandTag("UPDATE 1")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.CreateToken(context.Background(), 42, "token123")
	assert.NoError(t, err)
}

func TestCreateToken_InvalidUserID(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	err := repo.CreateToken(context.Background(), 0, "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
}

func TestCreateToken_EmptyToken(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	err := repo.CreateToken(context.Background(), 42, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty token")
}

func TestCreateToken_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(pgconn.CommandTag{}, errors.New("db error"))

	err := repo.CreateToken(context.Background(), 42, "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating token")
}

func TestCreateToken_NotFound(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockTag := pgconn.NewCommandTag("UPDATE 0")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.CreateToken(context.Background(), 42, "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no user found with")
}

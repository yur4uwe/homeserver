package tokenrepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	tokenrepo "homeserver/db/repositories/token_repo"
	"homeserver/test/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteToken_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockTag := pgconn.NewCommandTag("DELETE 1")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.DeleteToken(context.Background(), "token123")
	assert.NoError(t, err)
}

func TestDeleteToken_EmptyToken(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	err := repo.DeleteToken(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty token")
}

func TestDeleteToken_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(pgconn.CommandTag{}, errors.New("db error"))

	err := repo.DeleteToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting token")
}

func TestDeleteToken_NotFound(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockTag := pgconn.NewCommandTag("DELETE 0")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.DeleteToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no token found with value token123")
}

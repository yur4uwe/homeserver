package tokenrepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/constants/tconsts"
	"homeserver/db"
	tokenrepo "homeserver/db/repositories/token_repo"
	"homeserver/test"
	"homeserver/test/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyToken_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockOS := &mocks.MockOS{}
	defer test.MockFunction(&tokenrepo.Getenv, mockOS.Getenv, func() {})()
	defer test.MockFunction(&tokenrepo.Mkdir, mockOS.Mkdir, func() {})()

	mockOS.On("Mkdir", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("Getenv", mock.Anything).Return("/storage/")

	// Mock SELECT
	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 42
		*dest[1].(*string) = "john_doe"
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	// Mock UPDATE
	mockTag := pgconn.NewCommandTag("UPDATE 1")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.VerifyToken(context.Background(), "token123")
	assert.NoError(t, err)
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	err := repo.VerifyToken(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), tconsts.ErrEmptyToken)
}

func TestVerifyToken_SelectError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	err := repo.VerifyToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error verifying token")
}

func TestVerifyToken_InvalidUserID(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 0
		*dest[1].(*string) = "john_doe"
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	err := repo.VerifyToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID for token verification")
}

func TestVerifyToken_UpdateError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 42
		*dest[1].(*string) = "john_doe"
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(pgconn.CommandTag{}, errors.New("update error"))

	err := repo.VerifyToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating email verification status")
}

func TestVerifyToken_AffectedRowsMismatch(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 42
		*dest[1].(*string) = "john_doe"
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	mockTag := pgconn.NewCommandTag("UPDATE 0")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.VerifyToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "affected rows mismatch")
}

func TestVerifyToken_MkdirError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockOS := &mocks.MockOS{}
	defer test.MockFunction(&tokenrepo.Getenv, mockOS.Getenv, func() {})()
	defer test.MockFunction(&tokenrepo.Mkdir, mockOS.Mkdir, func() {})()

	mockOS.On("Mkdir", mock.Anything, mock.Anything).Return(errors.New("error creating directory"))
	mockOS.On("Getenv", mock.Anything).Return("/storage/")

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 42
		*dest[1].(*string) = "john_doe"
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	mockTag := pgconn.NewCommandTag("UPDATE 1")
	mockPool.On("Exec", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockTag, nil)

	err := repo.VerifyToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user directory")
}

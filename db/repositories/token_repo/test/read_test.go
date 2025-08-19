package tokenrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"homeserver/db"
	tokenrepo "homeserver/db/repositories/token_repo"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fieldDescs = []pgconn.FieldDescription{
	{Name: "id"},
	{Name: "username"},
	{Name: "email"},
	{Name: "auth_token"},
	{Name: "email_verified"},
	{Name: "created_at"},
}

func makeUser() *uschema.User {
	token := "token123"
	return &uschema.User{
		ID:            42,
		Username:      "john_doe",
		Password:      "hashed",
		Email:         "john@example.com",
		AuthToken:     &token,
		CreatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EmailVerified: true,
	}
}

func TestCheckToken_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 42
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	userID, err := repo.CheckToken(context.Background(), "token123")
	assert.NoError(t, err)
	assert.Equal(t, 42, userID)
}

func TestCheckToken_EmptyToken(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	userID, err := repo.CheckToken(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "empty token")
}

func TestCheckToken_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	userID, err := repo.CheckToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "error checking token")
}

func TestQueryToken_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	expected := makeUser()
	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{
			expected.ID, expected.Username, expected.Email,
			expected.AuthToken, expected.EmailVerified, expected.CreatedAt,
		}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	user, err := repo.QueryToken(context.Background(), "token123")
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.Username, user.Username)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.AuthToken, user.AuthToken)
	assert.Equal(t, expected.CreatedAt, user.CreatedAt)
	assert.Equal(t, expected.EmailVerified, user.EmailVerified)
}

func TestQueryToken_EmptyToken(t *testing.T) {
	mockDB := &db.Database{}
	repo := tokenrepo.NewTokenRepo(mockDB)

	user, err := repo.QueryToken(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "empty token")
}

func TestQueryToken_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := tokenrepo.NewTokenRepo(mockDB)

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("failed to insert user"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	user, err := repo.QueryToken(context.Background(), "token123")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "error querying token")
}

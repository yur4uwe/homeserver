package helpers

import (
	"homeserver/test/mocks"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScanRows(t *testing.T) {
	mockRows := new(mocks.MockRows)

	// Simulate returned columns
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
		{Name: "password"},
		{Name: "email"},
		{Name: "created_at"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	// Prepare expected scan values
	expectedUsername := "alice"
	expectedPassword := "secret"
	expectedEmail := "alice@example.com"
	expectedCreated := time.Now()

	// The scan args will be pointers to struct fields, so we capture them and set their values
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		*(ptrs[0].(*string)) = expectedUsername
		*(ptrs[1].(*string)) = expectedPassword
		*(ptrs[2].(*string)) = expectedEmail
		*(ptrs[3].(*time.Time)) = expectedCreated
	})

	type User struct {
		Uname    string    `db:"username"`
		Upass    string    `db:"password"`
		Uemail   string    `db:"email"`
		Ucreated time.Time `db:"created_at"`
	}

	user := &User{}
	err := ScanRows(mockRows, user)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsername, user.Uname)
	assert.Equal(t, expectedPassword, user.Upass)
	assert.Equal(t, expectedEmail, user.Uemail)
	assert.Equal(t, expectedCreated, user.Ucreated)
}

func TestScanRows_MissingTag(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
		{Name: "password"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		values := []any{"bob", "hunter2"}
		mocks.RunScan(ptrs, values)
	})

	type User struct {
		Uname string `db:"username"`
		// No tag for password
		Pass string
	}

	user := &User{}
	err := ScanRows(mockRows, user)
	assert.NoError(t, err)
	assert.Equal(t, "bob", user.Uname)
	assert.Equal(t, "", user.Pass) // Should remain zero value
}

func TestScanRows_ExtraColumn(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
		{Name: "password"},
		{Name: "extra_col"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		vals := []any{"carol", "hunter3", nil} // extra_col is ignored
		mocks.RunScan(ptrs, vals)
		// ptrs[2] is dummy, should not panic
	})

	type User struct {
		Uname string `db:"username"`
		Upass string `db:"password"`
	}

	user := &User{}
	err := ScanRows(mockRows, user)
	assert.NoError(t, err)
	assert.Equal(t, "carol", user.Uname)
	assert.Equal(t, "hunter3", user.Upass)
}

func TestScanRows_NoRows(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(false).Once()

	type User struct {
		Uname string `db:"username"`
	}

	user := &User{}
	err := ScanRows(mockRows, user)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestScanRows_ScanError(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(assert.AnError)

	type User struct {
		Uname string `db:"username"`
	}

	user := &User{}
	err := ScanRows(mockRows, user)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestScanRows_NilDest(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(nil)

	err := ScanRows(mockRows, nil)
	assert.Error(t, err)
}

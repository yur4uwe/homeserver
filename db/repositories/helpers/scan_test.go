package helpers

import (
	"homeserver/test/mocks"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScanRows_SingleRow(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
		{Name: "password"},
		{Name: "email"},
		{Name: "created_at"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	expectedUsername := "alice"
	expectedPassword := "secret"
	expectedEmail := "alice@example.com"
	expectedCreated := time.Now()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		vals := []any{expectedUsername, expectedPassword, expectedEmail, expectedCreated}
		mocks.RunScan(ptrs, vals)
	})

	type User struct {
		Uname    string    `db:"username"`
		Upass    string    `db:"password"`
		Uemail   string    `db:"email"`
		Ucreated time.Time `db:"created_at"`
	}

	var users []User
	err := ScanRows(mockRows, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, expectedUsername, users[0].Uname)
	assert.Equal(t, expectedPassword, users[0].Upass)
	assert.Equal(t, expectedEmail, users[0].Uemail)
	assert.Equal(t, expectedCreated, users[0].Ucreated)
}

func TestScanRows_MultipleRows(t *testing.T) {
	mockRows := new(mocks.MockRows)
	fieldDescs := []pgconn.FieldDescription{
		{Name: "username"},
	}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		mocks.RunScan(ptrs, []any{"bob"})
	})
	mockRows.On("Scan", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		ptrs := args.Get(0).([]any)
		mocks.RunScan(ptrs, []any{"carol"})
	})

	type User struct {
		Uname string `db:"username"`
	}

	var users []User
	err := ScanRows(mockRows, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "bob", users[0].Uname)
	assert.Equal(t, "carol", users[1].Uname)

	mockRows.AssertExpectations(t)
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
		mocks.RunScan(ptrs, []any{"bob", "hunter2"})
	})

	type User struct {
		Uname string `db:"username"`
		// No tag for password
		Pass string
	}

	var users []User
	err := ScanRows(mockRows, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "bob", users[0].Uname)
	assert.Equal(t, "", users[0].Pass) // Should remain zero value
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
		mocks.RunScan(ptrs, []any{"carol", "hunter3", nil}) // extra_col is ignored
	})

	type User struct {
		Uname string `db:"username"`
		Upass string `db:"password"`
	}

	var users []User
	err := ScanRows(mockRows, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "carol", users[0].Uname)
	assert.Equal(t, "hunter3", users[0].Upass)
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

	var users []User
	err := ScanRows(mockRows, &users)
	assert.NoError(t, err)
	assert.Len(t, users, 0)
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

	var users []User
	err := ScanRows(mockRows, &users)
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

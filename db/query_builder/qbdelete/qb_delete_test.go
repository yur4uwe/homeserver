package qbdelete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteQueryBuilder_Build(t *testing.T) {
	t.Run("Simple Delete", func(t *testing.T) {
		qb := New().From("users").Where("id = 1")
		expected := "DELETE FROM users WHERE id = 1"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Returning", func(t *testing.T) {
		qb := New().From("users").Where("id = 1").Returning("id")
		expected := "DELETE FROM users WHERE id = 1 RETURNING id"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Multiple Returning", func(t *testing.T) {
		qb := New().From("users").Where("id = 1").Returning("id").Returning("username")
		expected := "DELETE FROM users WHERE id = 1 RETURNING id, username"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Raw Clause", func(t *testing.T) {
		qb := New().From("users").Where("id = 1").Raw("CASCADE")
		expected := "DELETE FROM users WHERE id = 1 CASCADE"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete Without Where", func(t *testing.T) {
		qb := New().From("users")
		expected := "DELETE FROM users"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete Without Table Panics", func(t *testing.T) {
		qb := New().Where("id = 1")
		assert.PanicsWithValue(t, "No table specified for DELETE", func() { qb.Build() })
	})

	t.Run("Delete With Only Returning", func(t *testing.T) {
		qb := New().From("users").Returning("id")
		expected := "DELETE FROM users RETURNING id"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Only Raw", func(t *testing.T) {
		qb := New().From("users").Raw("CASCADE")
		expected := "DELETE FROM users CASCADE"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Where And Raw", func(t *testing.T) {
		qb := New().From("users").Where("id = 1").Raw("CASCADE")
		expected := "DELETE FROM users WHERE id = 1 CASCADE"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Delete With Where, Returning, and Raw", func(t *testing.T) {
		qb := New().From("users").Where("id = 1").Returning("id").Raw("CASCADE")
		expected := "DELETE FROM users WHERE id = 1 RETURNING id CASCADE"
		assert.Equal(t, expected, qb.Build())
	})
}

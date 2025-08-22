package querybuilders_test

import (
	"homeserver/db/query_builder/qbupdate"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateQueryBuilder_Build(t *testing.T) {
	t.Run("Simple Update", func(t *testing.T) {
		qb := qbupdate.New().
			From("users").
			Set("name").
			Where("id = $2")
		expected := "UPDATE users SET name = $1 WHERE id = $2"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Multiple Set Clauses", func(t *testing.T) {
		qb := qbupdate.New().
			From("users").
			Set("name").
			Set("email").
			Where("id = $3")
		expected := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("With Returning", func(t *testing.T) {
		qb := qbupdate.New().
			From("users").
			Set("name").
			Where("id = $2").
			Returning("id").
			Returning("name")
		expected := "UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("With Raw Clause", func(t *testing.T) {
		qb := qbupdate.New().
			From("users").
			Set("name").
			Raw("FOR UPDATE")
		expected := "UPDATE users SET name = $1 FOR UPDATE"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Panic on No Table", func(t *testing.T) {
		qb := qbupdate.New().
			Set("name")
		assert.PanicsWithValue(t, "No table specified for UPDATE", func() { qb.Build() })
	})

	t.Run("Panic on No Set Clauses", func(t *testing.T) {
		qb := qbupdate.New().
			From("users")
		assert.PanicsWithValue(t, "No SET columns specified for UPDATE", func() { qb.Build() })
	})
}

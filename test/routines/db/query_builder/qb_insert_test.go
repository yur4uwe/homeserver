package querybuilders_test

import (
	"homeserver/db/query_builder/qbinsert"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertQueryBuilder_Build(t *testing.T) {
	t.Run("Simple Insert", func(t *testing.T) {
		qb := qbinsert.New().
			Into("users").
			Columns("name", "email").
			ValuesNum(1)
		expected := "INSERT INTO users (name, email) VALUES ($1, $2)"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Multiple Rows Insert", func(t *testing.T) {
		qb := qbinsert.New().
			Into("users").
			Columns("name", "email").
			ValuesNum(2)
		expected := "INSERT INTO users (name, email) VALUES ($1, $2), ($3, $4)"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("With Returning", func(t *testing.T) {
		qb := qbinsert.New().
			Into("users").
			Columns("name", "email").
			ValuesNum(1).
			Returning("id").
			Returning("name")

		expected := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("With Raw Clause", func(t *testing.T) {
		qb := qbinsert.New().
			Into("users").
			Columns("name").
			ValuesNum(1).
			Raw("ON CONFLICT DO NOTHING")
		expected := "INSERT INTO users (name) VALUES ($1) ON CONFLICT DO NOTHING"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Panic on No Table", func(t *testing.T) {
		qb := qbinsert.New().
			Columns("name").
			ValuesNum(1)
		assert.PanicsWithValue(t, "No table specified for INSERT", func() { qb.Build() })
	})

	t.Run("Panic on No Columns", func(t *testing.T) {
		qb := qbinsert.New().
			Into("users").
			ValuesNum(1)
		assert.PanicsWithValue(t, "No columns specified for INSERT", func() { qb.Build() })
	})
}

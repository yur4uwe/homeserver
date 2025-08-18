package querybuilder

import (
	"homeserver/constants/qbconsts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("Insert Query", func(t *testing.T) {
		query := New().Insert("name", "age").Into("users").Values("John", 30).Build()
		assert.Equal(t, "INSERT INTO users (name, age) VALUES ($1, $2)", query)
	})

	t.Run("Insert Query - No Columns", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Insert().Into("users").Values().Build()
		})
	})

	t.Run("Select Query", func(t *testing.T) {
		query := New().Select("id", "name").From("users").Where("age > 18").Build()
		assert.Equal(t, "SELECT id, name FROM users WHERE age > 18", query)
	})

	t.Run("Select Query - No Columns", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Select().From("users").Build()
		})
	})

	t.Run("Select Query - No Table", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Select("id").Build()
		})
	})

	t.Run("Update Query", func(t *testing.T) {
		query := New().Update("users").Set("name", "age").Where("id = $1").Build()
		assert.Equal(t, "UPDATE users SET name = $1, age = $2 WHERE id = $1", query)
	})

	t.Run("Update Query - No Columns", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Update("users").Set().Build()
		})
	})

	t.Run("Update Query - No Table", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Set("name").Build()
		})
	})

	t.Run("Delete Query", func(t *testing.T) {
		query := New().Delete().From("users").Where("id = $1").Build()
		assert.Equal(t, "DELETE FROM users WHERE id = $1", query)
	})

	t.Run("Delete Query - No Table", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Delete().Build()
		})
	})

	t.Run("Delete Query - Not In Delete Mode", func(t *testing.T) {
		qb := New()
		assert.Panics(t, func() {
			qb.buildDelete()
		})
	})

	t.Run("Update Query - Not In Update Mode", func(t *testing.T) {
		qb := New()
		assert.Panics(t, func() {
			qb.buildUpdate()
		})
	})

	t.Run("Insert Query - Not In Insert Mode", func(t *testing.T) {
		qb := New()
		assert.Panics(t, func() {
			qb.buildInsert()
		})
	})

	t.Run("Select Query - Not In Select Mode", func(t *testing.T) {
		qb := New()
		assert.Panics(t, func() {
			qb.buildSelect()
		})
	})

	t.Run("Chained Where and OrderBy", func(t *testing.T) {
		qb := New().Select("id").From("users").Where("active = true").OrderBy(qbconsts.ASCENDING)
		// OrderBy is a no-op in buildSelect, but should not panic
		assert.NotPanics(t, func() { qb.Build() })
	})

	t.Run("Chained Limit and Offset", func(t *testing.T) {
		qb := New().Select("id").From("users").Limit(10).Offset(5)
		// Limit/Offset are no-ops in buildSelect, but should not panic
		assert.NotPanics(t, func() { qb.Build() })
	})

	t.Run("Returning Clause", func(t *testing.T) {
		qb := New().Insert("name").Into("users").Values("John").Returning("id")
		// Returning is a no-op in buildInsert, but should not panic
		assert.NotPanics(t, func() { qb.Build() })
	})
}

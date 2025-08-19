package querybuilder

import (
	"homeserver/constants/qbconsts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("Insert Query", func(t *testing.T) {
		query := New().Insert("username", "email").Into("users").Values("John", 30).Build()
		assert.Equal(t, "INSERT INTO users (username, email) VALUES ($1, $2)", query)
	})

	t.Run("Insert Query - No Columns", func(t *testing.T) {
		assert.Panics(t, func() {
			New().Insert().Into("users").Values().Build()
		})
	})

	t.Run("Select Query", func(t *testing.T) {
		query := New().Select("id", "username").From("users").Where("email = 'test@example.com'").Build()
		assert.Equal(t, "SELECT id, username FROM users WHERE email = 'test@example.com'", query)
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
		query := New().Update("users").Set("username", "email").Where("id = $1").Build()
		assert.Equal(t, "UPDATE users SET username = $1, email = $2 WHERE id = $1", query)
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
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "SELECT id FROM users WHERE active = true ORDER BY ASC", qb.Build())
	})

	t.Run("Chained Limit and Offset", func(t *testing.T) {
		qb := New().Select("id").From("users").Limit(10).Offset(5)
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "SELECT id FROM users LIMIT 10 OFFSET 5", qb.Build())
	})

	t.Run("Returning Clause", func(t *testing.T) {
		qb := New().Insert("username").Into("users").Values("John").Returning("id")
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "INSERT INTO users (username) VALUES ($1) RETURNING id", qb.Build())
	})

	t.Run("Returning Clause with Multiple Columns", func(t *testing.T) {
		qb := New().Insert("username", "email").Into("users").Values("John", "john@example.com").Returning("id", "created_at")
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, created_at", qb.Build())
	})

	t.Run("Returning Clause on Select", func(t *testing.T) {
		qb := New().Select("id").From("users")
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "SELECT id FROM users", qb.Build())
	})

	t.Run("Returning Clause on Update", func(t *testing.T) {
		qb := New().Update("users").Set("username").Where("id = $1").Returning("id", "username")
		assert.NotPanics(t, func() { qb.Build() })
		assert.Equal(t, "UPDATE users SET username = $1 WHERE id = $1 RETURNING id, username", qb.Build())
	})

	t.Run("Panic on Unknown Query Type", func(t *testing.T) {
		qb := New()
		qb.queryType = "UNKNOWN_TYPE"
		qb.table = "users"
		qb.columns = []string{"id"}
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Unknown Table", func(t *testing.T) {
		qb := New().Select("id").From("not_a_table")
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Empty Query Type", func(t *testing.T) {
		qb := New()
		qb.table = "users"
		qb.columns = []string{"id"}
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Empty Table", func(t *testing.T) {
		qb := New()
		qb.queryType = qbconsts.SELECT
		qb.columns = []string{"id"}
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Too Many Columns", func(t *testing.T) {
		qb := New().Select("id", "username", "password", "email", "auth_token", "created_at", "email_verified", "extra_column").From("users")
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Too Many Returning Columns", func(t *testing.T) {
		qb := New().Select("id").From("users").Returning("id", "username", "password", "email", "auth_token", "created_at", "email_verified", "extra_column")
		assert.Panics(t, func() { qb.Build() })
	})

	t.Run("Panic on Unknown Column", func(t *testing.T) {
		qb := New().Select("id", "not_a_field").From("users")
		assert.PanicsWithValue(t, "unknown column: not_a_field", func() { qb.Build() })
	})

	t.Run("Panic on Unknown Returning Column", func(t *testing.T) {
		qb := New().Select("id").From("users").Returning("not_a_field")
		assert.PanicsWithValue(t, "unknown returning column: not_a_field", func() { qb.Build() })
	})

	t.Run("Panic on Multiple Unknown Columns", func(t *testing.T) {
		qb := New().Select("id", "foo", "bar").From("users")
		assert.PanicsWithValue(t, "unknown column: foo", func() { qb.Build() })
	})

	t.Run("Panic on Multiple Unknown Returning Columns", func(t *testing.T) {
		qb := New().Select("id").From("users").Returning("foo", "bar")
		assert.PanicsWithValue(t, "unknown returning column: foo", func() { qb.Build() })
	})

	t.Run("No Panic on All Valid Columns and Returning", func(t *testing.T) {
		qb := New().Select("id", "username").From("users").Returning("id")
		assert.NotPanics(t, func() { qb.Build() })
	})
}

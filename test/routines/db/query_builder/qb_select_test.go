package querybuilders_test

import (
	"homeserver/constants/qbconsts"
	"homeserver/db/query_builder/qbselect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectQueryBuilder(t *testing.T) {
	t.Run("Simple Select", func(t *testing.T) {
		qb := qbselect.New().
			Col("id").
			Col("name").
			From("users").
			Where("age > 30")

		query := qb.Build()
		expected := "SELECT id, name FROM users WHERE age > 30"
		assert.Equal(t, expected, query)
	})

	t.Run("Joined Select", func(t *testing.T) {
		qb := qbselect.New().
			Col("users", "id").
			Col("users", "name").
			From("users").
			Join(qbconsts.InnerJoin, "orders", "users.id = orders.user_id").
			Col("orders", "order_id").
			Where("users.age > 30")

		query := qb.Build()
		expected := "SELECT users.id, users.name, orders.order_id FROM users INNER JOIN orders ON users.id = orders.user_id WHERE users.age > 30"
		assert.Equal(t, expected, query)
	})

	t.Run("Many Joins", func(t *testing.T) {
		qb := qbselect.New().
			Col("users", "id").
			Col("users", "name").
			From("users").
			Join(qbconsts.RightJoin, "orders", "users.id = orders.user_id").
			Col("orders", "order_id").
			Join(qbconsts.LeftJoin, "products", "orders.product_id = products.id").
			Col("products", "name").
			Where("users.age > 30")

		query := qb.Build()
		expected := "SELECT users.id, users.name, orders.order_id, products.name FROM users RIGHT JOIN orders ON users.id = orders.user_id LEFT JOIN products ON orders.product_id = products.id WHERE users.age > 30"
		assert.Equal(t, expected, query)
	})
	t.Run("OrderBy and Limit", func(t *testing.T) {
		qb := qbselect.New().
			Col("id").
			Col("name").
			From("users").
			OrderBy("name DESC").
			Limit(10)

		query := qb.Build()
		expected := "SELECT id, name FROM users ORDER BY name DESC LIMIT 10"
		assert.Equal(t, expected, query)
	})

	t.Run("Offset and GroupBy", func(t *testing.T) {
		qb := qbselect.New().
			Col("department").
			Col("", "COUNT(*)", "cnt").
			From("employees").
			GroupBy("department").
			Offset(5)

		query := qb.Build()
		expected := "SELECT department, COUNT(*) AS cnt FROM employees GROUP BY department OFFSET 5"
		assert.Equal(t, expected, query)
	})

	t.Run("Having Clause", func(t *testing.T) {
		qb := qbselect.New().
			Col("department").
			Col("", "COUNT(*)", "cnt").
			From("employees").
			GroupBy("department").
			Having("COUNT(*) > 3")

		query := qb.Build()
		expected := "SELECT department, COUNT(*) AS cnt FROM employees GROUP BY department HAVING COUNT(*) > 3"
		assert.Equal(t, expected, query)
	})

	t.Run("Distinct Select", func(t *testing.T) {
		qb := qbselect.New().
			Distinct().
			Col("email").
			From("users")

		expected := "SELECT DISTINCT email FROM users"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Raw SQL Clause (appended)", func(t *testing.T) {
		qb := qbselect.New().
			Col("id").
			From("users").
			Raw("FOR UPDATE")

		expected := "SELECT id FROM users FOR UPDATE"
		assert.Equal(t, expected, qb.Build())
	})

	t.Run("Columns", func(t *testing.T) {
		qb := qbselect.New().
			Columns("id", "name").
			From("users")

		query := qb.Build()
		expected := "SELECT id, name FROM users"
		assert.Equal(t, expected, query)
	})

	t.Run("Panic on Col with 0 args", func(t *testing.T) {
		qb := qbselect.New().From("users")
		assert.Panics(t, func() { qb.Col() })
	})

	t.Run("No Columns Specified", func(t *testing.T) {
		qb := qbselect.New().From("users")
		query := qb.Build()
		expected := "SELECT * FROM users"
		assert.Equal(t, expected, query)
	})

	t.Run("Returning", func(t *testing.T) {
		qb := qbselect.New().
			Col("id").
			Col("name").
			From("users").
			Where("age > 30").
			Ret("id").
			Ret("name")

		query := qb.Build()
		expected := "SELECT id, name FROM users WHERE age > 30 RETURNING id, name"
		assert.Equal(t, expected, query)
	})
}

package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("Insert Query", func(t *testing.T) {
		query := New().Insert("name", "age").Into("users").Values("John", 30).Build()
		assert.Equal(t, "INSERT INTO users (name, age) VALUES ($1, $2)", query)
	})
}

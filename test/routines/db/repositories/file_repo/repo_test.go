package filerepo_test

import (
	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileRepo(t *testing.T) {
	database := &db.Database{}
	repo := filerepo.NewFileRepo(database)

	assert.NotNil(t, repo, "Expected non-nil FileRepo instance")
}

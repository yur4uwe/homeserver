package tablewrappers

import (
	"context"
	"homeserver/db"
	"homeserver/db/wrappers/table_wrappers/tables"
)

type TableRepository struct {
	*db.Database
}

func (repo *TableRepository) ClearEntries(ctx context.Context) (string, error) {
	return tables.ClearEntries(ctx, repo.Database)
}

func (repo *TableRepository) DropTable(ctx context.Context, tableName string) error {
	return tables.DropTable(ctx, repo.Database, tableName)
}

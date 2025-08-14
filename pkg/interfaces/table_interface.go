package interfaces

import (
	"context"
	"homeserver/pkg/globals"
)

type TableInterface interface {
	CreateTable(ctx context.Context, table *globals.DBTable) error
	TableExists(ctx context.Context, tableName string) (bool, error)
	ClearEntries(ctx context.Context) (string, error)
	DropTable(ctx context.Context, tableName string) error
}

package setup

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
)

var required_db_tables = []*globals.DBTable{
	globals.DBUsersTable,
	globals.DBFilesTable,
}

func CheckDatabaseTables(connection *db.Database) error {
	ctx := context.Background()

	for _, table := range required_db_tables {
		exists, err := connection.TableExists(ctx, table.Name)
		if err != nil {
			return fmt.Errorf("error checking table %s: %v", table.Name, err)
		}

		if !exists {
			connection.CreateTable(ctx, table)
		}
	}

	fmt.Println("All required database tables exist.")
	return nil
}

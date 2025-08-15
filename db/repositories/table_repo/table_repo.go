package tablerepo

import (
	"context"
	"fmt"
	"homeserver/db"
)

type TableRepository struct {
	db *db.Database
}

func (repo *TableRepository) DropTable(ctx context.Context, table_name string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE;`, table_name)
	_, err := repo.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error dropping table %s: %v", table_name, err)
	}

	return nil
}

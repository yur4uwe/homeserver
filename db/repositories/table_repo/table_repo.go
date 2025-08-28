package tablerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InTableRepo interface {
	DropTable(ctx context.Context, table_name string) error
	TableExists(ctx context.Context, table_name string) (bool, error)
	CreateTable(ctx context.Context, table *dbconsts.DBTable) error
}

type TableRepo struct {
	pool *pgxpool.Pool
}

func NewTableRepo(pool *pgxpool.Pool) *TableRepo {
	return &TableRepo{
		pool: pool,
	}
}

func (repo *TableRepo) DropTable(ctx context.Context, table_name string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE;`, table_name)
	_, err := repo.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error dropping table %s: %v", table_name, err)
	}

	return nil
}

func (repo *TableRepo) TableExists(ctx context.Context, table_name string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s');`, table_name)
	var exists bool
	err := repo.pool.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if table %s exists: %v", table_name, err)
	}
	return exists, nil
}

func (repo *TableRepo) CreateTable(ctx context.Context, table *dbconsts.DBTable) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s);`, table.Name, table.SQLCreation)
	_, err := repo.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", table.Name, err)
	}
	return nil
}

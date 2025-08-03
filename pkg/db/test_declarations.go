package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestDatabase() (*Database, error) {
	testDBUrl := ""

	pool, err := pgxpool.New(context.Background(), testDBUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return &Database{pool: pool}, nil
}

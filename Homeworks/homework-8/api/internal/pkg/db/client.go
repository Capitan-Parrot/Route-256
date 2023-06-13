package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/homework8/config"
)

// NewDB билдер для Database
func NewDB(ctx context.Context) (*Database, error) {
	dsn := generateDsn()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return newDatabase(pool), nil
}

func generateDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBhost, config.DBport, config.DBuser, config.DBpassword, config.DBname)
}

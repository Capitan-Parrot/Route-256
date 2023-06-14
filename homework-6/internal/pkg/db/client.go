package db

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework6/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
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

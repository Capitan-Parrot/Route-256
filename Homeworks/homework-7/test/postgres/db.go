package postgres

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework7/internal/pkg/db"
	"strings"
	"sync"
	"testing"
)

type TDB struct {
	sync.Mutex
	DB *db.Database
}

func NewFromEnv() *TDB {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: database}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Lock()
	d.Truncate(ctx)
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Truncate(ctx)
}

func (d *TDB) Truncate(ctx context.Context) {
	var tables []string
	err := d.DB.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}
	if len(tables) == 0 {
		panic("no tables")
	}
	query := fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := d.DB.Exec(ctx, query); err != nil {
		panic(err)
	}
}

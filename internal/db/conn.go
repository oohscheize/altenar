package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPgxPool(ctx context.Context, user, password, dbname, host string, port int) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
	return pgxpool.New(ctx, dsn)
}

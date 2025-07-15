//go:build integration
// +build integration

package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
)

func StartPG(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("docker: %v", err)
	}

	res, err := pool.Run("postgres", "16",
		[]string{
			"POSTGRES_DB=casino",
			"POSTGRES_USER=casino",
			"POSTGRES_PASSWORD=casino",
		})
	if err != nil {
		tb.Fatalf("run pg: %v", err)
	}
	tb.Cleanup(func() { _ = pool.Purge(res) })

	dsn := fmt.Sprintf(
		"postgres://casino:casino@localhost:%s/casino?sslmode=disable",
		res.GetPort("5432/tcp"),
	)

	var pg *pgxpool.Pool
	pool.MaxWait = 200 * time.Second

	err = pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tmp, e := pgxpool.New(ctx, dsn)
		if e != nil {
			return e
		}
		for i := 0; i < 3; i++ {
			if e = tmp.Ping(ctx); e != nil {
				tmp.Close()
				return e
			}
			var x int
			if e = tmp.QueryRow(ctx, "SELECT 1").Scan(&x); e != nil {
				tmp.Close()
				return e
			}
			time.Sleep(200 * time.Millisecond)
		}
		pg = tmp
		return nil
	})
	if err != nil {
		tb.Fatalf("connect pg: %v", err)
	}

	_, _ = pg.Exec(context.Background(), `
	  CREATE TABLE IF NOT EXISTS transactions (
	    id SERIAL PRIMARY KEY,
	    user_id BIGINT NOT NULL,
	    transaction_type VARCHAR(10),
	    amount NUMERIC(12,2),
	    created_at TIMESTAMP NOT NULL DEFAULT NOW()
	  );`)

	return pg
}

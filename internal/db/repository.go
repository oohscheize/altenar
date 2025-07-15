package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Repositoryer interface {
	GetTransactions(ctx context.Context, userID *int64, txType *string) ([]Transaction, error)
	InsertTransaction(ctx context.Context, tx *Transaction) error
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertTransaction(ctx context.Context, t *Transaction) error {
	query := `
        INSERT INTO transactions (user_id, transaction_type, amount)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	return r.db.QueryRow(ctx, query, t.UserID, t.TransactionType, t.Amount).
		Scan(&t.ID, &t.CreatedAt)
}

func (r *Repository) GetTransactions(ctx context.Context, userID *int64, tType *string) ([]Transaction, error) {
	base := `SELECT id, user_id, transaction_type, amount, created_at FROM transactions`
	where := ""
	args := []interface{}{}
	argn := 1

	if userID != nil {
		where += fmt.Sprintf("user_id = $%d", argn)
		args = append(args, *userID)
		argn++
	}
	if tType != nil {
		if len(where) > 0 {
			where += " AND "
		}
		where += fmt.Sprintf("transaction_type = $%d", argn)
		args = append(args, *tType)
		argn++
	}
	if where != "" {
		base += " WHERE " + where
	}
	base += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.TransactionType, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

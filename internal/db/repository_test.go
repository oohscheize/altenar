package db

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
)

func newRepoMock(t *testing.T) (*Repository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("mock: %v", err)
	}
	return NewRepository(mock), mock
}

func TestRepository_InsertAndGet_NoFilters(t *testing.T) {
	repo, mock := newRepoMock(t)
	defer mock.Close()

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO transactions (user_id, transaction_type, amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`)).
		WithArgs(int64(7), "bet", 11.0).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).
			AddRow(int64(1), now))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, transaction_type, amount, created_at
		FROM transactions
		ORDER BY created_at DESC`)).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "user_id", "transaction_type", "amount", "created_at"}).
			AddRow(int64(1), int64(7), "bet", 11.0, now))

	tx := &Transaction{UserID: 7, TransactionType: "bet", Amount: 11}
	_ = repo.InsertTransaction(context.Background(), tx)

	list, _ := repo.GetTransactions(context.Background(), nil, nil)
	if len(list) != 1 || list[0].UserID != 7 {
		t.Fatalf("want 1 row with user 7, got %+v", list)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_GetTransactions_WithFilters(t *testing.T) {
	repo, mock := newRepoMock(t)
	defer mock.Close()

	uid := int64(42)
	typ := "win"
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, transaction_type, amount, created_at
		FROM transactions
		WHERE user_id = $1 AND transaction_type = $2
		ORDER BY created_at DESC`)).
		WithArgs(uid, typ).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "user_id", "transaction_type", "amount", "created_at"}).
			AddRow(int64(5), uid, typ, 99.9, now))

	list, err := repo.GetTransactions(context.Background(), &uid, &typ)
	if err != nil || len(list) != 1 || list[0].TransactionType != typ {
		t.Fatalf("filter err=%v list=%+v", err, list)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

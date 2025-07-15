//go:build integration
// +build integration

package db

import (
	"context"
	"testing"
	"time"
)

func TestRepository_InsertAndGet_PG(t *testing.T) {
	pg := StartPG(t)
	repo := NewRepository(pg)

	tx := &Transaction{UserID: 77, TransactionType: "bet", Amount: 50, CreatedAt: time.Now()}
	if err := repo.InsertTransaction(context.Background(), tx); err != nil {
		t.Fatalf("insert: %v", err)
	}

	res, err := repo.GetTransactions(context.Background(), nil, nil)
	if err != nil || len(res) != 1 || res[0].UserID != 77 {
		t.Fatalf("unexpected result: %+v / err %v", res, err)
	}
}

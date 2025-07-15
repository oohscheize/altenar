package db

import (
	"testing"
	"time"
)

func TestTransaction_IsValid(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name string
		tx   Transaction
		want bool
	}{
		{"valid bet", Transaction{UserID: 1, TransactionType: "bet", Amount: 10, CreatedAt: now}, true},
		{"valid win", Transaction{UserID: 2, TransactionType: "win", Amount: 50, CreatedAt: now}, true},
		{"zero user", Transaction{UserID: 0, TransactionType: "bet", Amount: 10, CreatedAt: now}, false},
		{"empty type", Transaction{UserID: 1, TransactionType: "", Amount: 10, CreatedAt: now}, false},
		{"invalid type", Transaction{UserID: 1, TransactionType: "foo", Amount: 10, CreatedAt: now}, false},
		{"zero amount", Transaction{UserID: 1, TransactionType: "bet", Amount: 0, CreatedAt: now}, false},
		{"negative amount", Transaction{UserID: 1, TransactionType: "bet", Amount: -5, CreatedAt: now}, false},
		{"zero time", Transaction{UserID: 1, TransactionType: "bet", Amount: 10, CreatedAt: time.Time{}}, false},
	}
	for _, c := range cases {
		if got := c.tx.IsValid(); got != c.want {
			t.Errorf("%s: expected %v, got %v", c.name, c.want, got)
		}
	}
}

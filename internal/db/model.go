package db

import "time"

type Transaction struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	TransactionType string    `db:"transaction_type" json:"transaction_type"`
	Amount          float64   `db:"amount" json:"amount"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func (t *Transaction) IsValid() bool {
	if t.UserID == 0 {
		return false
	}
	if t.TransactionType != "bet" && t.TransactionType != "win" {
		return false
	}
	if t.Amount <= 0 {
		return false
	}
	if t.CreatedAt.IsZero() {
		return false
	}
	return true
}

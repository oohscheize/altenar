package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/oohscheize/altenar/internal/db"
)

type mockRepo struct{ inserted []*db.Transaction }

func (m *mockRepo) InsertTransaction(_ context.Context, tx *db.Transaction) error {
	m.inserted = append(m.inserted, tx)
	return nil
}
func (m *mockRepo) GetTransactions(context.Context, *int64, *string) ([]db.Transaction, error) {
	return nil, nil
}

func TestConsumer_HandleValidAndInvalid(t *testing.T) {
	type tc struct {
		name    string
		msg     TransactionMessage
		wantIns bool
	}
	cases := []tc{
		{
			name:    "valid",
			msg:     TransactionMessage{UserID: 10, TransactionType: "bet", Amount: 9.9, CreatedAt: time.Now()},
			wantIns: true,
		},
		{
			name:    "invalid user_id",
			msg:     TransactionMessage{UserID: 0, TransactionType: "bet", Amount: 9.9, CreatedAt: time.Now()},
			wantIns: false,
		},
	}

	repo := &mockRepo{}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			b, _ := json.Marshal(cse.msg)
			var m TransactionMessage
			_ = json.Unmarshal(b, &m)

			tx := &db.Transaction{
				UserID:          m.UserID,
				TransactionType: m.TransactionType,
				Amount:          m.Amount,
				CreatedAt:       m.CreatedAt,
			}
			if tx.IsValid() {
				_ = repo.InsertTransaction(context.Background(), tx)
			}

			got := len(repo.inserted)
			if cse.wantIns && got != 1 {
				t.Fatalf("expected insert, repo len=%d", got)
			}
			if !cse.wantIns && got != 0 {
				t.Fatalf("expected skip, but repo len=%d", got)
			}
			repo.inserted = nil
		})
	}
}

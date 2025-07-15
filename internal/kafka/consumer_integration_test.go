//go:build integration
// +build integration

package kafka

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/oohscheize/altenar/internal/db"
	"github.com/segmentio/kafka-go"
)

type stubReader struct{ msgs [][]byte }

func (s *stubReader) ReadMessage(context.Context) (kafka.Message, error) {
	if len(s.msgs) == 0 {
		return kafka.Message{}, io.EOF
	}
	val := s.msgs[0]
	s.msgs = s.msgs[1:]
	return kafka.Message{Value: val}, nil
}

func TestConsumer_Run_PersistsValid(t *testing.T) {
	pg := db.StartPG(t)
	repo := db.NewRepository(pg)

	valid, _ := json.Marshal(TransactionMessage{UserID: 42, TransactionType: "win", Amount: 99.9, CreatedAt: time.Now()})
	invalid, _ := json.Marshal(TransactionMessage{UserID: 0, TransactionType: "bet", Amount: 1, CreatedAt: time.Now()})

	c := &Consumer{
		Reader:     &stubReader{msgs: [][]byte{valid, invalid}},
		Repository: repo,
	}
	c.Run(context.Background())

	list, _ := repo.GetTransactions(context.Background(), nil, nil)
	if len(list) != 1 || list[0].UserID != 42 {
		t.Fatalf("want 1 valid tx, got %+v", list)
	}
}

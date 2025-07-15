package kafka

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/oohscheize/altenar/internal/db"
	"github.com/segmentio/kafka-go"
)

type TransactionMessage struct {
	UserID          int64     `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}

type MessageReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}

type Consumer struct {
	Reader     MessageReader
	Repository db.Repositoryer
}

func NewConsumer(brokers []string, topic string, groupID string, repo *db.Repository) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{
		Reader:     reader,
		Repository: repo,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("consumer stopped:", ctx.Err())
			return
		default:
		}

		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				log.Println("kafka reader reached EOF, exiting")
				return
			}
			log.Printf("kafka read error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		var msg TransactionMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("invalid message: %v", err)
			continue
		}
		t := &db.Transaction{
			UserID:          msg.UserID,
			TransactionType: msg.TransactionType,
			Amount:          msg.Amount,
			CreatedAt:       msg.CreatedAt,
		}

		if !t.IsValid() {
			log.Printf("invalid transaction, skipping: %+v", t)
			continue
		}

		if err := c.Repository.InsertTransaction(ctx, t); err != nil {
			log.Printf("db insert error: %v", err)
		} else {
			log.Printf("transaction saved: %+v", t)
		}
	}
}

package main

import (
	"context"
	"log"

	"github.com/oohscheize/altenar/internal/config"
	"github.com/oohscheize/altenar/internal/db"
	"github.com/oohscheize/altenar/internal/kafka"
)

func main() {
	ctx := context.Background()
	dbpool, err := db.OpenPgxPool(
		ctx,
		config.PGUser,
		config.PGPassword,
		config.PGDB,
		config.PGHost,
		config.PGPort,
	)
	if err != nil {
		log.Fatalf("can't connect db: %v", err)
	}
	defer dbpool.Close()

	repo := db.NewRepository(dbpool)

	consumer := kafka.NewConsumer(
		config.KafkaBrokers,
		config.KafkaTopic,
		config.KafkaGroupID,
		repo,
	)
	log.Println("Kafka consumer started")
	consumer.Run(ctx)
}

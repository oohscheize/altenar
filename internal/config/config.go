package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PGHost     string
	PGPort     int
	PGDB       string
	PGUser     string
	PGPassword string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	APIPort string
)

func init() {
	_ = godotenv.Load()

	PGHost = getenv("PG_HOST", "localhost")
	PGPort = getenvInt("PG_PORT", 5432)
	PGDB = getenv("PG_DB", "casino")
	PGUser = getenv("PG_USER", "casino")
	PGPassword = getenv("PG_PASSWORD", "casino")
	
	KafkaBrokers = []string{getenv("KAFKA_BROKER", "localhost:9092")}
	KafkaTopic = getenv("KAFKA_TOPIC", "transactions")
	KafkaGroupID = getenv("KAFKA_GROUP_ID", "transaction-consumer")

	APIPort = getenv("API_PORT", "8080")
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getenvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("warning: env %s: %v, fallback to %d", key, err, fallback)
		return fallback
	}
	return i
}

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oohscheize/altenar/internal/api"
	"github.com/oohscheize/altenar/internal/config"
	"github.com/oohscheize/altenar/internal/db"
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
	apiHandler := api.NewAPI(repo)

	r := gin.Default()
	r.GET("/transactions", apiHandler.GetTransactions)

	log.Printf("API server started on :%s", config.APIPort)
	if err := r.Run(":" + config.APIPort); err != nil {
		log.Fatal(err)
	}
}

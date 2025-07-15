//go:build integration
// +build integration

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oohscheize/altenar/internal/db"
)

func TestGetTransactions_Integration(t *testing.T) {
	pg := db.StartPG(t)
	repo := db.NewRepository(pg)

	_ = repo.InsertTransaction(context.Background(),
		&db.Transaction{UserID: 55, TransactionType: "bet", Amount: 12, CreatedAt: time.Now()})

	gin.SetMode(gin.TestMode)
	h := NewAPI(repo)
	r := gin.New()
	r.GET("/transactions", h.GetTransactions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions?user_id=55", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var out []db.Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil || len(out) != 1 {
		t.Fatalf("body %s", w.Body.String())
	}
}

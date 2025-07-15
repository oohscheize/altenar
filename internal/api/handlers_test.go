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

type stubRepo struct{ list []db.Transaction }

func (s *stubRepo) GetTransactions(context.Context, *int64, *string) ([]db.Transaction, error) {
	return s.list, nil
}
func (s *stubRepo) InsertTransaction(context.Context, *db.Transaction) error { return nil }

func TestGetTransactions_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tr := db.Transaction{ID: 1, UserID: 42, TransactionType: "win", Amount: 77, CreatedAt: time.Now()}
	repo := &stubRepo{list: []db.Transaction{tr}}
	h := NewAPI(repo)

	r := gin.New()
	r.GET("/transactions", h.GetTransactions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions?user_id=42", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp []db.Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || len(resp) != 1 || resp[0].UserID != 42 {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

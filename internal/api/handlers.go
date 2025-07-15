package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oohscheize/altenar/internal/db"
)

type API struct {
	Repo db.Repositoryer
}

func NewAPI(repo db.Repositoryer) *API {
	return &API{Repo: repo}
}

func (a *API) GetTransactions(c *gin.Context) {
	var userID *int64
	if idStr := c.Query("user_id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = &id
	}
	var tType *string
	if tStr := c.Query("type"); tStr != "" {
		tType = &tStr
	}

	txs, err := a.Repo.GetTransactions(c.Request.Context(), userID, tType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, txs)
}

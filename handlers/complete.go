package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"product-catalog/go-app/controllers"
)

func CompleteHandler(c *gin.Context) {
	query := c.Query("q")
	results := controllers.CompleteQuery(query)
	if results == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query"})
		return
	}
	c.JSON(http.StatusOK, results)
}

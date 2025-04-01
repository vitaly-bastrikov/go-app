package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"product-catalog/go-app/controllers"
)

func PrecomputeHandler(c *gin.Context) {
	err := controllers.PrecomputeEmbeddings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to precompute embeddings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Embeddings precomputed successfully"})
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"product-catalog/go-app/controllers"
	"product-catalog/go-app/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("📦 Loading data...")
	controllers.LoadCatalogItems("data/nava_items.json")
	controllers.LoadCatalogItems("data/preferences.json")
	controllers.LoadCatalogItems("data/products.json")
	log.Println("✅ Data loaded")

	embeddingURL := os.Getenv("EMBEDDING_URL")
	log.Printf("🔗 EMBEDDING_URL = %s", embeddingURL)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/complete", handlers.CompleteHandler)
	router.POST("/precompute", handlers.PrecomputeHandler)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "🟢 Go app is running!")
	})

	log.Printf("🚀 Starting on port %s...", port)
	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

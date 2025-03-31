package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CatalogItem struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	InternalNotes string    `json:"internal_notes"`
	Tags          []string  `json:"tags"`
	Embedding     []float64 `json:"embedding"`
}

var navaItems []CatalogItem
var preferences []CatalogItem
var products []CatalogItem

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/complete", CompleteHandler)
	router.POST("/precompute", PrecomputeHandler)

	// Load from file
	navaItems = loadCatalogItems("nava_items.json")
	preferences = loadCatalogItems("preferences.json")
	products = loadCatalogItems("products.json")

	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}

func loadCatalogItems(filename string) []CatalogItem {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("❌ Failed to read %s: %v", filename, err)
	}
	var items []CatalogItem
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatalf("❌ Failed to parse %s: %v", filename, err)
	}
	log.Printf("✅ Loaded %d items from %s", len(items), filename)
	return items
}

func saveWithEmbeddings(filename string, items []CatalogItem) error {
	for i, item := range items {
		text := item.Name + " " + item.Description + " " + item.InternalNotes + " " + strings.Join(item.Tags, " ")
		embedding, err := GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to embed item %s: %v", item.ID, err)
			continue
		}
		items[i].Embedding = embedding
	}
	data, _ := json.MarshalIndent(items, "", "  ")
	return os.WriteFile(filename, data, 0644)
}

func PrecomputeHandler(c *gin.Context) {
	err1 := saveWithEmbeddings("nava_items.json", navaItems)
	err2 := saveWithEmbeddings("preferences.json", preferences)
	err3 := saveWithEmbeddings("products.json", products)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "❌ Failed to precompute some embeddings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Embeddings precomputed and saved"})
}

func GetEmbedding(text string) ([]float64, error) {
	type EmbedRequest struct {
		Text string `json:"text"`
	}
	type EmbedResponse struct {
		Embedding []float64 `json:"embedding"`
	}

	req := EmbedRequest{Text: text}
	body, _ := json.Marshal(req)

	resp, err := http.Post("http://localhost:8000/embed", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedding, nil
}

func cosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

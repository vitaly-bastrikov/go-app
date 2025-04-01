package controllers

import (
	"log"
	"strings"

	"product-catalog/go-app/gateway"
)

func PrecomputeEmbeddings() error {
	// Process nava items
	for i := range navaItems {
		text := navaItems[i].Name + " " + strings.Join(navaItems[i].Tags, " ")
		embedding, err := gateway.GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to embed nava item %s: %v", navaItems[i].ID, err)
			continue
		}
		navaItems[i].Embedding = embedding
	}

	// Process preferences
	for i := range preferences {
		text := preferences[i].Name + " " + strings.Join(preferences[i].Tags, " ")
		embedding, err := gateway.GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to embed preference %s: %v", preferences[i].ID, err)
			continue
		}
		preferences[i].Embedding = embedding
	}

	// Process products
	for i := range products {
		text := products[i].Name + " " + strings.Join(products[i].Tags, " ")
		embedding, err := gateway.GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to embed product %s: %v", products[i].ID, err)
			continue
		}
		products[i].Embedding = embedding
	}

	// Save all items with their new embeddings
	err1 := SaveWithEmbeddings("data/nava_items.json", navaItems)
	err2 := SaveWithEmbeddings("data/preferences.json", preferences)
	err3 := SaveWithEmbeddings("data/products.json", products)

	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("❌ Failed to save items with embeddings")
		return nil
	}

	log.Printf("✅ Successfully precomputed embeddings for %d nava items, %d preferences, and %d products",
		len(navaItems), len(preferences), len(products))
	return nil
}

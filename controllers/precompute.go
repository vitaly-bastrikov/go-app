package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"product-catalog/go-app/gateway"
)

func PrecomputeEmbeddings() error {
	log.Println("🔄 Starting precomputation...")

	// Process nava items
	for i := range navaItems {
		text := navaItems[i].Name
		if len(navaItems[i].Tags) > 0 {
			text += " " + navaItems[i].Tags[0]
		}

		embedding, err := gateway.GetEmbedding(text)
		if err != nil {
			log.Printf("❌ Failed to get embedding for nava item %s: %v", navaItems[i].Name, err)
			continue
		}
		navaItems[i].Embedding = embedding
		log.Printf("✅ Processed nava item: %s", navaItems[i].Name)
	}

	// Save nava items
	navaData, err := json.MarshalIndent(navaItems, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join("data", "nava_items.json"), navaData, 0644)
	if err != nil {
		return err
	}

	log.Printf("✅ Precomputation complete. Processed %d nava items", len(navaItems))
	return nil
}

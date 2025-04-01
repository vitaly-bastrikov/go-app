package controllers

import (
	"encoding/json"
	"log"
	"os"

	"product-catalog/go-app/entity"
)

var navaItems []entity.SearchItem

func LoadCatalogItems(filename string) []entity.SearchItem {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("❌ Failed to read %s: %v", filename, err)
		return nil
	}

	var items []entity.SearchItem
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("❌ Failed to parse %s: %v", filename, err)
		return nil
	}

	navaItems = items
	log.Printf("✅ Loaded %d items from %s", len(items), filename)
	return items
}

func SaveWithEmbeddings(filename string, items []entity.SearchItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

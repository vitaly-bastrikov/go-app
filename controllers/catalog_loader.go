package controllers

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"product-catalog/go-app/entity"
)

var (
	navaItems   []entity.CatalogItem
	preferences []entity.CatalogItem
	products    []entity.CatalogItem
)

func LoadCatalogItems(filename string) []entity.CatalogItem {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("❌ Failed to read %s: %v", filename, err)
		return nil
	}

	var items []entity.CatalogItem
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("❌ Failed to parse %s: %v", filename, err)
		return nil
	}

	// Store items in the appropriate global variable
	switch {
	case strings.Contains(filename, "nava"):
		navaItems = items
	case strings.Contains(filename, "preferences"):
		preferences = items
	case strings.Contains(filename, "products"):
		products = items
	}

	log.Printf("✅ Loaded %d items from %s", len(items), filename)
	return items
}

func SaveWithEmbeddings(filename string, items []entity.CatalogItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

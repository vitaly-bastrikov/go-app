package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"net/http"
	"sort"
	"strings"
)

func CompleteHandler(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query"})
		return
	}

	log.Printf("🔍 Query: %s", query)

	queryEmbedding, err := GetEmbedding(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to embed query"})
		return
	}

	type ScoredItem struct {
		Type  string  `json:"type"` // "nava", "preference", or "product"
		Name  string  `json:"name"`
		Score float64 `json:"score"`
		Boost bool    `json:"boost"`
	}

	var results []ScoredItem
	const threshold = 0.1
	const prefixBoost = 0.5
	const fuzzyBoost = 0.3
	queryLower := strings.ToLower(query)

	// 🔍 Nava Items
	for _, item := range navaItems {
		itemNameLower := strings.ToLower(item.Name)
		score := cosineSimilarity(queryEmbedding, item.Embedding)

		isPrefixMatch := false
		for _, word := range strings.Fields(itemNameLower) {
			if strings.HasPrefix(word, queryLower) {
				isPrefixMatch = true
				break
			}
		}
		isFuzzyMatch := fuzzy.MatchNormalizedFold(queryLower, itemNameLower)

		if isPrefixMatch {
			score += prefixBoost
		}
		if isFuzzyMatch {
			score += fuzzyBoost
		}

		if score >= threshold {
			results = append(results, ScoredItem{
				Type:  "nava",
				Name:  item.Name,
				Score: score,
				Boost: isPrefixMatch || isFuzzyMatch,
			})
		}
	}

	// 🔍 Preferences
	for _, pref := range preferences {
		nameLower := strings.ToLower(pref.Name)
		score := cosineSimilarity(queryEmbedding, pref.Embedding)

		isPrefixMatch := false
		for _, word := range strings.Fields(nameLower) {
			if strings.HasPrefix(word, queryLower) {
				isPrefixMatch = true
				break
			}
		}
		isFuzzyMatch := fuzzy.MatchNormalizedFold(queryLower, nameLower)

		if isPrefixMatch {
			score += prefixBoost
		}
		if isFuzzyMatch {
			score += fuzzyBoost
		}

		if score >= threshold {
			results = append(results, ScoredItem{
				Type:  "preference",
				Name:  pref.Name,
				Score: score,
				Boost: isPrefixMatch || isFuzzyMatch,
			})
		}
	}

	// 🔍 Products
	for _, p := range products {
		nameLower := strings.ToLower(p.Name)
		score := cosineSimilarity(queryEmbedding, p.Embedding)

		isPrefixMatch := false
		for _, word := range strings.Fields(nameLower) {
			if strings.HasPrefix(word, queryLower) {
				isPrefixMatch = true
				break
			}
		}
		isFuzzyMatch := fuzzy.MatchNormalizedFold(queryLower, nameLower)

		if isPrefixMatch {
			score += prefixBoost
		}
		if isFuzzyMatch {
			score += fuzzyBoost
		}

		if score >= threshold {
			results = append(results, ScoredItem{
				Type:  "product",
				Name:  p.Name,
				Score: score,
				Boost: isPrefixMatch || isFuzzyMatch,
			})
		}
	}

	// Sort results by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit to top 5
	if len(results) > 5 {
		results = results[:5]
	}

	c.JSON(http.StatusOK, results)
}

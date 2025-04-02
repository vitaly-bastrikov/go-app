package controllers

import (
	"log"
	"math"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"product-catalog/go-app/entity"
	"product-catalog/go-app/gateway"
)

func CompleteQuery(query string) []entity.ScoredItem {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	queryEmbedding, err := gateway.GetEmbedding(query)
	if err != nil {
		log.Printf("❌ Embedding failed: %v", err)
		return nil
	}

	return matchAllInParallel(queryEmbedding, strings.ToLower(query), navaItems)
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func matchAllInParallel(queryEmbedding []float64, queryLower string, allItems []entity.SearchItem) []entity.ScoredItem {
	const threshold = 0.1
	const namePrefixBoost = 0.5
	const nameFuzzyBoost = 0.3
	const tagPrefixBoost = 0.2
	const tagFuzzyBoost = 0.1

	numWorkers := runtime.NumCPU()
	itemCh := make(chan entity.SearchItem, len(allItems))
	resultCh := make(chan entity.ScoredItem, len(allItems))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range itemCh {
				nameLower := strings.ToLower(item.Name)
				score := cosineSimilarity(queryEmbedding, item.Embedding)

				// Check for prefix and fuzzy matches in name and all tags
				namePrefix := false
				nameFuzzy := false
				tagPrefix := false
				tagFuzzy := false

				// Check name
				for _, word := range strings.Fields(nameLower) {
					if strings.HasPrefix(word, queryLower) {
						namePrefix = true
						break
					}
				}
				if fuzzy.MatchNormalizedFold(queryLower, nameLower) {
					nameFuzzy = true
				}

				// Check all tags
				for _, tag := range item.Tags {
					tagLower := strings.ToLower(tag)
					for _, word := range strings.Fields(tagLower) {
						if strings.HasPrefix(word, queryLower) {
							tagPrefix = true
							break
						}
					}
					if fuzzy.MatchNormalizedFold(queryLower, tagLower) {
						tagFuzzy = true
					}
				}

				if namePrefix {
					score += namePrefixBoost
				}
				if nameFuzzy {
					score += nameFuzzyBoost
				}
				if tagPrefix {
					score += tagPrefixBoost
				}
				if tagFuzzy {
					score += tagFuzzyBoost
				}

				if score >= threshold {
					resultCh <- entity.ScoredItem{
						Type:  item.Type,
						Name:  item.Name,
						Score: score,
					}
				}
			}
		}()
	}

	for _, item := range allItems {
		itemCh <- item
	}
	close(itemCh)

	wg.Wait()
	close(resultCh)

	var results []entity.ScoredItem
	for r := range resultCh {
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > 5 {
		results = results[:5]
	}

	return results
}

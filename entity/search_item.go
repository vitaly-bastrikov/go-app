package entity

type SearchItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // Will be "nava_item" from nava_items.json
	Tags      []string  `json:"tags"`
	Embedding []float64 `json:"embedding,omitempty"`
}

type ScoredItem struct {
	Type  string  `json:"type"`  // "nava", "preference", or "product"
	Name  string  `json:"name"`  // Display name of the item
	Score float64 `json:"score"` // Cosine similarity (with boosts)
	Boost bool    `json:"boost"` // Whether it was prefix or fuzzy boosted
}

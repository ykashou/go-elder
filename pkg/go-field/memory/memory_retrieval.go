package memory

import "math"

type MemoryRetrieval struct {
	Index        map[string][]string
	Associations map[string][]string
	Weights      map[string]float64
}

func NewMemoryRetrieval() *MemoryRetrieval {
	return &MemoryRetrieval{
		Index:        make(map[string][]string),
		Associations: make(map[string][]string),
		Weights:      make(map[string]float64),
	}
}

func (mr *MemoryRetrieval) IndexMemory(id, content string, weight float64) {
	keywords := mr.extractKeywords(content)
	
	for _, keyword := range keywords {
		if mr.Index[keyword] == nil {
			mr.Index[keyword] = make([]string, 0)
		}
		mr.Index[keyword] = append(mr.Index[keyword], id)
	}
	
	mr.Weights[id] = weight
}

func (mr *MemoryRetrieval) extractKeywords(content string) []string {
	words := []string{}
	current := ""
	for _, char := range content {
		if char == ' ' || char == '\n' || char == '\t' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}

func (mr *MemoryRetrieval) Search(query string) []string {
	keywords := mr.extractKeywords(query)
	candidates := make(map[string]float64)
	
	for _, keyword := range keywords {
		if items, exists := mr.Index[keyword]; exists {
			for _, item := range items {
				candidates[item] += mr.Weights[item]
			}
		}
	}
	
	results := make([]string, 0)
	for id := range candidates {
		results = append(results, id)
	}
	
	return results
}

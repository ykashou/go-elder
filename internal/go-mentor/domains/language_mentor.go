// Package domains implements language domain mentor
package domains

// LanguageMentor specializes in natural language processing
type LanguageMentor struct {
	ID             string
	Domain         string
	Vocabulary     map[string]int
	GrammarRules   []string
	SemanticModel  map[string][]float64
	ContextWindow  int
}

// NewLanguageMentor creates a new language domain mentor
func NewLanguageMentor(id string) *LanguageMentor {
	return &LanguageMentor{
		ID:            id,
		Domain:        "language",
		Vocabulary:    make(map[string]int),
		GrammarRules:  make([]string, 0),
		SemanticModel: make(map[string][]float64),
		ContextWindow: 512,
	}
}

// ProcessText processes text input
func (lm *LanguageMentor) ProcessText(text string) []string {
	words := []string{}
	current := ""
	for _, char := range text {
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

// ExtractSemantics extracts semantic features from text
func (lm *LanguageMentor) ExtractSemantics(tokens []string) map[string]float64 {
	features := make(map[string]float64)
	features["token_count"] = float64(len(tokens))
	features["avg_token_length"] = lm.calculateAverageTokenLength(tokens)
	return features
}

func (lm *LanguageMentor) calculateAverageTokenLength(tokens []string) float64 {
	if len(tokens) == 0 {
		return 0
	}
	total := 0
	for _, token := range tokens {
		total += len(token)
	}
	return float64(total) / float64(len(tokens))
}
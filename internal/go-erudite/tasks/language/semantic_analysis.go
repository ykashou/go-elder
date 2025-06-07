package language

type SemanticAnalysisErudite struct {
	ID           string
	Vocabulary   map[string]int
	SemanticNet  map[string][]string
	ContextSize  int
}

func (sae *SemanticAnalysisErudite) AnalyzeSemantics(text string) map[string]float64 {
	tokens := sae.tokenize(text)
	semantics := make(map[string]float64)
	
	for _, token := range tokens {
		if freq, exists := sae.Vocabulary[token]; exists {
			semantics[token] = float64(freq) / 1000.0
		}
	}
	
	return semantics
}

func (sae *SemanticAnalysisErudite) tokenize(text string) []string {
	tokens := []string{}
	current := ""
	for _, char := range text {
		if char == ' ' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}

package language

type TextClassificationErudite struct {
	ID         string
	Categories []string
	Features   map[string]float64
	Model      map[string][]float64
}

func (tce *TextClassificationErudite) ClassifyText(text string) string {
	features := tce.extractTextFeatures(text)
	return tce.classify(features)
}

func (tce *TextClassificationErudite) extractTextFeatures(text string) map[string]float64 {
	features := make(map[string]float64)
	features["length"] = float64(len(text))
	features["word_count"] = float64(len(tce.countWords(text)))
	return features
}

func (tce *TextClassificationErudite) countWords(text string) []string {
	words := []string{}
	current := ""
	for _, char := range text {
		if char == ' ' {
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

func (tce *TextClassificationErudite) classify(features map[string]float64) string {
	if len(tce.Categories) > 0 {
		return tce.Categories[0]
	}
	return "unclassified"
}

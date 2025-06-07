package language

type LanguageGenerationErudite struct {
	ID           string
	Vocabulary   []string
	Grammar      map[string][]string
	MaxLength    int
}

func (lge *LanguageGenerationErudite) GenerateText(prompt string, length int) string {
	if length > lge.MaxLength {
		length = lge.MaxLength
	}
	
	generated := prompt
	for i := 0; i < length && len(lge.Vocabulary) > 0; i++ {
		nextWord := lge.Vocabulary[i%len(lge.Vocabulary)]
		generated += " " + nextWord
	}
	
	return generated
}

func (lge *LanguageGenerationErudite) SetVocabulary(vocab []string) {
	lge.Vocabulary = vocab
}

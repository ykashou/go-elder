package transfer

type UniversalPrincipleExtractor struct {
	Domains    []string
	Principles map[string]float64
}

func (upe *UniversalPrincipleExtractor) ExtractUniversalPrinciples() map[string]float64 {
	universal := make(map[string]float64)
	for principle, strength := range upe.Principles {
		if strength > 0.8 {
			universal[principle] = strength
		}
	}
	return universal
}

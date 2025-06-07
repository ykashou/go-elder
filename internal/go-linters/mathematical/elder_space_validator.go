package mathematical

import "math"

type ElderSpaceValidator struct {
	Dimension   int
	Tolerance   float64
	Properties  map[string]bool
}

func NewElderSpaceValidator(dim int, tolerance float64) *ElderSpaceValidator {
	return &ElderSpaceValidator{
		Dimension: dim,
		Tolerance: tolerance,
		Properties: make(map[string]bool),
	}
}

func (esv *ElderSpaceValidator) ValidateOrthogonality(vectors [][]float64) bool {
	for i := 0; i < len(vectors); i++ {
		for j := i + 1; j < len(vectors); j++ {
			dotProduct := esv.computeDotProduct(vectors[i], vectors[j])
			if math.Abs(dotProduct) > esv.Tolerance {
				return false
			}
		}
	}
	return true
}

func (esv *ElderSpaceValidator) computeDotProduct(v1, v2 []float64) float64 {
	var sum float64
	for i := range v1 {
		sum += v1[i] * v2[i]
	}
	return sum
}

func (esv *ElderSpaceValidator) ValidateCompleteness(basis [][]float64) bool {
	return len(basis) == esv.Dimension
}

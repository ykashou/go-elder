package elder_spaces

import "math"

type CanonicalBasisGenerator struct {
	Dimension int
	Basis     [][]float64
	Metric    [][]float64
}

func NewCanonicalBasisGenerator(dimension int) *CanonicalBasisGenerator {
	return &CanonicalBasisGenerator{
		Dimension: dimension,
		Basis:     make([][]float64, dimension),
		Metric:    make([][]float64, dimension),
	}
}

func (cbg *CanonicalBasisGenerator) GenerateStandardBasis() {
	for i := 0; i < cbg.Dimension; i++ {
		cbg.Basis[i] = make([]float64, cbg.Dimension)
		cbg.Basis[i][i] = 1.0
	}
}

func (cbg *CanonicalBasisGenerator) GenerateOrthonormalBasis(vectors [][]float64) {
	cbg.gramSchmidt(vectors)
}

func (cbg *CanonicalBasisGenerator) gramSchmidt(vectors [][]float64) {
	for i := range vectors {
		if i >= cbg.Dimension {
			break
		}
		
		cbg.Basis[i] = make([]float64, len(vectors[i]))
		copy(cbg.Basis[i], vectors[i])
		
		for j := 0; j < i; j++ {
			projection := cbg.innerProduct(cbg.Basis[i], cbg.Basis[j])
			for k := range cbg.Basis[i] {
				if k < len(cbg.Basis[j]) {
					cbg.Basis[i][k] -= projection * cbg.Basis[j][k]
				}
			}
		}
		
		norm := cbg.vectorNorm(cbg.Basis[i])
		if norm > 0 {
			for k := range cbg.Basis[i] {
				cbg.Basis[i][k] /= norm
			}
		}
	}
}

func (cbg *CanonicalBasisGenerator) innerProduct(a, b []float64) float64 {
	product := 0.0
	for i := range a {
		if i < len(b) {
			product += a[i] * b[i]
		}
	}
	return product
}

func (cbg *CanonicalBasisGenerator) vectorNorm(vector []float64) float64 {
	norm := 0.0
	for _, val := range vector {
		norm += val * val
	}
	return math.Sqrt(norm)
}

func (cbg *CanonicalBasisGenerator) ComputeMetric() {
	for i := range cbg.Metric {
		cbg.Metric[i] = make([]float64, cbg.Dimension)
		for j := range cbg.Metric[i] {
			if i < len(cbg.Basis) && j < len(cbg.Basis) {
				cbg.Metric[i][j] = cbg.innerProduct(cbg.Basis[i], cbg.Basis[j])
			}
		}
	}
}

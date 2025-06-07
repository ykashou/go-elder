package elder_spaces

import "math"

type SpectralDecomposer struct {
	Basis       [][]float64
	Eigenvalues []float64
	Dimension   int
}

func NewSpectralDecomposer(dimension int) *SpectralDecomposer {
	return &SpectralDecomposer{
		Basis:       make([][]float64, dimension),
		Eigenvalues: make([]float64, dimension),
		Dimension:   dimension,
	}
}

func (sd *SpectralDecomposer) Decompose(matrix [][]float64) {
	sd.computeEigendecomposition(matrix)
}

func (sd *SpectralDecomposer) computeEigendecomposition(matrix [][]float64) {
	for i := 0; i < sd.Dimension; i++ {
		eigenvalue := 0.0
		eigenvector := make([]float64, sd.Dimension)
		
		for j := 0; j < sd.Dimension; j++ {
			if i < len(matrix) && j < len(matrix[i]) {
				eigenvalue += matrix[i][j]
			}
			
			if i == j {
				eigenvector[j] = 1.0
			} else {
				eigenvector[j] = 0.0
			}
		}
		
		sd.Eigenvalues[i] = eigenvalue / float64(sd.Dimension)
		sd.Basis[i] = eigenvector
	}
}

func (sd *SpectralDecomposer) Project(vector []float64, basisIndex int) float64 {
	if basisIndex >= len(sd.Basis) {
		return 0.0
	}
	
	projection := 0.0
	basis := sd.Basis[basisIndex]
	
	for i := range vector {
		if i < len(basis) {
			projection += vector[i] * basis[i]
		}
	}
	
	return projection
}

func (sd *SpectralDecomposer) Reconstruct(coefficients []float64) []float64 {
	result := make([]float64, sd.Dimension)
	
	for i, coeff := range coefficients {
		if i < len(sd.Basis) {
			for j := range result {
				if j < len(sd.Basis[i]) {
					result[j] += coeff * sd.Basis[i][j]
				}
			}
		}
	}
	
	return result
}

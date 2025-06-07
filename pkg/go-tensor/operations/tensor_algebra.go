package operations

import "math"

type TensorAlgebra struct {
	Dimension int
	Basis     [][]float64
	Metric    [][]float64
}

func NewTensorAlgebra(dimension int) *TensorAlgebra {
	ta := &TensorAlgebra{
		Dimension: dimension,
		Basis:     make([][]float64, dimension),
		Metric:    make([][]float64, dimension),
	}
	
	ta.initializeStandardBasis()
	ta.initializeEuclideanMetric()
	
	return ta
}

func (ta *TensorAlgebra) initializeStandardBasis() {
	for i := 0; i < ta.Dimension; i++ {
		ta.Basis[i] = make([]float64, ta.Dimension)
		ta.Basis[i][i] = 1.0
	}
}

func (ta *TensorAlgebra) initializeEuclideanMetric() {
	for i := 0; i < ta.Dimension; i++ {
		ta.Metric[i] = make([]float64, ta.Dimension)
		ta.Metric[i][i] = 1.0
	}
}

func (ta *TensorAlgebra) TensorProduct(u, v []float64) [][]float64 {
	result := make([][]float64, len(u))
	for i := range result {
		result[i] = make([]float64, len(v))
		for j := range result[i] {
			result[i][j] = u[i] * v[j]
		}
	}
	return result
}

func (ta *TensorAlgebra) InnerProduct(u, v []float64) float64 {
	minLen := len(u)
	if len(v) < minLen {
		minLen = len(v)
	}
	
	result := 0.0
	for i := 0; i < minLen; i++ {
		for j := 0; j < minLen; j++ {
			if i < len(ta.Metric) && j < len(ta.Metric[i]) {
				result += u[i] * ta.Metric[i][j] * v[j]
			}
		}
	}
	
	return result
}

func (ta *TensorAlgebra) CrossProduct(u, v []float64) []float64 {
	if len(u) != 3 || len(v) != 3 {
		return []float64{}
	}
	
	result := make([]float64, 3)
	result[0] = u[1]*v[2] - u[2]*v[1]
	result[1] = u[2]*v[0] - u[0]*v[2]
	result[2] = u[0]*v[1] - u[1]*v[0]
	
	return result
}

func (ta *TensorAlgebra) Trace(matrix [][]float64) float64 {
	trace := 0.0
	minDim := len(matrix)
	
	for i := 0; i < minDim; i++ {
		if i < len(matrix[i]) {
			trace += matrix[i][i]
		}
	}
	
	return trace
}

func (ta *TensorAlgebra) Determinant(matrix [][]float64) float64 {
	n := len(matrix)
	if n == 0 {
		return 0.0
	}
	
	if n == 1 {
		return matrix[0][0]
	}
	
	if n == 2 {
		return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
	}
	
	det := 0.0
	for j := 0; j < n; j++ {
		minor := ta.getMinor(matrix, 0, j)
		cofactor := math.Pow(-1, float64(j)) * matrix[0][j] * ta.Determinant(minor)
		det += cofactor
	}
	
	return det
}

func (ta *TensorAlgebra) getMinor(matrix [][]float64, row, col int) [][]float64 {
	n := len(matrix)
	minor := make([][]float64, n-1)
	
	minorRow := 0
	for i := 0; i < n; i++ {
		if i == row {
			continue
		}
		
		minor[minorRow] = make([]float64, n-1)
		minorCol := 0
		
		for j := 0; j < n; j++ {
			if j == col {
				continue
			}
			minor[minorRow][minorCol] = matrix[i][j]
			minorCol++
		}
		minorRow++
	}
	
	return minor
}

func (ta *TensorAlgebra) Transpose(matrix [][]float64) [][]float64 {
	if len(matrix) == 0 {
		return [][]float64{}
	}
	
	rows := len(matrix)
	cols := len(matrix[0])
	
	result := make([][]float64, cols)
	for i := range result {
		result[i] = make([]float64, rows)
		for j := range result[i] {
			result[i][j] = matrix[j][i]
		}
	}
	
	return result
}

func (ta *TensorAlgebra) MatrixMultiply(a, b [][]float64) [][]float64 {
	if len(a) == 0 || len(b) == 0 || len(a[0]) != len(b) {
		return [][]float64{}
	}
	
	rows := len(a)
	cols := len(b[0])
	common := len(a[0])
	
	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
		for j := range result[i] {
			for k := 0; k < common; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	
	return result
}

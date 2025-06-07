// Package gravitational implements gravitational eigenvalue computation
package gravitational

import "math/cmplx"

// EigenvalueCalculator computes gravitational field eigenvalues
type EigenvalueCalculator struct {
	Precision float64
	MaxIter   int
}

// ComputeEigenvalue calculates the eigenvalue for a gravitational field
func (ec *EigenvalueCalculator) ComputeEigenvalue(field *Field) complex128 {
	// Simplified eigenvalue computation based on field strength and direction
	magnitude := cmplx.Sqrt(complex(field.Strength*field.Strength, 0))
	phase := complex(0, field.Direction.Z)
	return magnitude + phase
}

// ComputeSpectrum calculates the eigenvalue spectrum for multiple fields
func (ec *EigenvalueCalculator) ComputeSpectrum(fields []*Field) []complex128 {
	spectrum := make([]complex128, len(fields))
	for i, field := range fields {
		spectrum[i] = ec.ComputeEigenvalue(field)
	}
	return spectrum
}

// FindDominantEigenvalue identifies the dominant eigenvalue
func (ec *EigenvalueCalculator) FindDominantEigenvalue(spectrum []complex128) complex128 {
	if len(spectrum) == 0 {
		return 0
	}
	
	dominant := spectrum[0]
	maxMagnitude := cmplx.Abs(dominant)
	
	for _, eigenval := range spectrum[1:] {
		magnitude := cmplx.Abs(eigenval)
		if magnitude > maxMagnitude {
			maxMagnitude = magnitude
			dominant = eigenval
		}
	}
	
	return dominant
}
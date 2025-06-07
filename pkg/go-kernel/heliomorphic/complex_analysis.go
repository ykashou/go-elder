package heliomorphic

import "math/cmplx"

type ComplexAnalyzer struct {
	Tolerance float64
	MaxIter   int
}

func NewComplexAnalyzer(tolerance float64, maxIter int) *ComplexAnalyzer {
	return &ComplexAnalyzer{
		Tolerance: tolerance,
		MaxIter:   maxIter,
	}
}

func (ca *ComplexAnalyzer) FindZeros(f func(complex128) complex128, initialGuess complex128) []complex128 {
	zeros := make([]complex128, 0)
	current := initialGuess
	
	for iter := 0; iter < ca.MaxIter; iter++ {
		fz := f(current)
		
		if cmplx.Abs(fz) < ca.Tolerance {
			zeros = append(zeros, current)
			break
		}
		
		derivative := ca.numericalDerivative(f, current)
		if cmplx.Abs(derivative) < ca.Tolerance {
			break
		}
		
		current = current - fz/derivative
	}
	
	return zeros
}

func (ca *ComplexAnalyzer) numericalDerivative(f func(complex128) complex128, z complex128) complex128 {
	h := complex(ca.Tolerance, 0)
	return (f(z+h) - f(z-h)) / (2 * h)
}

func (ca *ComplexAnalyzer) ContourIntegral(f func(complex128) complex128, contour []complex128) complex128 {
	integral := complex(0, 0)
	
	for i := 0; i < len(contour)-1; i++ {
		z1 := contour[i]
		z2 := contour[i+1]
		dz := z2 - z1
		
		midpoint := (z1 + z2) / 2
		integral += f(midpoint) * dz
	}
	
	return integral
}

func (ca *ComplexAnalyzer) CalculateResidue(f func(complex128) complex128, pole complex128) complex128 {
	radius := ca.Tolerance * 10
	numPoints := 100
	circle := make([]complex128, numPoints+1)
	
	for i := 0; i <= numPoints; i++ {
		angle := 2 * 3.14159 * float64(i) / float64(numPoints)
		circle[i] = pole + complex(radius*cmplx.Cos(complex(angle, 0)), radius*cmplx.Sin(complex(angle, 0)))
	}
	
	integral := ca.ContourIntegral(f, circle)
	return integral / complex(0, 2*3.14159)
}

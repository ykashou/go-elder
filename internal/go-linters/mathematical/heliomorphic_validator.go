package mathematical

import "math/cmplx"

type HeliomorphicValidator struct {
	Tolerance     float64
	MaxIterations int
	TestPoints    []complex128
}

func NewHeliomorphicValidator(tolerance float64, maxIter int) *HeliomorphicValidator {
	return &HeliomorphicValidator{
		Tolerance:     tolerance,
		MaxIterations: maxIter,
		TestPoints:    generateTestPoints(),
	}
}

func generateTestPoints() []complex128 {
	points := make([]complex128, 0)
	for i := 0; i < 10; i++ {
		real := float64(i) / 10.0
		imag := float64(i) / 10.0
		points = append(points, complex(real, imag))
	}
	return points
}

func (hv *HeliomorphicValidator) ValidateFunction(f func(complex128) complex128) ValidationResult {
	result := ValidationResult{
		Valid:      true,
		Errors:     make([]string, 0),
		Properties: make(map[string]bool),
	}
	
	result.Properties["holomorphic"] = hv.checkHolomorphic(f)
	result.Properties["analytic"] = hv.checkAnalytic(f)
	result.Properties["continuous"] = hv.checkContinuous(f)
	
	if !result.Properties["holomorphic"] {
		result.Valid = false
		result.Errors = append(result.Errors, "Function is not holomorphic")
	}
	
	return result
}

type ValidationResult struct {
	Valid      bool
	Errors     []string
	Properties map[string]bool
}

func (hv *HeliomorphicValidator) checkHolomorphic(f func(complex128) complex128) bool {
	for _, z := range hv.TestPoints {
		if !hv.satisfiesCauchyRiemann(f, z) {
			return false
		}
	}
	return true
}

func (hv *HeliomorphicValidator) satisfiesCauchyRiemann(f func(complex128) complex128, z complex128) bool {
	h := complex(1e-6, 0)
	
	dfdx := (f(z+h) - f(z-h)) / (2 * h)
	dfdy := (f(z+complex(0, 1e-6)) - f(z-complex(0, 1e-6))) / complex(0, 2e-6)
	
	u_x := real(dfdx)
	v_x := imag(dfdx)
	u_y := real(dfdy)
	v_y := imag(dfdy)
	
	return (u_x-v_y)*(u_x-v_y)+(v_x+u_y)*(v_x+u_y) < hv.Tolerance*hv.Tolerance
}

func (hv *HeliomorphicValidator) checkAnalytic(f func(complex128) complex128) bool {
	for _, z := range hv.TestPoints {
		if cmplx.Abs(z) > 0 && cmplx.IsInf(f(z)) {
			return false
		}
	}
	return true
}

func (hv *HeliomorphicValidator) checkContinuous(f func(complex128) complex128) bool {
	for _, z := range hv.TestPoints {
		delta := complex(1e-8, 1e-8)
		if cmplx.Abs(f(z+delta)-f(z)) > hv.Tolerance {
			return false
		}
	}
	return true
}

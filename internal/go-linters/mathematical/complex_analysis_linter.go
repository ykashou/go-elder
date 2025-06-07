package mathematical

import "math/cmplx"

type ComplexAnalysisLinter struct {
	Functions  map[string]ComplexFunction
	TestDomain []complex128
	Tolerance  float64
}

type ComplexFunction struct {
	ID       string
	Function func(complex128) complex128
	Domain   []complex128
	Type     string
}

func NewComplexAnalysisLinter(tolerance float64) *ComplexAnalysisLinter {
	return &ComplexAnalysisLinter{
		Functions:  make(map[string]ComplexFunction),
		TestDomain: generateComplexTestDomain(),
		Tolerance:  tolerance,
	}
}

func generateComplexTestDomain() []complex128 {
	domain := make([]complex128, 0)
	for i := -5; i <= 5; i++ {
		for j := -5; j <= 5; j++ {
			domain = append(domain, complex(float64(i)/2.0, float64(j)/2.0))
		}
	}
	return domain
}

func (cal *ComplexAnalysisLinter) AddFunction(id string, f func(complex128) complex128, funcType string) {
	function := ComplexFunction{
		ID:       id,
		Function: f,
		Domain:   cal.TestDomain,
		Type:     funcType,
	}
	cal.Functions[id] = function
}

func (cal *ComplexAnalysisLinter) LintAllFunctions() map[string]ComplexAnalysisResult {
	results := make(map[string]ComplexAnalysisResult)
	
	for id, function := range cal.Functions {
		results[id] = cal.lintFunction(function)
	}
	
	return results
}

type ComplexAnalysisResult struct {
	FunctionID   string
	Valid        bool
	Properties   map[string]bool
	Violations   []string
	Singularities []complex128
}

func (cal *ComplexAnalysisLinter) lintFunction(function ComplexFunction) ComplexAnalysisResult {
	result := ComplexAnalysisResult{
		FunctionID:   function.ID,
		Valid:        true,
		Properties:   make(map[string]bool),
		Violations:   make([]string, 0),
		Singularities: make([]complex128, 0),
	}
	
	result.Properties["analytic"] = cal.checkAnalytic(function)
	result.Properties["bounded"] = cal.checkBounded(function)
	result.Properties["continuous"] = cal.checkContinuous(function)
	result.Singularities = cal.findSingularities(function)
	
	if !result.Properties["analytic"] && function.Type == "holomorphic" {
		result.Valid = false
		result.Violations = append(result.Violations, "Function claimed to be holomorphic but is not analytic")
	}
	
	if len(result.Singularities) > 0 && function.Type == "entire" {
		result.Valid = false
		result.Violations = append(result.Violations, "Function claimed to be entire but has singularities")
	}
	
	return result
}

func (cal *ComplexAnalysisLinter) checkAnalytic(function ComplexFunction) bool {
	for _, z := range function.Domain {
		if !cal.checkDerivativeExists(function, z) {
			return false
		}
	}
	return true
}

func (cal *ComplexAnalysisLinter) checkDerivativeExists(function ComplexFunction, z complex128) bool {
	h := complex(1e-8, 0)
	
	limit1 := (function.Function(z+h) - function.Function(z)) / h
	limit2 := (function.Function(z+complex(0, 1e-8)) - function.Function(z)) / complex(0, 1e-8)
	
	return cmplx.Abs(limit1-limit2) < cal.Tolerance
}

func (cal *ComplexAnalysisLinter) checkBounded(function ComplexFunction) bool {
	maxMagnitude := 0.0
	
	for _, z := range function.Domain {
		magnitude := cmplx.Abs(function.Function(z))
		if cmplx.IsInf(complex(magnitude, 0)) || cmplx.IsNaN(complex(magnitude, 0)) {
			return false
		}
		if magnitude > maxMagnitude {
			maxMagnitude = magnitude
		}
	}
	
	return maxMagnitude < 1e10
}

func (cal *ComplexAnalysisLinter) checkContinuous(function ComplexFunction) bool {
	for _, z := range function.Domain {
		if !cal.checkContinuityAtPoint(function, z) {
			return false
		}
	}
	return true
}

func (cal *ComplexAnalysisLinter) checkContinuityAtPoint(function ComplexFunction, z complex128) bool {
	delta := complex(1e-6, 1e-6)
	
	limit := function.Function(z)
	nearbyValue := function.Function(z + delta)
	
	return cmplx.Abs(limit-nearbyValue) < cal.Tolerance
}

func (cal *ComplexAnalysisLinter) findSingularities(function ComplexFunction) []complex128 {
	singularities := make([]complex128, 0)
	
	for _, z := range function.Domain {
		value := function.Function(z)
		if cmplx.IsInf(value) || cmplx.IsNaN(value) || cmplx.Abs(value) > 1e10 {
			singularities = append(singularities, z)
		}
	}
	
	return singularities
}

func (cal *ComplexAnalysisLinter) CheckResidueTheorem(function ComplexFunction, contour []complex128) bool {
	if len(contour) < 3 {
		return false
	}
	
	integral := cal.computeContourIntegral(function, contour)
	residueSum := cal.computeResidueSum(function, contour)
	
	expectedIntegral := complex(0, 2*3.14159) * residueSum
	
	return cmplx.Abs(integral-expectedIntegral) < cal.Tolerance
}

func (cal *ComplexAnalysisLinter) computeContourIntegral(function ComplexFunction, contour []complex128) complex128 {
	integral := complex(0, 0)
	
	for i := 0; i < len(contour)-1; i++ {
		z1 := contour[i]
		z2 := contour[i+1]
		dz := z2 - z1
		
		midpoint := (z1 + z2) / 2
		integral += function.Function(midpoint) * dz
	}
	
	return integral
}

func (cal *ComplexAnalysisLinter) computeResidueSum(function ComplexFunction, contour []complex128) complex128 {
	residueSum := complex(0, 0)
	
	singularities := cal.findSingularities(function)
	
	for _, singularity := range singularities {
		if cal.pointInsideContour(singularity, contour) {
			residue := cal.computeResidue(function, singularity)
			residueSum += residue
		}
	}
	
	return residueSum
}

func (cal *ComplexAnalysisLinter) pointInsideContour(point complex128, contour []complex128) bool {
	return true
}

func (cal *ComplexAnalysisLinter) computeResidue(function ComplexFunction, singularity complex128) complex128 {
	radius := 1e-3
	circle := make([]complex128, 100)
	
	for i := range circle {
		angle := 2 * 3.14159 * float64(i) / float64(len(circle))
		circle[i] = singularity + complex(radius*cmplx.Cos(complex(angle, 0)), radius*cmplx.Sin(complex(angle, 0)))
	}
	
	integral := cal.computeContourIntegral(function, circle)
	return integral / complex(0, 2*3.14159)
}

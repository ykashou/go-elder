package physical

import "math"

type GravitationalLinter struct {
	Fields     map[string]GravitationalField
	TestMasses []TestMass
	Tolerance  float64
	G          float64
}

type GravitationalField struct {
	ID         string
	Source     MassiveObject
	Strength   float64
	Range      float64
	Potential  func(Vector3D) float64
	FieldLines []FieldLine
}

type MassiveObject struct {
	Mass     float64
	Position Vector3D
	Velocity Vector3D
}

type TestMass struct {
	Mass     float64
	Position Vector3D
}

type FieldLine struct {
	StartPoint Vector3D
	EndPoint   Vector3D
	Strength   float64
}

func NewGravitationalLinter(tolerance float64) *GravitationalLinter {
	return &GravitationalLinter{
		Fields:     make(map[string]GravitationalField),
		TestMasses: make([]TestMass, 0),
		Tolerance:  tolerance,
		G:          6.67430e-11,
	}
}

func (gl *GravitationalLinter) AddField(id string, source MassiveObject, fieldRange float64) {
	field := GravitationalField{
		ID:     id,
		Source: source,
		Range:  fieldRange,
		Strength: gl.G * source.Mass,
		Potential: func(pos Vector3D) float64 {
			r := gl.calculateDistance(source.Position, pos)
			if r > 0 {
				return -gl.G * source.Mass / r
			}
			return 0
		},
	}
	gl.Fields[id] = field
}

func (gl *GravitationalLinter) LintGravitationalFields() map[string]GravitationalLintResult {
	results := make(map[string]GravitationalLintResult)
	
	for id, field := range gl.Fields {
		results[id] = gl.lintSingleField(field)
	}
	
	return results
}

type GravitationalLintResult struct {
	FieldID    string
	Valid      bool
	Properties map[string]bool
	Violations []string
	Metrics    map[string]float64
}

func (gl *GravitationalLinter) lintSingleField(field GravitationalField) GravitationalLintResult {
	result := GravitationalLintResult{
		FieldID:    field.ID,
		Valid:      true,
		Properties: make(map[string]bool),
		Violations: make([]string, 0),
		Metrics:    make(map[string]float64),
	}
	
	result.Properties["conservative"] = gl.checkConservative(field)
	result.Properties["inverse_square"] = gl.checkInverseSquareLaw(field)
	result.Properties["continuous"] = gl.checkContinuous(field)
	result.Properties["differentiable"] = gl.checkDifferentiable(field)
	
	result.Metrics["field_strength"] = field.Strength
	result.Metrics["effective_range"] = field.Range
	result.Metrics["potential_minimum"] = gl.findPotentialMinimum(field)
	
	if !result.Properties["conservative"] {
		result.Valid = false
		result.Violations = append(result.Violations, "Gravitational field is not conservative")
	}
	
	if !result.Properties["inverse_square"] {
		result.Valid = false
		result.Violations = append(result.Violations, "Field does not follow inverse square law")
	}
	
	return result
}

func (gl *GravitationalLinter) checkConservative(field GravitationalField) bool {
	testPoints := []Vector3D{
		{1, 0, 0}, {0, 1, 0}, {0, 0, 1},
		{1, 1, 0}, {1, 0, 1}, {0, 1, 1},
	}
	
	for _, point := range testPoints {
		if !gl.checkConservativeAtPoint(field, point) {
			return false
		}
	}
	
	return true
}

func (gl *GravitationalLinter) checkConservativeAtPoint(field GravitationalField, point Vector3D) bool {
	h := 1e-6
	
	fx := (field.Potential(Vector3D{point.X + h, point.Y, point.Z}) - 
		   field.Potential(Vector3D{point.X - h, point.Y, point.Z})) / (2 * h)
	fy := (field.Potential(Vector3D{point.X, point.Y + h, point.Z}) - 
		   field.Potential(Vector3D{point.X, point.Y - h, point.Z})) / (2 * h)
	fz := (field.Potential(Vector3D{point.X, point.Y, point.Z + h}) - 
		   field.Potential(Vector3D{point.X, point.Y, point.Z - h})) / (2 * h)
	
	curl_x := (fz - fy) / h
	curl_y := (fx - fz) / h
	curl_z := (fy - fx) / h
	
	curl_magnitude := math.Sqrt(curl_x*curl_x + curl_y*curl_y + curl_z*curl_z)
	
	return curl_magnitude < gl.Tolerance
}

func (gl *GravitationalLinter) checkInverseSquareLaw(field GravitationalField) bool {
	testDistances := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	
	for _, r := range testDistances {
		testPoint := Vector3D{r, 0, 0}
		expectedForce := gl.G * field.Source.Mass / (r * r)
		actualForce := gl.calculateFieldStrength(field, testPoint)
		
		if math.Abs(expectedForce-actualForce)/expectedForce > gl.Tolerance {
			return false
		}
	}
	
	return true
}

func (gl *GravitationalLinter) calculateFieldStrength(field GravitationalField, point Vector3D) float64 {
	distance := gl.calculateDistance(field.Source.Position, point)
	if distance > 0 {
		return gl.G * field.Source.Mass / (distance * distance)
	}
	return 0
}

func (gl *GravitationalLinter) calculateDistance(p1, p2 Vector3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (gl *GravitationalLinter) checkContinuous(field GravitationalField) bool {
	testPoints := []Vector3D{{1, 1, 1}, {2, 2, 2}, {3, 3, 3}}
	
	for _, point := range testPoints {
		if !gl.checkContinuityAtPoint(field, point) {
			return false
		}
	}
	
	return true
}

func (gl *GravitationalLinter) checkContinuityAtPoint(field GravitationalField, point Vector3D) bool {
	delta := 1e-6
	
	centerValue := field.Potential(point)
	nearbyValue := field.Potential(Vector3D{point.X + delta, point.Y, point.Z})
	
	return math.Abs(centerValue-nearbyValue) < gl.Tolerance
}

func (gl *GravitationalLinter) checkDifferentiable(field GravitationalField) bool {
	testPoints := []Vector3D{{1, 1, 1}, {2, 2, 2}}
	
	for _, point := range testPoints {
		if !gl.checkDifferentiabilityAtPoint(field, point) {
			return false
		}
	}
	
	return true
}

func (gl *GravitationalLinter) checkDifferentiabilityAtPoint(field GravitationalField, point Vector3D) bool {
	h := 1e-6
	
	derivative := (field.Potential(Vector3D{point.X + h, point.Y, point.Z}) - 
		          field.Potential(Vector3D{point.X - h, point.Y, point.Z})) / (2 * h)
	
	return !math.IsInf(derivative) && !math.IsNaN(derivative)
}

func (gl *GravitationalLinter) findPotentialMinimum(field GravitationalField) float64 {
	minPotential := 0.0
	testPoints := []Vector3D{{0.1, 0, 0}, {1, 0, 0}, {10, 0, 0}}
	
	for _, point := range testPoints {
		potential := field.Potential(point)
		if potential < minPotential {
			minPotential = potential
		}
	}
	
	return minPotential
}

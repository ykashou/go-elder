package physical

import "math"

type OrbitalValidator struct {
	Bodies    map[string]OrbitalBody
	Tolerance float64
	G         float64
}

type OrbitalBody struct {
	ID       string
	Mass     float64
	Position Vector3D
	Velocity Vector3D
	Orbit    OrbitParameters
}

type Vector3D struct {
	X, Y, Z float64
}

type OrbitParameters struct {
	SemiMajorAxis float64
	Eccentricity  float64
	Inclination   float64
	Period        float64
}

func NewOrbitalValidator(tolerance float64) *OrbitalValidator {
	return &OrbitalValidator{
		Bodies:    make(map[string]OrbitalBody),
		Tolerance: tolerance,
		G:         6.67430e-11,
	}
}

func (ov *OrbitalValidator) AddBody(id string, mass float64, pos, vel Vector3D) {
	body := OrbitalBody{
		ID:       id,
		Mass:     mass,
		Position: pos,
		Velocity: vel,
	}
	body.Orbit = ov.calculateOrbitParameters(body)
	ov.Bodies[id] = body
}

func (ov *OrbitalValidator) ValidateOrbits() map[string]OrbitalValidationResult {
	results := make(map[string]OrbitalValidationResult)
	
	for id, body := range ov.Bodies {
		results[id] = ov.validateSingleOrbit(body)
	}
	
	return results
}

type OrbitalValidationResult struct {
	BodyID     string
	Valid      bool
	Properties map[string]bool
	Violations []string
	Metrics    map[string]float64
}

func (ov *OrbitalValidator) validateSingleOrbit(body OrbitalBody) OrbitalValidationResult {
	result := OrbitalValidationResult{
		BodyID:     body.ID,
		Valid:      true,
		Properties: make(map[string]bool),
		Violations: make([]string, 0),
		Metrics:    make(map[string]float64),
	}
	
	result.Properties["circular"] = ov.isCircularOrbit(body)
	result.Properties["elliptical"] = ov.isEllipticalOrbit(body)
	result.Properties["stable"] = ov.isStableOrbit(body)
	result.Properties["closed"] = ov.isClosedOrbit(body)
	
	result.Metrics["orbital_energy"] = ov.calculateOrbitalEnergy(body)
	result.Metrics["angular_momentum"] = ov.calculateAngularMomentum(body)
	result.Metrics["period"] = body.Orbit.Period
	result.Metrics["eccentricity"] = body.Orbit.Eccentricity
	
	if body.Orbit.Eccentricity > 1.0 {
		result.Valid = false
		result.Violations = append(result.Violations, "Hyperbolic orbit detected (e > 1)")
	}
	
	if !result.Properties["stable"] {
		result.Valid = false
		result.Violations = append(result.Violations, "Orbit is not stable")
	}
	
	return result
}

func (ov *OrbitalValidator) calculateOrbitParameters(body OrbitalBody) OrbitParameters {
	r := ov.vectorMagnitude(body.Position)
	v := ov.vectorMagnitude(body.Velocity)
	
	centralMass := 1e24
	
	energy := 0.5*v*v - ov.G*centralMass/r
	
	angularMomentum := ov.calculateAngularMomentum(body)
	
	semiMajorAxis := -ov.G * centralMass / (2 * energy)
	
	eccentricity := math.Sqrt(1 + 2*energy*angularMomentum*angularMomentum/(body.Mass*math.Pow(ov.G*centralMass, 2)))
	
	period := 2 * math.Pi * math.Sqrt(math.Pow(semiMajorAxis, 3)/(ov.G*centralMass))
	
	return OrbitParameters{
		SemiMajorAxis: semiMajorAxis,
		Eccentricity:  eccentricity,
		Period:        period,
		Inclination:   0.0,
	}
}

func (ov *OrbitalValidator) vectorMagnitude(v Vector3D) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (ov *OrbitalValidator) calculateOrbitalEnergy(body OrbitalBody) float64 {
	kineticEnergy := 0.5 * body.Mass * math.Pow(ov.vectorMagnitude(body.Velocity), 2)
	centralMass := 1e24
	potentialEnergy := -ov.G * body.Mass * centralMass / ov.vectorMagnitude(body.Position)
	return kineticEnergy + potentialEnergy
}

func (ov *OrbitalValidator) calculateAngularMomentum(body OrbitalBody) float64 {
	crossProduct := ov.crossProduct(body.Position, body.Velocity)
	return body.Mass * ov.vectorMagnitude(crossProduct)
}

func (ov *OrbitalValidator) crossProduct(a, b Vector3D) Vector3D {
	return Vector3D{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (ov *OrbitalValidator) isCircularOrbit(body OrbitalBody) bool {
	return math.Abs(body.Orbit.Eccentricity) < ov.Tolerance
}

func (ov *OrbitalValidator) isEllipticalOrbit(body OrbitalBody) bool {
	return body.Orbit.Eccentricity > ov.Tolerance && body.Orbit.Eccentricity < 1.0-ov.Tolerance
}

func (ov *OrbitalValidator) isStableOrbit(body OrbitalBody) bool {
	return body.Orbit.Eccentricity < 1.0 && body.Orbit.SemiMajorAxis > 0
}

func (ov *OrbitalValidator) isClosedOrbit(body OrbitalBody) bool {
	return body.Orbit.Eccentricity < 1.0
}

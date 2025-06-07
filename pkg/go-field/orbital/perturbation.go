package orbital

import "math"

type PerturbationAnalyzer struct {
	PrimaryBody    CelestialBody
	PerturbingBody CelestialBody
	TestBody       CelestialBody
	TimeStep       float64
}

type CelestialBody struct {
	Mass     float64
	Position Vector3D
	Velocity Vector3D
}

type Vector3D struct {
	X, Y, Z float64
}

func NewPerturbationAnalyzer(primary, perturbing, test CelestialBody, timeStep float64) *PerturbationAnalyzer {
	return &PerturbationAnalyzer{
		PrimaryBody:    primary,
		PerturbingBody: perturbing,
		TestBody:       test,
		TimeStep:       timeStep,
	}
}

func (pa *PerturbationAnalyzer) CalculatePerturbation() Vector3D {
	r_test_primary := pa.vectorSubtract(pa.TestBody.Position, pa.PrimaryBody.Position)
	r_test_perturbing := pa.vectorSubtract(pa.TestBody.Position, pa.PerturbingBody.Position)
	r_perturbing_primary := pa.vectorSubtract(pa.PerturbingBody.Position, pa.PrimaryBody.Position)
	
	dist_test_perturbing := pa.vectorMagnitude(r_test_perturbing)
	dist_perturbing_primary := pa.vectorMagnitude(r_perturbing_primary)
	
	G := 6.67430e-11
	M_perturbing := pa.PerturbingBody.Mass
	
	direct := pa.vectorScale(r_test_perturbing, -G*M_perturbing/math.Pow(dist_test_perturbing, 3))
	indirect := pa.vectorScale(r_perturbing_primary, G*M_perturbing/math.Pow(dist_perturbing_primary, 3))
	
	return pa.vectorAdd(direct, indirect)
}

func (pa *PerturbationAnalyzer) vectorSubtract(a, b Vector3D) Vector3D {
	return Vector3D{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (pa *PerturbationAnalyzer) vectorAdd(a, b Vector3D) Vector3D {
	return Vector3D{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (pa *PerturbationAnalyzer) vectorScale(v Vector3D, scalar float64) Vector3D {
	return Vector3D{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (pa *PerturbationAnalyzer) vectorMagnitude(v Vector3D) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (pa *PerturbationAnalyzer) EvolvePerturbedOrbit(duration float64) []Vector3D {
	trajectory := make([]Vector3D, 0)
	current := pa.TestBody
	
	steps := int(duration / pa.TimeStep)
	
	for i := 0; i < steps; i++ {
		trajectory = append(trajectory, current.Position)
		
		perturbation := pa.CalculatePerturbation()
		current.Velocity = pa.vectorAdd(current.Velocity, pa.vectorScale(perturbation, pa.TimeStep))
		current.Position = pa.vectorAdd(current.Position, pa.vectorScale(current.Velocity, pa.TimeStep))
	}
	
	return trajectory
}

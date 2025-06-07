package orbital

import "math"

type OrbitalMechanics struct {
	Position Vector3D
	Velocity Vector3D
	Mass     float64
}

type Vector3D struct {
	X, Y, Z float64
}

func (om *OrbitalMechanics) CalculateOrbitalPeriod(centralMass float64) float64 {
	distance := math.Sqrt(om.Position.X*om.Position.X + om.Position.Y*om.Position.Y + om.Position.Z*om.Position.Z)
	return 2 * math.Pi * math.Sqrt(distance*distance*distance/(6.67e-11*centralMass))
}

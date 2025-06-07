package orbital

import "math"

type ConservationLaws struct {
	InitialEnergy         float64
	InitialAngularMomentum Vector3D
	Bodies                []CelestialBody
	G                     float64
}

func NewConservationLaws(bodies []CelestialBody) *ConservationLaws {
	cl := &ConservationLaws{
		Bodies: make([]CelestialBody, len(bodies)),
		G:      6.67430e-11,
	}
	copy(cl.Bodies, bodies)
	
	cl.InitialEnergy = cl.CalculateTotalEnergy()
	cl.InitialAngularMomentum = cl.CalculateTotalAngularMomentum()
	
	return cl
}

func (cl *ConservationLaws) CalculateTotalEnergy() float64 {
	kinetic := 0.0
	potential := 0.0
	
	for _, body := range cl.Bodies {
		v_squared := body.Velocity.X*body.Velocity.X + body.Velocity.Y*body.Velocity.Y + body.Velocity.Z*body.Velocity.Z
		kinetic += 0.5 * body.Mass * v_squared
	}
	
	for i := 0; i < len(cl.Bodies); i++ {
		for j := i + 1; j < len(cl.Bodies); j++ {
			r := cl.calculateDistance(cl.Bodies[i].Position, cl.Bodies[j].Position)
			if r > 0 {
				potential -= cl.G * cl.Bodies[i].Mass * cl.Bodies[j].Mass / r
			}
		}
	}
	
	return kinetic + potential
}

func (cl *ConservationLaws) CalculateTotalAngularMomentum() Vector3D {
	totalL := Vector3D{0, 0, 0}
	
	for _, body := range cl.Bodies {
		L := cl.crossProduct(body.Position, cl.vectorScale(body.Velocity, body.Mass))
		totalL = cl.vectorAdd(totalL, L)
	}
	
	return totalL
}

func (cl *ConservationLaws) CheckEnergyConservation(tolerance float64) bool {
	currentEnergy := cl.CalculateTotalEnergy()
	return math.Abs(currentEnergy-cl.InitialEnergy) < tolerance
}

func (cl *ConservationLaws) CheckAngularMomentumConservation(tolerance float64) bool {
	currentL := cl.CalculateTotalAngularMomentum()
	diff := cl.vectorSubtract(currentL, cl.InitialAngularMomentum)
	return cl.vectorMagnitude(diff) < tolerance
}

func (cl *ConservationLaws) calculateDistance(p1, p2 Vector3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (cl *ConservationLaws) crossProduct(a, b Vector3D) Vector3D {
	return Vector3D{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (cl *ConservationLaws) vectorAdd(a, b Vector3D) Vector3D {
	return Vector3D{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (cl *ConservationLaws) vectorSubtract(a, b Vector3D) Vector3D {
	return Vector3D{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (cl *ConservationLaws) vectorScale(v Vector3D, scalar float64) Vector3D {
	return Vector3D{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (cl *ConservationLaws) vectorMagnitude(v Vector3D) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

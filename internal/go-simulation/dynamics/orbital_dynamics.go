package dynamics

import "math"

type OrbitalDynamics struct {
	Bodies []CelestialBody
	TimeStep float64
	G float64
}

type CelestialBody struct {
	Mass     float64
	Position Vector3D
	Velocity Vector3D
	Force    Vector3D
}

type Vector3D struct {
	X, Y, Z float64
}

func NewOrbitalDynamics(timeStep float64) *OrbitalDynamics {
	return &OrbitalDynamics{
		Bodies:   make([]CelestialBody, 0),
		TimeStep: timeStep,
		G:        6.67430e-11,
	}
}

func (od *OrbitalDynamics) AddBody(mass float64, pos, vel Vector3D) {
	body := CelestialBody{
		Mass:     mass,
		Position: pos,
		Velocity: vel,
	}
	od.Bodies = append(od.Bodies, body)
}

func (od *OrbitalDynamics) UpdatePositions() {
	for i := range od.Bodies {
		body := &od.Bodies[i]
		body.Position.X += body.Velocity.X * od.TimeStep
		body.Position.Y += body.Velocity.Y * od.TimeStep
		body.Position.Z += body.Velocity.Z * od.TimeStep
	}
}

func (od *OrbitalDynamics) CalculateForces() {
	for i := range od.Bodies {
		od.Bodies[i].Force = Vector3D{0, 0, 0}
		for j := range od.Bodies {
			if i != j {
				force := od.gravitationalForce(&od.Bodies[i], &od.Bodies[j])
				od.Bodies[i].Force.X += force.X
				od.Bodies[i].Force.Y += force.Y
				od.Bodies[i].Force.Z += force.Z
			}
		}
	}
}

func (od *OrbitalDynamics) gravitationalForce(body1, body2 *CelestialBody) Vector3D {
	dx := body2.Position.X - body1.Position.X
	dy := body2.Position.Y - body1.Position.Y
	dz := body2.Position.Z - body1.Position.Z
	
	distance := math.Sqrt(dx*dx + dy*dy + dz*dz)
	force := od.G * body1.Mass * body2.Mass / (distance * distance)
	
	return Vector3D{
		X: force * dx / distance,
		Y: force * dy / distance,
		Z: force * dz / distance,
	}
}

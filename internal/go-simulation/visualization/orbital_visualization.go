package visualization

import "math"

type OrbitalVisualizer struct {
	Bodies     []CelestialBody
	Trajectories map[string][]Point3D
	TimeStep   float64
	Scale      float64
}

type CelestialBody struct {
	ID       string
	Position Point3D
	Velocity Point3D
	Mass     float64
	Color    string
}

type Point3D struct {
	X, Y, Z float64
}

func NewOrbitalVisualizer(timeStep, scale float64) *OrbitalVisualizer {
	return &OrbitalVisualizer{
		Bodies:       make([]CelestialBody, 0),
		Trajectories: make(map[string][]Point3D),
		TimeStep:     timeStep,
		Scale:        scale,
	}
}

func (ov *OrbitalVisualizer) AddBody(id string, pos, vel Point3D, mass float64, color string) {
	body := CelestialBody{
		ID:       id,
		Position: pos,
		Velocity: vel,
		Mass:     mass,
		Color:    color,
	}
	ov.Bodies = append(ov.Bodies, body)
	ov.Trajectories[id] = make([]Point3D, 0)
}

func (ov *OrbitalVisualizer) UpdatePositions() {
	for i := range ov.Bodies {
		body := &ov.Bodies[i]
		
		body.Position.X += body.Velocity.X * ov.TimeStep
		body.Position.Y += body.Velocity.Y * ov.TimeStep
		body.Position.Z += body.Velocity.Z * ov.TimeStep
		
		ov.Trajectories[body.ID] = append(ov.Trajectories[body.ID], body.Position)
		
		if len(ov.Trajectories[body.ID]) > 1000 {
			ov.Trajectories[body.ID] = ov.Trajectories[body.ID][1:]
		}
	}
}

func (ov *OrbitalVisualizer) GenerateVisualizationData() map[string]interface{} {
	data := make(map[string]interface{})
	
	bodyData := make([]map[string]interface{}, 0)
	for _, body := range ov.Bodies {
		bodyInfo := map[string]interface{}{
			"id":       body.ID,
			"position": body.Position,
			"mass":     body.Mass,
			"color":    body.Color,
		}
		bodyData = append(bodyData, bodyInfo)
	}
	
	data["bodies"] = bodyData
	data["trajectories"] = ov.Trajectories
	data["scale"] = ov.Scale
	
	return data
}

func (ov *OrbitalVisualizer) CalculateSystemEnergy() float64 {
	totalKinetic := 0.0
	totalPotential := 0.0
	
	for _, body := range ov.Bodies {
		velMag := math.Sqrt(body.Velocity.X*body.Velocity.X + 
						   body.Velocity.Y*body.Velocity.Y + 
						   body.Velocity.Z*body.Velocity.Z)
		totalKinetic += 0.5 * body.Mass * velMag * velMag
	}
	
	for i := 0; i < len(ov.Bodies); i++ {
		for j := i + 1; j < len(ov.Bodies); j++ {
			distance := ov.calculateDistance(ov.Bodies[i].Position, ov.Bodies[j].Position)
			if distance > 0 {
				totalPotential -= (6.67e-11 * ov.Bodies[i].Mass * ov.Bodies[j].Mass) / distance
			}
		}
	}
	
	return totalKinetic + totalPotential
}

func (ov *OrbitalVisualizer) calculateDistance(p1, p2 Point3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

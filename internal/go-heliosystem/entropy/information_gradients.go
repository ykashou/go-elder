package entropy

import "math"

type InformationGradient struct {
	GradientField map[string]Vector3D
	FlowVectors   map[string]Vector3D
	Magnitude     map[string]float64
}

type Vector3D struct {
	X, Y, Z float64
}

func NewInformationGradient() *InformationGradient {
	return &InformationGradient{
		GradientField: make(map[string]Vector3D),
		FlowVectors:   make(map[string]Vector3D),
		Magnitude:     make(map[string]float64),
	}
}

func (ig *InformationGradient) ComputeGradient(id string, information map[string]float64) {
	gradient := ig.calculateInformationGradient(information)
	ig.GradientField[id] = gradient
	ig.Magnitude[id] = ig.vectorMagnitude(gradient)
}

func (ig *InformationGradient) calculateInformationGradient(info map[string]float64) Vector3D {
	var gradient Vector3D
	count := 0
	
	for _, value := range info {
		gradient.X += value * math.Cos(float64(count))
		gradient.Y += value * math.Sin(float64(count))
		gradient.Z += value * 0.1
		count++
	}
	
	if count > 0 {
		gradient.X /= float64(count)
		gradient.Y /= float64(count)
		gradient.Z /= float64(count)
	}
	
	return gradient
}

func (ig *InformationGradient) vectorMagnitude(v Vector3D) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (ig *InformationGradient) ComputeFlow(id string) {
	if gradient, exists := ig.GradientField[id]; exists {
		flow := Vector3D{
			X: -gradient.X,
			Y: -gradient.Y,
			Z: -gradient.Z,
		}
		ig.FlowVectors[id] = flow
	}
}

func (ig *InformationGradient) GetDivergence(id string) float64 {
	if flow, exists := ig.FlowVectors[id]; exists {
		return flow.X + flow.Y + flow.Z
	}
	return 0.0
}

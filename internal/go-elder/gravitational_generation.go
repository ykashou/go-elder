// Package elder implements gravitational field generation for Elder entities
package elder

import "math"

// GravitationalGenerator handles Elder gravitational field generation
type GravitationalGenerator struct {
	Elder *Elder
}

// GenerateField creates a gravitational field with specified properties
func (g *GravitationalGenerator) GenerateField(strength float64, direction Vector3D) *GravitationalField {
	return &GravitationalField{
		Strength:  strength,
		Direction: direction,
		Range:     g.calculateRange(strength),
		Stability: g.calculateStability(strength, direction),
	}
}

// calculateRange computes the effective range of a gravitational field
func (g *GravitationalGenerator) calculateRange(strength float64) float64 {
	return math.Sqrt(strength) * 10.0
}

// calculateStability determines the stability coefficient of the field
func (g *GravitationalGenerator) calculateStability(strength float64, direction Vector3D) float64 {
	magnitude := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	return strength / (1.0 + magnitude)
}
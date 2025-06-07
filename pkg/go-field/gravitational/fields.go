// Package gravitational implements gravitational field operations
package gravitational

import "math"

// Field represents a gravitational field
type Field struct {
	Strength   float64
	Direction  Vector3D
	Range      float64
	Eigenvalue complex128
}

// Vector3D represents a 3D vector
type Vector3D struct {
	X, Y, Z float64
}

// FieldGenerator creates gravitational fields
type FieldGenerator struct {
	BaseStrength float64
	MaxRange     float64
}

// GenerateField creates a new gravitational field
func (fg *FieldGenerator) GenerateField(position Vector3D, mass float64) *Field {
	strength := fg.BaseStrength * mass
	direction := fg.calculateDirection(position)
	fieldRange := fg.calculateRange(strength)
	
	return &Field{
		Strength:  strength,
		Direction: direction,
		Range:     fieldRange,
		Eigenvalue: complex(strength, 0),
	}
}

// calculateDirection determines field direction from position
func (fg *FieldGenerator) calculateDirection(pos Vector3D) Vector3D {
	magnitude := math.Sqrt(pos.X*pos.X + pos.Y*pos.Y + pos.Z*pos.Z)
	if magnitude == 0 {
		return Vector3D{0, 0, 1}
	}
	return Vector3D{pos.X / magnitude, pos.Y / magnitude, pos.Z / magnitude}
}

// calculateRange computes effective field range
func (fg *FieldGenerator) calculateRange(strength float64) float64 {
	return math.Min(fg.MaxRange, strength*10.0)
}
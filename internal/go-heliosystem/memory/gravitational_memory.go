package memory

import "math"

type GravitationalMemory struct {
	FieldStrength map[string]float64
	MemoryFields  map[string]GravitationalField
	Capacity      float64
	Decay         float64
}

type GravitationalField struct {
	Position  Vector3D
	Strength  float64
	Direction Vector3D
	Data      []byte
}

type Vector3D struct {
	X, Y, Z float64
}

func NewGravitationalMemory(capacity, decay float64) *GravitationalMemory {
	return &GravitationalMemory{
		FieldStrength: make(map[string]float64),
		MemoryFields:  make(map[string]GravitationalField),
		Capacity:      capacity,
		Decay:         decay,
	}
}

func (gm *GravitationalMemory) StoreInField(id string, data []byte, position Vector3D) {
	strength := math.Min(gm.Capacity, float64(len(data))/1024.0)
	
	field := GravitationalField{
		Position:  position,
		Strength:  strength,
		Direction: gm.calculateDirection(position),
		Data:      data,
	}
	
	gm.MemoryFields[id] = field
	gm.FieldStrength[id] = strength
}

func (gm *GravitationalMemory) calculateDirection(pos Vector3D) Vector3D {
	magnitude := math.Sqrt(pos.X*pos.X + pos.Y*pos.Y + pos.Z*pos.Z)
	if magnitude == 0 {
		return Vector3D{0, 0, 1}
	}
	return Vector3D{pos.X / magnitude, pos.Y / magnitude, pos.Z / magnitude}
}

func (gm *GravitationalMemory) ApplyDecay() {
	for id, strength := range gm.FieldStrength {
		newStrength := strength * (1.0 - gm.Decay)
		if newStrength < 0.01 {
			delete(gm.FieldStrength, id)
			delete(gm.MemoryFields, id)
		} else {
			gm.FieldStrength[id] = newStrength
		}
	}
}

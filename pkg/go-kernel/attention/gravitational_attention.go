package attention

import "math"

type GravitationalAttention struct {
	Masses    map[string]float64
	Positions map[string]Vector3D
	G         float64
}

type Vector3D struct {
	X, Y, Z float64
}

func NewGravitationalAttention() *GravitationalAttention {
	return &GravitationalAttention{
		Masses:    make(map[string]float64),
		Positions: make(map[string]Vector3D),
		G:         6.67430e-11,
	}
}

func (ga *GravitationalAttention) ComputeAttention(query, key, value []float64) []float64 {
	attention := make([]float64, len(value))
	
	for i := range attention {
		force := ga.calculateGravitationalForce(query, key, i)
		weight := ga.normalizeForce(force)
		attention[i] = weight * value[i%len(value)]
	}
	
	return attention
}

func (ga *GravitationalAttention) calculateGravitationalForce(query, key []float64, index int) float64 {
	if index >= len(query) || index >= len(key) {
		return 0.0
	}
	
	mass1 := math.Abs(query[index]) + 1e-6
	mass2 := math.Abs(key[index]) + 1e-6
	distance := math.Abs(query[index]-key[index]) + 1e-6
	
	return ga.G * mass1 * mass2 / (distance * distance)
}

func (ga *GravitationalAttention) normalizeForce(force float64) float64 {
	return 1.0 / (1.0 + math.Exp(-force*1e10))
}

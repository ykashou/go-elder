package entropy

import "math"

type EntropyDynamics struct {
	CurrentEntropy float64
	EntropyHistory []float64
	EvolutionRate  float64
	MaxEntropy     float64
}

func NewEntropyDynamics(maxEntropy, rate float64) *EntropyDynamics {
	return &EntropyDynamics{
		EntropyHistory: make([]float64, 0),
		EvolutionRate:  rate,
		MaxEntropy:     maxEntropy,
	}
}

func (ed *EntropyDynamics) EvolveEntropy(deltaTime float64) {
	// Entropy evolution based on thermodynamic principles
	deltaEntropy := ed.EvolutionRate * deltaTime * (ed.MaxEntropy - ed.CurrentEntropy)
	ed.CurrentEntropy += deltaEntropy
	
	if ed.CurrentEntropy > ed.MaxEntropy {
		ed.CurrentEntropy = ed.MaxEntropy
	}
	
	ed.EntropyHistory = append(ed.EntropyHistory, ed.CurrentEntropy)
}

func (ed *EntropyDynamics) CalculateEntropyProduction() float64 {
	if len(ed.EntropyHistory) < 2 {
		return 0.0
	}
	
	current := ed.EntropyHistory[len(ed.EntropyHistory)-1]
	previous := ed.EntropyHistory[len(ed.EntropyHistory)-2]
	return current - previous
}

func (ed *EntropyDynamics) GetEntropyGradient() float64 {
	if len(ed.EntropyHistory) < 3 {
		return 0.0
	}
	
	n := len(ed.EntropyHistory)
	recent := ed.EntropyHistory[n-1]
	older := ed.EntropyHistory[n-3]
	return (recent - older) / 2.0
}

func (ed *EntropyDynamics) CalculateSystemOrder() float64 {
	if ed.MaxEntropy == 0 {
		return 1.0
	}
	return 1.0 - (ed.CurrentEntropy / ed.MaxEntropy)
}

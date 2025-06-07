package attention

import "math"

type PhaseAttention struct {
	PhaseShifts map[string]float64
	Frequencies map[string]float64
	Amplitudes  map[string]float64
}

func NewPhaseAttention() *PhaseAttention {
	return &PhaseAttention{
		PhaseShifts: make(map[string]float64),
		Frequencies: make(map[string]float64),
		Amplitudes:  make(map[string]float64),
	}
}

func (pa *PhaseAttention) ComputeAttention(query, key, value []float64, phase float64) []float64 {
	attention := make([]float64, len(value))
	
	for i := range attention {
		queryPhase := pa.calculatePhase(query, i, phase)
		keyPhase := pa.calculatePhase(key, i, phase)
		
		phaseDiff := queryPhase - keyPhase
		weight := math.Cos(phaseDiff)
		
		attention[i] = weight * value[i%len(value)]
	}
	
	return attention
}

func (pa *PhaseAttention) calculatePhase(vector []float64, index int, basePhase float64) float64 {
	if index < len(vector) {
		return basePhase + vector[index]*0.1
	}
	return basePhase
}

func (pa *PhaseAttention) SynchronizePhases(phases []float64) []float64 {
	avgPhase := 0.0
	for _, phase := range phases {
		avgPhase += phase
	}
	avgPhase /= float64(len(phases))
	
	synchronized := make([]float64, len(phases))
	for i, phase := range phases {
		synchronized[i] = avgPhase + 0.1*(phase-avgPhase)
	}
	
	return synchronized
}

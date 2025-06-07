package attention

import "math"

type ResonanceAttention struct {
	Frequencies map[string]float64
	Amplitudes  map[string]float64
	Phases      map[string]float64
	Damping     float64
}

func NewResonanceAttention(damping float64) *ResonanceAttention {
	return &ResonanceAttention{
		Frequencies: make(map[string]float64),
		Amplitudes:  make(map[string]float64),
		Phases:      make(map[string]float64),
		Damping:     damping,
	}
}

func (ra *ResonanceAttention) ComputeAttention(query, key, value []float64, time float64) []float64 {
	attention := make([]float64, len(value))
	
	for i := range attention {
		resonance := ra.calculateResonance(query, key, i, time)
		attention[i] = resonance * value[i%len(value)]
	}
	
	return attention
}

func (ra *ResonanceAttention) calculateResonance(query, key []float64, index int, time float64) float64 {
	if index >= len(query) || index >= len(key) {
		return 0.0
	}
	
	freq1 := math.Abs(query[index]) + 0.1
	freq2 := math.Abs(key[index]) + 0.1
	
	freqDiff := math.Abs(freq1 - freq2)
	resonanceStrength := math.Exp(-freqDiff * ra.Damping)
	
	oscillation := math.Sin(freq1*time) * math.Sin(freq2*time)
	
	return resonanceStrength * math.Abs(oscillation)
}

func (ra *ResonanceAttention) FindResonantFrequencies(signal []float64) []float64 {
	frequencies := make([]float64, 0)
	
	for i := 1; i < len(signal)-1; i++ {
		if signal[i] > signal[i-1] && signal[i] > signal[i+1] {
			freq := float64(i) / float64(len(signal))
			frequencies = append(frequencies, freq)
		}
	}
	
	return frequencies
}

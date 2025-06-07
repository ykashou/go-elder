package optimization

import "math"

type ResonanceOptimizer struct {
	Frequency    float64
	Amplitude    float64
	Phase        float64
	Damping      float64
	Oscillators  map[string]Oscillator
}

type Oscillator struct {
	Frequency float64
	Amplitude float64
	Phase     float64
}

func NewResonanceOptimizer(freq, amp, phase, damping float64) *ResonanceOptimizer {
	return &ResonanceOptimizer{
		Frequency:   freq,
		Amplitude:   amp,
		Phase:       phase,
		Damping:     damping,
		Oscillators: make(map[string]Oscillator),
	}
}

func (ro *ResonanceOptimizer) OptimizeResonance(objective func([]float64) float64, params []float64) []float64 {
	optimized := make([]float64, len(params))
	copy(optimized, params)
	
	for iter := 0; iter < 100; iter++ {
		t := float64(iter) * 0.1
		resonanceSignal := ro.generateResonanceSignal(t)
		
		for i := range optimized {
			perturbation := resonanceSignal * 0.01
			candidate := optimized[i] + perturbation
			
			testParams := make([]float64, len(optimized))
			copy(testParams, optimized)
			testParams[i] = candidate
			
			if objective(testParams) < objective(optimized) {
				optimized[i] = candidate
			}
		}
		
		ro.updateResonance(t)
	}
	
	return optimized
}

func (ro *ResonanceOptimizer) generateResonanceSignal(t float64) float64 {
	signal := ro.Amplitude * math.Cos(ro.Frequency*t + ro.Phase)
	return signal * math.Exp(-ro.Damping*t)
}

func (ro *ResonanceOptimizer) updateResonance(t float64) {
	ro.Phase += ro.Frequency * 0.1
	ro.Amplitude *= math.Exp(-ro.Damping * 0.1)
}

func (ro *ResonanceOptimizer) AddOscillator(id string, freq, amp, phase float64) {
	ro.Oscillators[id] = Oscillator{
		Frequency: freq,
		Amplitude: amp,
		Phase:     phase,
	}
}

func (ro *ResonanceOptimizer) SynchronizeOscillators(targetFreq float64) {
	for id, osc := range ro.Oscillators {
		osc.Frequency += (targetFreq - osc.Frequency) * 0.1
		ro.Oscillators[id] = osc
	}
}

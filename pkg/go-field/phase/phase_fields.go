package phase

import "math"

type PhaseField struct {
	ID         string
	Phase      complex128
	Frequency  float64
	Amplitude  float64
	Coherence  float64
	Evolution  func(float64) complex128
}

type PhaseFieldSystem struct {
	Fields      map[string]PhaseField
	Couplings   map[string][]string
	GlobalPhase complex128
}

func NewPhaseFieldSystem() *PhaseFieldSystem {
	return &PhaseFieldSystem{
		Fields:    make(map[string]PhaseField),
		Couplings: make(map[string][]string),
		GlobalPhase: complex(0, 0),
	}
}

func (pfs *PhaseFieldSystem) AddField(id string, freq, amp float64, initialPhase complex128) {
	field := PhaseField{
		ID:        id,
		Phase:     initialPhase,
		Frequency: freq,
		Amplitude: amp,
		Coherence: 1.0,
		Evolution: func(t float64) complex128 {
			return complex(amp*math.Cos(freq*t), amp*math.Sin(freq*t))
		},
	}
	pfs.Fields[id] = field
}

func (pfs *PhaseFieldSystem) EvolveFields(deltaTime float64) {
	for id, field := range pfs.Fields {
		newPhase := field.Evolution(deltaTime)
		field.Phase = newPhase
		pfs.Fields[id] = field
	}
	
	pfs.updateGlobalPhase()
}

func (pfs *PhaseFieldSystem) updateGlobalPhase() {
	totalPhase := complex(0, 0)
	count := 0
	
	for _, field := range pfs.Fields {
		totalPhase += field.Phase
		count++
	}
	
	if count > 0 {
		pfs.GlobalPhase = totalPhase / complex(float64(count), 0)
	}
}

func (pfs *PhaseFieldSystem) CalculateCoherence() float64 {
	if len(pfs.Fields) < 2 {
		return 1.0
	}
	
	totalCoherence := 0.0
	pairs := 0
	
	fieldList := make([]PhaseField, 0)
	for _, field := range pfs.Fields {
		fieldList = append(fieldList, field)
	}
	
	for i := 0; i < len(fieldList); i++ {
		for j := i + 1; j < len(fieldList); j++ {
			coherence := pfs.calculatePairCoherence(fieldList[i], fieldList[j])
			totalCoherence += coherence
			pairs++
		}
	}
	
	return totalCoherence / float64(pairs)
}

func (pfs *PhaseFieldSystem) calculatePairCoherence(field1, field2 PhaseField) float64 {
	phaseDiff := field1.Phase - field2.Phase
	return math.Abs(real(phaseDiff))
}

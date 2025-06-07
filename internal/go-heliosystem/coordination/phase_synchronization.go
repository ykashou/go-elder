package coordination

import "math"

type PhaseSynchronizer struct {
	Phases map[string]float64
	FrequencyRange [2]float64
	SyncThreshold float64
}

func NewPhaseSynchronizer(threshold float64) *PhaseSynchronizer {
	return &PhaseSynchronizer{
		Phases: make(map[string]float64),
		FrequencyRange: [2]float64{0.1, 10.0},
		SyncThreshold: threshold,
	}
}

func (ps *PhaseSynchronizer) RegisterPhase(entityID string, phase float64) {
	ps.Phases[entityID] = phase
}

func (ps *PhaseSynchronizer) SynchronizePhases() bool {
	if len(ps.Phases) < 2 {
		return true
	}
	
	avgPhase := ps.calculateAveragePhase()
	synchronized := true
	
	for entityID, phase := range ps.Phases {
		diff := math.Abs(phase - avgPhase)
		if diff > ps.SyncThreshold {
			synchronized = false
			ps.Phases[entityID] = avgPhase + (phase-avgPhase)*0.1
		}
	}
	
	return synchronized
}

func (ps *PhaseSynchronizer) calculateAveragePhase() float64 {
	sum := 0.0
	for _, phase := range ps.Phases {
		sum += phase
	}
	return sum / float64(len(ps.Phases))
}

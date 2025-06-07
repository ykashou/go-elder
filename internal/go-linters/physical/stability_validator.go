package physical

import "math"

type StabilityValidator struct {
	Systems     map[string]PhysicalSystem
	Tolerance   float64
	TimeWindow  float64
	Perturbations map[string]Perturbation
}

type PhysicalSystem struct {
	ID          string
	State       []float64
	Parameters  map[string]float64
	Dynamics    func([]float64, map[string]float64) []float64
	EnergyFunc  func([]float64) float64
}

type Perturbation struct {
	Magnitude float64
	Direction []float64
	Duration  float64
}

func NewStabilityValidator(tolerance, timeWindow float64) *StabilityValidator {
	return &StabilityValidator{
		Systems:       make(map[string]PhysicalSystem),
		Tolerance:     tolerance,
		TimeWindow:    timeWindow,
		Perturbations: make(map[string]Perturbation),
	}
}

func (sv *StabilityValidator) AddSystem(id string, initialState []float64, params map[string]float64) {
	system := PhysicalSystem{
		ID:         id,
		State:      make([]float64, len(initialState)),
		Parameters: make(map[string]float64),
	}
	copy(system.State, initialState)
	for k, v := range params {
		system.Parameters[k] = v
	}
	sv.Systems[id] = system
}

func (sv *StabilityValidator) ValidateStability(systemID string) StabilityResult {
	system := sv.Systems[systemID]
	
	result := StabilityResult{
		SystemID:      systemID,
		Stable:        true,
		StabilityType: "unknown",
		LyapunovExponent: 0.0,
		Violations:    make([]string, 0),
	}
	
	result.LyapunovExponent = sv.calculateLyapunovExponent(system)
	
	if result.LyapunovExponent > 0 {
		result.Stable = false
		result.StabilityType = "unstable"
		result.Violations = append(result.Violations, "Positive Lyapunov exponent indicates chaos")
	} else if result.LyapunovExponent < -sv.Tolerance {
		result.StabilityType = "asymptotically_stable"
	} else {
		result.StabilityType = "marginally_stable"
	}
	
	if !sv.checkBoundedness(system) {
		result.Stable = false
		result.Violations = append(result.Violations, "System trajectories are unbounded")
	}
	
	return result
}

type StabilityResult struct {
	SystemID         string
	Stable           bool
	StabilityType    string
	LyapunovExponent float64
	Violations       []string
}

func (sv *StabilityValidator) calculateLyapunovExponent(system PhysicalSystem) float64 {
	initialSeparation := 1e-8
	timeSteps := 100
	dt := sv.TimeWindow / float64(timeSteps)
	
	state1 := make([]float64, len(system.State))
	state2 := make([]float64, len(system.State))
	copy(state1, system.State)
	copy(state2, system.State)
	
	if len(state2) > 0 {
		state2[0] += initialSeparation
	}
	
	totalLyapunov := 0.0
	
	for t := 0; t < timeSteps; t++ {
		state1 = sv.evolveState(state1, system.Parameters, dt)
		state2 = sv.evolveState(state2, system.Parameters, dt)
		
		separation := sv.calculateSeparation(state1, state2)
		if separation > 0 {
			lyapunov := math.Log(separation / initialSeparation) / (float64(t+1) * dt)
			totalLyapunov += lyapunov
		}
		
		separation = math.Max(separation, 1e-12)
		factor := initialSeparation / separation
		for i := range state2 {
			state2[i] = state1[i] + (state2[i]-state1[i])*factor
		}
	}
	
	return totalLyapunov / float64(timeSteps)
}

func (sv *StabilityValidator) evolveState(state []float64, params map[string]float64, dt float64) []float64 {
	newState := make([]float64, len(state))
	for i := range state {
		newState[i] = state[i] + dt*0.1*(params["damping"]*state[i])
	}
	return newState
}

func (sv *StabilityValidator) calculateSeparation(state1, state2 []float64) float64 {
	separation := 0.0
	for i := range state1 {
		diff := state1[i] - state2[i]
		separation += diff * diff
	}
	return math.Sqrt(separation)
}

func (sv *StabilityValidator) checkBoundedness(system PhysicalSystem) bool {
	timeSteps := 1000
	dt := sv.TimeWindow / float64(timeSteps)
	state := make([]float64, len(system.State))
	copy(state, system.State)
	
	for t := 0; t < timeSteps; t++ {
		state = sv.evolveState(state, system.Parameters, dt)
		
		magnitude := 0.0
		for _, x := range state {
			magnitude += x * x
		}
		magnitude = math.Sqrt(magnitude)
		
		if magnitude > 1e6 {
			return false
		}
	}
	
	return true
}

package physical

import "math"

type ResonanceChecker struct {
	Oscillators map[string]Oscillator
	Tolerance   float64
	TimeWindow  float64
}

type Oscillator struct {
	ID        string
	Frequency float64
	Amplitude float64
	Phase     float64
	Damping   float64
	Driving   *DrivingForce
}

type DrivingForce struct {
	Frequency float64
	Amplitude float64
	Phase     float64
}

func NewResonanceChecker(tolerance, timeWindow float64) *ResonanceChecker {
	return &ResonanceChecker{
		Oscillators: make(map[string]Oscillator),
		Tolerance:   tolerance,
		TimeWindow:  timeWindow,
	}
}

func (rc *ResonanceChecker) AddOscillator(id string, freq, amp, phase, damping float64, driving *DrivingForce) {
	oscillator := Oscillator{
		ID:        id,
		Frequency: freq,
		Amplitude: amp,
		Phase:     phase,
		Damping:   damping,
		Driving:   driving,
	}
	rc.Oscillators[id] = oscillator
}

func (rc *ResonanceChecker) CheckResonances() map[string]ResonanceResult {
	results := make(map[string]ResonanceResult)
	
	for id, oscillator := range rc.Oscillators {
		results[id] = rc.checkOscillatorResonance(oscillator)
	}
	
	return results
}

type ResonanceResult struct {
	OscillatorID   string
	InResonance    bool
	ResonanceType  string
	QualityFactor  float64
	ResonantFreq   float64
	AmplificationRatio float64
	Violations     []string
}

func (rc *ResonanceChecker) checkOscillatorResonance(osc Oscillator) ResonanceResult {
	result := ResonanceResult{
		OscillatorID: osc.ID,
		InResonance:  false,
		ResonanceType: "none",
		Violations:   make([]string, 0),
	}
	
	result.QualityFactor = rc.calculateQualityFactor(osc)
	result.ResonantFreq = rc.findResonantFrequency(osc)
	
	if osc.Driving != nil {
		result.AmplificationRatio = rc.calculateAmplificationRatio(osc)
		
		freqDiff := math.Abs(osc.Driving.Frequency - result.ResonantFreq)
		if freqDiff < rc.Tolerance {
			result.InResonance = true
			result.ResonanceType = "driven_resonance"
		}
		
		if result.AmplificationRatio > 10.0 {
			result.InResonance = true
			if result.ResonanceType == "none" {
				result.ResonanceType = "amplitude_resonance"
			}
		}
	}
	
	if rc.checkParametricResonance(osc) {
		result.InResonance = true
		result.ResonanceType = "parametric_resonance"
	}
	
	if result.QualityFactor > 100 && result.InResonance {
		result.Violations = append(result.Violations, "High Q-factor may cause system instability")
	}
	
	if result.AmplificationRatio > 50 {
		result.Violations = append(result.Violations, "Excessive amplification detected")
	}
	
	return result
}

func (rc *ResonanceChecker) calculateQualityFactor(osc Oscillator) float64 {
	if osc.Damping > 0 {
		return osc.Frequency / (2 * osc.Damping)
	}
	return math.Inf(1)
}

func (rc *ResonanceChecker) findResonantFrequency(osc Oscillator) float64 {
	naturalFreq := osc.Frequency
	dampingCorrection := math.Sqrt(1 - osc.Damping*osc.Damping/(4*naturalFreq*naturalFreq))
	return naturalFreq * dampingCorrection
}

func (rc *ResonanceChecker) calculateAmplificationRatio(osc Oscillator) float64 {
	if osc.Driving == nil {
		return 1.0
	}
	
	omega := osc.Driving.Frequency
	omega0 := osc.Frequency
	gamma := osc.Damping
	
	denominator := math.Sqrt(math.Pow(omega0*omega0-omega*omega, 2) + math.Pow(2*gamma*omega, 2))
	if denominator > 0 {
		return osc.Driving.Amplitude / denominator
	}
	
	return math.Inf(1)
}

func (rc *ResonanceChecker) checkParametricResonance(osc Oscillator) bool {
	if osc.Driving == nil {
		return false
	}
	
	parametricFreq := 2 * osc.Frequency
	freqDiff := math.Abs(osc.Driving.Frequency - parametricFreq)
	
	return freqDiff < rc.Tolerance
}

func (rc *ResonanceChecker) AnalyzeCoupledOscillators() []CoupledResonanceResult {
	results := make([]CoupledResonanceResult, 0)
	
	oscillatorList := make([]Oscillator, 0)
	for _, osc := range rc.Oscillators {
		oscillatorList = append(oscillatorList, osc)
	}
	
	for i := 0; i < len(oscillatorList); i++ {
		for j := i + 1; j < len(oscillatorList); j++ {
			result := rc.checkCoupledResonance(oscillatorList[i], oscillatorList[j])
			if result.InResonance {
				results = append(results, result)
			}
		}
	}
	
	return results
}

type CoupledResonanceResult struct {
	Oscillator1   string
	Oscillator2   string
	InResonance   bool
	CouplingType  string
	BeatFrequency float64
	PhaseRelation string
}

func (rc *ResonanceChecker) checkCoupledResonance(osc1, osc2 Oscillator) CoupledResonanceResult {
	result := CoupledResonanceResult{
		Oscillator1: osc1.ID,
		Oscillator2: osc2.ID,
		InResonance: false,
	}
	
	freqDiff := math.Abs(osc1.Frequency - osc2.Frequency)
	result.BeatFrequency = freqDiff
	
	if freqDiff < rc.Tolerance {
		result.InResonance = true
		result.CouplingType = "frequency_lock"
	}
	
	freqRatio := osc1.Frequency / osc2.Frequency
	if rc.isSimpleRatio(freqRatio) {
		result.InResonance = true
		result.CouplingType = "harmonic_resonance"
	}
	
	phaseDiff := math.Abs(osc1.Phase - osc2.Phase)
	if phaseDiff < rc.Tolerance || math.Abs(phaseDiff-math.Pi) < rc.Tolerance {
		result.PhaseRelation = "synchronized"
	} else {
		result.PhaseRelation = "unsynchronized"
	}
	
	return result
}

func (rc *ResonanceChecker) isSimpleRatio(ratio float64) bool {
	simpleRatios := []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0}
	
	for _, simple := range simpleRatios {
		if math.Abs(ratio-simple) < rc.Tolerance {
			return true
		}
	}
	
	return false
}

func (rc *ResonanceChecker) SimulateResonanceEvolution(oscillatorID string, duration float64) []ResonanceSnapshot {
	oscillator := rc.Oscillators[oscillatorID]
	snapshots := make([]ResonanceSnapshot, 0)
	
	timeSteps := 1000
	dt := duration / float64(timeSteps)
	
	for t := 0; t < timeSteps; t++ {
		time := float64(t) * dt
		amplitude := rc.calculateAmplitudeAtTime(oscillator, time)
		phase := rc.calculatePhaseAtTime(oscillator, time)
		
		snapshot := ResonanceSnapshot{
			Time:      time,
			Amplitude: amplitude,
			Phase:     phase,
			Energy:    0.5 * amplitude * amplitude,
		}
		snapshots = append(snapshots, snapshot)
	}
	
	return snapshots
}

type ResonanceSnapshot struct {
	Time      float64
	Amplitude float64
	Phase     float64
	Energy    float64
}

func (rc *ResonanceChecker) calculateAmplitudeAtTime(osc Oscillator, t float64) float64 {
	if osc.Driving == nil {
		return osc.Amplitude * math.Exp(-osc.Damping*t) * math.Cos(osc.Frequency*t+osc.Phase)
	}
	
	transient := osc.Amplitude * math.Exp(-osc.Damping*t) * math.Cos(osc.Frequency*t+osc.Phase)
	steadyState := rc.calculateAmplificationRatio(osc) * math.Cos(osc.Driving.Frequency*t+osc.Driving.Phase)
	
	return transient + steadyState
}

func (rc *ResonanceChecker) calculatePhaseAtTime(osc Oscillator, t float64) float64 {
	if osc.Driving == nil {
		return osc.Phase + osc.Frequency*t
	}
	
	phaseShift := math.Atan2(2*osc.Damping*osc.Driving.Frequency, 
		                     osc.Frequency*osc.Frequency-osc.Driving.Frequency*osc.Driving.Frequency)
	
	return osc.Driving.Phase + osc.Driving.Frequency*t - phaseShift
}

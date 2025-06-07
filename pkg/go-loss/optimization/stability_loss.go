package optimization

import "math"

type StabilityLoss struct {
	LyapunovThreshold float64
	EnergyBound       float64
	PhaseSpace        [][]float64
	EnergyHistory     []float64
}

func NewStabilityLoss(lyapunovThreshold, energyBound float64) *StabilityLoss {
	return &StabilityLoss{
		LyapunovThreshold: lyapunovThreshold,
		EnergyBound:       energyBound,
		PhaseSpace:        make([][]float64, 0),
		EnergyHistory:     make([]float64, 0),
	}
}

func (sl *StabilityLoss) ComputeLoss(state, velocity []float64) float64 {
	sl.updatePhaseSpace(state, velocity)
	
	lyapunovLoss := sl.computeLyapunovLoss()
	energyLoss := sl.computeEnergyLoss()
	phaseLoss := sl.computePhaseStabilityLoss()
	
	return lyapunovLoss + energyLoss + phaseLoss
}

func (sl *StabilityLoss) updatePhaseSpace(state, velocity []float64) {
	phasePoint := make([]float64, len(state)+len(velocity))
	copy(phasePoint[:len(state)], state)
	copy(phasePoint[len(state):], velocity)
	
	sl.PhaseSpace = append(sl.PhaseSpace, phasePoint)
	
	energy := sl.computeEnergy(state, velocity)
	sl.EnergyHistory = append(sl.EnergyHistory, energy)
	
	maxHistory := 1000
	if len(sl.PhaseSpace) > maxHistory {
		sl.PhaseSpace = sl.PhaseSpace[1:]
	}
	if len(sl.EnergyHistory) > maxHistory {
		sl.EnergyHistory = sl.EnergyHistory[1:]
	}
}

func (sl *StabilityLoss) computeLyapunovLoss() float64 {
	if len(sl.PhaseSpace) < 10 {
		return 0.0
	}
	
	lyapunovExponent := sl.estimateLyapunovExponent()
	
	if lyapunovExponent > sl.LyapunovThreshold {
		excess := lyapunovExponent - sl.LyapunovThreshold
		return excess * excess
	}
	
	return 0.0
}

func (sl *StabilityLoss) estimateLyapunovExponent() float64 {
	if len(sl.PhaseSpace) < 2 {
		return 0.0
	}
	
	n := len(sl.PhaseSpace)
	divergenceSum := 0.0
	count := 0
	
	for i := 1; i < n; i++ {
		distance := sl.computeDistance(sl.PhaseSpace[i], sl.PhaseSpace[i-1])
		if distance > 0 {
			divergenceSum += math.Log(distance)
			count++
		}
	}
	
	if count == 0 {
		return 0.0
	}
	
	return divergenceSum / float64(count)
}

func (sl *StabilityLoss) computeEnergyLoss() float64 {
	if len(sl.EnergyHistory) == 0 {
		return 0.0
	}
	
	currentEnergy := sl.EnergyHistory[len(sl.EnergyHistory)-1]
	
	if currentEnergy > sl.EnergyBound {
		excess := currentEnergy - sl.EnergyBound
		return excess * excess
	}
	
	if len(sl.EnergyHistory) > 1 {
		energyVariation := sl.computeEnergyVariation()
		return energyVariation
	}
	
	return 0.0
}

func (sl *StabilityLoss) computeEnergyVariation() float64 {
	if len(sl.EnergyHistory) < 2 {
		return 0.0
	}
	
	variation := 0.0
	for i := 1; i < len(sl.EnergyHistory); i++ {
		diff := sl.EnergyHistory[i] - sl.EnergyHistory[i-1]
		variation += diff * diff
	}
	
	return variation / float64(len(sl.EnergyHistory)-1)
}

func (sl *StabilityLoss) computePhaseStabilityLoss() float64 {
	if len(sl.PhaseSpace) < 3 {
		return 0.0
	}
	
	attractorRadius := sl.estimateAttractorRadius()
	currentRadius := sl.computeCurrentRadius()
	
	if currentRadius > attractorRadius*2 {
		excess := currentRadius - attractorRadius*2
		return excess * excess
	}
	
	return 0.0
}

func (sl *StabilityLoss) estimateAttractorRadius() float64 {
	if len(sl.PhaseSpace) < 2 {
		return 1.0
	}
	
	centroid := sl.computeCentroid()
	maxDistance := 0.0
	
	for _, point := range sl.PhaseSpace {
		distance := sl.computeDistance(point, centroid)
		if distance > maxDistance {
			maxDistance = distance
		}
	}
	
	return maxDistance
}

func (sl *StabilityLoss) computeCurrentRadius() float64 {
	if len(sl.PhaseSpace) == 0 {
		return 0.0
	}
	
	centroid := sl.computeCentroid()
	current := sl.PhaseSpace[len(sl.PhaseSpace)-1]
	
	return sl.computeDistance(current, centroid)
}

func (sl *StabilityLoss) computeCentroid() []float64 {
	if len(sl.PhaseSpace) == 0 {
		return []float64{}
	}
	
	dimension := len(sl.PhaseSpace[0])
	centroid := make([]float64, dimension)
	
	for _, point := range sl.PhaseSpace {
		for i, val := range point {
			if i < len(centroid) {
				centroid[i] += val
			}
		}
	}
	
	for i := range centroid {
		centroid[i] /= float64(len(sl.PhaseSpace))
	}
	
	return centroid
}

func (sl *StabilityLoss) computeDistance(point1, point2 []float64) float64 {
	distance := 0.0
	minLen := len(point1)
	if len(point2) < minLen {
		minLen = len(point2)
	}
	
	for i := 0; i < minLen; i++ {
		diff := point1[i] - point2[i]
		distance += diff * diff
	}
	
	return math.Sqrt(distance)
}

func (sl *StabilityLoss) computeEnergy(state, velocity []float64) float64 {
	kinetic := 0.0
	potential := 0.0
	
	for _, v := range velocity {
		kinetic += v * v
	}
	kinetic *= 0.5
	
	for _, s := range state {
		potential += s * s
	}
	potential *= 0.5
	
	return kinetic + potential
}

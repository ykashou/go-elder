package optimization

import "math"

type ConvergenceLoss struct {
	TargetRate     float64
	ToleranceWindow float64
	HistoryLength  int
	History        []float64
}

func NewConvergenceLoss(targetRate, tolerance float64, historyLen int) *ConvergenceLoss {
	return &ConvergenceLoss{
		TargetRate:      targetRate,
		ToleranceWindow: tolerance,
		HistoryLength:   historyLen,
		History:         make([]float64, 0),
	}
}

func (cl *ConvergenceLoss) ComputeLoss(currentValue float64, gradient []float64) float64 {
	cl.updateHistory(currentValue)
	
	convergenceRateLoss := cl.computeConvergenceRateLoss()
	stabilityLoss := cl.computeStabilityLoss()
	gradientLoss := cl.computeGradientLoss(gradient)
	
	return convergenceRateLoss + stabilityLoss + gradientLoss
}

func (cl *ConvergenceLoss) updateHistory(value float64) {
	cl.History = append(cl.History, value)
	
	if len(cl.History) > cl.HistoryLength {
		cl.History = cl.History[1:]
	}
}

func (cl *ConvergenceLoss) computeConvergenceRateLoss() float64 {
	if len(cl.History) < 3 {
		return 0.0
	}
	
	n := len(cl.History)
	recent := cl.History[n-1]
	previous := cl.History[n-2]
	earlier := cl.History[n-3]
	
	actualRate := math.Abs(recent-previous) / math.Abs(previous-earlier)
	rateDifference := math.Abs(actualRate - cl.TargetRate)
	
	return rateDifference * rateDifference
}

func (cl *ConvergenceLoss) computeStabilityLoss() float64 {
	if len(cl.History) < 2 {
		return 0.0
	}
	
	variance := 0.0
	mean := 0.0
	
	for _, val := range cl.History {
		mean += val
	}
	mean /= float64(len(cl.History))
	
	for _, val := range cl.History {
		diff := val - mean
		variance += diff * diff
	}
	variance /= float64(len(cl.History))
	
	if variance > cl.ToleranceWindow*cl.ToleranceWindow {
		return variance - cl.ToleranceWindow*cl.ToleranceWindow
	}
	
	return 0.0
}

func (cl *ConvergenceLoss) computeGradientLoss(gradient []float64) float64 {
	gradientNorm := 0.0
	for _, g := range gradient {
		gradientNorm += g * g
	}
	gradientNorm = math.Sqrt(gradientNorm)
	
	if gradientNorm > 1.0 {
		return (gradientNorm - 1.0) * (gradientNorm - 1.0)
	}
	
	return 0.0
}

func (cl *ConvergenceLoss) IsConverged() bool {
	if len(cl.History) < 2 {
		return false
	}
	
	n := len(cl.History)
	recent := cl.History[n-1]
	previous := cl.History[n-2]
	
	return math.Abs(recent-previous) < cl.ToleranceWindow
}

func (cl *ConvergenceLoss) GetConvergenceRate() float64 {
	if len(cl.History) < 3 {
		return 0.0
	}
	
	n := len(cl.History)
	recent := cl.History[n-1]
	previous := cl.History[n-2]
	earlier := cl.History[n-3]
	
	if math.Abs(previous-earlier) < 1e-12 {
		return 0.0
	}
	
	return math.Abs(recent-previous) / math.Abs(previous-earlier)
}

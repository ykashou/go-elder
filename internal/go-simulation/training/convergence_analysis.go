package training

import "math"

type ConvergenceAnalyzer struct {
	LossHistory     []float64
	AccuracyHistory []float64
	LearningCurve   []float64
	Threshold       float64
	WindowSize      int
}

func NewConvergenceAnalyzer(threshold float64, windowSize int) *ConvergenceAnalyzer {
	return &ConvergenceAnalyzer{
		LossHistory:     make([]float64, 0),
		AccuracyHistory: make([]float64, 0),
		LearningCurve:   make([]float64, 0),
		Threshold:       threshold,
		WindowSize:      windowSize,
	}
}

func (ca *ConvergenceAnalyzer) RecordMetrics(loss, accuracy float64) {
	ca.LossHistory = append(ca.LossHistory, loss)
	ca.AccuracyHistory = append(ca.AccuracyHistory, accuracy)
	ca.updateLearningCurve()
}

func (ca *ConvergenceAnalyzer) updateLearningCurve() {
	if len(ca.LossHistory) > 0 {
		recentLoss := ca.LossHistory[len(ca.LossHistory)-1]
		ca.LearningCurve = append(ca.LearningCurve, recentLoss)
	}
}

func (ca *ConvergenceAnalyzer) CheckConvergence() bool {
	if len(ca.LossHistory) < ca.WindowSize {
		return false
	}
	
	recentLosses := ca.LossHistory[len(ca.LossHistory)-ca.WindowSize:]
	variance := ca.calculateVariance(recentLosses)
	
	return variance < ca.Threshold
}

func (ca *ConvergenceAnalyzer) calculateVariance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	mean := ca.calculateMean(values)
	variance := 0.0
	
	for _, value := range values {
		diff := value - mean
		variance += diff * diff
	}
	
	return variance / float64(len(values))
}

func (ca *ConvergenceAnalyzer) calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	
	return sum / float64(len(values))
}

func (ca *ConvergenceAnalyzer) PredictConvergenceTime() int {
	if len(ca.LossHistory) < 2 {
		return -1
	}
	
	recentGradient := ca.calculateGradient()
	if recentGradient >= 0 {
		return -1
	}
	
	currentLoss := ca.LossHistory[len(ca.LossHistory)-1]
	stepsToConvergence := int(math.Ceil(currentLoss / math.Abs(recentGradient)))
	
	return stepsToConvergence
}

func (ca *ConvergenceAnalyzer) calculateGradient() float64 {
	if len(ca.LossHistory) < 2 {
		return 0
	}
	
	current := ca.LossHistory[len(ca.LossHistory)-1]
	previous := ca.LossHistory[len(ca.LossHistory)-2]
	
	return current - previous
}

func (ca *ConvergenceAnalyzer) GetConvergenceReport() ConvergenceReport {
	return ConvergenceReport{
		Converged:         ca.CheckConvergence(),
		CurrentLoss:       ca.getCurrentLoss(),
		CurrentAccuracy:   ca.getCurrentAccuracy(),
		LossVariance:      ca.getLossVariance(),
		PredictedSteps:    ca.PredictConvergenceTime(),
		TotalEpochs:       len(ca.LossHistory),
	}
}

type ConvergenceReport struct {
	Converged       bool
	CurrentLoss     float64
	CurrentAccuracy float64
	LossVariance    float64
	PredictedSteps  int
	TotalEpochs     int
}

func (ca *ConvergenceAnalyzer) getCurrentLoss() float64 {
	if len(ca.LossHistory) > 0 {
		return ca.LossHistory[len(ca.LossHistory)-1]
	}
	return 0
}

func (ca *ConvergenceAnalyzer) getCurrentAccuracy() float64 {
	if len(ca.AccuracyHistory) > 0 {
		return ca.AccuracyHistory[len(ca.AccuracyHistory)-1]
	}
	return 0
}

func (ca *ConvergenceAnalyzer) getLossVariance() float64 {
	windowSize := ca.WindowSize
	if len(ca.LossHistory) < windowSize {
		windowSize = len(ca.LossHistory)
	}
	
	if windowSize > 0 {
		recent := ca.LossHistory[len(ca.LossHistory)-windowSize:]
		return ca.calculateVariance(recent)
	}
	
	return 0
}

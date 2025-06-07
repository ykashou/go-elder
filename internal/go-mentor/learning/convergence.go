package learning

type ConvergenceAnalyzer struct {
	LossHistory []float64
	Threshold   float64
}

func (ca *ConvergenceAnalyzer) CheckConvergence() bool {
	if len(ca.LossHistory) < 2 {
		return false
	}
	recent := ca.LossHistory[len(ca.LossHistory)-1]
	previous := ca.LossHistory[len(ca.LossHistory)-2]
	return (previous-recent)/previous < ca.Threshold
}

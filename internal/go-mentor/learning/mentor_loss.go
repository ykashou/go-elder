package learning

type MentorLossFunction struct {
	LossType   string
	Parameters map[string]float64
}

func (mlf *MentorLossFunction) CalculateLoss(predicted, actual []float64) float64 {
	var loss float64
	for i := range predicted {
		diff := predicted[i] - actual[i]
		loss += diff * diff
	}
	return loss / float64(len(predicted))
}

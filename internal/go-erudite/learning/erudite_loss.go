// Package learning implements erudite-specific loss functions
package learning

// EruditeLossFunction handles loss calculation for erudite entities
type EruditeLossFunction struct {
	LossType     string
	Regularization float64
	WeightDecay  float64
}

// CalculateTaskLoss computes task-specific loss
func (elf *EruditeLossFunction) CalculateTaskLoss(predicted, actual []float64) float64 {
	var loss float64
	for i := range predicted {
		diff := predicted[i] - actual[i]
		loss += diff * diff
	}
	return loss / float64(len(predicted))
}

// CalculateRegularizedLoss adds regularization to the loss
func (elf *EruditeLossFunction) CalculateRegularizedLoss(baseLoss float64, weights []float64) float64 {
	var regularization float64
	for _, weight := range weights {
		regularization += weight * weight
	}
	return baseLoss + elf.Regularization*regularization
}
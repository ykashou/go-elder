package elder

import "math"

type ElderLossFunction struct {
	Type        string
	Weight      float64
	Regularization float64
}

func (elf *ElderLossFunction) ComputeLoss(predicted, actual []float64) float64 {
	var loss float64
	switch elf.Type {
	case "mse":
		loss = elf.meanSquaredError(predicted, actual)
	case "mae":
		loss = elf.meanAbsoluteError(predicted, actual)
	default:
		loss = elf.meanSquaredError(predicted, actual)
	}
	return loss * elf.Weight
}

func (elf *ElderLossFunction) meanSquaredError(predicted, actual []float64) float64 {
	var sum float64
	for i := range predicted {
		diff := predicted[i] - actual[i]
		sum += diff * diff
	}
	return sum / float64(len(predicted))
}

func (elf *ElderLossFunction) meanAbsoluteError(predicted, actual []float64) float64 {
	var sum float64
	for i := range predicted {
		sum += math.Abs(predicted[i] - actual[i])
	}
	return sum / float64(len(predicted))
}

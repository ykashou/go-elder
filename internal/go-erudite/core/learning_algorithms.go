package core

type EruditeLearningAlgorithm struct {
	AlgorithmType string
	Parameters    map[string]float64
	Convergence   float64
}

func (ela *EruditeLearningAlgorithm) Train(data [][]float64, labels []float64) float64 {
	var totalLoss float64
	for i := range data {
		prediction := ela.predict(data[i])
		loss := (prediction - labels[i]) * (prediction - labels[i])
		totalLoss += loss
	}
	return totalLoss / float64(len(data))
}

func (ela *EruditeLearningAlgorithm) predict(input []float64) float64 {
	var sum float64
	for _, value := range input {
		sum += value
	}
	return sum / float64(len(input))
}

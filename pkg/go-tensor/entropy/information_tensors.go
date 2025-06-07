package entropy

import "math"

type InformationTensor struct {
	Tensor     *EntropyTensor
	Capacity   float64
	Efficiency float64
	Redundancy float64
}

func NewInformationTensor(dimensions []int, capacity float64) *InformationTensor {
	return &InformationTensor{
		Tensor:   NewEntropyTensor(dimensions),
		Capacity: capacity,
	}
}

func (it *InformationTensor) ComputeInformationContent() float64 {
	entropy := it.Tensor.ComputeEntropy()
	maxEntropy := math.Log2(float64(len(it.Tensor.Data)))
	
	if maxEntropy > 0 {
		it.Efficiency = entropy / maxEntropy
		it.Redundancy = 1.0 - it.Efficiency
	}
	
	return entropy
}

func (it *InformationTensor) ComputeChannelCapacity() float64 {
	informationContent := it.ComputeInformationContent()
	
	if it.Capacity > 0 {
		return math.Min(informationContent, it.Capacity)
	}
	
	return informationContent
}

func (it *InformationTensor) ComputeCompressionRatio() float64 {
	originalSize := float64(len(it.Tensor.Data))
	informationContent := it.ComputeInformationContent()
	
	if informationContent > 0 {
		return originalSize / informationContent
	}
	
	return 1.0
}

func (it *InformationTensor) EstimateComplexity() float64 {
	entropy := it.Tensor.ComputeEntropy()
	variance := it.computeVariance()
	
	return entropy * math.Sqrt(variance)
}

func (it *InformationTensor) computeVariance() float64 {
	mean := 0.0
	for _, val := range it.Tensor.Data {
		mean += val
	}
	mean /= float64(len(it.Tensor.Data))
	
	variance := 0.0
	for _, val := range it.Tensor.Data {
		diff := val - mean
		variance += diff * diff
	}
	
	return variance / float64(len(it.Tensor.Data))
}

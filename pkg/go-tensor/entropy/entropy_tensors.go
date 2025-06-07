package entropy

import "math"

type EntropyTensor struct {
	Dimensions []int
	Data       []float64
	Rank       int
	Shape      []int
}

func NewEntropyTensor(dimensions []int) *EntropyTensor {
	size := 1
	for _, dim := range dimensions {
		size *= dim
	}
	
	return &EntropyTensor{
		Dimensions: make([]int, len(dimensions)),
		Data:       make([]float64, size),
		Rank:       len(dimensions),
		Shape:      make([]int, len(dimensions)),
	}
}

func (et *EntropyTensor) ComputeEntropy() float64 {
	total := 0.0
	for _, val := range et.Data {
		total += val
	}
	
	if total == 0 {
		return 0.0
	}
	
	entropy := 0.0
	for _, val := range et.Data {
		if val > 0 {
			prob := val / total
			entropy -= prob * math.Log2(prob)
		}
	}
	
	return entropy
}

func (et *EntropyTensor) ComputeConditionalEntropy(conditioningTensor *EntropyTensor) float64 {
	jointEntropy := et.computeJointEntropy(conditioningTensor)
	conditioningEntropy := conditioningTensor.ComputeEntropy()
	
	return jointEntropy - conditioningEntropy
}

func (et *EntropyTensor) computeJointEntropy(other *EntropyTensor) float64 {
	minSize := len(et.Data)
	if len(other.Data) < minSize {
		minSize = len(other.Data)
	}
	
	jointData := make([]float64, minSize)
	total := 0.0
	
	for i := 0; i < minSize; i++ {
		jointData[i] = et.Data[i] * other.Data[i]
		total += jointData[i]
	}
	
	if total == 0 {
		return 0.0
	}
	
	entropy := 0.0
	for _, val := range jointData {
		if val > 0 {
			prob := val / total
			entropy -= prob * math.Log2(prob)
		}
	}
	
	return entropy
}

func (et *EntropyTensor) ComputeMutualInformation(other *EntropyTensor) float64 {
	entropyA := et.ComputeEntropy()
	entropyB := other.ComputeEntropy()
	jointEntropy := et.computeJointEntropy(other)
	
	return entropyA + entropyB - jointEntropy
}

func (et *EntropyTensor) ComputeKLDivergence(reference *EntropyTensor) float64 {
	minSize := len(et.Data)
	if len(reference.Data) < minSize {
		minSize = len(reference.Data)
	}
	
	totalP := 0.0
	totalQ := 0.0
	
	for i := 0; i < minSize; i++ {
		totalP += et.Data[i]
		totalQ += reference.Data[i]
	}
	
	if totalP == 0 || totalQ == 0 {
		return math.Inf(1)
	}
	
	kl := 0.0
	for i := 0; i < minSize; i++ {
		p := et.Data[i] / totalP
		q := reference.Data[i] / totalQ
		
		if p > 0 && q > 0 {
			kl += p * math.Log2(p/q)
		} else if p > 0 && q == 0 {
			return math.Inf(1)
		}
	}
	
	return kl
}

func (et *EntropyTensor) Normalize() {
	total := 0.0
	for _, val := range et.Data {
		total += val
	}
	
	if total > 0 {
		for i := range et.Data {
			et.Data[i] /= total
		}
	}
}

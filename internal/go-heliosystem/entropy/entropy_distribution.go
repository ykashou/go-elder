package entropy

import "math"

type EntropyDistribution struct {
	HierarchyLevels map[int]float64
	TotalEntropy    float64
	Distribution    map[string]float64
	FlowRates       map[string]float64
}

func NewEntropyDistribution() *EntropyDistribution {
	return &EntropyDistribution{
		HierarchyLevels: make(map[int]float64),
		Distribution:    make(map[string]float64),
		FlowRates:       make(map[string]float64),
	}
}

func (ed *EntropyDistribution) SetLevelEntropy(level int, entropy float64) {
	ed.HierarchyLevels[level] = entropy
	ed.recalculateTotal()
}

func (ed *EntropyDistribution) recalculateTotal() {
	total := 0.0
	for _, entropy := range ed.HierarchyLevels {
		total += entropy
	}
	ed.TotalEntropy = total
}

func (ed *EntropyDistribution) DistributeEntropy() {
	if ed.TotalEntropy == 0 {
		return
	}
	
	for level, entropy := range ed.HierarchyLevels {
		proportion := entropy / ed.TotalEntropy
		ed.Distribution[fmt.Sprintf("level_%d", level)] = proportion
	}
}

func (ed *EntropyDistribution) CalculateInformationContent(level int) float64 {
	if entropy, exists := ed.HierarchyLevels[level]; exists {
		return -math.Log2(entropy + 1e-10)
	}
	return 0.0
}

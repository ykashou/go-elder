package gravitational

import "math"

type StabilityAnalyzer struct {
	Fields []Field
	Threshold float64
}

func (sa *StabilityAnalyzer) AnalyzeStability() float64 {
	var totalStability float64
	for _, field := range sa.Fields {
		stability := math.Min(1.0, field.Strength/field.Range)
		totalStability += stability
	}
	return totalStability / float64(len(sa.Fields))
}

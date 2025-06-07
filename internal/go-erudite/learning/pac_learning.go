// Package learning implements PAC learning bounds
package learning

import "math"

// PACLearningBounds implements PAC learning theory bounds
type PACLearningBounds struct {
	Confidence    float64
	Accuracy      float64
	VCDimension   int
	SampleSize    int
}

// CalculateSampleComplexity computes required sample size for PAC learning
func (plb *PACLearningBounds) CalculateSampleComplexity() int {
	epsilon := 1.0 - plb.Accuracy
	delta := 1.0 - plb.Confidence
	
	// PAC bound: m >= (1/epsilon) * (ln(|H|) + ln(1/delta))
	// Using VC dimension approximation
	vcTerm := float64(plb.VCDimension) * math.Log(2.0*math.E*float64(plb.SampleSize)/float64(plb.VCDimension))
	deltaLog := math.Log(1.0 / delta)
	
	return int(math.Ceil((vcTerm + deltaLog) / epsilon))
}

// ValidatePACBounds checks if current learning satisfies PAC bounds
func (plb *PACLearningBounds) ValidatePACBounds(empiricalError float64) bool {
	bound := plb.calculateGeneralizationBound()
	return empiricalError <= bound
}

func (plb *PACLearningBounds) calculateGeneralizationBound() float64 {
	vcTerm := float64(plb.VCDimension) / float64(plb.SampleSize)
	return math.Sqrt(vcTerm * math.Log(2.0/1.0-plb.Confidence))
}
package hierarchical

import "math"

type CrossLevelLoss struct {
	InformationFlow    float64
	HierarchyIntegrity float64
	CausalConsistency  float64
	TemporalCoherence  float64
}

func NewCrossLevelLoss() *CrossLevelLoss {
	return &CrossLevelLoss{
		InformationFlow:    1.0,
		HierarchyIntegrity: 0.9,
		CausalConsistency:  0.8,
		TemporalCoherence:  0.7,
	}
}

func (cll *CrossLevelLoss) ComputeLoss(elderState []float64, mentorStates, eruditeStates [][]float64) float64 {
	infoFlowLoss := cll.computeInformationFlowLoss(elderState, mentorStates, eruditeStates)
	integrityLoss := cll.computeHierarchyIntegrityLoss(elderState, mentorStates, eruditeStates)
	causalLoss := cll.computeCausalConsistencyLoss(elderState, mentorStates, eruditeStates)
	temporalLoss := cll.computeTemporalCoherenceLoss(elderState, mentorStates, eruditeStates)
	
	totalLoss := cll.InformationFlow*infoFlowLoss +
		cll.HierarchyIntegrity*integrityLoss +
		cll.CausalConsistency*causalLoss +
		cll.TemporalCoherence*temporalLoss
	
	return totalLoss
}

func (cll *CrossLevelLoss) computeInformationFlowLoss(elderState []float64, mentorStates, eruditeStates [][]float64) float64 {
	elderEntropy := cll.computeEntropy(elderState)
	
	mentorEntropy := 0.0
	for _, mentorState := range mentorStates {
		mentorEntropy += cll.computeEntropy(mentorState)
	}
	mentorEntropy /= float64(len(mentorStates))
	
	eruditeEntropy := 0.0
	for _, eruditeState := range eruditeStates {
		eruditeEntropy += cll.computeEntropy(eruditeState)
	}
	eruditeEntropy /= float64(len(eruditeStates))
	
	expectedFlow := elderEntropy - mentorEntropy - eruditeEntropy
	return math.Abs(expectedFlow)
}

func (cll *CrossLevelLoss) computeHierarchyIntegrityLoss(elderState []float64, mentorStates, eruditeStates [][]float64) float64 {
	loss := 0.0
	
	elderNorm := cll.computeNorm(elderState)
	
	for _, mentorState := range mentorStates {
		mentorNorm := cll.computeNorm(mentorState)
		if mentorNorm > elderNorm {
			loss += (mentorNorm - elderNorm) * (mentorNorm - elderNorm)
		}
	}
	
	for _, eruditeState := range eruditeStates {
		eruditeNorm := cll.computeNorm(eruditeState)
		avgMentorNorm := 0.0
		for _, mentorState := range mentorStates {
			avgMentorNorm += cll.computeNorm(mentorState)
		}
		avgMentorNorm /= float64(len(mentorStates))
		
		if eruditeNorm > avgMentorNorm {
			loss += (eruditeNorm - avgMentorNorm) * (eruditeNorm - avgMentorNorm)
		}
	}
	
	return loss
}

func (cll *CrossLevelLoss) computeCausalConsistencyLoss(elderState []float64, mentorStates, eruditeStates [][]float64) float64 {
	loss := 0.0
	
	for i, mentorState := range mentorStates {
		elderInfluence := cll.computeInfluence(elderState, mentorState)
		
		mentorInfluenceSum := 0.0
		for j, eruditeState := range eruditeStates {
			if j%len(mentorStates) == i {
				mentorInfluence := cll.computeInfluence(mentorState, eruditeState)
				mentorInfluenceSum += mentorInfluence
			}
		}
		
		influenceRatio := mentorInfluenceSum / elderInfluence
		if influenceRatio < 0.5 || influenceRatio > 2.0 {
			loss += math.Abs(influenceRatio - 1.0)
		}
	}
	
	return loss / float64(len(mentorStates))
}

func (cll *CrossLevelLoss) computeTemporalCoherenceLoss(elderState []float64, mentorStates, eruditeStates [][]float64) float64 {
	loss := 0.0
	
	elderVariance := cll.computeVariance(elderState)
	
	mentorVarianceSum := 0.0
	for _, mentorState := range mentorStates {
		mentorVarianceSum += cll.computeVariance(mentorState)
	}
	avgMentorVariance := mentorVarianceSum / float64(len(mentorStates))
	
	eruditeVarianceSum := 0.0
	for _, eruditeState := range eruditeStates {
		eruditeVarianceSum += cll.computeVariance(eruditeState)
	}
	avgEruditeVariance := eruditeVarianceSum / float64(len(eruditeStates))
	
	varianceRatio1 := avgMentorVariance / elderVariance
	varianceRatio2 := avgEruditeVariance / avgMentorVariance
	
	loss += math.Abs(varianceRatio1 - 1.0) + math.Abs(varianceRatio2 - 1.0)
	
	return loss
}

func (cll *CrossLevelLoss) computeEntropy(state []float64) float64 {
	entropy := 0.0
	for _, val := range state {
		if val > 0 {
			entropy -= val * math.Log2(val)
		}
	}
	return entropy
}

func (cll *CrossLevelLoss) computeNorm(state []float64) float64 {
	norm := 0.0
	for _, val := range state {
		norm += val * val
	}
	return math.Sqrt(norm)
}

func (cll *CrossLevelLoss) computeInfluence(source, target []float64) float64 {
	influence := 0.0
	minLen := len(source)
	if len(target) < minLen {
		minLen = len(target)
	}
	
	for i := 0; i < minLen; i++ {
		influence += source[i] * target[i]
	}
	
	return math.Abs(influence)
}

func (cll *CrossLevelLoss) computeVariance(state []float64) float64 {
	mean := 0.0
	for _, val := range state {
		mean += val
	}
	mean /= float64(len(state))
	
	variance := 0.0
	for _, val := range state {
		diff := val - mean
		variance += diff * diff
	}
	
	return variance / float64(len(state))
}

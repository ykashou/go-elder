package hierarchical

import "math"

type ElderMentorLoss struct {
	CoordinationWeight  float64
	AlignmentWeight     float64
	EfficiencyWeight    float64
	StabilityWeight     float64
	HierarchyLevels     int
}

func NewElderMentorLoss(levels int) *ElderMentorLoss {
	return &ElderMentorLoss{
		CoordinationWeight: 1.0,
		AlignmentWeight:    0.8,
		EfficiencyWeight:   0.6,
		StabilityWeight:    0.9,
		HierarchyLevels:    levels,
	}
}

func (eml *ElderMentorLoss) ComputeLoss(elderState, mentorStates [][]float64) float64 {
	coordinationLoss := eml.computeCoordinationLoss(elderState, mentorStates)
	alignmentLoss := eml.computeAlignmentLoss(elderState, mentorStates)
	efficiencyLoss := eml.computeEfficiencyLoss(mentorStates)
	stabilityLoss := eml.computeStabilityLoss(elderState)
	
	totalLoss := eml.CoordinationWeight*coordinationLoss +
		eml.AlignmentWeight*alignmentLoss +
		eml.EfficiencyWeight*efficiencyLoss +
		eml.StabilityWeight*stabilityLoss
	
	return totalLoss
}

func (eml *ElderMentorLoss) computeCoordinationLoss(elderState []float64, mentorStates [][]float64) float64 {
	loss := 0.0
	
	for _, mentorState := range mentorStates {
		deviation := eml.computeDeviation(elderState, mentorState)
		loss += deviation * deviation
	}
	
	return loss / float64(len(mentorStates))
}

func (eml *ElderMentorLoss) computeAlignmentLoss(elderState []float64, mentorStates [][]float64) float64 {
	loss := 0.0
	
	for i, mentorState := range mentorStates {
		for j := i + 1; j < len(mentorStates); j++ {
			otherMentor := mentorStates[j]
			misalignment := eml.computeDeviation(mentorState, otherMentor)
			loss += misalignment
		}
	}
	
	if len(mentorStates) > 1 {
		pairs := len(mentorStates) * (len(mentorStates) - 1) / 2
		loss /= float64(pairs)
	}
	
	return loss
}

func (eml *ElderMentorLoss) computeEfficiencyLoss(mentorStates [][]float64) float64 {
	loss := 0.0
	
	for _, mentorState := range mentorStates {
		energy := 0.0
		for _, val := range mentorState {
			energy += val * val
		}
		loss += energy
	}
	
	return loss / float64(len(mentorStates))
}

func (eml *ElderMentorLoss) computeStabilityLoss(elderState []float64) float64 {
	variance := 0.0
	mean := 0.0
	
	for _, val := range elderState {
		mean += val
	}
	mean /= float64(len(elderState))
	
	for _, val := range elderState {
		diff := val - mean
		variance += diff * diff
	}
	
	return math.Sqrt(variance / float64(len(elderState)))
}

func (eml *ElderMentorLoss) computeDeviation(state1, state2 []float64) float64 {
	deviation := 0.0
	minLen := len(state1)
	if len(state2) < minLen {
		minLen = len(state2)
	}
	
	for i := 0; i < minLen; i++ {
		diff := state1[i] - state2[i]
		deviation += diff * diff
	}
	
	return math.Sqrt(deviation)
}

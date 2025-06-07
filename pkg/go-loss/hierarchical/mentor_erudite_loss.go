package hierarchical

import "math"

type MentorEruditeLoss struct {
	SupervisionWeight   float64
	SpecializationWeight float64
	ConvergenceWeight   float64
	DiversityWeight     float64
}

func NewMentorEruditeLoss() *MentorEruditeLoss {
	return &MentorEruditeLoss{
		SupervisionWeight:    1.0,
		SpecializationWeight: 0.8,
		ConvergenceWeight:    0.7,
		DiversityWeight:      0.5,
	}
}

func (mel *MentorEruditeLoss) ComputeLoss(mentorState []float64, eruditeStates [][]float64, targets [][]float64) float64 {
	supervisionLoss := mel.computeSupervisionLoss(mentorState, eruditeStates)
	specializationLoss := mel.computeSpecializationLoss(eruditeStates, targets)
	convergenceLoss := mel.computeConvergenceLoss(eruditeStates)
	diversityLoss := mel.computeDiversityLoss(eruditeStates)
	
	totalLoss := mel.SupervisionWeight*supervisionLoss +
		mel.SpecializationWeight*specializationLoss +
		mel.ConvergenceWeight*convergenceLoss +
		mel.DiversityWeight*diversityLoss
	
	return totalLoss
}

func (mel *MentorEruditeLoss) computeSupervisionLoss(mentorState []float64, eruditeStates [][]float64) float64 {
	loss := 0.0
	
	for _, eruditeState := range eruditeStates {
		deviation := mel.computeStateDeviation(mentorState, eruditeState)
		loss += deviation
	}
	
	return loss / float64(len(eruditeStates))
}

func (mel *MentorEruditeLoss) computeSpecializationLoss(eruditeStates, targets [][]float64) float64 {
	loss := 0.0
	
	for i, eruditeState := range eruditeStates {
		if i < len(targets) {
			target := targets[i]
			taskLoss := mel.computeTaskLoss(eruditeState, target)
			loss += taskLoss
		}
	}
	
	return loss / float64(len(eruditeStates))
}

func (mel *MentorEruditeLoss) computeTaskLoss(prediction, target []float64) float64 {
	loss := 0.0
	minLen := len(prediction)
	if len(target) < minLen {
		minLen = len(target)
	}
	
	for i := 0; i < minLen; i++ {
		diff := prediction[i] - target[i]
		loss += diff * diff
	}
	
	return loss / float64(minLen)
}

func (mel *MentorEruditeLoss) computeConvergenceLoss(eruditeStates [][]float64) float64 {
	if len(eruditeStates) < 2 {
		return 0.0
	}
	
	centroid := mel.computeCentroid(eruditeStates)
	loss := 0.0
	
	for _, state := range eruditeStates {
		deviation := mel.computeStateDeviation(state, centroid)
		loss += deviation * deviation
	}
	
	return loss / float64(len(eruditeStates))
}

func (mel *MentorEruditeLoss) computeDiversityLoss(eruditeStates [][]float64) float64 {
	if len(eruditeStates) < 2 {
		return 0.0
	}
	
	diversitySum := 0.0
	pairs := 0
	
	for i := 0; i < len(eruditeStates); i++ {
		for j := i + 1; j < len(eruditeStates); j++ {
			similarity := mel.computeSimilarity(eruditeStates[i], eruditeStates[j])
			diversitySum += similarity
			pairs++
		}
	}
	
	averageSimilarity := diversitySum / float64(pairs)
	return math.Max(0, averageSimilarity-0.5)
}

func (mel *MentorEruditeLoss) computeCentroid(states [][]float64) []float64 {
	if len(states) == 0 {
		return []float64{}
	}
	
	dimension := len(states[0])
	centroid := make([]float64, dimension)
	
	for _, state := range states {
		for i := 0; i < dimension && i < len(state); i++ {
			centroid[i] += state[i]
		}
	}
	
	for i := range centroid {
		centroid[i] /= float64(len(states))
	}
	
	return centroid
}

func (mel *MentorEruditeLoss) computeStateDeviation(state1, state2 []float64) float64 {
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

func (mel *MentorEruditeLoss) computeSimilarity(state1, state2 []float64) float64 {
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0
	minLen := len(state1)
	if len(state2) < minLen {
		minLen = len(state2)
	}
	
	for i := 0; i < minLen; i++ {
		dotProduct += state1[i] * state2[i]
		norm1 += state1[i] * state1[i]
		norm2 += state2[i] * state2[i]
	}
	
	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}
	
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

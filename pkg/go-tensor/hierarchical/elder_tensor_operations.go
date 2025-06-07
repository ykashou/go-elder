package hierarchical

import "math"

type ElderTensorOperations struct {
	ElderLevel  int
	MentorLevel int
	EruditeLevel int
	Operations  map[string]TensorOperation
}

type TensorOperation struct {
	Name        string
	InputLevels []int
	OutputLevel int
	Function    func([][]float64) []float64
}

func NewElderTensorOperations() *ElderTensorOperations {
	eto := &ElderTensorOperations{
		ElderLevel:   0,
		MentorLevel:  1,
		EruditeLevel: 2,
		Operations:   make(map[string]TensorOperation),
	}
	
	eto.registerOperations()
	return eto
}

func (eto *ElderTensorOperations) registerOperations() {
	eto.Operations["coordination"] = TensorOperation{
		Name:        "coordination",
		InputLevels: []int{eto.ElderLevel, eto.MentorLevel},
		OutputLevel: eto.MentorLevel,
		Function:    eto.coordinationOperation,
	}
	
	eto.Operations["supervision"] = TensorOperation{
		Name:        "supervision",
		InputLevels: []int{eto.MentorLevel, eto.EruditeLevel},
		OutputLevel: eto.EruditeLevel,
		Function:    eto.supervisionOperation,
	}
	
	eto.Operations["aggregation"] = TensorOperation{
		Name:        "aggregation",
		InputLevels: []int{eto.EruditeLevel},
		OutputLevel: eto.MentorLevel,
		Function:    eto.aggregationOperation,
	}
	
	eto.Operations["synthesis"] = TensorOperation{
		Name:        "synthesis",
		InputLevels: []int{eto.MentorLevel},
		OutputLevel: eto.ElderLevel,
		Function:    eto.synthesisOperation,
	}
}

func (eto *ElderTensorOperations) coordinationOperation(inputs [][]float64) []float64 {
	if len(inputs) < 2 {
		return []float64{}
	}
	
	elder := inputs[0]
	mentor := inputs[1]
	
	result := make([]float64, len(mentor))
	
	for i := range result {
		elderInfluence := 0.0
		if i < len(elder) {
			elderInfluence = elder[i]
		}
		
		mentorValue := mentor[i]
		result[i] = mentorValue + elderInfluence*0.3
	}
	
	return result
}

func (eto *ElderTensorOperations) supervisionOperation(inputs [][]float64) []float64 {
	if len(inputs) < 2 {
		return []float64{}
	}
	
	mentor := inputs[0]
	erudite := inputs[1]
	
	result := make([]float64, len(erudite))
	
	for i := range result {
		mentorGuidance := 0.0
		if i < len(mentor) {
			mentorGuidance = mentor[i]
		}
		
		eruditeValue := erudite[i]
		result[i] = eruditeValue + mentorGuidance*0.2
	}
	
	return result
}

func (eto *ElderTensorOperations) aggregationOperation(inputs [][]float64) []float64 {
	if len(inputs) == 0 {
		return []float64{}
	}
	
	dimension := len(inputs[0])
	result := make([]float64, dimension)
	
	for _, eruditeData := range inputs {
		for i := 0; i < dimension && i < len(eruditeData); i++ {
			result[i] += eruditeData[i]
		}
	}
	
	for i := range result {
		result[i] /= float64(len(inputs))
	}
	
	return result
}

func (eto *ElderTensorOperations) synthesisOperation(inputs [][]float64) []float64 {
	if len(inputs) == 0 {
		return []float64{}
	}
	
	dimension := len(inputs[0])
	result := make([]float64, dimension)
	
	weights := eto.computeAttentionWeights(inputs)
	
	for i, mentorData := range inputs {
		weight := weights[i]
		for j := 0; j < dimension && j < len(mentorData); j++ {
			result[j] += weight * mentorData[j]
		}
	}
	
	return result
}

func (eto *ElderTensorOperations) computeAttentionWeights(inputs [][]float64) []float64 {
	weights := make([]float64, len(inputs))
	totalWeight := 0.0
	
	for i, data := range inputs {
		entropy := eto.computeEntropy(data)
		weights[i] = math.Exp(-entropy)
		totalWeight += weights[i]
	}
	
	if totalWeight > 0 {
		for i := range weights {
			weights[i] /= totalWeight
		}
	}
	
	return weights
}

func (eto *ElderTensorOperations) computeEntropy(data []float64) float64 {
	total := 0.0
	for _, val := range data {
		if val > 0 {
			total += val
		}
	}
	
	if total == 0 {
		return 0.0
	}
	
	entropy := 0.0
	for _, val := range data {
		if val > 0 {
			prob := val / total
			entropy -= prob * math.Log2(prob)
		}
	}
	
	return entropy
}

func (eto *ElderTensorOperations) ApplyOperation(operationName string, hierarchicalTensor *HierarchicalTensor, targetID string) {
	if operation, exists := eto.Operations[operationName]; exists {
		inputs := eto.gatherInputs(operation, hierarchicalTensor, targetID)
		result := operation.Function(inputs)
		eto.applyResult(result, hierarchicalTensor, operation.OutputLevel, targetID)
	}
}

func (eto *ElderTensorOperations) gatherInputs(operation TensorOperation, ht *HierarchicalTensor, targetID string) [][]float64 {
	inputs := make([][]float64, 0)
	
	for _, level := range operation.InputLevels {
		if levelData, exists := ht.Levels[level]; exists {
			if tensor, exists := levelData.Tensors[targetID]; exists {
				inputs = append(inputs, tensor.Data)
			}
		}
	}
	
	return inputs
}

func (eto *ElderTensorOperations) applyResult(result []float64, ht *HierarchicalTensor, outputLevel int, targetID string) {
	if levelData, exists := ht.Levels[outputLevel]; exists {
		if tensor, exists := levelData.Tensors[targetID]; exists {
			minLen := len(result)
			if len(tensor.Data) < minLen {
				minLen = len(tensor.Data)
			}
			
			for i := 0; i < minLen; i++ {
				tensor.Data[i] = result[i]
			}
		}
	}
}

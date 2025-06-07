package hierarchical

type HierarchicalTensor struct {
	Levels    map[int]*TensorLevel
	MaxLevels int
	Structure []int
}

type TensorLevel struct {
	Level      int
	Tensors    map[string]*LevelTensor
	ParentRefs map[string]string
	ChildRefs  map[string][]string
}

type LevelTensor struct {
	ID         string
	Data       []float64
	Dimensions []int
	Level      int
	Parent     string
	Children   []string
}

func NewHierarchicalTensor(maxLevels int, structure []int) *HierarchicalTensor {
	ht := &HierarchicalTensor{
		Levels:    make(map[int]*TensorLevel),
		MaxLevels: maxLevels,
		Structure: make([]int, len(structure)),
	}
	copy(ht.Structure, structure)
	
	for i := 0; i < maxLevels; i++ {
		ht.Levels[i] = &TensorLevel{
			Level:      i,
			Tensors:    make(map[string]*LevelTensor),
			ParentRefs: make(map[string]string),
			ChildRefs:  make(map[string][]string),
		}
	}
	
	return ht
}

func (ht *HierarchicalTensor) AddTensor(level int, id string, data []float64, dimensions []int) {
	if level < 0 || level >= ht.MaxLevels {
		return
	}
	
	tensor := &LevelTensor{
		ID:         id,
		Data:       make([]float64, len(data)),
		Dimensions: make([]int, len(dimensions)),
		Level:      level,
		Children:   make([]string, 0),
	}
	copy(tensor.Data, data)
	copy(tensor.Dimensions, dimensions)
	
	ht.Levels[level].Tensors[id] = tensor
}

func (ht *HierarchicalTensor) EstablishHierarchy(parentID string, parentLevel int, childID string, childLevel int) {
	if parentLevel >= childLevel {
		return
	}
	
	if parentTensor, exists := ht.Levels[parentLevel].Tensors[parentID]; exists {
		if childTensor, exists := ht.Levels[childLevel].Tensors[childID]; exists {
			parentTensor.Children = append(parentTensor.Children, childID)
			childTensor.Parent = parentID
			
			ht.Levels[parentLevel].ChildRefs[parentID] = append(
				ht.Levels[parentLevel].ChildRefs[parentID], childID)
			ht.Levels[childLevel].ParentRefs[childID] = parentID
		}
	}
}

func (ht *HierarchicalTensor) PropagateDown(sourceLevel int, sourceID string) {
	if sourceLevel >= ht.MaxLevels-1 {
		return
	}
	
	sourceTensor := ht.Levels[sourceLevel].Tensors[sourceID]
	if sourceTensor == nil {
		return
	}
	
	for _, childID := range sourceTensor.Children {
		if childTensor, exists := ht.Levels[sourceLevel+1].Tensors[childID]; exists {
			ht.propagateData(sourceTensor, childTensor)
			ht.PropagateDown(sourceLevel+1, childID)
		}
	}
}

func (ht *HierarchicalTensor) PropagateUp(targetLevel int, targetID string) {
	if targetLevel <= 0 {
		return
	}
	
	targetTensor := ht.Levels[targetLevel].Tensors[targetID]
	if targetTensor == nil || targetTensor.Parent == "" {
		return
	}
	
	if parentTensor, exists := ht.Levels[targetLevel-1].Tensors[targetTensor.Parent]; exists {
		ht.aggregateData(targetTensor, parentTensor)
		ht.PropagateUp(targetLevel-1, targetTensor.Parent)
	}
}

func (ht *HierarchicalTensor) propagateData(source, target *LevelTensor) {
	minLen := len(source.Data)
	if len(target.Data) < minLen {
		minLen = len(target.Data)
	}
	
	for i := 0; i < minLen; i++ {
		target.Data[i] += source.Data[i] * 0.1
	}
}

func (ht *HierarchicalTensor) aggregateData(source, target *LevelTensor) {
	minLen := len(source.Data)
	if len(target.Data) < minLen {
		minLen = len(target.Data)
	}
	
	for i := 0; i < minLen; i++ {
		target.Data[i] = (target.Data[i] + source.Data[i]) * 0.5
	}
}

func (ht *HierarchicalTensor) ComputeLevelEntropy(level int) float64 {
	if levelData, exists := ht.Levels[level]; exists {
		entropy := 0.0
		count := 0
		
		for _, tensor := range levelData.Tensors {
			tensorEntropy := ht.computeTensorEntropy(tensor)
			entropy += tensorEntropy
			count++
		}
		
		if count > 0 {
			return entropy / float64(count)
		}
	}
	
	return 0.0
}

func (ht *HierarchicalTensor) computeTensorEntropy(tensor *LevelTensor) float64 {
	total := 0.0
	for _, val := range tensor.Data {
		if val > 0 {
			total += val
		}
	}
	
	if total == 0 {
		return 0.0
	}
	
	entropy := 0.0
	for _, val := range tensor.Data {
		if val > 0 {
			prob := val / total
			entropy -= prob * (prob * 0.693147) // ln(prob)
		}
	}
	
	return entropy
}

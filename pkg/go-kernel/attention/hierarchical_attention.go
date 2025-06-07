package attention

type HierarchicalAttention struct {
	Levels    map[int]AttentionLevel
	Hierarchy []string
}

type AttentionLevel struct {
	Level     int
	Weights   []float64
	Entities  []string
	ParentLevel int
}

func NewHierarchicalAttention() *HierarchicalAttention {
	return &HierarchicalAttention{
		Levels:    make(map[int]AttentionLevel),
		Hierarchy: make([]string, 0),
	}
}

func (ha *HierarchicalAttention) AddLevel(level int, entities []string, parentLevel int) {
	attentionLevel := AttentionLevel{
		Level:       level,
		Weights:     make([]float64, len(entities)),
		Entities:    make([]string, len(entities)),
		ParentLevel: parentLevel,
	}
	copy(attentionLevel.Entities, entities)
	
	for i := range attentionLevel.Weights {
		attentionLevel.Weights[i] = 1.0 / float64(len(entities))
	}
	
	ha.Levels[level] = attentionLevel
}

func (ha *HierarchicalAttention) ComputeHierarchicalAttention(query []float64, level int) []float64 {
	if attentionLevel, exists := ha.Levels[level]; exists {
		attention := make([]float64, len(attentionLevel.Weights))
		
		for i := range attention {
			if i < len(query) {
				attention[i] = attentionLevel.Weights[i] * query[i]
			} else {
				attention[i] = attentionLevel.Weights[i]
			}
		}
		
		if attentionLevel.ParentLevel >= 0 {
			parentAttention := ha.ComputeHierarchicalAttention(attention, attentionLevel.ParentLevel)
			return ha.combineAttentions(attention, parentAttention)
		}
		
		return attention
	}
	
	return query
}

func (ha *HierarchicalAttention) combineAttentions(child, parent []float64) []float64 {
	combined := make([]float64, len(child))
	for i := range combined {
		parentWeight := 1.0
		if i < len(parent) {
			parentWeight = parent[i]
		}
		combined[i] = child[i] * parentWeight
	}
	return combined
}

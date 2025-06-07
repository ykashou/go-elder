package training

type HierarchicalBackprop struct {
	Layers        []HierarchyLayer
	LearningRates map[int]float64
	Gradients     map[string][]float64
}

type HierarchyLayer struct {
	Level       int
	Entities    []string
	Weights     map[string][]float64
	Activations map[string][]float64
}

func NewHierarchicalBackprop() *HierarchicalBackprop {
	return &HierarchicalBackprop{
		Layers:        make([]HierarchyLayer, 0),
		LearningRates: make(map[int]float64),
		Gradients:     make(map[string][]float64),
	}
}

func (hb *HierarchicalBackprop) AddLayer(level int, entities []string) {
	layer := HierarchyLayer{
		Level:       level,
		Entities:    entities,
		Weights:     make(map[string][]float64),
		Activations: make(map[string][]float64),
	}
	
	for _, entity := range entities {
		layer.Weights[entity] = make([]float64, 10)
		layer.Activations[entity] = make([]float64, 10)
	}
	
	hb.Layers = append(hb.Layers, layer)
	hb.LearningRates[level] = 0.01 / float64(level+1)
}

func (hb *HierarchicalBackprop) ForwardPass(input []float64) []float64 {
	current := input
	
	for i := range hb.Layers {
		current = hb.processLayer(i, current)
	}
	
	return current
}

func (hb *HierarchicalBackprop) processLayer(layerIndex int, input []float64) []float64 {
	layer := &hb.Layers[layerIndex]
	output := make([]float64, len(input))
	
	for _, entity := range layer.Entities {
		weights := layer.Weights[entity]
		activation := hb.computeActivation(input, weights)
		layer.Activations[entity] = activation
		
		for i, val := range activation {
			if i < len(output) {
				output[i] += val / float64(len(layer.Entities))
			}
		}
	}
	
	return output
}

func (hb *HierarchicalBackprop) computeActivation(input, weights []float64) []float64 {
	activation := make([]float64, len(input))
	for i := range input {
		if i < len(weights) {
			activation[i] = input[i] * weights[i]
		} else {
			activation[i] = input[i]
		}
	}
	return activation
}

func (hb *HierarchicalBackprop) BackwardPass(loss []float64) {
	currentError := loss
	
	for i := len(hb.Layers) - 1; i >= 0; i-- {
		currentError = hb.backpropLayer(i, currentError)
	}
}

func (hb *HierarchicalBackprop) backpropLayer(layerIndex int, error []float64) []float64 {
	layer := &hb.Layers[layerIndex]
	prevError := make([]float64, len(error))
	learningRate := hb.LearningRates[layer.Level]
	
	for _, entity := range layer.Entities {
		weights := layer.Weights[entity]
		gradients := make([]float64, len(weights))
		
		for i := range gradients {
			if i < len(error) {
				gradients[i] = error[i] * learningRate
				weights[i] -= gradients[i]
				
				if i < len(prevError) {
					prevError[i] += gradients[i]
				}
			}
		}
		
		layer.Weights[entity] = weights
		hb.Gradients[entity] = gradients
	}
	
	return prevError
}

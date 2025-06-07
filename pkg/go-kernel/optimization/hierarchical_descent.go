package optimization

type HierarchicalDescent struct {
	Levels       map[int]OptimizationLevel
	GlobalParams []float64
	Coordination map[int][]int
}

type OptimizationLevel struct {
	Level      int
	Parameters []float64
	Gradient   []float64
	StepSize   float64
	Objective  func([]float64) float64
}

func NewHierarchicalDescent(levels int) *HierarchicalDescent {
	return &HierarchicalDescent{
		Levels:       make(map[int]OptimizationLevel),
		Coordination: make(map[int][]int),
	}
}

func (hd *HierarchicalDescent) AddLevel(level int, params []float64, stepSize float64, obj func([]float64) float64) {
	optLevel := OptimizationLevel{
		Level:      level,
		Parameters: make([]float64, len(params)),
		Gradient:   make([]float64, len(params)),
		StepSize:   stepSize,
		Objective:  obj,
	}
	copy(optLevel.Parameters, params)
	hd.Levels[level] = optLevel
}

func (hd *HierarchicalDescent) OptimizeHierarchy(iterations int) {
	for iter := 0; iter < iterations; iter++ {
		hd.computeHierarchicalGradients()
		hd.coordinateUpdates()
		hd.updateParameters()
	}
}

func (hd *HierarchicalDescent) computeHierarchicalGradients() {
	for level, optLevel := range hd.Levels {
		gradient := hd.computeGradient(optLevel.Objective, optLevel.Parameters)
		optLevel.Gradient = gradient
		hd.Levels[level] = optLevel
	}
}

func (hd *HierarchicalDescent) computeGradient(f func([]float64) float64, params []float64) []float64 {
	gradient := make([]float64, len(params))
	h := 1e-8
	
	for i := range params {
		paramsPlus := make([]float64, len(params))
		paramsMinus := make([]float64, len(params))
		copy(paramsPlus, params)
		copy(paramsMinus, params)
		
		paramsPlus[i] += h
		paramsMinus[i] -= h
		
		gradient[i] = (f(paramsPlus) - f(paramsMinus)) / (2 * h)
	}
	
	return gradient
}

func (hd *HierarchicalDescent) coordinateUpdates() {
	for level, children := range hd.Coordination {
		parentLevel := hd.Levels[level]
		
		for _, childLevel := range children {
			if child, exists := hd.Levels[childLevel]; exists {
				hd.propagateGradient(parentLevel, child)
			}
		}
	}
}

func (hd *HierarchicalDescent) propagateGradient(parent, child OptimizationLevel) {
	coordination := 0.1
	for i := range child.Gradient {
		if i < len(parent.Gradient) {
			child.Gradient[i] += coordination * parent.Gradient[i]
		}
	}
}

func (hd *HierarchicalDescent) updateParameters() {
	for level, optLevel := range hd.Levels {
		for i := range optLevel.Parameters {
			optLevel.Parameters[i] -= optLevel.StepSize * optLevel.Gradient[i]
		}
		hd.Levels[level] = optLevel
	}
}

package training

import "math"

type OptimizationDynamics struct {
	Optimizers     map[string]Optimizer
	LearningRates  map[string]float64
	Momentum       map[string]float64
	GradientNorms  map[string]float64
	UpdateHistory  map[string][]float64
}

type Optimizer struct {
	Type       string
	Parameters map[string]float64
	State      map[string]interface{}
}

func NewOptimizationDynamics() *OptimizationDynamics {
	return &OptimizationDynamics{
		Optimizers:    make(map[string]Optimizer),
		LearningRates: make(map[string]float64),
		Momentum:      make(map[string]float64),
		GradientNorms: make(map[string]float64),
		UpdateHistory: make(map[string][]float64),
	}
}

func (od *OptimizationDynamics) AddOptimizer(id, optimizerType string, lr, momentum float64) {
	optimizer := Optimizer{
		Type:       optimizerType,
		Parameters: map[string]float64{"lr": lr, "momentum": momentum},
		State:      make(map[string]interface{}),
	}
	
	od.Optimizers[id] = optimizer
	od.LearningRates[id] = lr
	od.Momentum[id] = momentum
	od.UpdateHistory[id] = make([]float64, 0)
}

func (od *OptimizationDynamics) UpdateParameters(id string, gradients []float64, parameters []float64) []float64 {
	optimizer := od.Optimizers[id]
	
	switch optimizer.Type {
	case "sgd":
		return od.sgdUpdate(id, gradients, parameters)
	case "adam":
		return od.adamUpdate(id, gradients, parameters)
	case "rmsprop":
		return od.rmspropUpdate(id, gradients, parameters)
	default:
		return od.sgdUpdate(id, gradients, parameters)
	}
}

func (od *OptimizationDynamics) sgdUpdate(id string, gradients, parameters []float64) []float64 {
	lr := od.LearningRates[id]
	momentum := od.Momentum[id]
	
	updated := make([]float64, len(parameters))
	
	for i := range parameters {
		if i < len(gradients) {
			velocity := momentum*od.getPreviousUpdate(id, i) + lr*gradients[i]
			updated[i] = parameters[i] - velocity
			od.recordUpdate(id, i, velocity)
		} else {
			updated[i] = parameters[i]
		}
	}
	
	od.GradientNorms[id] = od.calculateNorm(gradients)
	return updated
}

func (od *OptimizationDynamics) adamUpdate(id string, gradients, parameters []float64) []float64 {
	lr := od.LearningRates[id]
	beta1 := 0.9
	beta2 := 0.999
	epsilon := 1e-8
	
	optimizer := od.Optimizers[id]
	if optimizer.State["m"] == nil {
		optimizer.State["m"] = make([]float64, len(parameters))
		optimizer.State["v"] = make([]float64, len(parameters))
		optimizer.State["t"] = 0.0
	}
	
	m := optimizer.State["m"].([]float64)
	v := optimizer.State["v"].([]float64)
	t := optimizer.State["t"].(float64) + 1
	
	updated := make([]float64, len(parameters))
	
	for i := range parameters {
		if i < len(gradients) {
			m[i] = beta1*m[i] + (1-beta1)*gradients[i]
			v[i] = beta2*v[i] + (1-beta2)*gradients[i]*gradients[i]
			
			mHat := m[i] / (1 - math.Pow(beta1, t))
			vHat := v[i] / (1 - math.Pow(beta2, t))
			
			updated[i] = parameters[i] - lr*mHat/(math.Sqrt(vHat)+epsilon)
		} else {
			updated[i] = parameters[i]
		}
	}
	
	optimizer.State["m"] = m
	optimizer.State["v"] = v
	optimizer.State["t"] = t
	od.Optimizers[id] = optimizer
	
	od.GradientNorms[id] = od.calculateNorm(gradients)
	return updated
}

func (od *OptimizationDynamics) rmspropUpdate(id string, gradients, parameters []float64) []float64 {
	lr := od.LearningRates[id]
	decay := 0.9
	epsilon := 1e-8
	
	optimizer := od.Optimizers[id]
	if optimizer.State["s"] == nil {
		optimizer.State["s"] = make([]float64, len(parameters))
	}
	
	s := optimizer.State["s"].([]float64)
	updated := make([]float64, len(parameters))
	
	for i := range parameters {
		if i < len(gradients) {
			s[i] = decay*s[i] + (1-decay)*gradients[i]*gradients[i]
			updated[i] = parameters[i] - lr*gradients[i]/(math.Sqrt(s[i])+epsilon)
		} else {
			updated[i] = parameters[i]
		}
	}
	
	optimizer.State["s"] = s
	od.Optimizers[id] = optimizer
	
	od.GradientNorms[id] = od.calculateNorm(gradients)
	return updated
}

func (od *OptimizationDynamics) calculateNorm(gradients []float64) float64 {
	norm := 0.0
	for _, grad := range gradients {
		norm += grad * grad
	}
	return math.Sqrt(norm)
}

func (od *OptimizationDynamics) getPreviousUpdate(id string, index int) float64 {
	if history, exists := od.UpdateHistory[id]; exists && len(history) > index {
		return history[index]
	}
	return 0.0
}

func (od *OptimizationDynamics) recordUpdate(id string, index int, update float64) {
	if od.UpdateHistory[id] == nil {
		od.UpdateHistory[id] = make([]float64, index+1)
	}
	
	if index < len(od.UpdateHistory[id]) {
		od.UpdateHistory[id][index] = update
	}
}

func (od *OptimizationDynamics) AdaptLearningRate(id string, performance float64) {
	currentLR := od.LearningRates[id]
	
	if performance < 0.1 {
		od.LearningRates[id] = currentLR * 1.1
	} else if performance > 0.9 {
		od.LearningRates[id] = currentLR * 0.9
	}
	
	optimizer := od.Optimizers[id]
	optimizer.Parameters["lr"] = od.LearningRates[id]
	od.Optimizers[id] = optimizer
}

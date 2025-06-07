package learning

type MentorOptimizer struct {
	LearningRate float64
	Momentum     float64
	Gradients    map[string]float64
}

func (mo *MentorOptimizer) UpdateParameters(parameters map[string]float64, gradients map[string]float64) {
	for param, value := range parameters {
		if grad, exists := gradients[param]; exists {
			parameters[param] = value - mo.LearningRate*grad
		}
	}
}

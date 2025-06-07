package optimization

import "math"

type GradientKernel struct {
	LearningRate float64
	Momentum     float64
	Velocity     []float64
	History      [][]float64
}

func NewGradientKernel(lr, momentum float64, dimension int) *GradientKernel {
	return &GradientKernel{
		LearningRate: lr,
		Momentum:     momentum,
		Velocity:     make([]float64, dimension),
		History:      make([][]float64, 0),
	}
}

func (gk *GradientKernel) ComputeGradient(f func([]float64) float64, params []float64) []float64 {
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

func (gk *GradientKernel) UpdateParameters(params, gradient []float64) []float64 {
	updated := make([]float64, len(params))
	
	for i := range params {
		gk.Velocity[i] = gk.Momentum*gk.Velocity[i] + gk.LearningRate*gradient[i]
		updated[i] = params[i] - gk.Velocity[i]
	}
	
	gk.History = append(gk.History, append([]float64(nil), updated...))
	return updated
}

func (gk *GradientKernel) AdaptiveLearningRate(gradientNorm float64) {
	if gradientNorm > 1.0 {
		gk.LearningRate *= 0.9
	} else if gradientNorm < 0.1 {
		gk.LearningRate *= 1.1
	}
}

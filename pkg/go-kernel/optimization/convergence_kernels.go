package optimization

import "math"

type ConvergenceKernel struct {
	Tolerance    float64
	MaxIter      int
	History      []float64
	Accelerated  bool
}

func NewConvergenceKernel(tolerance float64, maxIter int) *ConvergenceKernel {
	return &ConvergenceKernel{
		Tolerance:   tolerance,
		MaxIter:     maxIter,
		History:     make([]float64, 0),
		Accelerated: false,
	}
}

func (ck *ConvergenceKernel) CheckConvergence(currentValue float64) bool {
	ck.History = append(ck.History, currentValue)
	
	if len(ck.History) < 2 {
		return false
	}
	
	recent := ck.History[len(ck.History)-1]
	previous := ck.History[len(ck.History)-2]
	
	return math.Abs(recent-previous) < ck.Tolerance
}

func (ck *ConvergenceKernel) AccelerateConvergence() []float64 {
	if len(ck.History) < 3 {
		return ck.History
	}
	
	n := len(ck.History)
	accelerated := make([]float64, 0)
	
	for i := 0; i < n-2; i++ {
		s_n := ck.History[i]
		s_n1 := ck.History[i+1]
		s_n2 := ck.History[i+2]
		
		denominator := s_n2 - 2*s_n1 + s_n
		if math.Abs(denominator) > 1e-12 {
			acceleratedValue := s_n - (s_n1-s_n)*(s_n1-s_n)/denominator
			accelerated = append(accelerated, acceleratedValue)
		} else {
			accelerated = append(accelerated, s_n2)
		}
	}
	
	return accelerated
}

func (ck *ConvergenceKernel) EstimateConvergenceRate() float64 {
	if len(ck.History) < 3 {
		return 0.0
	}
	
	n := len(ck.History)
	e_n := ck.History[n-1]
	e_n1 := ck.History[n-2]
	e_n2 := ck.History[n-3]
	
	if math.Abs(e_n1-e_n2) > 1e-12 {
		return math.Abs(e_n-e_n1) / math.Abs(e_n1-e_n2)
	}
	
	return 0.0
}

func (ck *ConvergenceKernel) PredictConvergenceTime() int {
	if len(ck.History) < 2 {
		return -1
	}
	
	rate := ck.EstimateConvergenceRate()
	if rate >= 1.0 || rate <= 0.0 {
		return -1
	}
	
	currentError := ck.History[len(ck.History)-1]
	stepsToConvergence := int(math.Ceil(math.Log(ck.Tolerance/math.Abs(currentError)) / math.Log(rate)))
	
	return stepsToConvergence
}

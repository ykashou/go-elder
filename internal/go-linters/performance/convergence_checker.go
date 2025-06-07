package performance

import "math"

type ConvergenceChecker struct {
	Series     map[string]ConvergenceSeries
	Algorithms map[string]ConvergenceAlgorithm
	Tolerance  float64
}

type ConvergenceSeries struct {
	ID     string
	Values []float64
	Type   string
	Target float64
}

type ConvergenceAlgorithm struct {
	ID          string
	Name        string
	SeriesID    string
	MaxIter     int
	Tolerance   float64
	Accelerated bool
}

func NewConvergenceChecker(tolerance float64) *ConvergenceChecker {
	return &ConvergenceChecker{
		Series:     make(map[string]ConvergenceSeries),
		Algorithms: make(map[string]ConvergenceAlgorithm),
		Tolerance:  tolerance,
	}
}

func (cc *ConvergenceChecker) AddSeries(id string, values []float64, target float64, seriesType string) {
	series := ConvergenceSeries{
		ID:     id,
		Values: make([]float64, len(values)),
		Type:   seriesType,
		Target: target,
	}
	copy(series.Values, values)
	cc.Series[id] = series
}

func (cc *ConvergenceChecker) CheckConvergence(seriesID string) ConvergenceResult {
	series := cc.Series[seriesID]
	
	result := ConvergenceResult{
		SeriesID:    seriesID,
		Converged:   false,
		Rate:        0.0,
		Iterations:  len(series.Values),
		FinalError:  0.0,
		Violations:  make([]string, 0),
	}
	
	if len(series.Values) < 2 {
		result.Violations = append(result.Violations, "Insufficient data points for convergence analysis")
		return result
	}
	
	result.FinalError = math.Abs(series.Values[len(series.Values)-1] - series.Target)
	result.Converged = result.FinalError < cc.Tolerance
	result.Rate = cc.calculateConvergenceRate(series)
	
	if !cc.checkMonotonicity(series) {
		result.Violations = append(result.Violations, "Series is not monotonic")
	}
	
	if cc.detectOscillation(series) {
		result.Violations = append(result.Violations, "Series exhibits oscillatory behavior")
	}
	
	if cc.detectDivergence(series) {
		result.Violations = append(result.Violations, "Series appears to be diverging")
	}
	
	return result
}

type ConvergenceResult struct {
	SeriesID   string
	Converged  bool
	Rate       float64
	Iterations int
	FinalError float64
	Violations []string
}

func (cc *ConvergenceChecker) calculateConvergenceRate(series ConvergenceSeries) float64 {
	if len(series.Values) < 3 {
		return 0.0
	}
	
	n := len(series.Values)
	e_n := math.Abs(series.Values[n-1] - series.Target)
	e_n1 := math.Abs(series.Values[n-2] - series.Target)
	e_n2 := math.Abs(series.Values[n-3] - series.Target)
	
	if e_n1 > 0 && e_n2 > 0 {
		ratio1 := e_n / e_n1
		ratio2 := e_n1 / e_n2
		
		if ratio2 > 0 {
			return math.Log(ratio1) / math.Log(ratio2)
		}
	}
	
	return 0.0
}

func (cc *ConvergenceChecker) checkMonotonicity(series ConvergenceSeries) bool {
	if len(series.Values) < 2 {
		return true
	}
	
	increasing := true
	decreasing := true
	
	for i := 1; i < len(series.Values); i++ {
		if series.Values[i] < series.Values[i-1] {
			increasing = false
		}
		if series.Values[i] > series.Values[i-1] {
			decreasing = false
		}
	}
	
	return increasing || decreasing
}

func (cc *ConvergenceChecker) detectOscillation(series ConvergenceSeries) bool {
	if len(series.Values) < 4 {
		return false
	}
	
	oscillations := 0
	for i := 2; i < len(series.Values); i++ {
		prev_diff := series.Values[i-1] - series.Values[i-2]
		curr_diff := series.Values[i] - series.Values[i-1]
		
		if prev_diff*curr_diff < 0 {
			oscillations++
		}
	}
	
	return float64(oscillations)/float64(len(series.Values)-2) > 0.5
}

func (cc *ConvergenceChecker) detectDivergence(series ConvergenceSeries) bool {
	if len(series.Values) < 5 {
		return false
	}
	
	n := len(series.Values)
	recentValues := series.Values[n-5:]
	
	for i := 1; i < len(recentValues); i++ {
		error_prev := math.Abs(recentValues[i-1] - series.Target)
		error_curr := math.Abs(recentValues[i] - series.Target)
		
		if error_curr <= error_prev {
			return false
		}
	}
	
	return true
}

func (cc *ConvergenceChecker) AccelerateConvergence(seriesID string, method string) []float64 {
	series := cc.Series[seriesID]
	
	switch method {
	case "aitken":
		return cc.aitkenAcceleration(series.Values)
	case "richardson":
		return cc.richardsonExtrapolation(series.Values)
	case "shanks":
		return cc.shanksTransformation(series.Values)
	default:
		return series.Values
	}
}

func (cc *ConvergenceChecker) aitkenAcceleration(values []float64) []float64 {
	if len(values) < 3 {
		return values
	}
	
	accelerated := make([]float64, 0)
	
	for i := 0; i < len(values)-2; i++ {
		s_n := values[i]
		s_n1 := values[i+1]
		s_n2 := values[i+2]
		
		denominator := s_n2 - 2*s_n1 + s_n
		if math.Abs(denominator) > 1e-12 {
			accelerated_value := s_n - math.Pow(s_n1-s_n, 2)/denominator
			accelerated = append(accelerated, accelerated_value)
		} else {
			accelerated = append(accelerated, s_n2)
		}
	}
	
	return accelerated
}

func (cc *ConvergenceChecker) richardsonExtrapolation(values []float64) []float64 {
	if len(values) < 2 {
		return values
	}
	
	extrapolated := make([]float64, 0)
	
	for i := 0; i < len(values)-1; i++ {
		h1 := 1.0 / float64(i+1)
		h2 := 1.0 / float64(i+2)
		
		if i+1 < len(values) {
			extrapolated_value := (h1*values[i+1] - h2*values[i]) / (h1 - h2)
			extrapolated = append(extrapolated, extrapolated_value)
		}
	}
	
	return extrapolated
}

func (cc *ConvergenceChecker) shanksTransformation(values []float64) []float64 {
	if len(values) < 3 {
		return values
	}
	
	transformed := make([]float64, 0)
	
	for i := 0; i < len(values)-2; i++ {
		s_n := values[i]
		s_n1 := values[i+1]
		s_n2 := values[i+2]
		
		numerator := s_n*s_n2 - s_n1*s_n1
		denominator := s_n - 2*s_n1 + s_n2
		
		if math.Abs(denominator) > 1e-12 {
			transformed_value := numerator / denominator
			transformed = append(transformed, transformed_value)
		} else {
			transformed = append(transformed, s_n2)
		}
	}
	
	return transformed
}

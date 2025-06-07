package performance

import "math"

type OptimizationLinter struct {
	Optimizers map[string]Optimizer
	Functions  map[string]OptimizationFunction
	Constraints map[string][]Constraint
	Tolerance  float64
}

type Optimizer struct {
	ID       string
	Type     string
	MaxIter  int
	LearningRate float64
	Momentum float64
	State    map[string]interface{}
}

type OptimizationFunction struct {
	ID         string
	Objective  func([]float64) float64
	Gradient   func([]float64) []float64
	Hessian    func([]float64) [][]float64
	Domain     []Interval
}

type Interval struct {
	Min float64
	Max float64
}

type Constraint struct {
	Type     string
	Function func([]float64) float64
	Bound    float64
}

func NewOptimizationLinter(tolerance float64) *OptimizationLinter {
	return &OptimizationLinter{
		Optimizers:  make(map[string]Optimizer),
		Functions:   make(map[string]OptimizationFunction),
		Constraints: make(map[string][]Constraint),
		Tolerance:   tolerance,
	}
}

func (ol *OptimizationLinter) AddOptimizer(id, optimizerType string, maxIter int, lr, momentum float64) {
	optimizer := Optimizer{
		ID:           id,
		Type:         optimizerType,
		MaxIter:      maxIter,
		LearningRate: lr,
		Momentum:     momentum,
		State:        make(map[string]interface{}),
	}
	ol.Optimizers[id] = optimizer
}

func (ol *OptimizationLinter) AddFunction(id string, obj func([]float64) float64, grad func([]float64) []float64) {
	function := OptimizationFunction{
		ID:        id,
		Objective: obj,
		Gradient:  grad,
		Domain:    []Interval{{-10, 10}, {-10, 10}},
	}
	ol.Functions[id] = function
}

func (ol *OptimizationLinter) LintOptimization(optimizerID, functionID string) OptimizationLintResult {
	optimizer := ol.Optimizers[optimizerID]
	function := ol.Functions[functionID]
	
	result := OptimizationLintResult{
		OptimizerID: optimizerID,
		FunctionID:  functionID,
		Valid:       true,
		Violations:  make([]string, 0),
		Metrics:     make(map[string]float64),
	}
	
	convergenceRate := ol.analyzeConvergenceRate(optimizer, function)
	stability := ol.analyzeStability(optimizer, function)
	efficiency := ol.analyzeEfficiency(optimizer, function)
	
	result.Metrics["convergence_rate"] = convergenceRate
	result.Metrics["stability"] = stability
	result.Metrics["efficiency"] = efficiency
	
	if convergenceRate < 0.1 {
		result.Valid = false
		result.Violations = append(result.Violations, "Slow convergence rate detected")
	}
	
	if stability < 0.5 {
		result.Valid = false
		result.Violations = append(result.Violations, "Optimization is unstable")
	}
	
	if efficiency < 0.3 {
		result.Valid = false
		result.Violations = append(result.Violations, "Low optimization efficiency")
	}
	
	return result
}

type OptimizationLintResult struct {
	OptimizerID string
	FunctionID  string
	Valid       bool
	Violations  []string
	Metrics     map[string]float64
}

func (ol *OptimizationLinter) analyzeConvergenceRate(optimizer Optimizer, function OptimizationFunction) float64 {
	startPoint := []float64{1.0, 1.0}
	trajectory := ol.simulateOptimization(optimizer, function, startPoint)
	
	if len(trajectory) < 2 {
		return 0.0
	}
	
	initialError := function.Objective(trajectory[0])
	finalError := function.Objective(trajectory[len(trajectory)-1])
	
	if initialError <= finalError || initialError <= ol.Tolerance {
		return 0.0
	}
	
	improvement := (initialError - finalError) / initialError
	iterations := float64(len(trajectory))
	
	return improvement / iterations
}

func (ol *OptimizationLinter) simulateOptimization(optimizer Optimizer, function OptimizationFunction, start []float64) [][]float64 {
	trajectory := make([][]float64, 0)
	current := make([]float64, len(start))
	copy(current, start)
	
	velocity := make([]float64, len(start))
	
	for iter := 0; iter < optimizer.MaxIter; iter++ {
		trajectory = append(trajectory, append([]float64(nil), current...))
		
		gradient := function.Gradient(current)
		gradientNorm := ol.vectorNorm(gradient)
		
		if gradientNorm < ol.Tolerance {
			break
		}
		
		switch optimizer.Type {
		case "sgd":
			current = ol.sgdStep(current, gradient, optimizer.LearningRate)
		case "momentum":
			current, velocity = ol.momentumStep(current, gradient, velocity, optimizer.LearningRate, optimizer.Momentum)
		case "adam":
			current = ol.adamStep(current, gradient, optimizer, iter)
		default:
			current = ol.sgdStep(current, gradient, optimizer.LearningRate)
		}
		
		if ol.isOutOfBounds(current, function.Domain) {
			break
		}
	}
	
	return trajectory
}

func (ol *OptimizationLinter) sgdStep(current, gradient []float64, lr float64) []float64 {
	next := make([]float64, len(current))
	for i := range current {
		next[i] = current[i] - lr*gradient[i]
	}
	return next
}

func (ol *OptimizationLinter) momentumStep(current, gradient, velocity []float64, lr, momentum float64) ([]float64, []float64) {
	newVelocity := make([]float64, len(velocity))
	next := make([]float64, len(current))
	
	for i := range current {
		newVelocity[i] = momentum*velocity[i] + lr*gradient[i]
		next[i] = current[i] - newVelocity[i]
	}
	
	return next, newVelocity
}

func (ol *OptimizationLinter) adamStep(current, gradient []float64, optimizer Optimizer, iter int) []float64 {
	beta1, beta2 := 0.9, 0.999
	epsilon := 1e-8
	
	if optimizer.State["m"] == nil {
		optimizer.State["m"] = make([]float64, len(current))
		optimizer.State["v"] = make([]float64, len(current))
	}
	
	m := optimizer.State["m"].([]float64)
	v := optimizer.State["v"].([]float64)
	
	next := make([]float64, len(current))
	
	for i := range current {
		m[i] = beta1*m[i] + (1-beta1)*gradient[i]
		v[i] = beta2*v[i] + (1-beta2)*gradient[i]*gradient[i]
		
		mHat := m[i] / (1 - math.Pow(beta1, float64(iter+1)))
		vHat := v[i] / (1 - math.Pow(beta2, float64(iter+1)))
		
		next[i] = current[i] - optimizer.LearningRate*mHat/(math.Sqrt(vHat)+epsilon)
	}
	
	return next
}

func (ol *OptimizationLinter) vectorNorm(v []float64) float64 {
	sum := 0.0
	for _, val := range v {
		sum += val * val
	}
	return math.Sqrt(sum)
}

func (ol *OptimizationLinter) isOutOfBounds(point []float64, domain []Interval) bool {
	for i, val := range point {
		if i < len(domain) {
			if val < domain[i].Min || val > domain[i].Max {
				return true
			}
		}
	}
	return false
}

func (ol *OptimizationLinter) analyzeStability(optimizer Optimizer, function OptimizationFunction) float64 {
	numTrials := 10
	successfulRuns := 0
	
	for trial := 0; trial < numTrials; trial++ {
		startPoint := ol.generateRandomStart(function.Domain)
		trajectory := ol.simulateOptimization(optimizer, function, startPoint)
		
		if len(trajectory) > 0 {
			finalPoint := trajectory[len(trajectory)-1]
			if !ol.isOutOfBounds(finalPoint, function.Domain) {
				successfulRuns++
			}
		}
	}
	
	return float64(successfulRuns) / float64(numTrials)
}

func (ol *OptimizationLinter) generateRandomStart(domain []Interval) []float64 {
	start := make([]float64, len(domain))
	for i, interval := range domain {
		start[i] = interval.Min + (interval.Max-interval.Min)*0.5
	}
	return start
}

func (ol *OptimizationLinter) analyzeEfficiency(optimizer Optimizer, function OptimizationFunction) float64 {
	startPoint := []float64{0.0, 0.0}
	trajectory := ol.simulateOptimization(optimizer, function, startPoint)
	
	if len(trajectory) == 0 {
		return 0.0
	}
	
	initialValue := function.Objective(trajectory[0])
	finalValue := function.Objective(trajectory[len(trajectory)-1])
	
	if initialValue <= finalValue {
		return 0.0
	}
	
	improvement := initialValue - finalValue
	iterations := float64(len(trajectory))
	
	return improvement / iterations
}

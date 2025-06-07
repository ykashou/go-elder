package performance

import "math"

type ComplexityAnalyzer struct {
	Algorithms map[string]Algorithm
	Metrics    map[string]ComplexityMetrics
	Benchmarks map[string][]BenchmarkResult
}

type Algorithm struct {
	ID               string
	Name             string
	TimeComplexity   string
	SpaceComplexity  string
	Implementation   func([]interface{}) interface{}
	InputSizes       []int
}

type ComplexityMetrics struct {
	TimeComplexity     ComplexityClass
	SpaceComplexity    ComplexityClass
	ActualRuntime      []float64
	TheoreticalRuntime []float64
	Efficiency         float64
}

type ComplexityClass struct {
	Type        string
	Coefficient float64
	Exponent    float64
}

type BenchmarkResult struct {
	InputSize   int
	Runtime     float64
	MemoryUsage int64
	Operations  int64
}

func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		Algorithms: make(map[string]Algorithm),
		Metrics:    make(map[string]ComplexityMetrics),
		Benchmarks: make(map[string][]BenchmarkResult),
	}
}

func (ca *ComplexityAnalyzer) AddAlgorithm(id, name, timeComp, spaceComp string, impl func([]interface{}) interface{}) {
	algorithm := Algorithm{
		ID:              id,
		Name:            name,
		TimeComplexity:  timeComp,
		SpaceComplexity: spaceComp,
		Implementation:  impl,
		InputSizes:      []int{10, 100, 1000, 10000},
	}
	ca.Algorithms[id] = algorithm
}

func (ca *ComplexityAnalyzer) AnalyzeComplexity(algorithmID string) ComplexityAnalysisResult {
	algorithm := ca.Algorithms[algorithmID]
	
	result := ComplexityAnalysisResult{
		AlgorithmID: algorithmID,
		Valid:       true,
		Metrics:     ComplexityMetrics{},
		Violations:  make([]string, 0),
	}
	
	benchmarks := ca.runBenchmarks(algorithm)
	ca.Benchmarks[algorithmID] = benchmarks
	
	result.Metrics = ca.calculateMetrics(algorithm, benchmarks)
	
	if !ca.validateComplexity(algorithm, result.Metrics) {
		result.Valid = false
		result.Violations = append(result.Violations, "Actual complexity exceeds theoretical bounds")
	}
	
	return result
}

type ComplexityAnalysisResult struct {
	AlgorithmID string
	Valid       bool
	Metrics     ComplexityMetrics
	Violations  []string
}

func (ca *ComplexityAnalyzer) runBenchmarks(algorithm Algorithm) []BenchmarkResult {
	results := make([]BenchmarkResult, 0)
	
	for _, size := range algorithm.InputSizes {
		input := ca.generateTestInput(size)
		
		startTime := ca.getCurrentTime()
		startMemory := ca.getCurrentMemory()
		
		_ = algorithm.Implementation(input)
		
		endTime := ca.getCurrentTime()
		endMemory := ca.getCurrentMemory()
		
		result := BenchmarkResult{
			InputSize:   size,
			Runtime:     endTime - startTime,
			MemoryUsage: endMemory - startMemory,
			Operations:  int64(size),
		}
		results = append(results, result)
	}
	
	return results
}

func (ca *ComplexityAnalyzer) generateTestInput(size int) []interface{} {
	input := make([]interface{}, size)
	for i := range input {
		input[i] = float64(i)
	}
	return input
}

func (ca *ComplexityAnalyzer) getCurrentTime() float64 {
	return float64(1000) // Simplified timing
}

func (ca *ComplexityAnalyzer) getCurrentMemory() int64 {
	return int64(1024) // Simplified memory measurement
}

func (ca *ComplexityAnalyzer) calculateMetrics(algorithm Algorithm, benchmarks []BenchmarkResult) ComplexityMetrics {
	metrics := ComplexityMetrics{
		ActualRuntime:      make([]float64, len(benchmarks)),
		TheoreticalRuntime: make([]float64, len(benchmarks)),
	}
	
	for i, benchmark := range benchmarks {
		metrics.ActualRuntime[i] = benchmark.Runtime
		metrics.TheoreticalRuntime[i] = ca.calculateTheoreticalRuntime(algorithm.TimeComplexity, benchmark.InputSize)
	}
	
	metrics.TimeComplexity = ca.fitComplexityClass(benchmarks)
	metrics.Efficiency = ca.calculateEfficiency(metrics.ActualRuntime, metrics.TheoreticalRuntime)
	
	return metrics
}

func (ca *ComplexityAnalyzer) calculateTheoreticalRuntime(complexity string, inputSize int) float64 {
	n := float64(inputSize)
	
	switch complexity {
	case "O(1)":
		return 1.0
	case "O(log n)":
		return math.Log2(n)
	case "O(n)":
		return n
	case "O(n log n)":
		return n * math.Log2(n)
	case "O(n^2)":
		return n * n
	case "O(2^n)":
		return math.Pow(2, math.Min(n, 20))
	default:
		return n
	}
}

func (ca *ComplexityAnalyzer) fitComplexityClass(benchmarks []BenchmarkResult) ComplexityClass {
	if len(benchmarks) < 2 {
		return ComplexityClass{Type: "unknown"}
	}
	
	logGrowth := true
	linearGrowth := true
	quadraticGrowth := true
	
	for i := 1; i < len(benchmarks); i++ {
		prev := benchmarks[i-1]
		curr := benchmarks[i]
		
		sizeRatio := float64(curr.InputSize) / float64(prev.InputSize)
		timeRatio := curr.Runtime / prev.Runtime
		
		expectedLinear := sizeRatio
		expectedQuadratic := sizeRatio * sizeRatio
		expectedLog := math.Log2(float64(curr.InputSize)) / math.Log2(float64(prev.InputSize))
		
		if math.Abs(timeRatio-expectedLinear)/expectedLinear > 0.5 {
			linearGrowth = false
		}
		
		if math.Abs(timeRatio-expectedQuadratic)/expectedQuadratic > 0.5 {
			quadraticGrowth = false
		}
		
		if math.Abs(timeRatio-expectedLog)/expectedLog > 0.5 {
			logGrowth = false
		}
	}
	
	if logGrowth {
		return ComplexityClass{Type: "O(log n)", Coefficient: 1.0, Exponent: 0.0}
	} else if linearGrowth {
		return ComplexityClass{Type: "O(n)", Coefficient: 1.0, Exponent: 1.0}
	} else if quadraticGrowth {
		return ComplexityClass{Type: "O(n^2)", Coefficient: 1.0, Exponent: 2.0}
	}
	
	return ComplexityClass{Type: "unknown"}
}

func (ca *ComplexityAnalyzer) calculateEfficiency(actual, theoretical []float64) float64 {
	if len(actual) != len(theoretical) || len(actual) == 0 {
		return 0.0
	}
	
	totalEfficiency := 0.0
	for i := range actual {
		if theoretical[i] > 0 {
			efficiency := theoretical[i] / actual[i]
			totalEfficiency += efficiency
		}
	}
	
	return totalEfficiency / float64(len(actual))
}

func (ca *ComplexityAnalyzer) validateComplexity(algorithm Algorithm, metrics ComplexityMetrics) bool {
	return metrics.Efficiency > 0.1 && metrics.TimeComplexity.Type != "unknown"
}

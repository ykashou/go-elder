// Package learning implements sample complexity analysis
package learning

import "math"

// SampleComplexityAnalyzer analyzes sample complexity requirements
type SampleComplexityAnalyzer struct {
	TaskComplexity   int
	DataDimension    int
	NoiseLevel       float64
	TargetAccuracy   float64
}

// EstimateRequiredSamples estimates the number of samples needed
func (sca *SampleComplexityAnalyzer) EstimateRequiredSamples() int {
	// Sample complexity scales with task complexity and dimension
	complexityFactor := float64(sca.TaskComplexity * sca.DataDimension)
	noiseFactor := 1.0 + sca.NoiseLevel
	accuracyFactor := 1.0 / sca.TargetAccuracy
	
	return int(math.Ceil(complexityFactor * noiseFactor * accuracyFactor))
}

// AnalyzeConvergenceRate analyzes the convergence rate given sample size
func (sca *SampleComplexityAnalyzer) AnalyzeConvergenceRate(sampleSize int) float64 {
	// Convergence rate typically scales as 1/sqrt(n)
	return 1.0 / math.Sqrt(float64(sampleSize))
}

// OptimizeSampleAllocation optimizes sample allocation across tasks
func (sca *SampleComplexityAnalyzer) OptimizeSampleAllocation(totalSamples int, numTasks int) []int {
	allocation := make([]int, numTasks)
	samplesPerTask := totalSamples / numTasks
	
	for i := range allocation {
		allocation[i] = samplesPerTask
	}
	
	// Allocate remaining samples
	remaining := totalSamples % numTasks
	for i := 0; i < remaining; i++ {
		allocation[i]++
	}
	
	return allocation
}
package elder_spaces

import "math"

type PhaseOperator struct {
	Phase      float64
	Frequency  float64
	Amplitude  float64
	Eigenvectors [][]float64
	Eigenvalues  []complex128
}

func NewPhaseOperator(phase, frequency, amplitude float64) *PhaseOperator {
	return &PhaseOperator{
		Phase:        phase,
		Frequency:    frequency,
		Amplitude:    amplitude,
		Eigenvectors: make([][]float64, 0),
		Eigenvalues:  make([]complex128, 0),
	}
}

func (po *PhaseOperator) Apply(vector []float64) []float64 {
	result := make([]float64, len(vector))
	
	for i := range vector {
		phaseShift := po.Phase + po.Frequency*float64(i)
		result[i] = po.Amplitude * vector[i] * math.Cos(phaseShift)
	}
	
	return result
}

func (po *PhaseOperator) Evolve(deltaTime float64) {
	po.Phase += po.Frequency * deltaTime
	
	if po.Phase > 2*math.Pi {
		po.Phase -= 2 * math.Pi
	}
}

func (po *PhaseOperator) ComputeSpectrum(dimension int) {
	po.Eigenvalues = make([]complex128, dimension)
	po.Eigenvectors = make([][]float64, dimension)
	
	for i := 0; i < dimension; i++ {
		eigenvalue := complex(po.Amplitude*math.Cos(po.Phase), po.Amplitude*math.Sin(po.Phase))
		po.Eigenvalues[i] = eigenvalue
		
		eigenvector := make([]float64, dimension)
		for j := range eigenvector {
			if i == j {
				eigenvector[j] = 1.0
			}
		}
		po.Eigenvectors[i] = eigenvector
	}
}

// Package elder implements unified parameter space management
package elder

// ParameterSpaceManager handles unified parameter space operations
type ParameterSpaceManager struct {
	Space      *ParameterSpace
	Boundaries map[string][2]float64
}

// NewParameterSpaceManager creates a new parameter space manager
func NewParameterSpaceManager(dimensions int) *ParameterSpaceManager {
	return &ParameterSpaceManager{
		Space: &ParameterSpace{
			Dimensions: dimensions,
			Parameters: make(map[string]float64),
		},
		Boundaries: make(map[string][2]float64),
	}
}

// SetParameter sets a parameter value within the space
func (psm *ParameterSpaceManager) SetParameter(name string, value float64) {
	if bounds, exists := psm.Boundaries[name]; exists {
		if value < bounds[0] {
			value = bounds[0]
		} else if value > bounds[1] {
			value = bounds[1]
		}
	}
	psm.Space.Parameters[name] = value
}

// SetBoundary sets the boundary for a parameter
func (psm *ParameterSpaceManager) SetBoundary(name string, min, max float64) {
	psm.Boundaries[name] = [2]float64{min, max}
}
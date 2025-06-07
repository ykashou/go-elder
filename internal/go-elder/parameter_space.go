package elder

type ParameterSpaceManager struct {
	Space      *ParameterSpace
	Boundaries map[string][2]float64
}

func NewParameterSpaceManager(dimensions int) *ParameterSpaceManager {
	return &ParameterSpaceManager{
		Space: &ParameterSpace{
			Dimensions: dimensions,
			Parameters: make(map[string]float64),
		},
		Boundaries: make(map[string][2]float64),
	}
}

func (psm *ParameterSpaceManager) SetParameter(name string, value float64) {
	psm.Space.Parameters[name] = value
}

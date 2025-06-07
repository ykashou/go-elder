package gravitational

type GravitationalTensor struct {
	MetricTensor   [][]float64
	RiemannTensor  [][][][]float64
	EinsteinTensor [][]float64
	Dimension      int
}

func NewGravitationalTensor(dim int) *GravitationalTensor {
	return &GravitationalTensor{
		MetricTensor:   make([][]float64, dim),
		EinsteinTensor: make([][]float64, dim),
		Dimension:      dim,
	}
}

func (gt *GravitationalTensor) ComputeCurvature() [][]float64 {
	curvature := make([][]float64, gt.Dimension)
	for i := range curvature {
		curvature[i] = make([]float64, gt.Dimension)
	}
	return curvature
}

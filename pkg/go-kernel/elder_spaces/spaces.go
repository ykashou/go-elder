package elder_spaces

type ElderSpace struct {
	Dimension   int
	Basis       [][]float64
	Operations  map[string]func([]float64, []float64) []float64
	Metric      [][]float64
}

func NewElderSpace(dim int) *ElderSpace {
	return &ElderSpace{
		Dimension:  dim,
		Basis:      make([][]float64, dim),
		Operations: make(map[string]func([]float64, []float64) []float64),
		Metric:     make([][]float64, dim),
	}
}

func (es *ElderSpace) AddOperation(name string, op func([]float64, []float64) []float64) {
	es.Operations[name] = op
}

func (es *ElderSpace) ApplyOperation(name string, v1, v2 []float64) []float64 {
	if op, exists := es.Operations[name]; exists {
		return op(v1, v2)
	}
	return nil
}

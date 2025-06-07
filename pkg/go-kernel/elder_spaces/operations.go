package elder_spaces

type ElderOperation struct {
	Name     string
	Operator func([]float64, []float64) []float64
	Identity []float64
}

type OperationRegistry struct {
	Operations map[string]ElderOperation
}

func NewOperationRegistry() *OperationRegistry {
	registry := &OperationRegistry{
		Operations: make(map[string]ElderOperation),
	}
	
	registry.registerDefaultOperations()
	return registry
}

func (or *OperationRegistry) registerDefaultOperations() {
	or.Operations["direct_sum"] = ElderOperation{
		Name: "direct_sum",
		Operator: func(a, b []float64) []float64 {
			result := make([]float64, len(a)+len(b))
			copy(result[:len(a)], a)
			copy(result[len(a):], b)
			return result
		},
		Identity: []float64{0},
	}
	
	or.Operations["tensor_product"] = ElderOperation{
		Name: "tensor_product",
		Operator: func(a, b []float64) []float64 {
			result := make([]float64, len(a)*len(b))
			for i := range a {
				for j := range b {
					result[i*len(b)+j] = a[i] * b[j]
				}
			}
			return result
		},
		Identity: []float64{1},
	}
	
	or.Operations["star_product"] = ElderOperation{
		Name: "star_product",
		Operator: func(a, b []float64) []float64 {
			result := make([]float64, len(a))
			for i := range result {
				if i < len(b) {
					result[i] = a[i] * b[i]
				} else {
					result[i] = a[i]
				}
			}
			return result
		},
		Identity: []float64{1},
	}
}

func (or *OperationRegistry) Apply(operationName string, a, b []float64) []float64 {
	if op, exists := or.Operations[operationName]; exists {
		return op.Operator(a, b)
	}
	return a
}

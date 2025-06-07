package operations

import "math"

type TensorOperator struct {
	Operations map[string]Operation
}

type Operation struct {
	Name     string
	Function func([][]float64) []float64
	Arity    int
}

func NewTensorOperator() *TensorOperator {
	to := &TensorOperator{
		Operations: make(map[string]Operation),
	}
	
	to.registerOperations()
	return to
}

func (to *TensorOperator) registerOperations() {
	to.Operations["add"] = Operation{
		Name:     "add",
		Function: to.addOperation,
		Arity:    2,
	}
	
	to.Operations["multiply"] = Operation{
		Name:     "multiply",
		Function: to.multiplyOperation,
		Arity:    2,
	}
	
	to.Operations["contract"] = Operation{
		Name:     "contract",
		Function: to.contractOperation,
		Arity:    2,
	}
	
	to.Operations["outer"] = Operation{
		Name:     "outer",
		Function: to.outerProductOperation,
		Arity:    2,
	}
	
	to.Operations["transform"] = Operation{
		Name:     "transform",
		Function: to.transformOperation,
		Arity:    1,
	}
}

func (to *TensorOperator) addOperation(tensors [][]float64) []float64 {
	if len(tensors) < 2 {
		return []float64{}
	}
	
	tensor1 := tensors[0]
	tensor2 := tensors[1]
	
	minLen := len(tensor1)
	if len(tensor2) < minLen {
		minLen = len(tensor2)
	}
	
	result := make([]float64, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = tensor1[i] + tensor2[i]
	}
	
	return result
}

func (to *TensorOperator) multiplyOperation(tensors [][]float64) []float64 {
	if len(tensors) < 2 {
		return []float64{}
	}
	
	tensor1 := tensors[0]
	tensor2 := tensors[1]
	
	minLen := len(tensor1)
	if len(tensor2) < minLen {
		minLen = len(tensor2)
	}
	
	result := make([]float64, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = tensor1[i] * tensor2[i]
	}
	
	return result
}

func (to *TensorOperator) contractOperation(tensors [][]float64) []float64 {
	if len(tensors) < 2 {
		return []float64{}
	}
	
	tensor1 := tensors[0]
	tensor2 := tensors[1]
	
	minLen := len(tensor1)
	if len(tensor2) < minLen {
		minLen = len(tensor2)
	}
	
	contraction := 0.0
	for i := 0; i < minLen; i++ {
		contraction += tensor1[i] * tensor2[i]
	}
	
	return []float64{contraction}
}

func (to *TensorOperator) outerProductOperation(tensors [][]float64) []float64 {
	if len(tensors) < 2 {
		return []float64{}
	}
	
	tensor1 := tensors[0]
	tensor2 := tensors[1]
	
	result := make([]float64, len(tensor1)*len(tensor2))
	
	for i, val1 := range tensor1 {
		for j, val2 := range tensor2 {
			result[i*len(tensor2)+j] = val1 * val2
		}
	}
	
	return result
}

func (to *TensorOperator) transformOperation(tensors [][]float64) []float64 {
	if len(tensors) < 1 {
		return []float64{}
	}
	
	tensor := tensors[0]
	result := make([]float64, len(tensor))
	
	for i, val := range tensor {
		result[i] = math.Tanh(val)
	}
	
	return result
}

func (to *TensorOperator) Apply(operationName string, tensors [][]float64) []float64 {
	if operation, exists := to.Operations[operationName]; exists {
		if len(tensors) >= operation.Arity {
			return operation.Function(tensors)
		}
	}
	return []float64{}
}

func (to *TensorOperator) ComputeNorm(tensor []float64) float64 {
	sum := 0.0
	for _, val := range tensor {
		sum += val * val
	}
	return math.Sqrt(sum)
}

func (to *TensorOperator) Normalize(tensor []float64) []float64 {
	norm := to.ComputeNorm(tensor)
	if norm == 0 {
		return tensor
	}
	
	result := make([]float64, len(tensor))
	for i, val := range tensor {
		result[i] = val / norm
	}
	
	return result
}

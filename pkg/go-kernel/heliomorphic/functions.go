package heliomorphic

import "math/cmplx"

type HeliomorphicFunction struct {
	Coefficients []complex128
	Domain       complex128
	Range        complex128
}

func (hf *HeliomorphicFunction) Evaluate(z complex128) complex128 {
	result := complex(0, 0)
	power := complex(1, 0)
	for _, coeff := range hf.Coefficients {
		result += coeff * power
		power *= z
	}
	return result
}

func (hf *HeliomorphicFunction) Derivative(z complex128) complex128 {
	result := complex(0, 0)
	power := complex(1, 0)
	for i, coeff := range hf.Coefficients[1:] {
		result += complex(float64(i+1), 0) * coeff * power
		power *= z
	}
	return result
}

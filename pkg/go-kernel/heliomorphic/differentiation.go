package heliomorphic

import "math/cmplx"

type HeliomorphicDifferentiator struct {
	StepSize float64
	Order    int
}

func NewHeliomorphicDifferentiator(stepSize float64, order int) *HeliomorphicDifferentiator {
	return &HeliomorphicDifferentiator{
		StepSize: stepSize,
		Order:    order,
	}
}

func (hd *HeliomorphicDifferentiator) Differentiate(f func(complex128) complex128, z complex128) complex128 {
	h := complex(hd.StepSize, 0)
	
	switch hd.Order {
	case 1:
		return (f(z+h) - f(z-h)) / (2 * h)
	case 2:
		return (f(z+h) - 2*f(z) + f(z-h)) / (h * h)
	default:
		return (f(z+h) - f(z-h)) / (2 * h)
	}
}

func (hd *HeliomorphicDifferentiator) PartialDerivative(f func(complex128) complex128, z complex128, direction complex128) complex128 {
	h := hd.StepSize
	dir := direction / cmplx.Abs(direction)
	step := complex(h, 0) * dir
	
	return (f(z+step) - f(z-step)) / (2 * step)
}

package isomorphisms

import "math/cmplx"

type ElderHeliomorphicMapping struct {
	ElderSpace      [][]float64
	HeliomorphicSpace map[string]func(complex128) complex128
	Mappings        map[string]string
}

func NewElderHeliomorphicMapping() *ElderHeliomorphicMapping {
	return &ElderHeliomorphicMapping{
		ElderSpace:        make([][]float64, 0),
		HeliomorphicSpace: make(map[string]func(complex128) complex128),
		Mappings:          make(map[string]string),
	}
}

func (ehm *ElderHeliomorphicMapping) MapElderToHeliomorphic(elderVector []float64, id string) func(complex128) complex128 {
	return func(z complex128) complex128 {
		result := complex(0, 0)
		for i, val := range elderVector {
			power := complex(float64(i), 0)
			result += complex(val, 0) * cmplx.Pow(z, power)
		}
		return result
	}
}

func (ehm *ElderHeliomorphicMapping) MapHeliomorphicToElder(f func(complex128) complex128, samples []complex128) []float64 {
	elderVector := make([]float64, len(samples))
	
	for i, z := range samples {
		value := f(z)
		elderVector[i] = real(value)
	}
	
	return elderVector
}

func (ehm *ElderHeliomorphicMapping) VerifyIsomorphism(elderID, heliomorphicID string) bool {
	if elderIndex := ehm.findElderIndex(elderID); elderIndex >= 0 {
		if heliomorphicFunc, exists := ehm.HeliomorphicSpace[heliomorphicID]; exists {
			return ehm.checkStructurePreservation(elderIndex, heliomorphicFunc)
		}
	}
	return false
}

func (ehm *ElderHeliomorphicMapping) findElderIndex(id string) int {
	for i := range ehm.ElderSpace {
		if len(ehm.ElderSpace[i]) > 0 {
			return i
		}
	}
	return -1
}

func (ehm *ElderHeliomorphicMapping) checkStructurePreservation(elderIndex int, f func(complex128) complex128) bool {
	testPoints := []complex128{complex(1, 0), complex(0, 1), complex(1, 1)}
	
	for _, z := range testPoints {
		if cmplx.IsInf(f(z)) || cmplx.IsNaN(f(z)) {
			return false
		}
	}
	
	return true
}

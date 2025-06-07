package gravitational

type FieldPhaseCoupling struct {
	CouplingStrength float64
	PhaseShift       float64
	Tensor           [][]complex128
}

func (fpc *FieldPhaseCoupling) CalculateCoupling(field1, field2 *Field) complex128 {
	return complex(field1.Strength*field2.Strength*fpc.CouplingStrength, fpc.PhaseShift)
}

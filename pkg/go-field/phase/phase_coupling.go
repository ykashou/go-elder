package phase

import "math/cmplx"

type PhaseCoupling struct {
	SourceField string
	TargetField string
	Strength    float64
	Type        string
}

type CouplingMatrix struct {
	Couplings map[string]map[string]float64
	Fields    []string
}

func NewCouplingMatrix(fields []string) *CouplingMatrix {
	cm := &CouplingMatrix{
		Couplings: make(map[string]map[string]float64),
		Fields:    make([]string, len(fields)),
	}
	copy(cm.Fields, fields)
	
	for _, field := range fields {
		cm.Couplings[field] = make(map[string]float64)
	}
	
	return cm
}

func (cm *CouplingMatrix) SetCoupling(source, target string, strength float64) {
	if cm.Couplings[source] == nil {
		cm.Couplings[source] = make(map[string]float64)
	}
	cm.Couplings[source][target] = strength
}

func (cm *CouplingMatrix) CalculateCoupledEvolution(fields map[string]PhaseField, deltaTime float64) map[string]complex128 {
	newPhases := make(map[string]complex128)
	
	for _, fieldID := range cm.Fields {
		field := fields[fieldID]
		coupledPhase := field.Phase
		
		for targetID, strength := range cm.Couplings[fieldID] {
			if targetField, exists := fields[targetID]; exists {
				coupling := complex(strength*deltaTime, 0) * targetField.Phase
				coupledPhase += coupling
			}
		}
		
		newPhases[fieldID] = coupledPhase
	}
	
	return newPhases
}

func (cm *CouplingMatrix) CalculateCouplingEnergy(fields map[string]PhaseField) float64 {
	energy := 0.0
	
	for sourceID, targets := range cm.Couplings {
		sourceField := fields[sourceID]
		
		for targetID, strength := range targets {
			targetField := fields[targetID]
			coupling := sourceField.Phase * cmplx.Conj(targetField.Phase)
			energy += strength * real(coupling)
		}
	}
	
	return energy
}

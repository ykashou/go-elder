package physical

import "math"

type ConservationChecker struct {
	EnergyTolerance    float64
	MomentumTolerance  float64
	AngularTolerance   float64
}

func NewConservationChecker() *ConservationChecker {
	return &ConservationChecker{
		EnergyTolerance:   1e-6,
		MomentumTolerance: 1e-6,
		AngularTolerance:  1e-6,
	}
}

func (cc *ConservationChecker) CheckEnergyConservation(initialEnergy, finalEnergy float64) bool {
	return math.Abs(initialEnergy-finalEnergy) < cc.EnergyTolerance
}

func (cc *ConservationChecker) CheckMomentumConservation(initialMomentum, finalMomentum []float64) bool {
	for i := range initialMomentum {
		if math.Abs(initialMomentum[i]-finalMomentum[i]) > cc.MomentumTolerance {
			return false
		}
	}
	return true
}

type ConservationReport struct {
	EnergyConserved   bool
	MomentumConserved bool
	AngularConserved  bool
	Violations        []string
}

func (cc *ConservationChecker) GenerateReport() *ConservationReport {
	return &ConservationReport{
		Violations: make([]string, 0),
	}
}

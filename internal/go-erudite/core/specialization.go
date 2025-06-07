package core

type SpecializationManager struct {
	Domain         string
	Specializations map[string]float64
	Expertise      float64
}

func (sm *SpecializationManager) DevelopSpecialization(area string, intensity float64) {
	if sm.Specializations == nil {
		sm.Specializations = make(map[string]float64)
	}
	sm.Specializations[area] = intensity
}

func (sm *SpecializationManager) GetExpertiseLevel(area string) float64 {
	if level, exists := sm.Specializations[area]; exists {
		return level
	}
	return 0.0
}

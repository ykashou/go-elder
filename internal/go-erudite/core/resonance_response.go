package core

type ResonanceResponseMechanism struct {
	FrequencyRange [2]float64
	Sensitivity    float64
	Response       map[float64]float64
}

func (rrm *ResonanceResponseMechanism) RespondToResonance(frequency float64) float64 {
	if frequency >= rrm.FrequencyRange[0] && frequency <= rrm.FrequencyRange[1] {
		return rrm.Sensitivity * frequency
	}
	return 0.0
}

func (rrm *ResonanceResponseMechanism) CalibrateResponse(frequency, expectedResponse float64) {
	if rrm.Response == nil {
		rrm.Response = make(map[float64]float64)
	}
	rrm.Response[frequency] = expectedResponse
}

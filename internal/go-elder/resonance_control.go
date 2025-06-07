package elder

type ResonanceController struct {
	Elder           *Elder
	ResonanceFields map[string]float64
	Frequency       float64
	Amplitude       float64
}

func (rc *ResonanceController) InitializeResonance(frequency, amplitude float64) {
	rc.Frequency = frequency
	rc.Amplitude = amplitude
	rc.ResonanceFields = make(map[string]float64)
}

func (rc *ResonanceController) ModulateResonance(fieldID string, modulation float64) {
	rc.ResonanceFields[fieldID] = modulation
}

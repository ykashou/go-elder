package audio

type AudioEventDetectionErudite struct {
	ID         string
	EventTypes []string
	Threshold  float64
}

func (aede *AudioEventDetectionErudite) DetectEvents(audioSignal []float64) []string {
	events := []string{}
	energy := aede.calculateEnergy(audioSignal)
	if energy > aede.Threshold {
		events = append(events, "high_energy_event")
	}
	return events
}

func (aede *AudioEventDetectionErudite) calculateEnergy(signal []float64) float64 {
	var energy float64
	for _, sample := range signal {
		energy += sample * sample
	}
	return energy / float64(len(signal))
}

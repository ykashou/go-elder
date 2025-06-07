package orbital

import "math"

type ResonanceAnalyzer struct {
	Bodies     []OrbitalBody
	Resonances map[string]Resonance
	Tolerance  float64
}

type OrbitalBody struct {
	ID           string
	Mass         float64
	SemiMajorAxis float64
	Period       float64
	Eccentricity float64
}

type Resonance struct {
	Body1    string
	Body2    string
	Ratio    [2]int
	Strength float64
	Type     string
}

func NewResonanceAnalyzer(tolerance float64) *ResonanceAnalyzer {
	return &ResonanceAnalyzer{
		Bodies:     make([]OrbitalBody, 0),
		Resonances: make(map[string]Resonance),
		Tolerance:  tolerance,
	}
}

func (ra *ResonanceAnalyzer) AddBody(id string, mass, semiMajorAxis, eccentricity float64) {
	period := ra.calculatePeriod(semiMajorAxis, mass)
	body := OrbitalBody{
		ID:            id,
		Mass:          mass,
		SemiMajorAxis: semiMajorAxis,
		Period:        period,
		Eccentricity:  eccentricity,
	}
	ra.Bodies = append(ra.Bodies, body)
}

func (ra *ResonanceAnalyzer) calculatePeriod(semiMajorAxis, mass float64) float64 {
	G := 6.67430e-11
	centralMass := 1e24
	return 2 * math.Pi * math.Sqrt(math.Pow(semiMajorAxis, 3)/(G*(centralMass+mass)))
}

func (ra *ResonanceAnalyzer) DetectResonances() []Resonance {
	resonances := make([]Resonance, 0)
	
	for i := 0; i < len(ra.Bodies); i++ {
		for j := i + 1; j < len(ra.Bodies); j++ {
			body1 := ra.Bodies[i]
			body2 := ra.Bodies[j]
			
			if resonance := ra.checkResonance(body1, body2); resonance != nil {
				resonances = append(resonances, *resonance)
			}
		}
	}
	
	return resonances
}

func (ra *ResonanceAnalyzer) checkResonance(body1, body2 OrbitalBody) *Resonance {
	ratio := body1.Period / body2.Period
	
	simpleRatios := [][2]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}, {2, 5}, {3, 5}}
	
	for _, simpleRatio := range simpleRatios {
		expectedRatio := float64(simpleRatio[0]) / float64(simpleRatio[1])
		if math.Abs(ratio-expectedRatio) < ra.Tolerance {
			strength := ra.calculateResonanceStrength(body1, body2, simpleRatio)
			return &Resonance{
				Body1:    body1.ID,
				Body2:    body2.ID,
				Ratio:    simpleRatio,
				Strength: strength,
				Type:     "mean_motion",
			}
		}
	}
	
	return nil
}

func (ra *ResonanceAnalyzer) calculateResonanceStrength(body1, body2 OrbitalBody, ratio [2]int) float64 {
	massRatio := (body1.Mass + body2.Mass) / 1e24
	eccentricityFactor := math.Sqrt(body1.Eccentricity*body1.Eccentricity + body2.Eccentricity*body2.Eccentricity)
	return massRatio * eccentricityFactor * 0.1
}

package entropy

import "math"

type FieldEntropyCalculator struct {
	Fields    map[string]Field
	Tolerance float64
}

type Field struct {
	ID      string
	Values  []float64
	Density map[float64]float64
	Volume  float64
}

func NewFieldEntropyCalculator(tolerance float64) *FieldEntropyCalculator {
	return &FieldEntropyCalculator{
		Fields:    make(map[string]Field),
		Tolerance: tolerance,
	}
}

func (fec *FieldEntropyCalculator) AddField(id string, values []float64, volume float64) {
	field := Field{
		ID:      id,
		Values:  make([]float64, len(values)),
		Density: make(map[float64]float64),
		Volume:  volume,
	}
	copy(field.Values, values)
	
	field.Density = fec.calculateDensity(values)
	fec.Fields[id] = field
}

func (fec *FieldEntropyCalculator) calculateDensity(values []float64) map[float64]float64 {
	density := make(map[float64]float64)
	total := float64(len(values))
	
	for _, value := range values {
		rounded := math.Round(value/fec.Tolerance) * fec.Tolerance
		density[rounded]++
	}
	
	for key := range density {
		density[key] /= total
	}
	
	return density
}

func (fec *FieldEntropyCalculator) CalculateEntropy(fieldID string) float64 {
	field := fec.Fields[fieldID]
	entropy := 0.0
	
	for _, probability := range field.Density {
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}
	
	return entropy
}

func (fec *FieldEntropyCalculator) CalculateRelativeEntropy(field1ID, field2ID string) float64 {
	field1 := fec.Fields[field1ID]
	field2 := fec.Fields[field2ID]
	
	relativeEntropy := 0.0
	
	for value, p1 := range field1.Density {
		if p2, exists := field2.Density[value]; exists && p2 > 0 {
			relativeEntropy += p1 * math.Log2(p1/p2)
		}
	}
	
	return relativeEntropy
}

func (fec *FieldEntropyCalculator) CalculateMutualInformation(field1ID, field2ID string) float64 {
	entropy1 := fec.CalculateEntropy(field1ID)
	entropy2 := fec.CalculateEntropy(field2ID)
	jointEntropy := fec.calculateJointEntropy(field1ID, field2ID)
	
	return entropy1 + entropy2 - jointEntropy
}

func (fec *FieldEntropyCalculator) calculateJointEntropy(field1ID, field2ID string) float64 {
	field1 := fec.Fields[field1ID]
	field2 := fec.Fields[field2ID]
	
	jointDensity := make(map[[2]float64]float64)
	total := float64(len(field1.Values))
	
	for i := range field1.Values {
		if i < len(field2.Values) {
			key := [2]float64{
				math.Round(field1.Values[i]/fec.Tolerance) * fec.Tolerance,
				math.Round(field2.Values[i]/fec.Tolerance) * fec.Tolerance,
			}
			jointDensity[key]++
		}
	}
	
	entropy := 0.0
	for _, count := range jointDensity {
		probability := count / total
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}
	
	return entropy
}

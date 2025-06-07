package isomorphisms

type ParameterMapping struct {
	SourceSpace map[string]float64
	TargetSpace map[string]float64
	MappingRules map[string]MappingRule
}

type MappingRule struct {
	SourceParam string
	TargetParam string
	Transform   func(float64) float64
	Inverse     func(float64) float64
}

func NewParameterMapping() *ParameterMapping {
	return &ParameterMapping{
		SourceSpace:  make(map[string]float64),
		TargetSpace:  make(map[string]float64),
		MappingRules: make(map[string]MappingRule),
	}
}

func (pm *ParameterMapping) AddMapping(sourceParam, targetParam string, transform, inverse func(float64) float64) {
	rule := MappingRule{
		SourceParam: sourceParam,
		TargetParam: targetParam,
		Transform:   transform,
		Inverse:     inverse,
	}
	pm.MappingRules[sourceParam] = rule
}

func (pm *ParameterMapping) MapForward(sourceParams map[string]float64) map[string]float64 {
	targetParams := make(map[string]float64)
	
	for sourceParam, value := range sourceParams {
		if rule, exists := pm.MappingRules[sourceParam]; exists {
			transformedValue := rule.Transform(value)
			targetParams[rule.TargetParam] = transformedValue
		}
	}
	
	return targetParams
}

func (pm *ParameterMapping) MapBackward(targetParams map[string]float64) map[string]float64 {
	sourceParams := make(map[string]float64)
	
	for _, rule := range pm.MappingRules {
		if value, exists := targetParams[rule.TargetParam]; exists {
			originalValue := rule.Inverse(value)
			sourceParams[rule.SourceParam] = originalValue
		}
	}
	
	return sourceParams
}

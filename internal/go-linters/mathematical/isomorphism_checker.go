package mathematical

type IsomorphismChecker struct {
	SourceStructure map[string]interface{}
	TargetStructure map[string]interface{}
	Mappings       map[string]string
	Tolerance      float64
}

func NewIsomorphismChecker(tolerance float64) *IsomorphismChecker {
	return &IsomorphismChecker{
		SourceStructure: make(map[string]interface{}),
		TargetStructure: make(map[string]interface{}),
		Mappings:        make(map[string]string),
		Tolerance:       tolerance,
	}
}

func (ic *IsomorphismChecker) SetStructures(source, target map[string]interface{}) {
	ic.SourceStructure = source
	ic.TargetStructure = target
}

func (ic *IsomorphismChecker) CheckIsomorphism() IsomorphismResult {
	result := IsomorphismResult{
		IsIsomorphic: true,
		Mappings:     make(map[string]string),
		Violations:   make([]string, 0),
	}
	
	if len(ic.SourceStructure) != len(ic.TargetStructure) {
		result.IsIsomorphic = false
		result.Violations = append(result.Violations, "Structures have different cardinality")
		return result
	}
	
	mapping := ic.findMapping()
	if mapping == nil {
		result.IsIsomorphic = false
		result.Violations = append(result.Violations, "No valid mapping found")
		return result
	}
	
	result.Mappings = mapping
	
	if !ic.preservesOperation(mapping) {
		result.IsIsomorphic = false
		result.Violations = append(result.Violations, "Mapping does not preserve operations")
	}
	
	return result
}

type IsomorphismResult struct {
	IsIsomorphic bool
	Mappings     map[string]string
	Violations   []string
}

func (ic *IsomorphismChecker) findMapping() map[string]string {
	mapping := make(map[string]string)
	used := make(map[string]bool)
	
	for sourceKey := range ic.SourceStructure {
		found := false
		for targetKey := range ic.TargetStructure {
			if !used[targetKey] && ic.elementsCompatible(sourceKey, targetKey) {
				mapping[sourceKey] = targetKey
				used[targetKey] = true
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	
	return mapping
}

func (ic *IsomorphismChecker) elementsCompatible(sourceKey, targetKey string) bool {
	sourceVal := ic.SourceStructure[sourceKey]
	targetVal := ic.TargetStructure[targetKey]
	
	return ic.valuesCompatible(sourceVal, targetVal)
}

func (ic *IsomorphismChecker) valuesCompatible(source, target interface{}) bool {
	switch s := source.(type) {
	case float64:
		if t, ok := target.(float64); ok {
			return (s-t)*(s-t) < ic.Tolerance*ic.Tolerance
		}
	case string:
		if t, ok := target.(string); ok {
			return len(s) == len(t)
		}
	case []float64:
		if t, ok := target.([]float64); ok {
			return len(s) == len(t)
		}
	}
	return false
}

func (ic *IsomorphismChecker) preservesOperation(mapping map[string]string) bool {
	for sourceKey1, targetKey1 := range mapping {
		for sourceKey2, targetKey2 := range mapping {
			if sourceKey1 != sourceKey2 {
				sourceResult := ic.applyOperation(sourceKey1, sourceKey2, ic.SourceStructure)
				targetResult := ic.applyOperation(targetKey1, targetKey2, ic.TargetStructure)
				
				if !ic.valuesCompatible(sourceResult, targetResult) {
					return false
				}
			}
		}
	}
	return true
}

func (ic *IsomorphismChecker) applyOperation(key1, key2 string, structure map[string]interface{}) interface{} {
	val1 := structure[key1]
	val2 := structure[key2]
	
	switch v1 := val1.(type) {
	case float64:
		if v2, ok := val2.(float64); ok {
			return v1 + v2
		}
	case []float64:
		if v2, ok := val2.([]float64); ok {
			if len(v1) == len(v2) {
				result := make([]float64, len(v1))
				for i := range v1 {
					result[i] = v1[i] + v2[i]
				}
				return result
			}
		}
	}
	
	return nil
}

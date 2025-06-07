// Package transfer implements domain mapping protocols
package transfer

// DomainMappingProtocol handles mappings between different domains
type DomainMappingProtocol struct {
	MappingRules map[string]MappingRule
	Domains      []string
}

// MappingRule defines how concepts map between domains
type MappingRule struct {
	SourceConcept string
	TargetConcept string
	Confidence    float64
	Bidirectional bool
}

// CreateDomainMapping creates a new domain mapping
func (dmp *DomainMappingProtocol) CreateDomainMapping(source, target string, confidence float64, bidirectional bool) {
	rule := MappingRule{
		SourceConcept: source,
		TargetConcept: target,
		Confidence:    confidence,
		Bidirectional: bidirectional,
	}
	dmp.MappingRules[source] = rule
	
	if bidirectional {
		reverseRule := MappingRule{
			SourceConcept: target,
			TargetConcept: source,
			Confidence:    confidence,
			Bidirectional: true,
		}
		dmp.MappingRules[target] = reverseRule
	}
}

// GetMapping retrieves mapping for a concept
func (dmp *DomainMappingProtocol) GetMapping(concept string) *MappingRule {
	if rule, exists := dmp.MappingRules[concept]; exists {
		return &rule
	}
	return nil
}
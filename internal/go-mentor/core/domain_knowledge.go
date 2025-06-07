// Package core implements domain-specific knowledge management
package core

// DomainKnowledgeManager handles domain-specific knowledge operations
type DomainKnowledgeManager struct {
	Knowledge *DomainKnowledge
}

// AddConcept adds a new concept to the domain knowledge
func (dkm *DomainKnowledgeManager) AddConcept(name string, concept interface{}) {
	dkm.Knowledge.Concepts[name] = concept
}

// LinkConcepts creates relationships between concepts
func (dkm *DomainKnowledgeManager) LinkConcepts(concept1, concept2 string) {
	if dkm.Knowledge.Relationships[concept1] == nil {
		dkm.Knowledge.Relationships[concept1] = make([]string, 0)
	}
	dkm.Knowledge.Relationships[concept1] = append(dkm.Knowledge.Relationships[concept1], concept2)
}

// GetRelatedConcepts returns concepts related to a given concept
func (dkm *DomainKnowledgeManager) GetRelatedConcepts(concept string) []string {
	return dkm.Knowledge.Relationships[concept]
}

// AddPrinciple adds a domain principle
func (dkm *DomainKnowledgeManager) AddPrinciple(principle string) {
	dkm.Knowledge.Principles = append(dkm.Knowledge.Principles, principle)
}
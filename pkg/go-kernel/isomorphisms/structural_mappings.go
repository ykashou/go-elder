package isomorphisms

type StructuralMapping struct {
	SourceStructure map[string]StructureElement
	TargetStructure map[string]StructureElement
	Correspondences map[string]string
}

type StructureElement struct {
	ID          string
	Type        string
	Properties  map[string]interface{}
	Connections []string
}

func NewStructuralMapping() *StructuralMapping {
	return &StructuralMapping{
		SourceStructure: make(map[string]StructureElement),
		TargetStructure: make(map[string]StructureElement),
		Correspondences: make(map[string]string),
	}
}

func (sm *StructuralMapping) AddSourceElement(id, elementType string, properties map[string]interface{}) {
	element := StructureElement{
		ID:          id,
		Type:        elementType,
		Properties:  make(map[string]interface{}),
		Connections: make([]string, 0),
	}
	
	for k, v := range properties {
		element.Properties[k] = v
	}
	
	sm.SourceStructure[id] = element
}

func (sm *StructuralMapping) AddTargetElement(id, elementType string, properties map[string]interface{}) {
	element := StructureElement{
		ID:          id,
		Type:        elementType,
		Properties:  make(map[string]interface{}),
		Connections: make([]string, 0),
	}
	
	for k, v := range properties {
		element.Properties[k] = v
	}
	
	sm.TargetStructure[id] = element
}

func (sm *StructuralMapping) EstablishCorrespondence(sourceID, targetID string) bool {
	sourceElement, sourceExists := sm.SourceStructure[sourceID]
	targetElement, targetExists := sm.TargetStructure[targetID]
	
	if !sourceExists || !targetExists {
		return false
	}
	
	if sourceElement.Type != targetElement.Type {
		return false
	}
	
	sm.Correspondences[sourceID] = targetID
	return true
}

func (sm *StructuralMapping) VerifyStructuralIsomorphism() bool {
	if len(sm.SourceStructure) != len(sm.TargetStructure) {
		return false
	}
	
	for sourceID, targetID := range sm.Correspondences {
		if !sm.verifyElementCorrespondence(sourceID, targetID) {
			return false
		}
	}
	
	return true
}

func (sm *StructuralMapping) verifyElementCorrespondence(sourceID, targetID string) bool {
	sourceElement := sm.SourceStructure[sourceID]
	targetElement := sm.TargetStructure[targetID]
	
	if len(sourceElement.Connections) != len(targetElement.Connections) {
		return false
	}
	
	for _, sourceConnection := range sourceElement.Connections {
		if targetConnection, exists := sm.Correspondences[sourceConnection]; exists {
			if !sm.containsConnection(targetElement.Connections, targetConnection) {
				return false
			}
		} else {
			return false
		}
	}
	
	return true
}

func (sm *StructuralMapping) containsConnection(connections []string, target string) bool {
	for _, connection := range connections {
		if connection == target {
			return true
		}
	}
	return false
}

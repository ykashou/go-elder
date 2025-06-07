package memory

type FieldBasedMemory struct {
	Fields       map[string]MemoryField
	Associations map[string][]string
	FieldTypes   map[string]string
}

type MemoryField struct {
	ID          string
	Type        string
	Strength    float64
	Content     interface{}
	Connections []string
}

func NewFieldBasedMemory() *FieldBasedMemory {
	return &FieldBasedMemory{
		Fields:       make(map[string]MemoryField),
		Associations: make(map[string][]string),
		FieldTypes:   make(map[string]string),
	}
}

func (fbm *FieldBasedMemory) CreateField(id, fieldType string, strength float64, content interface{}) {
	field := MemoryField{
		ID:          id,
		Type:        fieldType,
		Strength:    strength,
		Content:     content,
		Connections: make([]string, 0),
	}
	
	fbm.Fields[id] = field
	fbm.FieldTypes[id] = fieldType
}

func (fbm *FieldBasedMemory) AssociateFields(field1, field2 string) {
	if fbm.Associations[field1] == nil {
		fbm.Associations[field1] = make([]string, 0)
	}
	fbm.Associations[field1] = append(fbm.Associations[field1], field2)
	
	if fbm.Associations[field2] == nil {
		fbm.Associations[field2] = make([]string, 0)
	}
	fbm.Associations[field2] = append(fbm.Associations[field2], field1)
}

func (fbm *FieldBasedMemory) RetrieveByAssociation(fieldID string) []MemoryField {
	associated := make([]MemoryField, 0)
	if associations, exists := fbm.Associations[fieldID]; exists {
		for _, assocID := range associations {
			if field, exists := fbm.Fields[assocID]; exists {
				associated = append(associated, field)
			}
		}
	}
	return associated
}

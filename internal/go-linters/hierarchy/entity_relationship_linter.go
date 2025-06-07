package hierarchy

type EntityRelationshipLinter struct {
	Entities      map[string]Entity
	Relationships map[string][]Relationship
	Rules         []ValidationRule
}

type Entity struct {
	ID       string
	Type     string
	Level    int
	Parent   string
	Children []string
}

type Relationship struct {
	From string
	To   string
	Type string
}

type ValidationRule struct {
	Name        string
	Description string
	Validator   func(Entity, []Entity) bool
}

func NewEntityRelationshipLinter() *EntityRelationshipLinter {
	return &EntityRelationshipLinter{
		Entities:      make(map[string]Entity),
		Relationships: make(map[string][]Relationship),
		Rules:         make([]ValidationRule, 0),
	}
}

func (erl *EntityRelationshipLinter) ValidateHierarchy() []string {
	violations := []string{}
	
	for _, entity := range erl.Entities {
		if !erl.validateEntityLevel(entity) {
			violations = append(violations, "Invalid hierarchy level for entity: "+entity.ID)
		}
	}
	
	return violations
}

func (erl *EntityRelationshipLinter) validateEntityLevel(entity Entity) bool {
	if entity.Parent == "" {
		return entity.Level == 0
	}
	
	parent, exists := erl.Entities[entity.Parent]
	if !exists {
		return false
	}
	
	return entity.Level == parent.Level + 1
}

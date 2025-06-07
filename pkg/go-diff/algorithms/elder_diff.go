package algorithms

type ElderDifferential struct {
	SourceState map[string]interface{}
	TargetState map[string]interface{}
	Changes     []Change
}

type Change struct {
	Type     string
	Path     string
	OldValue interface{}
	NewValue interface{}
}

func (ed *ElderDifferential) ComputeDifference() []Change {
	changes := []Change{}
	for key, newVal := range ed.TargetState {
		if oldVal, exists := ed.SourceState[key]; exists {
			if oldVal != newVal {
				changes = append(changes, Change{
					Type:     "modified",
					Path:     key,
					OldValue: oldVal,
					NewValue: newVal,
				})
			}
		} else {
			changes = append(changes, Change{
				Type:     "added",
				Path:     key,
				NewValue: newVal,
			})
		}
	}
	ed.Changes = changes
	return changes
}

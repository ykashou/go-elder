package coordination

type HierarchyController struct {
	Levels map[int][]string
	ControlMatrix [][]float64
	ActiveEntities map[string]bool
}

func NewHierarchyController() *HierarchyController {
	return &HierarchyController{
		Levels: make(map[int][]string),
		ActiveEntities: make(map[string]bool),
	}
}

func (hc *HierarchyController) RegisterEntity(level int, entityID string) {
	if hc.Levels[level] == nil {
		hc.Levels[level] = make([]string, 0)
	}
	hc.Levels[level] = append(hc.Levels[level], entityID)
	hc.ActiveEntities[entityID] = true
}

func (hc *HierarchyController) ControlHierarchy() {
	for level, entities := range hc.Levels {
		hc.coordinateLevel(level, entities)
	}
}

func (hc *HierarchyController) coordinateLevel(level int, entities []string) {
	for _, entity := range entities {
		hc.ActiveEntities[entity] = true
	}
}

package architecture

type HierarchicalMapping struct {
	Levels map[int][]string
	Mappings map[string]string
}

func (hm *HierarchicalMapping) MapLevels(level int, entities []string) {
	if hm.Levels == nil {
		hm.Levels = make(map[int][]string)
	}
	hm.Levels[level] = entities
}

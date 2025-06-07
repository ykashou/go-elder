package memory

type ElderMemoryMap struct {
	MemoryRegions map[string]MemoryRegion
	TotalCapacity int64
	UsedCapacity  int64
	Mappings      map[string]string
}

type MemoryRegion struct {
	ID       string
	Size     int64
	Type     string
	Data     []byte
	MetaData map[string]interface{}
}

func NewElderMemoryMap(capacity int64) *ElderMemoryMap {
	return &ElderMemoryMap{
		MemoryRegions: make(map[string]MemoryRegion),
		TotalCapacity: capacity,
		Mappings:      make(map[string]string),
	}
}

func (emm *ElderMemoryMap) AllocateRegion(id string, size int64, regionType string) bool {
	if emm.UsedCapacity+size > emm.TotalCapacity {
		return false
	}
	
	region := MemoryRegion{
		ID:       id,
		Size:     size,
		Type:     regionType,
		Data:     make([]byte, size),
		MetaData: make(map[string]interface{}),
	}
	
	emm.MemoryRegions[id] = region
	emm.UsedCapacity += size
	return true
}

func (emm *ElderMemoryMap) AccessRegion(id string) *MemoryRegion {
	if region, exists := emm.MemoryRegions[id]; exists {
		return &region
	}
	return nil
}

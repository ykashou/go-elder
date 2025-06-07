package memory

import "math"

type InfiniteMemorySystem struct {
	Layers      map[int]MemoryLayer
	Compression map[int]float64
	MaxLayers   int
	GrowthRate  float64
}

type MemoryLayer struct {
	Level    int
	Capacity int64
	Used     int64
	Data     map[string][]byte
}

func NewInfiniteMemorySystem(maxLayers int, growthRate float64) *InfiniteMemorySystem {
	return &InfiniteMemorySystem{
		Layers:      make(map[int]MemoryLayer),
		Compression: make(map[int]float64),
		MaxLayers:   maxLayers,
		GrowthRate:  growthRate,
	}
}

func (ims *InfiniteMemorySystem) Store(id string, data []byte) bool {
	layer := ims.findBestLayer(len(data))
	
	if layer == -1 {
		layer = ims.createNewLayer()
	}
	
	memLayer := ims.Layers[layer]
	memLayer.Data[id] = make([]byte, len(data))
	copy(memLayer.Data[id], data)
	memLayer.Used += int64(len(data))
	ims.Layers[layer] = memLayer
	
	return true
}

func (ims *InfiniteMemorySystem) findBestLayer(dataSize int) int {
	for level, layer := range ims.Layers {
		if layer.Used+int64(dataSize) <= layer.Capacity {
			return level
		}
	}
	return -1
}

func (ims *InfiniteMemorySystem) createNewLayer() int {
	newLevel := len(ims.Layers)
	if newLevel >= ims.MaxLayers {
		return 0
	}
	
	capacity := int64(math.Pow(2, float64(newLevel+10)))
	layer := MemoryLayer{
		Level:    newLevel,
		Capacity: capacity,
		Used:     0,
		Data:     make(map[string][]byte),
	}
	
	ims.Layers[newLevel] = layer
	ims.Compression[newLevel] = 1.0 / math.Pow(ims.GrowthRate, float64(newLevel))
	
	return newLevel
}

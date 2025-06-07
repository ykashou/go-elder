package memory

import "math"

type GravitationalMemoryField struct {
	MemoryNodes map[string]MemoryNode
	FieldStrength float64
	DecayRate     float64
	Capacity      int64
}

type MemoryNode struct {
	ID       string
	Position Vector3D
	Data     []byte
	Weight   float64
	Created  float64
}

type Vector3D struct {
	X, Y, Z float64
}

func NewGravitationalMemoryField(strength, decay float64, capacity int64) *GravitationalMemoryField {
	return &GravitationalMemoryField{
		MemoryNodes:   make(map[string]MemoryNode),
		FieldStrength: strength,
		DecayRate:     decay,
		Capacity:      capacity,
	}
}

func (gmf *GravitationalMemoryField) StoreMemory(id string, data []byte, position Vector3D) bool {
	if int64(len(gmf.MemoryNodes)) >= gmf.Capacity {
		gmf.evictOldestNode()
	}
	
	weight := gmf.calculateWeight(data, position)
	node := MemoryNode{
		ID:       id,
		Position: position,
		Data:     make([]byte, len(data)),
		Weight:   weight,
		Created:  gmf.getCurrentTime(),
	}
	copy(node.Data, data)
	
	gmf.MemoryNodes[id] = node
	return true
}

func (gmf *GravitationalMemoryField) RetrieveMemory(queryPosition Vector3D, radius float64) []MemoryNode {
	results := make([]MemoryNode, 0)
	
	for _, node := range gmf.MemoryNodes {
		distance := gmf.calculateDistance(queryPosition, node.Position)
		if distance <= radius {
			strength := gmf.calculateFieldStrength(distance, node.Weight)
			if strength > 0.1 {
				results = append(results, node)
			}
		}
	}
	
	return results
}

func (gmf *GravitationalMemoryField) calculateWeight(data []byte, position Vector3D) float64 {
	dataWeight := float64(len(data)) / 1024.0
	positionWeight := math.Sqrt(position.X*position.X + position.Y*position.Y + position.Z*position.Z)
	return gmf.FieldStrength * dataWeight / (1.0 + positionWeight)
}

func (gmf *GravitationalMemoryField) calculateDistance(p1, p2 Vector3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (gmf *GravitationalMemoryField) calculateFieldStrength(distance, weight float64) float64 {
	if distance == 0 {
		return weight
	}
	return weight / (distance * distance)
}

func (gmf *GravitationalMemoryField) getCurrentTime() float64 {
	return 1000.0 // Simplified time
}

func (gmf *GravitationalMemoryField) evictOldestNode() {
	var oldestID string
	var oldestTime float64 = math.Inf(1)
	
	for id, node := range gmf.MemoryNodes {
		if node.Created < oldestTime {
			oldestTime = node.Created
			oldestID = id
		}
	}
	
	if oldestID != "" {
		delete(gmf.MemoryNodes, oldestID)
	}
}

func (gmf *GravitationalMemoryField) ApplyDecay(deltaTime float64) {
	toDelete := make([]string, 0)
	
	for id, node := range gmf.MemoryNodes {
		node.Weight *= math.Exp(-gmf.DecayRate * deltaTime)
		if node.Weight < 0.01 {
			toDelete = append(toDelete, id)
		} else {
			gmf.MemoryNodes[id] = node
		}
	}
	
	for _, id := range toDelete {
		delete(gmf.MemoryNodes, id)
	}
}

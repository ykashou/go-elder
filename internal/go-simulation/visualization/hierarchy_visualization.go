package visualization

type HierarchyVisualizer struct {
	Entities     map[string]HierarchyEntity
	Connections  []Connection
	Levels       map[int][]string
	LayoutType   string
}

type HierarchyEntity struct {
	ID       string
	Level    int
	Type     string
	Position Point2D
	Size     float64
	Color    string
	Status   string
}

type Point2D struct {
	X, Y float64
}

type Connection struct {
	From     string
	To       string
	Type     string
	Strength float64
}

func NewHierarchyVisualizer(layoutType string) *HierarchyVisualizer {
	return &HierarchyVisualizer{
		Entities:    make(map[string]HierarchyEntity),
		Connections: make([]Connection, 0),
		Levels:      make(map[int][]string),
		LayoutType:  layoutType,
	}
}

func (hv *HierarchyVisualizer) AddEntity(id string, level int, entityType, color, status string) {
	entity := HierarchyEntity{
		ID:     id,
		Level:  level,
		Type:   entityType,
		Size:   hv.calculateEntitySize(level),
		Color:  color,
		Status: status,
	}
	
	hv.Entities[id] = entity
	
	if hv.Levels[level] == nil {
		hv.Levels[level] = make([]string, 0)
	}
	hv.Levels[level] = append(hv.Levels[level], id)
}

func (hv *HierarchyVisualizer) calculateEntitySize(level int) float64 {
	baseSize := 10.0
	return baseSize + float64(3-level)*5.0
}

func (hv *HierarchyVisualizer) AddConnection(from, to, connectionType string, strength float64) {
	connection := Connection{
		From:     from,
		To:       to,
		Type:     connectionType,
		Strength: strength,
	}
	hv.Connections = append(hv.Connections, connection)
}

func (hv *HierarchyVisualizer) CalculateLayout() {
	switch hv.LayoutType {
	case "tree":
		hv.calculateTreeLayout()
	case "circular":
		hv.calculateCircularLayout()
	case "force":
		hv.calculateForceLayout()
	default:
		hv.calculateTreeLayout()
	}
}

func (hv *HierarchyVisualizer) calculateTreeLayout() {
	levelHeight := 100.0
	
	for level, entities := range hv.Levels {
		entityWidth := 800.0 / float64(len(entities))
		yPos := float64(level) * levelHeight
		
		for i, entityID := range entities {
			entity := hv.Entities[entityID]
			entity.Position = Point2D{
				X: float64(i)*entityWidth + entityWidth/2,
				Y: yPos,
			}
			hv.Entities[entityID] = entity
		}
	}
}

func (hv *HierarchyVisualizer) calculateCircularLayout() {
	centerX, centerY := 400.0, 300.0
	
	for level, entities := range hv.Levels {
		radius := float64(level+1) * 80.0
		angleStep := 6.28318 / float64(len(entities))
		
		for i, entityID := range entities {
			angle := float64(i) * angleStep
			entity := hv.Entities[entityID]
			entity.Position = Point2D{
				X: centerX + radius*math.Cos(angle),
				Y: centerY + radius*math.Sin(angle),
			}
			hv.Entities[entityID] = entity
		}
	}
}

func (hv *HierarchyVisualizer) calculateForceLayout() {
	iterations := 100
	repulsion := 1000.0
	attraction := 0.1
	damping := 0.9
	
	forces := make(map[string]Point2D)
	
	for iter := 0; iter < iterations; iter++ {
		for id := range hv.Entities {
			forces[id] = Point2D{0, 0}
		}
		
		for id1, entity1 := range hv.Entities {
			for id2, entity2 := range hv.Entities {
				if id1 != id2 {
					dx := entity1.Position.X - entity2.Position.X
					dy := entity1.Position.Y - entity2.Position.Y
					distance := math.Sqrt(dx*dx + dy*dy)
					
					if distance > 0 {
						force := repulsion / (distance * distance)
						forces[id1] = Point2D{
							X: forces[id1].X + force*dx/distance,
							Y: forces[id1].Y + force*dy/distance,
						}
					}
				}
			}
		}
		
		for _, conn := range hv.Connections {
			entity1 := hv.Entities[conn.From]
			entity2 := hv.Entities[conn.To]
			
			dx := entity2.Position.X - entity1.Position.X
			dy := entity2.Position.Y - entity1.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			
			if distance > 0 {
				force := attraction * distance * conn.Strength
				forces[conn.From] = Point2D{
					X: forces[conn.From].X + force*dx/distance,
					Y: forces[conn.From].Y + force*dy/distance,
				}
				forces[conn.To] = Point2D{
					X: forces[conn.To].X - force*dx/distance,
					Y: forces[conn.To].Y - force*dy/distance,
				}
			}
		}
		
		for id, entity := range hv.Entities {
			entity.Position.X += forces[id].X * damping
			entity.Position.Y += forces[id].Y * damping
			hv.Entities[id] = entity
		}
	}
}

func (hv *HierarchyVisualizer) GenerateVisualizationData() map[string]interface{} {
	data := make(map[string]interface{})
	
	entityData := make([]map[string]interface{}, 0)
	for _, entity := range hv.Entities {
		entityInfo := map[string]interface{}{
			"id":       entity.ID,
			"level":    entity.Level,
			"type":     entity.Type,
			"position": entity.Position,
			"size":     entity.Size,
			"color":    entity.Color,
			"status":   entity.Status,
		}
		entityData = append(entityData, entityInfo)
	}
	
	connectionData := make([]map[string]interface{}, 0)
	for _, conn := range hv.Connections {
		connInfo := map[string]interface{}{
			"from":     conn.From,
			"to":       conn.To,
			"type":     conn.Type,
			"strength": conn.Strength,
		}
		connectionData = append(connectionData, connInfo)
	}
	
	data["entities"] = entityData
	data["connections"] = connectionData
	data["layout_type"] = hv.LayoutType
	
	return data
}

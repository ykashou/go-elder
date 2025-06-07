package visualization

import "math"

type FieldVisualizer struct {
	GravitationalFields []GravField
	GridResolution      int
	BoundingBox         BoundingBox
	FieldLines          []FieldLine
}

type GravField struct {
	Position  Point3D
	Strength  float64
	Direction Point3D
	Range     float64
}

type BoundingBox struct {
	MinX, MaxX float64
	MinY, MaxY float64
	MinZ, MaxZ float64
}

type FieldLine struct {
	Points    []Point3D
	Strength  float64
	Direction Point3D
}

func NewFieldVisualizer(resolution int, bbox BoundingBox) *FieldVisualizer {
	return &FieldVisualizer{
		GravitationalFields: make([]GravField, 0),
		GridResolution:      resolution,
		BoundingBox:         bbox,
		FieldLines:          make([]FieldLine, 0),
	}
}

func (fv *FieldVisualizer) AddGravitationalField(pos Point3D, strength float64, dir Point3D, fieldRange float64) {
	field := GravField{
		Position:  pos,
		Strength:  strength,
		Direction: dir,
		Range:     fieldRange,
	}
	fv.GravitationalFields = append(fv.GravitationalFields, field)
}

func (fv *FieldVisualizer) GenerateFieldLines() {
	fv.FieldLines = make([]FieldLine, 0)
	
	stepX := (fv.BoundingBox.MaxX - fv.BoundingBox.MinX) / float64(fv.GridResolution)
	stepY := (fv.BoundingBox.MaxY - fv.BoundingBox.MinY) / float64(fv.GridResolution)
	
	for i := 0; i < fv.GridResolution; i++ {
		for j := 0; j < fv.GridResolution; j++ {
			startPoint := Point3D{
				X: fv.BoundingBox.MinX + float64(i)*stepX,
				Y: fv.BoundingBox.MinY + float64(j)*stepY,
				Z: 0,
			}
			
			fieldLine := fv.traceFieldLine(startPoint)
			if len(fieldLine.Points) > 1 {
				fv.FieldLines = append(fv.FieldLines, fieldLine)
			}
		}
	}
}

func (fv *FieldVisualizer) traceFieldLine(start Point3D) FieldLine {
	line := FieldLine{
		Points: make([]Point3D, 0),
	}
	
	current := start
	stepSize := 0.1
	maxSteps := 100
	
	for step := 0; step < maxSteps; step++ {
		line.Points = append(line.Points, current)
		
		fieldVector := fv.calculateFieldAtPoint(current)
		magnitude := math.Sqrt(fieldVector.X*fieldVector.X + 
							  fieldVector.Y*fieldVector.Y + 
							  fieldVector.Z*fieldVector.Z)
		
		if magnitude < 1e-6 {
			break
		}
		
		current.X += (fieldVector.X / magnitude) * stepSize
		current.Y += (fieldVector.Y / magnitude) * stepSize
		current.Z += (fieldVector.Z / magnitude) * stepSize
		
		if !fv.isInBounds(current) {
			break
		}
	}
	
	return line
}

func (fv *FieldVisualizer) calculateFieldAtPoint(point Point3D) Point3D {
	totalField := Point3D{0, 0, 0}
	
	for _, field := range fv.GravitationalFields {
		distance := fv.calculateDistance(point, field.Position)
		
		if distance > 0 && distance <= field.Range {
			fieldStrength := field.Strength / (distance * distance)
			
			direction := Point3D{
				X: (field.Position.X - point.X) / distance,
				Y: (field.Position.Y - point.Y) / distance,
				Z: (field.Position.Z - point.Z) / distance,
			}
			
			totalField.X += fieldStrength * direction.X
			totalField.Y += fieldStrength * direction.Y
			totalField.Z += fieldStrength * direction.Z
		}
	}
	
	return totalField
}

func (fv *FieldVisualizer) calculateDistance(p1, p2 Point3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (fv *FieldVisualizer) isInBounds(point Point3D) bool {
	return point.X >= fv.BoundingBox.MinX && point.X <= fv.BoundingBox.MaxX &&
		   point.Y >= fv.BoundingBox.MinY && point.Y <= fv.BoundingBox.MaxY &&
		   point.Z >= fv.BoundingBox.MinZ && point.Z <= fv.BoundingBox.MaxZ
}

func (fv *FieldVisualizer) GenerateVisualizationData() map[string]interface{} {
	data := make(map[string]interface{})
	
	fieldData := make([]map[string]interface{}, 0)
	for _, field := range fv.GravitationalFields {
		fieldInfo := map[string]interface{}{
			"position":  field.Position,
			"strength":  field.Strength,
			"direction": field.Direction,
			"range":     field.Range,
		}
		fieldData = append(fieldData, fieldInfo)
	}
	
	lineData := make([]map[string]interface{}, 0)
	for _, line := range fv.FieldLines {
		lineInfo := map[string]interface{}{
			"points":    line.Points,
			"strength":  line.Strength,
			"direction": line.Direction,
		}
		lineData = append(lineData, lineInfo)
	}
	
	data["fields"] = fieldData
	data["field_lines"] = lineData
	data["bounding_box"] = fv.BoundingBox
	
	return data
}

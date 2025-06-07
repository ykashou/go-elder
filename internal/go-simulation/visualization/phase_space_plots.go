package visualization

import "math"

type PhaseSpacePlotter struct {
	Trajectories map[string]PhaseTrajectory
	Dimensions   int
	TimeStep     float64
	MaxPoints    int
}

type PhaseTrajectory struct {
	ID     string
	Points []PhasePoint
	Color  string
	Style  string
}

type PhasePoint struct {
	Position []float64
	Velocity []float64
	Time     float64
	Energy   float64
}

func NewPhaseSpacePlotter(dimensions int, timeStep float64, maxPoints int) *PhaseSpacePlotter {
	return &PhaseSpacePlotter{
		Trajectories: make(map[string]PhaseTrajectory),
		Dimensions:   dimensions,
		TimeStep:     timeStep,
		MaxPoints:    maxPoints,
	}
}

func (psp *PhaseSpacePlotter) AddTrajectory(id, color, style string) {
	trajectory := PhaseTrajectory{
		ID:     id,
		Points: make([]PhasePoint, 0),
		Color:  color,
		Style:  style,
	}
	psp.Trajectories[id] = trajectory
}

func (psp *PhaseSpacePlotter) AddPoint(trajectoryID string, position, velocity []float64, time float64) {
	trajectory := psp.Trajectories[trajectoryID]
	
	energy := psp.calculateEnergy(position, velocity)
	
	point := PhasePoint{
		Position: make([]float64, len(position)),
		Velocity: make([]float64, len(velocity)),
		Time:     time,
		Energy:   energy,
	}
	
	copy(point.Position, position)
	copy(point.Velocity, velocity)
	
	trajectory.Points = append(trajectory.Points, point)
	
	if len(trajectory.Points) > psp.MaxPoints {
		trajectory.Points = trajectory.Points[1:]
	}
	
	psp.Trajectories[trajectoryID] = trajectory
}

func (psp *PhaseSpacePlotter) calculateEnergy(position, velocity []float64) float64 {
	kinetic := 0.0
	for _, v := range velocity {
		kinetic += v * v
	}
	kinetic *= 0.5
	
	potential := 0.0
	for _, p := range position {
		potential += p * p
	}
	potential *= 0.5
	
	return kinetic + potential
}

func (psp *PhaseSpacePlotter) GeneratePhasePortrait() map[string]interface{} {
	data := make(map[string]interface{})
	
	trajectoryData := make([]map[string]interface{}, 0)
	
	for _, trajectory := range psp.Trajectories {
		points := make([]map[string]interface{}, 0)
		
		for _, point := range trajectory.Points {
			pointData := map[string]interface{}{
				"position": point.Position,
				"velocity": point.Velocity,
				"time":     point.Time,
				"energy":   point.Energy,
			}
			points = append(points, pointData)
		}
		
		trajData := map[string]interface{}{
			"id":     trajectory.ID,
			"points": points,
			"color":  trajectory.Color,
			"style":  trajectory.Style,
		}
		trajectoryData = append(trajectoryData, trajData)
	}
	
	data["trajectories"] = trajectoryData
	data["dimensions"] = psp.Dimensions
	data["phase_analysis"] = psp.analyzePhaseSpace()
	
	return data
}

func (psp *PhaseSpacePlotter) analyzePhaseSpace() map[string]interface{} {
	analysis := make(map[string]interface{})
	
	totalEnergy := 0.0
	totalPoints := 0
	energyVariance := 0.0
	
	for _, trajectory := range psp.Trajectories {
		for _, point := range trajectory.Points {
			totalEnergy += point.Energy
			totalPoints++
		}
	}
	
	avgEnergy := totalEnergy / float64(totalPoints)
	
	for _, trajectory := range psp.Trajectories {
		for _, point := range trajectory.Points {
			diff := point.Energy - avgEnergy
			energyVariance += diff * diff
		}
	}
	
	energyVariance /= float64(totalPoints)
	
	analysis["average_energy"] = avgEnergy
	analysis["energy_variance"] = energyVariance
	analysis["energy_stability"] = 1.0 / (1.0 + energyVariance)
	analysis["total_trajectories"] = len(psp.Trajectories)
	analysis["total_points"] = totalPoints
	
	attractors := psp.findAttractors()
	analysis["attractors"] = attractors
	
	return analysis
}

func (psp *PhaseSpacePlotter) findAttractors() []map[string]interface{} {
	attractors := make([]map[string]interface{}, 0)
	
	for _, trajectory := range psp.Trajectories {
		if len(trajectory.Points) < 10 {
			continue
		}
		
		recentPoints := trajectory.Points[len(trajectory.Points)-10:]
		
		centerPos := make([]float64, psp.Dimensions)
		centerVel := make([]float64, psp.Dimensions)
		
		for _, point := range recentPoints {
			for i := 0; i < psp.Dimensions && i < len(point.Position); i++ {
				centerPos[i] += point.Position[i]
			}
			for i := 0; i < psp.Dimensions && i < len(point.Velocity); i++ {
				centerVel[i] += point.Velocity[i]
			}
		}
		
		for i := range centerPos {
			centerPos[i] /= float64(len(recentPoints))
		}
		for i := range centerVel {
			centerVel[i] /= float64(len(recentPoints))
		}
		
		variance := 0.0
		for _, point := range recentPoints {
			for i := 0; i < psp.Dimensions && i < len(point.Position); i++ {
				diff := point.Position[i] - centerPos[i]
				variance += diff * diff
			}
		}
		variance /= float64(len(recentPoints))
		
		if variance < 0.1 {
			attractor := map[string]interface{}{
				"trajectory_id": trajectory.ID,
				"center_pos":    centerPos,
				"center_vel":    centerVel,
				"variance":      variance,
				"type":          "fixed_point",
			}
			attractors = append(attractors, attractor)
		}
	}
	
	return attractors
}

func (psp *PhaseSpacePlotter) ExportData(format string) map[string]interface{} {
	export := make(map[string]interface{})
	
	switch format {
	case "csv":
		export["format"] = "csv"
		export["data"] = psp.generateCSVData()
	case "json":
		export["format"] = "json"
		export["data"] = psp.GeneratePhasePortrait()
	default:
		export["format"] = "json"
		export["data"] = psp.GeneratePhasePortrait()
	}
	
	return export
}

func (psp *PhaseSpacePlotter) generateCSVData() []string {
	csvLines := make([]string, 0)
	
	header := "trajectory_id,time,energy"
	for i := 0; i < psp.Dimensions; i++ {
		header += fmt.Sprintf(",pos_%d,vel_%d", i, i)
	}
	csvLines = append(csvLines, header)
	
	for _, trajectory := range psp.Trajectories {
		for _, point := range trajectory.Points {
			line := fmt.Sprintf("%s,%.6f,%.6f", trajectory.ID, point.Time, point.Energy)
			
			for i := 0; i < psp.Dimensions; i++ {
				if i < len(point.Position) {
					line += fmt.Sprintf(",%.6f", point.Position[i])
				} else {
					line += ",0.0"
				}
				if i < len(point.Velocity) {
					line += fmt.Sprintf(",%.6f", point.Velocity[i])
				} else {
					line += ",0.0"
				}
			}
			
			csvLines = append(csvLines, line)
		}
	}
	
	return csvLines
}

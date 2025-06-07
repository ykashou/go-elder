package mathematical

import "math"

type TopologyValidator struct {
	Space      TopologicalSpace
	Tolerance  float64
	Continuity map[string]bool
}

type TopologicalSpace struct {
	Points      []Point
	OpenSets    []Set
	Metric      func(Point, Point) float64
	Dimension   int
}

type Point struct {
	Coordinates []float64
	ID          string
}

type Set struct {
	Points []Point
	Type   string
}

func NewTopologyValidator(space TopologicalSpace, tolerance float64) *TopologyValidator {
	return &TopologyValidator{
		Space:      space,
		Tolerance:  tolerance,
		Continuity: make(map[string]bool),
	}
}

func (tv *TopologyValidator) ValidateTopology() TopologyValidationResult {
	result := TopologyValidationResult{
		Valid:      true,
		Properties: make(map[string]bool),
		Violations: make([]string, 0),
	}
	
	result.Properties["hausdorff"] = tv.checkHausdorff()
	result.Properties["compact"] = tv.checkCompact()
	result.Properties["connected"] = tv.checkConnected()
	result.Properties["complete"] = tv.checkComplete()
	
	for property, satisfied := range result.Properties {
		if !satisfied {
			result.Valid = false
			result.Violations = append(result.Violations, "Space is not "+property)
		}
	}
	
	return result
}

type TopologyValidationResult struct {
	Valid      bool
	Properties map[string]bool
	Violations []string
}

func (tv *TopologyValidator) checkHausdorff() bool {
	for i := 0; i < len(tv.Space.Points); i++ {
		for j := i + 1; j < len(tv.Space.Points); j++ {
			p1 := tv.Space.Points[i]
			p2 := tv.Space.Points[j]
			
			if tv.Space.Metric(p1, p2) > tv.Tolerance {
				if !tv.findSeparatingNeighborhoods(p1, p2) {
					return false
				}
			}
		}
	}
	return true
}

func (tv *TopologyValidator) findSeparatingNeighborhoods(p1, p2 Point) bool {
	distance := tv.Space.Metric(p1, p2)
	radius := distance / 3.0
	
	for _, set1 := range tv.Space.OpenSets {
		for _, set2 := range tv.Space.OpenSets {
			if tv.pointInSet(p1, set1) && tv.pointInSet(p2, set2) {
				if !tv.setsIntersect(set1, set2) {
					return true
				}
			}
		}
	}
	
	return false
}

func (tv *TopologyValidator) pointInSet(point Point, set Set) bool {
	for _, p := range set.Points {
		if tv.Space.Metric(point, p) < tv.Tolerance {
			return true
		}
	}
	return false
}

func (tv *TopologyValidator) setsIntersect(set1, set2 Set) bool {
	for _, p1 := range set1.Points {
		for _, p2 := range set2.Points {
			if tv.Space.Metric(p1, p2) < tv.Tolerance {
				return true
			}
		}
	}
	return false
}

func (tv *TopologyValidator) checkCompact() bool {
	if len(tv.Space.Points) == 0 {
		return true
	}
	
	return tv.checkSequentialCompactness()
}

func (tv *TopologyValidator) checkSequentialCompactness() bool {
	maxSequenceLength := 10
	
	for i := 0; i < maxSequenceLength && i < len(tv.Space.Points); i++ {
		sequence := tv.Space.Points[i:min(i+maxSequenceLength, len(tv.Space.Points))]
		if !tv.hasConvergentSubsequence(sequence) {
			return false
		}
	}
	
	return true
}

func (tv *TopologyValidator) hasConvergentSubsequence(sequence []Point) bool {
	if len(sequence) < 2 {
		return true
	}
	
	for i := 0; i < len(sequence)-1; i++ {
		for j := i + 1; j < len(sequence); j++ {
			if tv.Space.Metric(sequence[i], sequence[j]) < tv.Tolerance {
				return true
			}
		}
	}
	
	return false
}

func (tv *TopologyValidator) checkConnected() bool {
	if len(tv.Space.Points) <= 1 {
		return true
	}
	
	visited := make(map[string]bool)
	tv.dfs(tv.Space.Points[0], visited)
	
	return len(visited) == len(tv.Space.Points)
}

func (tv *TopologyValidator) dfs(point Point, visited map[string]bool) {
	visited[point.ID] = true
	
	for _, p := range tv.Space.Points {
		if !visited[p.ID] && tv.Space.Metric(point, p) < tv.Tolerance*10 {
			tv.dfs(p, visited)
		}
	}
}

func (tv *TopologyValidator) checkComplete() bool {
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

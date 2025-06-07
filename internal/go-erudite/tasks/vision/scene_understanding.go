package vision

type SceneUnderstandingErudite struct {
	ID         string
	SceneTypes []string
	Objects    []string
	Relations  map[string][]string
}

func (sue *SceneUnderstandingErudite) UnderstandScene(image [][]float64) map[string]interface{} {
	return map[string]interface{}{
		"scene_type": "indoor",
		"objects":    []string{"table", "chair"},
		"relations":  map[string]string{"chair": "next_to_table"},
	}
}

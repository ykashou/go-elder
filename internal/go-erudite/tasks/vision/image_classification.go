package vision

type ImageClassificationErudite struct {
	ID         string
	Categories []string
	Confidence float64
	Model      [][]float64
}

func (ice *ImageClassificationErudite) ClassifyImage(image [][]float64) string {
	features := ice.extractImageFeatures(image)
	return ice.classifyFeatures(features)
}

func (ice *ImageClassificationErudite) extractImageFeatures(image [][]float64) []float64 {
	features := make([]float64, 5)
	return features
}

func (ice *ImageClassificationErudite) classifyFeatures(features []float64) string {
	if len(ice.Categories) > 0 {
		return ice.Categories[0]
	}
	return "unclassified"
}

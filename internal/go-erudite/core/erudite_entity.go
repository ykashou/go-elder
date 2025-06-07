package core

type EruditeEntity struct {
	ID             string
	TaskType       string
	Specialization string
	Performance    float64
	LearningRate   float64
}

func NewEruditeEntity(id, taskType, specialization string) *EruditeEntity {
	return &EruditeEntity{
		ID:             id,
		TaskType:       taskType,
		Specialization: specialization,
		Performance:    0.0,
		LearningRate:   0.01,
	}
}

func (ee *EruditeEntity) UpdatePerformance(delta float64) {
	ee.Performance += delta * ee.LearningRate
}

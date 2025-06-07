package elder

type InformationCapacityController struct {
	Elder        *Elder
	MaxCapacity  float64
	UsedCapacity float64
	Channels     map[string]float64
}

func NewInformationCapacityController(maxCapacity float64) *InformationCapacityController {
	return &InformationCapacityController{
		MaxCapacity: maxCapacity,
		Channels:    make(map[string]float64),
	}
}

func (icc *InformationCapacityController) AllocateCapacity(channelID string, capacity float64) bool {
	if icc.UsedCapacity+capacity > icc.MaxCapacity {
		return false
	}
	icc.Channels[channelID] = capacity
	icc.UsedCapacity += capacity
	return true
}

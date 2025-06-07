package coordination

type InformationFlowManager struct {
	Channels map[string]Channel
	FlowRates map[string]float64
	Capacity float64
}

type Channel struct {
	Source string
	Target string
	Bandwidth float64
	Active bool
}

func NewInformationFlowManager(capacity float64) *InformationFlowManager {
	return &InformationFlowManager{
		Channels: make(map[string]Channel),
		FlowRates: make(map[string]float64),
		Capacity: capacity,
	}
}

func (ifm *InformationFlowManager) CreateChannel(id, source, target string, bandwidth float64) {
	ifm.Channels[id] = Channel{
		Source: source,
		Target: target,
		Bandwidth: bandwidth,
		Active: true,
	}
}

func (ifm *InformationFlowManager) ManageFlow() {
	totalFlow := 0.0
	for id, channel := range ifm.Channels {
		if channel.Active {
			flow := channel.Bandwidth * 0.8
			ifm.FlowRates[id] = flow
			totalFlow += flow
		}
	}
}

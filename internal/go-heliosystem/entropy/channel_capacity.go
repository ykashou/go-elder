package entropy

import "math"

type ChannelCapacity struct {
	Channels    map[string]Channel
	NoiseLevel  float64
	Bandwidth   float64
	SignalPower float64
}

type Channel struct {
	ID           string
	Capacity     float64
	Utilization  float64
	ErrorRate    float64
	Throughput   float64
}

func NewChannelCapacity(bandwidth, signalPower, noise float64) *ChannelCapacity {
	return &ChannelCapacity{
		Channels:    make(map[string]Channel),
		NoiseLevel:  noise,
		Bandwidth:   bandwidth,
		SignalPower: signalPower,
	}
}

func (cc *ChannelCapacity) CalculateShannonCapacity() float64 {
	snr := cc.SignalPower / cc.NoiseLevel
	return cc.Bandwidth * math.Log2(1+snr)
}

func (cc *ChannelCapacity) CreateChannel(id string) {
	capacity := cc.CalculateShannonCapacity()
	
	channel := Channel{
		ID:          id,
		Capacity:    capacity,
		Utilization: 0.0,
		ErrorRate:   cc.calculateErrorRate(),
		Throughput:  0.0,
	}
	
	cc.Channels[id] = channel
}

func (cc *ChannelCapacity) calculateErrorRate() float64 {
	snr := cc.SignalPower / cc.NoiseLevel
	return 1.0 / (1.0 + snr)
}

func (cc *ChannelCapacity) UpdateUtilization(channelID string, utilization float64) {
	if channel, exists := cc.Channels[channelID]; exists {
		channel.Utilization = math.Min(1.0, utilization)
		channel.Throughput = channel.Capacity * channel.Utilization * (1.0 - channel.ErrorRate)
		cc.Channels[channelID] = channel
	}
}

func (cc *ChannelCapacity) GetTotalThroughput() float64 {
	total := 0.0
	for _, channel := range cc.Channels {
		total += channel.Throughput
	}
	return total
}

func (cc *ChannelCapacity) OptimizeChannels() {
	for id, channel := range cc.Channels {
		if channel.Utilization > 0.8 {
			channel.Capacity *= 1.1
			cc.Channels[id] = channel
		}
	}
}

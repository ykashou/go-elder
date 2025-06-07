package config

type SimulationConfig struct {
	Engine      EngineConfig      `json:"engine"`
	Physics     PhysicsConfig     `json:"physics"`
	Integration IntegrationConfig `json:"integration"`
	Output      OutputConfig      `json:"output"`
}

type EngineConfig struct {
	Type        string  `json:"type"`
	Precision   string  `json:"precision"`
	Parallel    bool    `json:"parallel"`
	Threads     int     `json:"threads"`
	MemoryLimit int64   `json:"memory_limit"`
}

type PhysicsConfig struct {
	Gravity         bool    `json:"gravity"`
	Electromagnetic bool    `json:"electromagnetic"`
	Resonance       bool    `json:"resonance"`
	Damping         float64 `json:"damping"`
	StabilityCheck  bool    `json:"stability_check"`
}

type IntegrationConfig struct {
	Method   string  `json:"method"`
	TimeStep float64 `json:"time_step"`
	MaxSteps int     `json:"max_steps"`
	Adaptive bool    `json:"adaptive"`
}

type OutputConfig struct {
	Format     string   `json:"format"`
	Fields     []string `json:"fields"`
	Frequency  int      `json:"frequency"`
	Compress   bool     `json:"compress"`
	Directory  string   `json:"directory"`
}

package config

type AnalysisConfig struct {
	Targets    []string          `json:"targets"`
	Metrics    MetricsConfig     `json:"metrics"`
	Stability  StabilityConfig   `json:"stability"`
	Performance PerformanceConfig `json:"performance"`
	Output     AnalysisOutputConfig `json:"output"`
}

type MetricsConfig struct {
	Convergence    bool `json:"convergence"`
	Energy         bool `json:"energy"`
	Entropy        bool `json:"entropy"`
	Information    bool `json:"information"`
	Complexity     bool `json:"complexity"`
}

type StabilityConfig struct {
	LyapunovExponent bool    `json:"lyapunov_exponent"`
	PhaseAnalysis    bool    `json:"phase_analysis"`
	Perturbation     bool    `json:"perturbation"`
	Tolerance        float64 `json:"tolerance"`
	TimeWindow       float64 `json:"time_window"`
}

type PerformanceConfig struct {
	Memory       bool `json:"memory"`
	CPU          bool `json:"cpu"`
	Throughput   bool `json:"throughput"`
	Latency      bool `json:"latency"`
	Scalability  bool `json:"scalability"`
}

type AnalysisOutputConfig struct {
	Format      string   `json:"format"`
	Details     string   `json:"details"`
	Plots       bool     `json:"plots"`
	Interactive bool     `json:"interactive"`
	Export      []string `json:"export"`
}

package config

type ElderConfig struct {
	System      SystemConfig      `json:"system"`
	Entities    EntitiesConfig    `json:"entities"`
	Fields      FieldsConfig      `json:"fields"`
	Simulation  SimulationConfig  `json:"simulation"`
	Logging     LoggingConfig     `json:"logging"`
}

type SystemConfig struct {
	MaxMemory        int64   `json:"max_memory"`
	MaxProcessors    int     `json:"max_processors"`
	GravitationalG   float64 `json:"gravitational_g"`
	TimeStep         float64 `json:"time_step"`
	Tolerance        float64 `json:"tolerance"`
}

type EntitiesConfig struct {
	ElderCount   int `json:"elder_count"`
	MentorCount  int `json:"mentor_count"`
	EruditeCount int `json:"erudite_count"`
}

type FieldsConfig struct {
	MaxFieldStrength float64 `json:"max_field_strength"`
	FieldRange       float64 `json:"field_range"`
	DecayRate        float64 `json:"decay_rate"`
}

type SimulationConfig struct {
	MaxDuration     float64 `json:"max_duration"`
	OutputInterval  float64 `json:"output_interval"`
	CheckpointFreq  int     `json:"checkpoint_frequency"`
}

type LoggingConfig struct {
	Level      string `json:"level"`
	OutputFile string `json:"output_file"`
	Verbose    bool   `json:"verbose"`
}

func DefaultElderConfig() *ElderConfig {
	return &ElderConfig{
		System: SystemConfig{
			MaxMemory:      8 * 1024 * 1024 * 1024,
			MaxProcessors:  8,
			GravitationalG: 6.67430e-11,
			TimeStep:       0.01,
			Tolerance:      1e-8,
		},
		Entities: EntitiesConfig{
			ElderCount:   1,
			MentorCount:  3,
			EruditeCount: 9,
		},
		Fields: FieldsConfig{
			MaxFieldStrength: 1000.0,
			FieldRange:      100.0,
			DecayRate:       0.01,
		},
		Simulation: SimulationConfig{
			MaxDuration:    1000.0,
			OutputInterval: 1.0,
			CheckpointFreq: 100,
		},
		Logging: LoggingConfig{
			Level:      "INFO",
			OutputFile: "elder.log",
			Verbose:    false,
		},
	}
}

package config

type TransferConfig struct {
	Source   DomainConfig   `json:"source"`
	Target   DomainConfig   `json:"target"`
	Method   MethodConfig   `json:"method"`
	Mapping  MappingConfig  `json:"mapping"`
	Validation ValidationConfig `json:"validation"`
}

type DomainConfig struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Features    []string           `json:"features"`
	Dimensions  int                `json:"dimensions"`
	Parameters  map[string]float64 `json:"parameters"`
}

type MethodConfig struct {
	Type           string  `json:"type"`
	Isomorphic     bool    `json:"isomorphic"`
	Hierarchical   bool    `json:"hierarchical"`
	Resonance      bool    `json:"resonance"`
	Confidence     float64 `json:"confidence"`
	MaxIterations  int     `json:"max_iterations"`
}

type MappingConfig struct {
	AutoDetect     bool               `json:"auto_detect"`
	ManualMappings map[string]string  `json:"manual_mappings"`
	Weights        map[string]float64 `json:"weights"`
	Bidirectional  bool               `json:"bidirectional"`
}

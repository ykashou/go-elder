package config

type TrainingConfig struct {
	Model       ModelConfig       `json:"model"`
	Optimizer   OptimizerConfig   `json:"optimizer"`
	Data        DataConfig        `json:"data"`
	Validation  ValidationConfig  `json:"validation"`
	Checkpoints CheckpointConfig  `json:"checkpoints"`
}

type ModelConfig struct {
	Architecture string             `json:"architecture"`
	Layers       []LayerConfig      `json:"layers"`
	Parameters   map[string]float64 `json:"parameters"`
	Regularization RegularizationConfig `json:"regularization"`
}

type LayerConfig struct {
	Type       string             `json:"type"`
	Size       int                `json:"size"`
	Activation string             `json:"activation"`
	Parameters map[string]float64 `json:"parameters"`
}

type RegularizationConfig struct {
	L1         float64 `json:"l1"`
	L2         float64 `json:"l2"`
	Dropout    float64 `json:"dropout"`
	BatchNorm  bool    `json:"batch_norm"`
}

type OptimizerConfig struct {
	Type         string  `json:"type"`
	LearningRate float64 `json:"learning_rate"`
	Momentum     float64 `json:"momentum"`
	Beta1        float64 `json:"beta1"`
	Beta2        float64 `json:"beta2"`
	Epsilon      float64 `json:"epsilon"`
}

type DataConfig struct {
	BatchSize    int      `json:"batch_size"`
	Shuffle      bool     `json:"shuffle"`
	Augment      bool     `json:"augment"`
	Preprocessing []string `json:"preprocessing"`
}

type ValidationConfig struct {
	SplitRatio float64 `json:"split_ratio"`
	Frequency  int     `json:"frequency"`
	EarlyStopping bool `json:"early_stopping"`
	Patience   int     `json:"patience"`
}

type CheckpointConfig struct {
	Enabled   bool   `json:"enabled"`
	Frequency int    `json:"frequency"`
	Directory string `json:"directory"`
	KeepLast  int    `json:"keep_last"`
}

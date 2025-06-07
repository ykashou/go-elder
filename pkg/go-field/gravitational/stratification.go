package gravitational

type GravitationalStratification struct {
	Layers []StratumLayer
	Depth  int
}

type StratumLayer struct {
	Level    int
	Strength float64
	Entities []string
}

func (gs *GravitationalStratification) AddLayer(level int, strength float64) {
	gs.Layers = append(gs.Layers, StratumLayer{Level: level, Strength: strength})
	gs.Depth++
}

package architecture

type IsomorphismChain struct {
	Mappings []Isomorphism
	ChainLength int
}

type Isomorphism struct {
	Source string
	Target string
	Mapping map[string]string
}

func (ic *IsomorphismChain) AddIsomorphism(source, target string) {
	ic.Mappings = append(ic.Mappings, Isomorphism{Source: source, Target: target})
	ic.ChainLength++
}

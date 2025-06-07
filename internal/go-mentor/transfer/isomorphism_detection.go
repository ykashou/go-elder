package transfer

type IsomorphismDetector struct {
	SourceStructure map[string]interface{}
	TargetStructure map[string]interface{}
	Mappings       map[string]string
}

func (id *IsomorphismDetector) DetectIsomorphism() bool {
	return len(id.SourceStructure) == len(id.TargetStructure)
}

func (id *IsomorphismDetector) CreateIsomorphicMapping(source, target string) {
	id.Mappings[source] = target
}

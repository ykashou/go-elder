package isomorphisms

type DomainIsomorphism struct {
	SourceDomain string
	TargetDomain string
	ForwardMap   map[string]string
	BackwardMap  map[string]string
	Verified     bool
}

func NewDomainIsomorphism(source, target string) *DomainIsomorphism {
	return &DomainIsomorphism{
		SourceDomain: source,
		TargetDomain: target,
		ForwardMap:   make(map[string]string),
		BackwardMap:  make(map[string]string),
		Verified:     false,
	}
}

func (di *DomainIsomorphism) AddMapping(sourceElement, targetElement string) {
	di.ForwardMap[sourceElement] = targetElement
	di.BackwardMap[targetElement] = sourceElement
}

func (di *DomainIsomorphism) VerifyBijection() bool {
	forwardCount := len(di.ForwardMap)
	backwardCount := len(di.BackwardMap)
	
	if forwardCount != backwardCount {
		return false
	}
	
	for source, target := range di.ForwardMap {
		if backwardTarget, exists := di.BackwardMap[target]; !exists || backwardTarget != source {
			return false
		}
	}
	
	di.Verified = true
	return true
}

func (di *DomainIsomorphism) MapElement(element string, direction string) (string, bool) {
	switch direction {
	case "forward":
		if target, exists := di.ForwardMap[element]; exists {
			return target, true
		}
	case "backward":
		if source, exists := di.BackwardMap[element]; exists {
			return source, true
		}
	}
	return "", false
}

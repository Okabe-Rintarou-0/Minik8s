package selector

type Type byte

const (
	Random Type = iota
	MinimumCpuUtility
	MinimumMemoryUtility
	MinimumNumPods
	MaximumNumPods
)

var DefaultFactory = &Factory{}

type Factory struct{}

func (f *Factory) NewSelector(selectorType Type) Selector {
	switch selectorType {
	case Random:
		return random()
	case MinimumCpuUtility:
		return minimumCpuUtility()
	case MinimumMemoryUtility:
		return minimumMemoryUtility()
	case MinimumNumPods:
		return minimumNumPods()
	case MaximumNumPods:
		return maximumNumPods()
	}

	panic("Invalid Selector Type!")
}

package calculator

func NewDefaultRegistry() *Registry {
	return NewRegistry(
		Add{},
		Subtract{},
		Multiply{},
		Divide{},
		Power{},
		Sqrt{},
		Percentage{},
	)
}

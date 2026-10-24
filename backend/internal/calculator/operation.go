package calculator

type Arity struct {
	Min int
	Max int
}

type Operation interface {
	Name() string
	Arity() Arity
	Apply(operands []float64) (float64, error)
}

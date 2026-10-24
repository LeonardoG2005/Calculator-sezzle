package calculator

type Subtract struct{}

func (Subtract) Name() string { return "subtract" }
func (Subtract) Arity() Arity { return Arity{Min: 2, Max: 2} }
func (Subtract) Apply(values []float64) (float64, error) {
	return values[0] - values[1], nil
}

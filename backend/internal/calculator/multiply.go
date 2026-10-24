package calculator
type Multiply struct{}

func (Multiply) Name() string { return "multiply" }
func (Multiply) Arity() Arity { return Arity{Min: 2, Max: 2} }
func (Multiply) Apply(values []float64) (float64, error) {
	return values[0] * values[1], nil
}

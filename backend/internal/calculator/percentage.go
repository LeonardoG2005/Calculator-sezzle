package calculator

type Percentage struct{}

func (Percentage) Name() string { return "percentage" }
func (Percentage) Arity() Arity { return Arity{Min: 2, Max: 2} }
func (Percentage) Apply(values []float64) (float64, error) {
	return values[0] * values[1] / 100, nil
}

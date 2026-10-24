package calculator
type Add struct{}

func (Add) Name() string       { return "add" }
func (Add) Arity() Arity       { return Arity{Min: 2, Max: 2} }
func (Add) Apply(values []float64) (float64, error) {
	return values[0] + values[1], nil
}

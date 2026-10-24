package calculator
type Divide struct{}

func (Divide) Name() string { return "divide" }
func (Divide) Arity() Arity { return Arity{Min: 2, Max: 2} }
func (Divide) Apply(values []float64) (float64, error) {
	if values[1] == 0 {
		return 0, ErrDivisionByZero
	}
	return values[0] / values[1], nil
}

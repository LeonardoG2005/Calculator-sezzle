package calculator
import "math"

type Power struct{}

func (Power) Name() string { return "power" }
func (Power) Arity() Arity { return Arity{Min: 2, Max: 2} }
func (Power) Apply(values []float64) (float64, error) {
	if values[0] < 0 && values[1] != math.Trunc(values[1]) {
		return 0, ErrMathDomain
	}
	return math.Pow(values[0], values[1]), nil
}

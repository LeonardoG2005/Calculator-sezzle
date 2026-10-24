package calculator
import "math"

type Sqrt struct{}

func (Sqrt) Name() string { return "sqrt" }
func (Sqrt) Arity() Arity { return Arity{Min: 1, Max: 1} }
func (Sqrt) Apply(values []float64) (float64, error) {
	if values[0] < 0 {
		return 0, ErrMathDomain
	}
	return math.Sqrt(values[0]), nil
}

package calculator
import "testing"

func TestBasicOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		values    []float64
		want      float64
	}{
		{"add", Add{}, []float64{2, 3}, 5},
		{"subtract", Subtract{}, []float64{9, 4}, 5},
		{"multiply", Multiply{}, []float64{6, 7}, 42},
		{"divide", Divide{}, []float64{10, 4}, 2.5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.operation.Apply(test.values)
			if err != nil {
				t.Fatalf("Apply() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("Apply() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := (Divide{}).Apply([]float64{1, 0})
	if err != ErrDivisionByZero {
		t.Fatalf("Apply() error = %v, want %v", err, ErrDivisionByZero)
	}
}

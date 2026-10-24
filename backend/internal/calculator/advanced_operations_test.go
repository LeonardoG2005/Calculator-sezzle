package calculator
import "testing"

func TestAdvancedOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		values    []float64
		want      float64
	}{
		{"power", Power{}, []float64{2, 10}, 1024},
		{"square root", Sqrt{}, []float64{16}, 4},
		{"percentage", Percentage{}, []float64{20, 50}, 10},
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

func TestAdvancedDomainErrors(t *testing.T) {
	if _, err := (Sqrt{}).Apply([]float64{-1}); err != ErrMathDomain {
		t.Fatalf("Sqrt() error = %v, want %v", err, ErrMathDomain)
	}
	if _, err := (Power{}).Apply([]float64{-2, 0.5}); err != ErrMathDomain {
		t.Fatalf("Power() error = %v, want %v", err, ErrMathDomain)
	}
}

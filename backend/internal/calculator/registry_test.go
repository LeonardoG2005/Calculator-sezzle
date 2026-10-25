package calculator

import (
	"errors"
	"testing"
)

func TestDefaultRegistry(t *testing.T) {
	registry := NewDefaultRegistry()
	if len(registry.List()) != 7 {
		t.Fatalf("List() returned %d operations, want 7", len(registry.List()))
	}

	result, err := registry.Calculate("add", []float64{2, 3})
	if err != nil || result != 5 {
		t.Fatalf("Calculate(add) = %v, %v; want 5, nil", result, err)
	}
}

func TestRegistryErrors(t *testing.T) {
	registry := NewDefaultRegistry()
	tests := []struct {
		name      string
		operation string
		operands  []float64
		want      error
	}{
		{"unknown operation", "unknown", []float64{1}, ErrUnknownOperation},
		{"too few operands", "add", []float64{1}, ErrInvalidOperandCount},
		{"too many operands", "sqrt", []float64{1, 2}, ErrInvalidOperandCount},
		{"overflow", "power", []float64{10, 400}, ErrResultOutOfRange},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := registry.Calculate(test.operation, test.operands)
			if !errors.Is(err, test.want) {
				t.Fatalf("Calculate() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestOperationMetadata(t *testing.T) {
	operations := []Operation{Add{}, Subtract{}, Multiply{}, Divide{}, Power{}, Sqrt{}, Percentage{}}
	for _, operation := range operations {
		if operation.Name() == "" || operation.Arity().Min < 1 || operation.Arity().Max < operation.Arity().Min {
			t.Fatalf("invalid metadata for %T", operation)
		}
	}
}

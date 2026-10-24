package calculator

import "math"

type Registry struct {
	operations map[string]Operation
}

func NewRegistry(operations ...Operation) *Registry {
	registry := &Registry{operations: make(map[string]Operation, len(operations))}
	for _, operation := range operations {
		registry.operations[operation.Name()] = operation
	}
	return registry
}

func (r *Registry) List() []Operation {
	operations := make([]Operation, 0, len(r.operations))
	for _, operation := range r.operations {
		operations = append(operations, operation)
	}
	return operations
}

func (r *Registry) Calculate(name string, operands []float64) (float64, error) {
	operation, ok := r.operations[name]
	if !ok {
		return 0, ErrUnknownOperation
	}
	arity := operation.Arity()
	if len(operands) < arity.Min || (arity.Max >= 0 && len(operands) > arity.Max) {
		return 0, ErrInvalidOperandCount
	}
	result, err := operation.Apply(operands)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrResultOutOfRange
	}
	return result, nil
}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"calculator/internal/calculator"
)

const maxRequestBodyBytes = 1024

type Handlers struct {
	registry *calculator.Registry
}

func NewHandlers(registry *calculator.Registry) *Handlers {
	return &Handlers{registry: registry}
}

func (h *Handlers) Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) Operations(w http.ResponseWriter, _ *http.Request) {
	operations := h.registry.List()
	infos := make([]OperationInfo, len(operations))
	for i, operation := range operations {
		arity := operation.Arity()
		infos[i] = OperationInfo{
			Name:  operation.Name(),
			Arity: ArityInfo{Min: arity.Min, Max: arity.Max},
		}
	}
	writeJSON(w, http.StatusOK, OperationsResponse{Operations: infos})
}

func (h *Handlers) Calculate(w http.ResponseWriter, r *http.Request) {
	contentType, _, contentTypeError := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if contentTypeError != nil || contentType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, ErrCodeUnsupportedMediaType,
			"Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request CalculateRequest
	if err := decoder.Decode(&request); err != nil {
		var maxError *http.MaxBytesError
		if errors.As(err, &maxError) {
			writeError(w, http.StatusRequestEntityTooLarge, ErrCodeRequestTooLarge,
				"request body too large")
			return
		}
		if strings.Contains(err.Error(), "unknown field") {
			writeError(w, http.StatusBadRequest, ErrCodeUnknownField,
				"request body contains an unknown field")
		} else {
			writeError(w, http.StatusBadRequest, ErrCodeInvalidJSON,
				"request body is not valid JSON")
		}
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, ErrCodeInvalidJSON,
			"request body must contain a single JSON object")
		return
	}
	if strings.TrimSpace(request.Operation) == "" {
		writeError(w, http.StatusBadRequest, ErrCodeMissingField, "operation is required")
		return
	}

	operands := make([]float64, len(request.Operands))
	for index, operand := range request.Operands {
		if operand == nil {
			writeError(w, http.StatusBadRequest, ErrCodeInvalidJSON,
				fmt.Sprintf("operand at index %d must be a number, not null", index))
			return
		}
		operands[index] = *operand
	}

	result, err := h.registry.Calculate(request.Operation, operands)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrUnknownOperation):
		writeError(w, http.StatusBadRequest, ErrCodeUnknownOperation, err.Error())
	case errors.Is(err, calculator.ErrInvalidOperandCount):
		writeError(w, http.StatusBadRequest, ErrCodeInvalidOperandCount, err.Error())
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, ErrCodeDivisionByZero, err.Error())
	case errors.Is(err, calculator.ErrMathDomain):
		writeError(w, http.StatusUnprocessableEntity, ErrCodeMathDomainError, err.Error())
	case errors.Is(err, calculator.ErrResultOutOfRange):
		writeError(w, http.StatusUnprocessableEntity, ErrCodeResultOutOfRange, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

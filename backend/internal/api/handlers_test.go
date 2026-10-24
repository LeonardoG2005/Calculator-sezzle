package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator/internal/calculator"
)

func TestCalculate(t *testing.T) {

	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
		wantCode    string
	}{
		{"success", `{"operation":"add","operands":[2,3]}`, "application/json; charset=utf-8", http.StatusOK, ""},
		{"missing content type", `{"operation":"add","operands":[2,3]}`, "", http.StatusUnsupportedMediaType, string(ErrCodeUnsupportedMediaType)},
		{"unknown operation", `{"operation":"modulo","operands":[2,3]}`, "application/json", http.StatusBadRequest, string(ErrCodeUnknownOperation)},
		{"invalid count", `{"operation":"sqrt","operands":[2,3]}`, "application/json", http.StatusBadRequest, string(ErrCodeInvalidOperandCount)},
		{"division by zero", `{"operation":"divide","operands":[2,0]}`, "application/json", http.StatusUnprocessableEntity, string(ErrCodeDivisionByZero)},
		{"unknown field", `{"operation":"add","operands":[2,3],"extra":true}`, "application/json", http.StatusBadRequest, string(ErrCodeUnknownField)},
		{"trailing data", `{"operation":"add","operands":[2,3]} {}`, "application/json", http.StatusBadRequest, string(ErrCodeInvalidJSON)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(test.body))
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			response := httptest.NewRecorder()
			NewHandlers(calculator.NewDefaultRegistry()).Calculate(response, request)
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.wantCode == "" {
				return
			}
			var payload ErrorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if payload.Error.Code != test.wantCode {
				t.Fatalf("error code = %q, want %q", payload.Error.Code, test.wantCode)
			}
		})
	}
}

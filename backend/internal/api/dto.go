package api

type CalculateRequest struct {
	Operation string     `json:"operation"`
	Operands  []*float64 `json:"operands"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OperationsResponse struct {
	Operations []OperationInfo `json:"operations"`
}

type OperationInfo struct {
	Name  string    `json:"name"`
	Arity ArityInfo `json:"arity"`
}

type ArityInfo struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

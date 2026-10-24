package api

type ErrorCode string

const (
	ErrCodeInvalidJSON          ErrorCode = "INVALID_JSON"
	ErrCodeUnknownField         ErrorCode = "UNKNOWN_FIELD"
	ErrCodeMissingField         ErrorCode = "MISSING_FIELD"
	ErrCodeUnknownOperation     ErrorCode = "UNKNOWN_OPERATION"
	ErrCodeInvalidOperandCount  ErrorCode = "INVALID_OPERAND_COUNT"
	ErrCodeUnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
	ErrCodeRequestTooLarge      ErrorCode = "REQUEST_TOO_LARGE"
	ErrCodeDivisionByZero       ErrorCode = "DIVISION_BY_ZERO"
	ErrCodeMathDomainError      ErrorCode = "MATH_DOMAIN_ERROR"
	ErrCodeResultOutOfRange     ErrorCode = "RESULT_OUT_OF_RANGE"
	ErrCodeNotFound             ErrorCode = "NOT_FOUND"
	ErrCodeMethodNotAllowed     ErrorCode = "METHOD_NOT_ALLOWED"
	ErrCodeInternal             ErrorCode = "INTERNAL_ERROR"
)

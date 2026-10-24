package calculator
import "errors"

var (
	ErrUnknownOperation    = errors.New("unknown operation")
	ErrInvalidOperandCount = errors.New("invalid operand count")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrMathDomain           = errors.New("operation is undefined for the given operands")
	ErrResultOutOfRange     = errors.New("result is not a finite number")
)

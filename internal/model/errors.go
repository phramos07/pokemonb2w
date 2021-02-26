package model

// Exported error constants
const (
	ErrorDefault           = "DEFAULT_ERROR_TYPE"
	ErrorUnprocessableJSON = "UNPROCESSABLE_JSON_ERROR_TYPE"
	ErrorBadRequest        = "BAD_REQUEST_ERROR_TYPE"
)

// CustomError Interface
type CustomError interface {
	error
	ErrorType() string
}

type requestError struct {
	CustomError

	errStr  string
	errType string
}

// NewrequestError Makes new custom request with given status code and str
func NewrequestError(errStr string, errType string) CustomError {
	return &requestError{
		errType: errType,
		errStr:  errStr,
	}
}

func (r *requestError) Error() string {
	return r.errStr
}

func (r *requestError) StatusCode() string {
	return r.errType
}

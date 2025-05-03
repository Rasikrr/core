package errors

type CoreError struct {
	Code    int
	Message string
}

func (e *CoreError) Error() string {
	return e.Message
}

func (e *CoreError) StatusCode() int {
	return e.Code
}

func NewError(message string, code int) *CoreError {
	return &CoreError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrNotFound            = NewError("not found", 404)
	ErrBadRequest          = NewError("bad request", 400)
	ErrInternalServerError = NewError("internal server error", 500)
	ErrUnauthorized        = NewError("unauthorized", 401)
	ErrForbidden           = NewError("forbidden", 403)
	ErrConflict            = NewError("conflict", 409)
	ErrBadRequestBody      = NewError("bad request body", 400)
	ErrNotImplemented      = NewError("not implemented", 501)
	ErrInvalidToken        = NewError("invalid token", 401)
)

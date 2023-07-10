package errs

import "net/http"

type AppError struct {
	Code    int
	Message string
	Errors  any
}

func (e AppError) Error() string {
	return e.Message
}

func ErrBadRequest() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid data provided",
	}
}

func ErrUnexpected() *AppError {
	return &AppError{
		Code:    http.StatusBadGateway,
		Message: "Unexpected error",
	}
}

func ErrUnauthorized() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
}

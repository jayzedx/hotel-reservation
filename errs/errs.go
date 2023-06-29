package errs

type AppError struct {
	Code    int
	Message string
	Errors  any
}

func (e AppError) Error() string {
	return e.Message
}

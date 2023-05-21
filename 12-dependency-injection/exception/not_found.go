package exception

type NotFoundError struct {
	Error string
}

func NewNotfoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

package responses

type Error struct {
	StatusCode int
	Message    string
}

func NewError(code int, message string) Error {
	return Error{
		StatusCode: code,
		Message:    message,
	}
}

// Code returns the HTTP response code.
func (e Error) Code() int {
	return e.StatusCode
}

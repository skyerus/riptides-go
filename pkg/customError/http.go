package customError

type Error interface {
	OriginalError() error
	Code() int
	Message() string
}

type HttpError struct {
	code int
	message string
	error error
}

func NewHttpError(code int, message string, err error) Error {
	return &HttpError{code, message, err}
}

func NewGenericHttpError(err error) Error {
	return &HttpError{500, "Oops, something went wrong. Please try again later.", err}
}

func (e *HttpError) OriginalError() error {
	return e.error
}

func (e *HttpError) Code() int  {
	return e.code
}

func (e *HttpError) Message() string {
	return e.message
}

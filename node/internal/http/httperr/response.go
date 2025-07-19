package httperr

// error to write to http response
type Response struct {
	code    int
	details string
}

func NewResponse(code int, details string) *Response {
	return &Response{code: code, details: details}
}

func (e *Response) Error() string {
	return e.details
}

func (e *Response) Code() int {
	return e.code
}

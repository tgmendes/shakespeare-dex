package web

// HTTPError is used to return http client requests errors.
type HTTPError struct {
	StatusCode int
	Err        error
}

// Error implements the error interface, returning the http error message.
func (e *HTTPError) Error() string {
	return e.Err.Error()
}

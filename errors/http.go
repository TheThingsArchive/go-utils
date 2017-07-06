package errors

import "net/http"

// HTTPStatusCode returns the corresponding http status code from an error type
func (t Type) HTTPStatusCode() int {
	switch t {
	case InvalidArgument:
	case OutOfRange:
		return http.StatusBadRequest

	case NotFound:
		return http.StatusNotFound

	case Conflict:
	case AlreadyExists:
		return http.StatusConflict

	case Unauthorized:
		return http.StatusUnauthorized

	case PermissionDenied:
		return http.StatusForbidden

	case Timeout:
		return http.StatusRequestTimeout

	case NotImplemented:
		return http.StatusNotImplemented

	case TemporarilyUnavailable:
		return http.StatusBadGateway

	case PermanentlyUnavailable:
		return http.StatusGone

	case Unknown:
	case Internal:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

// HTTPStatusCode returns the HTTP status code for the given error
// or 500 if it doesn't know
func HTTPStatusCode(err error) int {
	e, ok := err.(Error)
	if ok {
		return e.Type().HTTPStatusCode()
	}

	return http.StatusInternalServerError
}

// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"encoding/json"
	"net/http"
)

// CodeHeader is the header where the error code will be stored
const CodeHeader = "X-TTN-Error-Code"

// HTTPStatusCode returns the corresponding http status code from an error type
func (t Type) HTTPStatusCode() int {
	switch t {
	case Canceled:
		return http.StatusRequestTimeout
	case InvalidArgument:
		return http.StatusBadRequest
	case OutOfRange:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
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
	case ResourceExhausted:
		return http.StatusForbidden
	case Internal:
		return http.StatusInternalServerError
	case Unknown:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

// HTTPStatusCode returns the HTTP status code for the given error or 500 if it doesn't know
func HTTPStatusCode(err error) int {
	e, ok := err.(Error)
	if ok {
		return e.Type().HTTPStatusCode()
	}

	return http.StatusInternalServerError
}

// HTTPStatusToType infers the error Type from a HTTP Status code
func HTTPStatusToType(status int) Type {
	switch status {
	case http.StatusBadRequest:
		return InvalidArgument
	case http.StatusNotFound:
		return NotFound
	case http.StatusConflict:
		return Conflict
	case http.StatusUnauthorized:
		return Unauthorized
	case http.StatusForbidden:
		return PermissionDenied
	case http.StatusRequestTimeout:
		return Timeout
	case http.StatusNotImplemented:
		return NotImplemented
	case http.StatusBadGateway:
	case http.StatusServiceUnavailable:
		return TemporarilyUnavailable
	case http.StatusGone:
		return PermanentlyUnavailable
	case http.StatusTooManyRequests:
		return ResourceExhausted
	case http.StatusInternalServerError:
		return Unknown
	}
	return Unknown
}

// FromHTTP parses the http.Response and returns the corresponding
// If the response is not an error (eg. 200 OK), it returns nil
func FromHTTP(resp *http.Response) Error {
	if resp.StatusCode < 399 {
		return nil
	}

	typ := HTTPStatusToType(resp.StatusCode)

	out := &impl{
		message: typ.String(),
		code:    parseCode(resp.Header.Get(CodeHeader)),
		typ:     typ,
	}

	// try to decode the error from the body
	defer resp.Body.Close()
	j := new(jsonError)
	err := json.NewDecoder(resp.Body).Decode(j)
	if err == nil {
		out.message = j.Message
		out.code = j.Code
		out.typ = j.Type
		out.attributes = j.Attributes
	}

	return out
}

// ToHTTP writes the error to the http response
func ToHTTP(in error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if err, ok := in.(Error); ok {
		w.Header().Set(CodeHeader, err.Code().String())
		w.WriteHeader(err.Type().HTTPStatusCode())
		return json.NewEncoder(w).Encode(toJSON(err))
	}

	w.WriteHeader(http.StatusInternalServerError)
	return json.NewEncoder(w).Encode(&jsonError{
		Message: in.Error(),
		Type:    Unknown,
	})
}

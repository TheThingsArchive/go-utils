package errors

import (
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TypeFromHTTPStatus(status int) Type {
	switch status {
	case http.StatusInternalServerError:
		return Internal
	case http.StatusBadRequest, http.StatusMethodNotAllowed, http.StatusLengthRequired, http.StatusNotAcceptable, http.StatusPreconditionFailed:
		return InvalidArgument
	case http.StatusForbidden:
		return Forbidden
	case http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusServiceUnavailable:
		return Unavailable
	case http.StatusConflict:
		return AlreadyExists
	case http.StatusUnauthorized, http.StatusNetworkAuthenticationRequired:
		return Unauthenticated
	case http.StatusNotFound:
		return NotFound
	case http.StatusNotImplemented:
		return NotImplemented
	}

	return Unknown
}

func HTTPStatusFromType(typ Type) int {
	switch typ {
	case Internal, Unknown:
		return http.StatusInternalServerError
	case InvalidArgument:
		return http.StatusBadRequest
	case Forbidden:
		return http.StatusForbidden
	case Unavailable:
		return http.StatusServiceUnavailable
	case AlreadyExists:
		return http.StatusConflict
	case Unauthenticated:
		return http.StatusUnauthorized
	case NotFound:
		return http.StatusNotFound
	case NotImplemented:
		return http.StatusNotImplemented
	}

	return http.StatusInternalServerError
}

func FromHTTPStatus(status int) Error {
	typ := TypeFromHTTPStatus(status)
	return Err(typ, string(typ))
}

func ToHTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	if te, ok := err.(Typeable); ok {
		return HTTPStatusFromType(te.Type())
	}

	code := grpc.Code(err)
	if code != codes.Unknown {
		return HTTPStatusFromType(TypeFromGRPCCode(code))
	}

	return http.StatusInternalServerError
}

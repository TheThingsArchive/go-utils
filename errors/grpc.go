package errors

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GRPCCode returns the corresponding http status code from an error type
func (t Type) GRPCCode() codes.Code {
	// TODO
	return codes.Unknown
}

// GRPCCode returns the corresponding http status code from an error
func GRPCCode(err error) codes.Code {
	e, ok := err.(Error)
	if ok {
		return e.Type().GRPCCode()
	}

	return grpc.Code(err)
}

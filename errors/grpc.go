package errors

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TypeFromGRPCCode(code codes.Code) Type {
	switch code {
	case codes.AlreadyExists:
		return AlreadyExists
	case codes.InvalidArgument, codes.FailedPrecondition:
		return InvalidArgument
	case codes.NotFound:
		return NotFound
	case codes.PermissionDenied:
		return Forbidden
	case codes.Unauthenticated:
		return Unauthenticated
	case codes.Internal:
		return Internal
	case codes.Unknown:
		return Unknown
	}

	return Unknown
}

func GRPCCodeFromType(typ Type) codes.Code {
	switch typ {
	case AlreadyExists:
		return codes.AlreadyExists
	case InvalidArgument:
		return codes.InvalidArgument
	case NotFound:
		return codes.NotFound
	case Forbidden:
		return codes.PermissionDenied
	case Unauthenticated:
		return codes.Unauthenticated
	case Internal:
		return codes.Internal
	case Unknown:
		return codes.Unknown
	}

	return codes.Unknown
}

func FromGRPCError(err error) Error {
	if err == nil {
		return nil
	}

	code := grpc.Code(err)
	desc := grpc.ErrorDesc(err)
	typ := TypeFromGRPCCode(code)
	return New(typ, desc)
}

func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	code := grpc.Code(err)
	if code != codes.Unknown {
		return err
	}

	if te, ok := err.(Typeable); ok {
		return grpc.Errorf(GRPCCodeFromType(te.Type()), err.Error())
	}

	return grpc.Errorf(codes.Unknown, err.Error())
}

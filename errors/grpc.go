package errors

import (
	"encoding/json"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GRPCCode returns the corresponding http status code from an error type
func (t Type) GRPCCode() codes.Code {
	switch t {
	case InvalidArgument:
		return codes.InvalidArgument
	case OutOfRange:
		return codes.OutOfRange
	case NotFound:
		return codes.NotFound
	case Conflict:
	case AlreadyExists:
		return codes.AlreadyExists
	case Unauthorized:
		return codes.Unauthenticated
	case PermissionDenied:
		return codes.PermissionDenied
	case Timeout:
		return codes.DeadlineExceeded
	case NotImplemented:
		return codes.Unimplemented
	case TemporarilyUnavailable:
		return codes.Unavailable
	case PermanentlyUnavailable:
		return codes.FailedPrecondition
	case Canceled:
		return codes.Canceled
	case ResourceExhausted:
		return codes.ResourceExhausted
	case Internal:
	case Unknown:
		return codes.Unknown
	}

	return codes.Unknown
}

func GRPCCodeToType(code codes.Code) Type {
	switch code {
	case codes.InvalidArgument:
		return InvalidArgument
	case codes.OutOfRange:
		return OutOfRange
	case codes.NotFound:
		return NotFound
	case codes.AlreadyExists:
		return AlreadyExists
	case codes.Unauthenticated:
		return Unauthorized
	case codes.PermissionDenied:
		return PermissionDenied
	case codes.DeadlineExceeded:
		return Timeout
	case codes.Unimplemented:
		return NotImplemented
	case codes.Unavailable:
		return TemporarilyUnavailable
	case codes.FailedPrecondition:
		return PermanentlyUnavailable
	case codes.Canceled:
		return Canceled
	case codes.ResourceExhausted:
		return ResourceExhausted
	case codes.Unknown:
		return Unknown
	}
	return Unknown
}

// GRPCCode returns the corresponding http status code from an error
func GRPCCode(err error) codes.Code {
	e, ok := err.(Error)
	if ok {
		return e.Type().GRPCCode()
	}

	return grpc.Code(err)
}

var regex = regexp.MustCompile(`.*desc = (.*) \(e:(\d+)\) attributes = (.*)`)
var format = "%s (e:%v) attributes = %s"

// FromGRPC parses a gRPC error and returns an Error
func FromGRPC(in error) Error {
	out := &impl{
		Imessage: grpc.ErrorDesc(in),
		Ityp:     GRPCCodeToType(grpc.Code(in)),
		Icode:    Code(0),
	}

	matches := regex.FindStringSubmatch(in.Error())

	if len(matches) < 4 {
		return out
	}

	// set message
	out.Imessage = matches[1]
	out.Icode = parseCode(matches[2])
	_ = json.Unmarshal([]byte(matches[3]), &out.Iattributes)

	got := Get(Code(out.Icode))
	if got == nil {
		return out
	}

	// Todo: find attributes
	return got.New(out.Iattributes)
}

// ToGRPC turns an error into a gRPC error
func ToGRPC(in error) error {
	if err, ok := in.(Error); ok {
		attrs, _ := json.Marshal(err.Attributes())
		return grpc.Errorf(err.Type().GRPCCode(), format, err.Error(), err.Code(), attrs)
	}

	return grpc.Errorf(codes.Unknown, in.Error())
}

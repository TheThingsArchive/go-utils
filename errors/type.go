package errors

type Type uint8

const (
	// Unknown is the type of unknown or unexpected errors
	Unknown Type = iota

	// Internal is the type of internal errors
	Internal

	// InvalidArgument is the type of errors that result from an invalid argument
	// in a request
	InvalidArgument

	// OutOfRange is the type of errors that result from an out of range request
	OutOfRange

	// NotFound is the type of errors that result from an entity that is not found
	// or not accessible
	NotFound

	// Conflict is the type of errors that result from a conflict
	Conflict

	// AlreadyExists is the type of errors that result from a conflict where the
	// updated/created entity already exists
	AlreadyExists

	// Unauthorized is the type of errors where the request is unauthorized where
	// it should be
	Unauthorized

	// PermissionDenied is the type of errors where the request was authorized but
	// did not grant access to the requested entity
	PermissionDenied

	// Timeout is the type of errors that are a result of a process taking too
	// long to complete
	Timeout

	// NotImplemented is the type of errors that result from a requested action
	// that is not (yet) implemented
	NotImplemented

	// TemporarilyUnavailable is the type of errors that result from a service
	// being temporarily unavailable (down)
	TemporarilyUnavailable

	// PermanentlyUnavailable is the type of errors that result from an action
	// that has been deprecated and is no longer available
	PermanentlyUnavailable
)

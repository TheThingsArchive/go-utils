package errors

type Type string

const (
	// Unknown means that the error type is unknown
	Unknown Type = "unknown"

	// Internal means something went wrong on the server side
	Internal Type = "internal"

	// Unavailable means a service the action depends on is unavailable
	Unavailable Type = "unavailable"

	// Unavailable means a service the action depends on is unavailable
	TemporaryUnavailable Type = "temporary unavailable"

	// AlreadyExists means that the entity already exists
	AlreadyExists Type = "already exists"

	// NotImplemented means action is known but not implemented yet
	NotImplemented Type = "not implemented"

	// Unauthenticated means the client has made no authentication, while it
	// should have
	Unauthenticated Type = "unauthenticated"

	// Forbidden means that the client is authenticated by not allowed to perform
	// the action
	Forbidden Type = "forbidden"

	// InvalidArgument means that an argument was not valid
	InvalidArgument Type = "invalid argument"

	// OutOfRange
	OutOfRange Type = "out of range"

	// NotFound means that the entity does not exist
	NotFound Type = "not found"
)

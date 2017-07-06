package errors

import (
	"fmt"
	"strings"
)

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

// string representations of the Types
// keep up to date with the iota
var str = []string{
	"Unknown",
	"Internal",
	"Invalid argument",
	"Out of range",
	"Not found",
	"Conflict",
	"Already exists",
	"Unauthorized",
	"Permission denied",
	"Timeout",
	"Not implemented",
	"Temporarily unavailable",
	"Permanently unavailable",
}

// String implements stringer
func (t Type) String() string {
	return str[t]
}

// MarshalText implements TextMarsheler
func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements TextUnmarsheler
func (t *Type) UnmarshalText(text []byte) error {
	enum := strings.ToLower(string(text))
	for i, typ := range str {
		if enum == strings.ToLower(typ) {
			*t = Type(i)
			return nil
		}
	}

	return fmt.Errorf("Invalid event type")
}

// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

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

	// Canceled indicates the operation was canceled (typically by the caller)
	Canceled

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted
)

// String implements stringer
func (t Type) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Internal:
		return "Internal"
	case InvalidArgument:
		return "Invalid argument"
	case OutOfRange:
		return "Out of range"
	case NotFound:
		return "Not found"
	case Conflict:
		return "Conflict"
	case AlreadyExists:
		return "Already exists"
	case Unauthorized:
		return "Unauthorized"
	case PermissionDenied:
		return "Permission denied"
	case Timeout:
		return "Timeout"
	case NotImplemented:
		return "Not implemented"
	case TemporarilyUnavailable:
		return "Temporarily unavailable"
	case PermanentlyUnavailable:
		return "Permanently unavailable"
	case Canceled:
		return "Canceled"
	case ResourceExhausted:
		return "Resource exhausted"
	default:
		return "Unknown"
	}
}

// MarshalText implements TextMarsheler
func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements TextUnmarsheler
func (t *Type) UnmarshalText(text []byte) error {
	e, err := fromString(string(text))
	if err != nil {
		return err
	}

	*t = e

	return nil
}

func fromString(str string) (Type, error) {
	enum := strings.ToLower(str)
	switch enum {
	case "unknown":
		return Unknown, nil
	case "internal":
		return Internal, nil
	case "invalid argument":
		return InvalidArgument, nil
	case "out of range":
		return OutOfRange, nil
	case "not found":
		return NotFound, nil
	case "conflict":
		return Conflict, nil
	case "already exists":
		return AlreadyExists, nil
	case "unauthorized":
		return Unauthorized, nil
	case "permission denied":
		return PermissionDenied, nil
	case "timeout":
		return Timeout, nil
	case "not implemented":
		return NotImplemented, nil
	case "temporarily unavailable":
		return TemporarilyUnavailable, nil
	case "permanently unavailable":
		return PermanentlyUnavailable, nil
	case "canceled":
		return Canceled, nil
	case "resource exhausted":
		return ResourceExhausted, nil
	default:
		return Unknown, fmt.Errorf("Invalid event type")
	}
}

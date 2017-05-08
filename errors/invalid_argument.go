// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// NewErrInvalidArgument returns a new errInvalidArgument for the given entitiy
func NewErrInvalidArgument(what, why string) error {
	return &errInvalidArgument{base: newBase(""), what: what, why: why}
}

// errInvalidArgument indicates that an entity already exists
type errInvalidArgument struct {
	base
	what string
	why  string
}

func (err errInvalidArgument) Code() codes.Code {
	return codes.InvalidArgument
}

func (err errInvalidArgument) Error() string {
	return fmt.Sprintf(`invalid argument "%s": %s`, err.what, err.why)
}

func (err errInvalidArgument) Fields() Fields {
	fields := err.base.Fields()
	fields["what"] = err.what
	fields["why"] = err.why
	return fields
}

func newErrInvalidArgumentFrom(msg string, f Fields) (error, bool) {
	what, ok := f["what"]
	why, _ := f["why"]
	if ok {
		return NewErrInvalidArgument(what.(string), why.(string)), true
	}
	return New(msg), false
}

// IsInvalidArgument returns whether error type is InvalidArgument
func IsInvalidArgument(err error) bool {
	return FindCode(FindCause(err)) == codes.InvalidArgument
}

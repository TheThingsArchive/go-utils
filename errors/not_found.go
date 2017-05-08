// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// NewErrNotFound returns a new errNotFound
func NewErrNotFound(what string) error {
	return errNotFound{base: newBase(""), what: what}
}

type errNotFound struct {
	base
	what string
}

func (err errNotFound) Code() codes.Code {
	return codes.NotFound
}

func (err errNotFound) Error() string {
	return fmt.Sprintf("%s not found", err.what)
}

func (err errNotFound) Fields() Fields {
	fields := err.base.Fields()
	fields["what"] = err.what
	return fields
}

func newErrNotFoundFrom(msg string, f Fields) (error, bool) {
	if what, ok := f["what"]; ok {
		return NewErrNotFound(what.(string)), true
	}
	return New(msg), false
}

// IsNotFound returns whether error type is NotFound
func IsNotFound(err error) bool {
	return FindCode(FindCause(err)) == codes.NotFound
}

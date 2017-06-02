// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import "fmt"
import "google.golang.org/grpc/codes"

// NewErrAlreadyExists returns a new errAlreadyExists
func NewErrAlreadyExists(what string) error {
	return errAlreadyExists{base: newBase(""), what: what}
}

type errAlreadyExists struct {
	base
	what string
}

func (err errAlreadyExists) Code() codes.Code {
	return codes.AlreadyExists
}

func (err errAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", err.what)
}

func (err errAlreadyExists) Fields() Fields {
	fields := err.base.Fields()
	fields["what"] = err.what
	return fields
}

func newErrAlreadyExistsFrom(msg string, f Fields) (error, bool) {
	if what, ok := f["what"]; ok {
		return NewErrAlreadyExists(what.(string)), true
	}
	return New(msg), false
}

// IsAlreadyExists returns whether error type is AlreadyExists
func IsAlreadyExists(err error) bool {
	return FindCode(FindCause(err)) == codes.AlreadyExists
}

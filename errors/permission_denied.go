// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// NewErrPermissionDenied returns a new errPermissionDenied
func NewErrPermissionDenied(why string) error {
	return errPermissionDenied{base: newBase(""), why: why}
}

type errPermissionDenied struct {
	base
	why string
}

func (err errPermissionDenied) Code() codes.Code {
	return codes.PermissionDenied
}

func (err errPermissionDenied) Error() string {
	return fmt.Sprintf("permission denied: %s", err.why)
}

func (err errPermissionDenied) Fields() Fields {
	fields := err.base.Fields()
	fields["why"] = err.why
	return fields
}

func newErrPermissionDeniedFrom(msg string, f Fields) (error, bool) {
	if why, ok := f["why"]; ok {
		return NewErrPermissionDenied(why.(string)), true
	}
	return New(msg), false
}

// IsPermissionDenied returns whether error type is PermissionDenied
func IsPermissionDenied(err error) bool {
	return FindCode(FindCause(err)) == codes.PermissionDenied
}

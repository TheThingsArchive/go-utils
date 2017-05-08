// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// NewErrInternal returns a new errInternal
func NewErrInternal(text string) error {
	return errInternal{base: newBase(text)}
}

type errInternal struct {
	base
}

func (err errInternal) Code() codes.Code {
	return codes.FailedPrecondition
}

func (err errInternal) Error() string {
	return fmt.Sprintf("internal error: %s", err.message)
}

// IsInternal returns whether error type is Internal
func IsInternal(err error) bool {
	return FindCode(FindCause(err)) == codes.FailedPrecondition
}

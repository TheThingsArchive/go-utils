// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"google.golang.org/grpc/codes"
)

type hasCode interface {
	Code() codes.Code
}

// FindCode finds the code of an error
func FindCode(err error) codes.Code {
	if err, ok := err.(hasCode); ok {
		return err.Code()
	}
	return codes.Unknown
}

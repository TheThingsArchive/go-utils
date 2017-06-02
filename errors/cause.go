// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

type hasCause interface {
	Cause() error
}

// FindCause finds the cause of an error
func FindCause(err error) error {
	if err, ok := err.(hasCause); ok {
		return err.Cause()
	}
	return err
}

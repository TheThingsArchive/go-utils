// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import "errors"

const causeKey = "cause"

// Causer is the type of errors that can have a cause
type Causer interface {
	Cause() error
}

// Cause returns the cause of an error
func Cause(err Error) error {
	attributes := err.Attributes()
	if attributes == nil {
		return nil
	}

	cause, ok := attributes[causeKey]
	if !ok {
		return nil
	}

	switch v := cause.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return nil
	}
}

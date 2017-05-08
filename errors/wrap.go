// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

type wrapped struct {
	base
}

// Wrap an existing error
func Wrap(err error, text string) error {
	wrapped := &wrapped{base: newBase(text)}
	wrapped.base.cause = err
	return wrapped
}

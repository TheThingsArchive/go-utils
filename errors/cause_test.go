// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"testing"

	s "github.com/smartystreets/assertions"
)

func TestCause(t *testing.T) {
	a := s.New(t)

	var err, wrapped error

	err = errors.New("err")
	wrapped = Wrap(err, "wrapped")
	a.So(FindCause(wrapped), s.ShouldEqual, err)
}

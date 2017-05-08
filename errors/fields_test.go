// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"testing"

	s "github.com/smartystreets/assertions"
)

func TestFields(t *testing.T) {
	a := s.New(t)

	any, err := (Fields{
		"foo":           "bar",
		"num":           42,
		"baz":           true,
		"not-supported": []string{"sorry"},
	}).Proto()
	a.So(err, s.ShouldBeNil)

	fields, err := FieldsFromProto(any)
	a.So(err, s.ShouldBeNil)
	a.So(fields, s.ShouldHaveLength, 3)
	a.So(fields["foo"], s.ShouldEqual, "bar")
	a.So(fields["num"], s.ShouldEqual, 42)
	a.So(fields["baz"], s.ShouldBeTrue)
}

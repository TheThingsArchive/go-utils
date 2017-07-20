// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions"
)

// code creates Codes for testing
var code = Range(10000, 11000)

func TestHTTP(t *testing.T) {
	a := assertions.New(t)

	d := &ErrDescriptor{
		MessageFormat: "You do not have access to app with id {app_id}",
		Code:          code(77),
		Type:          PermissionDenied,
		registered:    true,
	}

	attributes := Attributes{
		"app_id": "foo",
		"count":  42,
	}

	err := d.New(attributes)

	w := httptest.NewRecorder()
	e := ToHTTP(err, w)
	a.So(e, assertions.ShouldBeNil)

	resp := w.Result()

	got := FromHTTP(resp)
	a.So(got.Code(), assertions.ShouldEqual, err.Code())
	a.So(got.Type(), assertions.ShouldEqual, err.Type())
	a.So(got.Error(), assertions.ShouldEqual, err.Error())
	a.So(got.Attributes()["app_id"], assertions.ShouldResemble, attributes["app_id"])
	a.So(got.Attributes()["count"], assertions.ShouldAlmostEqual, attributes["count"])
}

func TestToUnspecifiedHTTP(t *testing.T) {
	a := assertions.New(t)

	err := errors.New("A random error")

	w := httptest.NewRecorder()
	e := ToHTTP(err, w)
	a.So(e, assertions.ShouldBeNil)

	resp := w.Result()

	got := FromHTTP(resp)
	a.So(got.Code(), assertions.ShouldEqual, NoCode)
	a.So(got.Type(), assertions.ShouldEqual, Unknown)
	a.So(got.Error(), assertions.ShouldEqual, err.Error())
	a.So(got.Attributes(), assertions.ShouldBeNil)
}

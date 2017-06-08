package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
)

func TestDescriptor(t *testing.T) {
	a := assertions.New(t)

	d := &ErrDescriptor{
		MessageFormat: "You do not have access to app with id {app_id}",
		Code:          77,
		Type:          PermissionDenied,
		registered:    true,
	}

	attributes := Attributes{
		"app_id": "foo",
	}
	err := New(d, attributes)

	a.So(err.Error(), assertions.ShouldEqual, "You do not have access to app with id foo")
	a.So(err.Code(), assertions.ShouldEqual, d.Code)
	a.So(err.Type(), assertions.ShouldEqual, d.Type)
	a.So(err.Attributes(), assertions.ShouldResemble, attributes)
}

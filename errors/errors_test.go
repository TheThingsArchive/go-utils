package errors

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/assertions"
)

func TestErrorWithFields(t *testing.T) {
	a := New(t)

	err := Err(Unknown, "Error")

	a.So(err.Type(), ShouldEqual, Unknown)
	a.So(err.Error(), ShouldEqual, "Error")
	a.So(err.Fields(), ShouldResemble, map[string]interface{}{})

	// adding fields to an error returns them as well
	a.So(err.WithField("foo", "bar").Fields(), ShouldResemble, map[string]interface{}{
		"foo": "bar",
	})

	// the original error should have no fields
	a.So(err.Fields(), ShouldResemble, map[string]interface{}{})

	fmt.Println(ToGRPCError(err))
}

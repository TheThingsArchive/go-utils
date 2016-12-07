package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	s "github.com/smartystreets/assertions"
)

func TestErrorWithFields(t *testing.T) {
	a := s.New(t)

	err := New(Unknown, "Error")

	a.So(err.Type(), s.ShouldEqual, Unknown)
	a.So(err.Error(), s.ShouldEqual, "Error")
	a.So(err.Fields(), s.ShouldResemble, map[string]interface{}{})

	// adding fields to an error returns them as well
	a.So(err.WithField("foo", "bar").Fields(), s.ShouldResemble, map[string]interface{}{
		"foo": "bar",
	})

	// the original error should have no fields
	a.So(err.Fields(), s.ShouldResemble, map[string]interface{}{})

	fmt.Println(ToGRPCError(err))
}

func TestErrorCause(t *testing.T) {
	a := s.New(t)

	root := New(Unknown, "Foo")
	first := New(Unknown, "Bar").WithCause(root)
	second := New(Unknown, "Baz").WithCause(first)
	third := New(Unknown, "Qux").WithCause(second)

	a.So(third.Cause(), s.ShouldEqual, second)
	a.So(second.Cause(), s.ShouldEqual, first)
	a.So(first.Cause(), s.ShouldEqual, root)

	a.So(RootCause(third), s.ShouldEqual, root)
	a.So(RootCause(second), s.ShouldEqual, root)
	a.So(RootCause(first), s.ShouldEqual, root)
	a.So(RootCause(root), s.ShouldEqual, root)
}

func TestFrom(t *testing.T) {
	a := s.New(t)

	// io is out of range
	{
		err := From(io.EOF)
		a.So(err.Type(), s.ShouldEqual, OutOfRange)
		a.So(err.Error(), s.ShouldEqual, io.EOF.Error())
	}

	// plain error is uknnown
	{
		err := From(errors.New("foo"))
		a.So(err.Type(), s.ShouldEqual, Unknown)
		a.So(err.Error(), s.ShouldEqual, "foo")
	}

	// parse grpc code
	{
		err := From(grpc.Errorf(codes.Unauthenticated, "derp"))
		a.So(err.Type(), s.ShouldEqual, Unauthenticated)
		a.So(err.Error(), s.ShouldEqual, "derp")
	}
}

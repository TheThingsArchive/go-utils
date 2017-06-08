package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
)

func TestTypeString(t *testing.T) {
	a := assertions.New(t)

	a.So(Unknown.String(), assertions.ShouldEqual, "Unknown")
	a.So(Timeout.String(), assertions.ShouldEqual, "Timeout")
}

func TestTypeMarshal(t *testing.T) {
	a := assertions.New(t)

	text, err := Unknown.MarshalText()
	a.So(err, assertions.ShouldBeNil)
	a.So(text, assertions.ShouldResemble, []byte("Unknown"))
}

func TestTypeUnmarshal(t *testing.T) {
	a := assertions.New(t)

	var typ Type
	err := typ.UnmarshalText([]byte("Temporarily unavailable"))
	a.So(err, assertions.ShouldBeNil)
	a.So(typ, assertions.ShouldEqual, TemporarilyUnavailable)

	err = typ.UnmarshalText([]byte("temporarily unavailable"))
	a.So(err, assertions.ShouldBeNil)
	a.So(typ, assertions.ShouldEqual, TemporarilyUnavailable)
}

package entries

import (
	"testing"

	"github.com/apex/log"
	"github.com/smartystreets/assertions"
)

func TestEntries(t *testing.T) {
	a := assertions.New(t)

	e := New()

	e.HandleLog(&log.Entry{
		Message: "foo",
		Fields: log.Fields{
			"a": 10,
		},
	})

	a.So(len(e.Entries), assertions.ShouldEqual, 1)
	a.So(e.Entries[0].Message, assertions.ShouldEqual, "foo")
	a.So(e.Entries[0].Fields["a"], assertions.ShouldEqual, 10)

	e.HandleLog(&log.Entry{
		Message: "bar",
		Fields: log.Fields{
			"b": 20,
		},
	})

	a.So(len(e.Entries), assertions.ShouldEqual, 2)
	a.So(e.Entries[1].Message, assertions.ShouldEqual, "bar")
	a.So(e.Entries[1].Fields["b"], assertions.ShouldEqual, 20)

}

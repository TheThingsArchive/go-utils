package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
)

func TestFormat(t *testing.T) {
	a := assertions.New(t)

	format := "{foo} - {bar} - {nil} - {list} - {map}"
	{
		res := Format(format, Attributes{
			"foo":  10,
			"bar":  "bar",
			"list": []int{1, 2, 3},
			"map":  map[string]int{"ok": 1},
		}, defaultValueFormatter)
		a.So(res, assertions.ShouldEqual, "10 - bar - <nil> - [1 2 3] - map[ok:1]")
	}
}

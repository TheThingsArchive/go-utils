package namespaced

import (
	"testing"

	. "github.com/smartystreets/assertions"
)

func TestNamespaces(t *testing.T) {
	a := New(t)

	ns := &ns{}

	// empty namespace is accepted
	{
		ns.Set([]string{})

		a.So(ns.IsEnabled("a"), ShouldBeFalse)
		a.So(ns.IsEnabled("b"), ShouldBeFalse)
		a.So(ns.IsEnabled("c"), ShouldBeFalse)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}

	// add some namespaces
	{
		ns.Set([]string{
			"a",
			"b",
		})

		a.So(ns.IsEnabled("a"), ShouldBeTrue)
		a.So(ns.IsEnabled("b"), ShouldBeTrue)
		a.So(ns.IsEnabled("c"), ShouldBeFalse)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}

	// * accepts everything
	{
		ns.Set([]string{
			"*",
		})

		a.So(ns.IsEnabled("a"), ShouldBeTrue)
		a.So(ns.IsEnabled("b"), ShouldBeTrue)
		a.So(ns.IsEnabled("c"), ShouldBeTrue)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}

	// negation wins from *
	{
		ns.Set([]string{
			"*",
			"-a",
		})

		a.So(ns.IsEnabled("a"), ShouldBeFalse)
		a.So(ns.IsEnabled("b"), ShouldBeTrue)
		a.So(ns.IsEnabled("c"), ShouldBeTrue)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}

	// order should not matter
	{
		ns.Set([]string{
			"-a",
			"-b",
			"*",
		})

		a.So(ns.IsEnabled("a"), ShouldBeFalse)
		a.So(ns.IsEnabled("b"), ShouldBeFalse)
		a.So(ns.IsEnabled("c"), ShouldBeTrue)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}

	// negation always wins
	{
		ns.Set([]string{
			"a",
			"-a",
		})

		a.So(ns.IsEnabled("a"), ShouldBeFalse)
		a.So(ns.IsEnabled(""), ShouldBeTrue)
	}
}

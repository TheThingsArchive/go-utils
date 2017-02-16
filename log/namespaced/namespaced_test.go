package namespaced

import (
	"testing"

	"github.com/TheThingsNetwork/go-utils/handlers/entries"
	apx "github.com/TheThingsNetwork/go-utils/log/apex"
	apex "github.com/apex/log"
	. "github.com/smartystreets/assertions"
)

func TestNamespaced(t *testing.T) {
	a := New(t)

	handler := entries.New()

	ctx := Wrap(apx.Wrap(&apex.Logger{
		Level:   apex.DebugLevel,
		Handler: handler,
	}))

	// should just log messages without a namespace
	{
		ctx.Info("message")
		a.So(len(handler.Entries), ShouldEqual, 1)
	}

	// should just log messages without a namespace
	// even when a namespace is set
	{
		ctx.SetNamespaces("foo")
		ctx.Info("message 2")
		a.So(len(handler.Entries), ShouldEqual, 2)
	}

	// correctly namespaced loggers should log
	{
		foo := Namespace(ctx, "foo")
		foo.Info("message 3")
		a.So(len(handler.Entries), ShouldEqual, 3)
	}

	// incorrectly namespaced loggers should not log
	{
		bar := Namespace(ctx, "bar")
		bar.Info("message 4 (bar) should be ignored")
		a.So(len(handler.Entries), ShouldEqual, 3)

		// set the namspaces to include bar and log bar
		ctx.SetNamespaces("foo", "bar")
		bar.Info("message 5 (bar)")
		a.So(len(handler.Entries), ShouldEqual, 4)

		// set the namspaces to include bar and log bar
		ctx.SetNamespaces()
		bar.Info("message 6 (bar) should be ignored")
		a.So(len(handler.Entries), ShouldEqual, 4)
	}

	for _, entry := range handler.Entries {
		// fmt.Println(entry.Message)
		a.So(entry.Message, ShouldNotEndWith, "should be ignored")
		a.So(entry.Fields, ShouldNotContainKey, "namespace")
	}
}

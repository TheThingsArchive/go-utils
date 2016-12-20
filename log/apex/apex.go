package apex

import (
	"os"

	"github.com/TheThingsNetwork/go-utils/handlers/cli"
	"github.com/TheThingsNetwork/go-utils/log"
	apex "github.com/apex/log"
)

// Stdout logging apex/log
func Stdout() log.Interface {
	return Wrap(&apex.Logger{
		Level:   apex.InfoLevel,
		Handler: cli.New(os.Stdout),
	})
}

// Wrap apex.Interface
func Wrap(ctx apex.Interface) log.Interface {
	return &apexInterfaceWrapper{ctx}
}

type apexInterfaceWrapper struct {
	apex.Interface
}

func (w *apexInterfaceWrapper) WithField(k string, v interface{}) log.Interface {
	return &apexEntryWrapper{w.Interface.WithField(k, v)}
}

func (w *apexInterfaceWrapper) WithFields(fields log.Fields) log.Interface {
	return &apexEntryWrapper{w.Interface.WithFields(apex.Fields(fields))}
}

func (w *apexInterfaceWrapper) WithError(err error) log.Interface {
	return &apexEntryWrapper{w.Interface.WithError(err)}
}

type apexEntryWrapper struct {
	*apex.Entry
}

func (w *apexEntryWrapper) WithField(k string, v interface{}) log.Interface {
	return &apexEntryWrapper{w.Entry.WithField(k, v)}
}

func (w *apexEntryWrapper) WithFields(fields log.Fields) log.Interface {
	return &apexEntryWrapper{w.Entry.WithFields(apex.Fields(fields))}
}

func (w *apexEntryWrapper) WithError(err error) log.Interface {
	return &apexEntryWrapper{w.Entry.WithError(err)}
}

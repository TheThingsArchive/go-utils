package apex

import (
	"os"

	"github.com/TheThingsNetwork/go-utils/handlers/cli"
	"github.com/TheThingsNetwork/go-utils/log"
	apex "github.com/apex/log"
)

// Stdout logging apex/log
func Stdout() *apexInterfaceWrapper {
	return Wrap(&apex.Logger{
		Level:   apex.InfoLevel,
		Handler: cli.New(os.Stdout),
	})
}

const (
	DebugLevel = apex.DebugLevel
	InfoLevel  = apex.InfoLevel
	WarnLevel  = apex.WarnLevel
	ErrorLevel = apex.ErrorLevel
	FatalLevel = apex.FatalLevel
)

var (
	ParseLevel     = apex.ParseLevel
	MustParseLevel = apex.MustParseLevel
)

// Wrap apex.Interface
func Wrap(ctx *apex.Logger) *apexInterfaceWrapper {
	return &apexInterfaceWrapper{ctx}
}

type apexInterfaceWrapper struct {
	*apex.Logger
}

func (w *apexInterfaceWrapper) WithField(k string, v interface{}) log.Interface {
	return &apexEntryWrapper{w.Logger.WithField(k, v)}
}

func (w *apexInterfaceWrapper) WithFields(fields log.Fields) log.Interface {
	return &apexEntryWrapper{w.Logger.WithFields(apex.Fields(fields))}
}

func (w *apexInterfaceWrapper) WithError(err error) log.Interface {
	return &apexEntryWrapper{w.Logger.WithError(err)}
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

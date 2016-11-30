package apex

import (
	"github.com/TheThingsNetwork/go-utils/log"
	apex "github.com/apex/log"
)

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

func (w *apexInterfaceWrapper) WithFields(fields map[string]interface{}) log.Interface {
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

func (w *apexEntryWrapper) WithFields(fields map[string]interface{}) log.Interface {
	return &apexEntryWrapper{w.Entry.WithFields(apex.Fields(fields))}
}

func (w *apexEntryWrapper) WithError(err error) log.Interface {
	return &apexEntryWrapper{w.Entry.WithError(err)}
}

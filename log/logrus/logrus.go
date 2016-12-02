package logrus

import (
	"github.com/Sirupsen/logrus"
	"github.com/TheThingsNetwork/go-utils/log"
)

// StandardLogrus wraps the standard Logrus Logger into a Logger
func StandardLogrus() log.Interface {
	return Wrap(logrus.StandardLogger())
}

// Wrap logrus.Logger
func Wrap(logger *logrus.Logger) log.Interface {
	return &logrusEntryWrapper{logrus.NewEntry(logger)}
}

type logrusEntryWrapper struct {
	*logrus.Entry
}

func (w *logrusEntryWrapper) Debug(msg string) {
	w.Entry.Debug(msg)
}

func (w *logrusEntryWrapper) Info(msg string) {
	w.Entry.Info(msg)
}

func (w *logrusEntryWrapper) Warn(msg string) {
	w.Entry.Warn(msg)
}

func (w *logrusEntryWrapper) Error(msg string) {
	w.Entry.Error(msg)
}

func (w *logrusEntryWrapper) Fatal(msg string) {
	w.Entry.Fatal(msg)
}

func (w *logrusEntryWrapper) WithError(err error) log.Interface {
	return &logrusEntryWrapper{w.Entry.WithError(err)}
}

func (w *logrusEntryWrapper) WithField(k string, v interface{}) log.Interface {
	return &logrusEntryWrapper{w.Entry.WithField(k, v)}
}

func (w *logrusEntryWrapper) WithFields(fields log.Fields) log.Interface {
	return &logrusEntryWrapper{w.Entry.WithFields(logrus.Fields(fields))}
}

package filtered

import (
	"testing"

	"github.com/TheThingsNetwork/go-utils/log"
	. "github.com/smartystreets/assertions"
)

// noopLogger just does nothing
type noopLogger struct{}

func (l noopLogger) Debug(msg string)                            {}
func (l noopLogger) Info(msg string)                             {}
func (l noopLogger) Warn(msg string)                             {}
func (l noopLogger) Error(msg string)                            {}
func (l noopLogger) Fatal(msg string)                            {}
func (l noopLogger) Debugf(msg string, v ...interface{})         {}
func (l noopLogger) Infof(msg string, v ...interface{})          {}
func (l noopLogger) Warnf(msg string, v ...interface{})          {}
func (l noopLogger) Errorf(msg string, v ...interface{})         {}
func (l noopLogger) Fatalf(msg string, v ...interface{})         {}
func (l noopLogger) WithField(string, interface{}) log.Interface { return l }
func (l noopLogger) WithFields(log.Fields) log.Interface         { return l }
func (l noopLogger) WithError(error) log.Interface               { return l }

// FieldsLogger is a logger that stores which fields it has been passed
type FieldsLogger struct {
	log.Interface
	Fields map[string]interface{}
}

// WithField saves the fields and calls the wrapped loggers WithField
func (f *FieldsLogger) WithField(k string, v interface{}) log.Interface {
	flds := make(map[string]interface{}, len(f.Fields)+1)
	for k, v := range f.Fields {
		flds[k] = v
	}

	// add the new field
	flds[k] = v

	return &FieldsLogger{
		Interface: f.Interface.WithField(k, v),
		Fields:    flds,
	}
}

// WithFields saves the fields and calls the wrapped loggers WithFields
func (f *FieldsLogger) WithFields(fields log.Fields) log.Interface {
	flds := make(map[string]interface{}, len(f.Fields)+1)
	for k, v := range f.Fields {
		flds[k] = v
	}

	// add the new field
	for k, v := range fields {
		flds[k] = v
	}

	return &FieldsLogger{
		Interface: f.Interface.WithFields(fields),
		Fields:    flds,
	}
}

func TestFieldsLogger(t *testing.T) {
	a := New(t)

	logger := &FieldsLogger{
		Interface: &noopLogger{},
	}

	logger = logger.WithField("foo", 42).(*FieldsLogger)

	a.So(logger.Fields, ShouldContainKey, "foo")
	a.So(logger.Fields["foo"], ShouldEqual, 42)

	logger = logger.WithFields(log.Fields{
		"bar": "lol",
		"baz": true,
	}).(*FieldsLogger)

	a.So(logger.Fields, ShouldContainKey, "foo")
	a.So(logger.Fields["foo"], ShouldEqual, 42)

	a.So(logger.Fields, ShouldContainKey, "bar")
	a.So(logger.Fields["bar"], ShouldEqual, "lol")

	a.So(logger.Fields, ShouldContainKey, "baz")
	a.So(logger.Fields["baz"], ShouldEqual, true)
}

func TestFilterSensitive(t *testing.T) {
	a := New(t)

	wrapped := Wrap(&FieldsLogger{
		Interface: &noopLogger{},
	}, DefaultSensitiveFilter)

	// should elide passwords
	{
		w := wrapped.WithField("password", "secret")
		fields := w.(*Filtered).Interface.(*FieldsLogger).Fields
		a.So(fields, ShouldContainKey, "password")
		a.So(fields["password"], ShouldEqual, defaultElided)
	}

	// should not elide other stuff
	{
		w := wrapped.WithField("foo", "bar")
		fields := w.(*Filtered).Interface.(*FieldsLogger).Fields
		a.So(fields, ShouldContainKey, "foo")
		a.So(fields["foo"], ShouldEqual, "bar")
	}

	// should work the same with more fields
	{
		w := wrapped.WithFields(log.Fields{
			"bar":   "baz",
			"token": "secret",
		})
		fields := w.(*Filtered).Interface.(*FieldsLogger).Fields
		a.So(fields, ShouldContainKey, "bar")
		a.So(fields["bar"], ShouldEqual, "baz")

		a.So(fields, ShouldContainKey, "token")
		a.So(fields["token"], ShouldEqual, defaultElided)
	}
}

func TestFilterQuery(t *testing.T) {
	a := New(t)

	wrapped := Wrap(&FieldsLogger{
		Interface: &noopLogger{},
	}, DefaultQueryFilter)

	w := wrapped.WithField("query", map[string][]string{
		"password": []string{"secret"},
		"foo":      []string{"bar"},
	})

	fields := w.(*Filtered).Interface.(*FieldsLogger).Fields

	a.So(fields, ShouldContainKey, "query")

	query := fields["query"].(map[string]interface{})

	a.So(query, ShouldContainKey, "password")
	a.So(query["password"], ShouldResemble, []interface{}{defaultElided})

	a.So(query, ShouldContainKey, "foo")
	a.So(query["foo"], ShouldResemble, []interface{}{"bar"})
}

func TestFilteredFields(t *testing.T) {
	a := New(t)

	wrapped := Wrap(&FieldsLogger{
		Interface: &noopLogger{},
	}, DefaultQueryFilter)

	fst := wrapped.WithField("foo", "bar")
	snd := wrapped.WithField("quu", "qux")

	fstFields := fst.(*Filtered).Interface.(*FieldsLogger).Fields
	sndFields := snd.(*Filtered).Interface.(*FieldsLogger).Fields

	a.So(fstFields, ShouldContainKey, "foo")
	a.So(fstFields, ShouldNotContainKey, "quu")

	a.So(sndFields, ShouldNotContainKey, "foo")
	a.So(sndFields, ShouldContainKey, "quu")

	fstbis := fst.WithFields(log.Fields{
		"foobis": "barbis",
	})
	sndbis := snd.WithFields(log.Fields{
		"quubis": "quxbis",
	})

	fstbisFields := fstbis.(*Filtered).Interface.(*FieldsLogger).Fields
	sndbisFields := sndbis.(*Filtered).Interface.(*FieldsLogger).Fields

	a.So(fstbisFields, ShouldContainKey, "foo")
	a.So(fstbisFields, ShouldNotContainKey, "quu")
	a.So(fstbisFields, ShouldContainKey, "foobis")
	a.So(fstbisFields, ShouldNotContainKey, "quubis")

	a.So(sndbisFields, ShouldNotContainKey, "foo")
	a.So(sndbisFields, ShouldContainKey, "quu")
	a.So(sndbisFields, ShouldNotContainKey, "foobix")
	a.So(sndbisFields, ShouldContainKey, "quubis")
}

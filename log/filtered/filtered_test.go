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
		wrapped.WithField("password", "secret")

		fields := wrapped.Interface.(*FieldsLogger).Fields
		a.So(fields, ShouldContainKey, "password")
		a.So(fields["password"], ShouldEqual, defaultElided)
	}

	// should not elide other stuff
	{
		wrapped.WithField("foo", "bar")
		fields := wrapped.Interface.(*FieldsLogger).Fields
		a.So(fields, ShouldContainKey, "foo")
		a.So(fields["foo"], ShouldEqual, "bar")
	}

	// should work the same with more fields
	{
		wrapped.WithFields(log.Fields{
			"bar":   "baz",
			"token": "secret",
		})
		fields := wrapped.Interface.(*FieldsLogger).Fields
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

	wrapped.WithField("query", map[string][]string{
		"password": []string{"secret"},
		"foo":      []string{"bar"},
	})

	fields := wrapped.Interface.(*FieldsLogger).Fields

	a.So(fields, ShouldContainKey, "query")

	query := fields["query"].(map[string]interface{})

	a.So(query, ShouldContainKey, "password")
	a.So(query["password"], ShouldResemble, []interface{}{defaultElided})

	a.So(query, ShouldContainKey, "foo")
	a.So(query["foo"], ShouldResemble, []interface{}{"bar"})
}

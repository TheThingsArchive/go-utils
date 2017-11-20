// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package prometheus

import (
	"errors"
	"testing"

	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

func TestPrometheus(t *testing.T) {
	a := New(t)

	wrapped := Wrap(&noopLogger{})

	withFields := wrapped.
		WithError(errors.New("a")).
		WithField("foo", "bar").
		WithFields(log.Fields{"bar": "baz"})

	ch := make(chan prometheus.Metric, 10)

	collect := func() []*dto.Metric {
		logCounter.Collect(ch)
		metrics := make([]*dto.Metric, 0, len(ch))
		for i := len(ch); i > 0; i-- {
			metric := <-ch
			m := new(dto.Metric)
			metric.Write(m)
			metrics = append(metrics, m)
		}
		return metrics
	}

	withFields.Debug("foo")
	withFields.Debugf("foo %d", 42)

	metrics := collect()
	a.So(metrics, ShouldHaveLength, 1)

	withFields.Info("foo")
	withFields.Infof("foo %d", 42)

	metrics = collect()
	a.So(metrics, ShouldHaveLength, 2)

	withFields.Warn("foo")
	withFields.Warnf("foo %d", 42)

	metrics = collect()
	a.So(metrics, ShouldHaveLength, 3)

	withFields.Error("foo")
	withFields.Errorf("foo %d", 42)

	metrics = collect()
	a.So(metrics, ShouldHaveLength, 4)

	withFields.Fatal("foo")
	withFields.Fatalf("foo %d", 42)

	metrics = collect()
	a.So(metrics, ShouldHaveLength, 5)
}

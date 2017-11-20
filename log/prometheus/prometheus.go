// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package prometheus

import (
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/prometheus/client_golang/prometheus"
)

var logCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "ttn",
		Subsystem: "log",
		Name:      "messages_total",
		Help:      "Total log messages.",
	},
	[]string{"level"},
)

func init() {
	prometheus.MustRegister(logCounter)
}

type prom struct {
	log.Interface
}

const (
	debug = "debug"
	info  = "info"
	warn  = "warn"
	err   = "error"
	fatal = "fatal"
)

func (p prom) Debug(msg string) {
	logCounter.WithLabelValues(debug).Inc()
	p.Interface.Debug(msg)
}
func (p prom) Info(msg string) {
	logCounter.WithLabelValues(info).Inc()
	p.Interface.Info(msg)
}
func (p prom) Warn(msg string) {
	logCounter.WithLabelValues(warn).Inc()
	p.Interface.Warn(msg)
}
func (p prom) Error(msg string) {
	logCounter.WithLabelValues(err).Inc()
	p.Interface.Error(msg)
}
func (p prom) Fatal(msg string) {
	logCounter.WithLabelValues(fatal).Inc()
	p.Interface.Fatal(msg)
}
func (p prom) Debugf(msg string, v ...interface{}) {
	logCounter.WithLabelValues(debug).Inc()
	p.Interface.Debugf(msg, v)
}
func (p prom) Infof(msg string, v ...interface{}) {
	logCounter.WithLabelValues(info).Inc()
	p.Interface.Infof(msg, v)
}
func (p prom) Warnf(msg string, v ...interface{}) {
	logCounter.WithLabelValues(warn).Inc()
	p.Interface.Warnf(msg, v)
}
func (p prom) Errorf(msg string, v ...interface{}) {
	logCounter.WithLabelValues(err).Inc()
	p.Interface.Errorf(msg, v)
}
func (p prom) Fatalf(msg string, v ...interface{}) {
	logCounter.WithLabelValues(fatal).Inc()
	p.Interface.Fatalf(msg, v)
}
func (p prom) WithField(k string, v interface{}) log.Interface {
	return prom{p.Interface.WithField(k, v)}
}
func (p prom) WithFields(f log.Fields) log.Interface {
	return prom{p.Interface.WithFields(f)}
}
func (p prom) WithError(err error) log.Interface {
	return prom{p.Interface.WithError(err)}
}

// Wrap wraps an existing logger, counting logs by level
func Wrap(logger log.Interface) log.Interface {
	return prom{logger}
}

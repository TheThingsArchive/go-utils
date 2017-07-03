// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package influx

import (
	"fmt"
	"sync"
	"time"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

// DefaultScalingInterval represents default scaling interval(time to wait before a new writer is started/killed) used by batching writer.
const DefaultScalingInterval = 500 * time.Millisecond

func newBatchPoints(bpConf influxdb.BatchPointsConfig) influxdb.BatchPoints {
	bp, err := influxdb.NewBatchPoints(bpConf)
	if err != nil {
		// Can only happen if there's an error in the code
		panic(fmt.Errorf("Invalid batch point configuration: %s", err))
	}
	return bp
}

// BatchPointWriter writes influxdb.BatchPoints to Influx database.
type BatchPointWriter interface {
	Write(bp influxdb.BatchPoints) error
}

// PointWriter writes *influxdb.Point to Influx database.
type PointWriter interface {
	Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error
}

type singlePointWriter struct {
	log    ttnlog.Interface
	writer BatchPointWriter
}

// NewSinglePointWriter creates new PointWriter, which writes points one-by-one
func NewSinglePointWriter(log ttnlog.Interface, w BatchPointWriter) PointWriter {
	return &singlePointWriter{
		log:    log,
		writer: w,
	}
}

func (w *singlePointWriter) Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error {
	bp := newBatchPoints(bpConf)
	bp.AddPoint(p)
	return w.writer.Write(bp)
}

type batchPoint struct {
	*influxdb.Point
	errch chan error
}

// pushError reports err(can be nil) to the waiter.
// It does not block and must only be executed once.
func (p *batchPoint) pushError(err error) {
	p.errch <- err
	close(p.errch)
}

func writeInBatches(log ttnlog.Interface, w BatchPointWriter, bpConf influxdb.BatchPointsConfig, scalingInterval time.Duration, ch <-chan *batchPoint) {
	log = log.WithField("config", bpConf)

	var points []*batchPoint
	for {
		select {
		case p := <-ch:
			points = append(points, p)
		default:
			if len(points) == 0 {
				select {
				case p := <-ch:
					points = append(points, p)
					continue
				case <-time.After(scalingInterval):
					log.Info("Removing batch writer instance")
					return
				}
			}

			bp := newBatchPoints(bpConf)
			for _, p := range points {
				bp.AddPoint(p.Point)
			}

			log.WithField("num", len(points)).Debug("Writing batch of points to Influx")
			err := w.Write(bp)
			for _, p := range points {
				go p.pushError(err)
			}
			points = points[:0]
		}
	}
}

type batchingWriter struct {
	log             ttnlog.Interface
	scalingInterval time.Duration

	writer     BatchPointWriter
	mutex      sync.RWMutex
	pointChans map[influxdb.BatchPointsConfig]chan *batchPoint
}

// NewBatchingWriter creates new PointWriter, which writes points in batches and scales automatically according to scalingInterval.
func NewBatchingWriter(log ttnlog.Interface, w BatchPointWriter, scalingInterval time.Duration) PointWriter {
	return &batchingWriter{
		log:             log,
		writer:          w,
		scalingInterval: scalingInterval,
		pointChans:      make(map[influxdb.BatchPointsConfig]chan *batchPoint),
	}
}

func (w *batchingWriter) Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error {
	w.mutex.RLock()
	ch, ok := w.pointChans[bpConf]
	w.mutex.RUnlock()
	if !ok {
		w.mutex.Lock()
		ch, ok = w.pointChans[bpConf]
		if !ok {
			ch = make(chan *batchPoint)
			w.pointChans[bpConf] = ch
			go writeInBatches(w.log, w.writer, bpConf, w.scalingInterval, ch)
		}
		w.mutex.Unlock()
	}

	point := &batchPoint{
		Point: p,
		errch: make(chan error, 1),
	}
	select {
	case ch <- point:
	case <-time.After(w.scalingInterval):
		w.log.WithField("config", bpConf).Info("Creating additional batch writer instance")
		go writeInBatches(w.log, w.writer, bpConf, w.scalingInterval, ch)
		ch <- point
	}
	return <-point.errch
}

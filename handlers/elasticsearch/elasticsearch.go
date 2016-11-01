// Copyright © 2016 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package elasticsearch

import (
	"fmt"
	"io"
	stdlog "log"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/tj/go-elastic/batch"
)

// Elasticsearch interface.
type Elasticsearch interface {
	Bulk(io.Reader) error
}

// Config for handler.
type Config struct {
	BufferSize int           // BufferSize is the number of logs to buffer before flush (default: 100)
	Client     Elasticsearch // Client for ES
	Prefix     string        // logs will be prefixed with this
}

// defaults applies defaults to the config.
func (c *Config) defaults() {
	if c.BufferSize == 0 {
		c.BufferSize = 100
	}

	if c.Prefix == "" {
		c.Prefix = "logs"
	}
}

// Handler implementation.
type Handler struct {
	*Config

	mu    sync.Mutex
	batch *batch.Batch
}

// indexName returns the index for the configured
func (h *Handler) indexName() string {
	return fmt.Sprintf("%s-%s", h.Config.Prefix, time.Now().Format("06-01-02"))
}

// New handler with BufferSize
func New(config *Config) *Handler {
	config.defaults()
	return &Handler{
		Config: config,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.batch == nil {
		h.batch = &batch.Batch{
			Elastic: h.Client,
			Index:   h.indexName(),
			Type:    "log",
		}
	}

	// Map fields
	for k, v := range e.Fields {
		switch t := v.(type) {
		case []byte: // addresses and EUIs are []byte
			e.Fields[k] = fmt.Sprintf("%X", t)
		case [21]byte: // bundle IDs [21]byte
			e.Fields[k] = fmt.Sprintf("%X-%X-%X-%X", t[0], t[1:9], t[9:17], t[17:])
		}
	}
	e.Timestamp = e.Timestamp.UTC()

	h.batch.Add(e)

	if h.batch.Size() >= h.BufferSize {
		go h.flush(h.batch)
		h.batch = nil
	}

	return nil
}

// Flush the current `batch`.
func (h *Handler) Flush() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.batch != nil {
		go h.flush(h.batch)
		h.batch = nil
	}
}

// flush the given `batch` asynchronously.
func (h *Handler) flush(batch *batch.Batch) {
	size := batch.Size()
	if err := batch.Flush(); err != nil {
		stdlog.Printf("log/elastic: failed to flush %d logs: %s", size, err)
	}
}

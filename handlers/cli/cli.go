// Copyright © 2016 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/apex/log"
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 90
)

// Colors mapping.
var Colors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// by sorts projects by call count.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu       sync.Mutex
	Writer   io.Writer
	UseColor bool
}

// colorTermSubstrings contains a list of substrings that indicate support for terminal colors
var colorTermSubstrings = []string{
	"color",
	"xterm",
}

// New handler.
func New(w io.Writer) *Handler {
	var useColor bool
	if os.Getenv("COLORTERM") != "" {
		useColor = true
	}
	if term := os.Getenv("TERM"); term != "" {
		for _, substring := range colorTermSubstrings {
			if strings.Contains(term, substring) {
				useColor = true
				break
			}
		}
	}
	return &Handler{
		Writer:   w,
		UseColor: useColor,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]

	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.UseColor {
		fmt.Fprintf(h.Writer, "\033[%dm%6s\033[0m %-40s", color, level, e.Message)
	} else {
		fmt.Fprintf(h.Writer, "%6s %-40s", level, e.Message)
	}

	for _, f := range fields {
		var value interface{}
		switch t := f.Value.(type) {
		case []byte: // addresses and EUIs are []byte
			value = fmt.Sprintf("%X", t)
		case [21]byte: // bundle IDs [21]byte
			value = fmt.Sprintf("%X-%X-%X-%X", t[0], t[1:9], t[9:17], t[17:])
		default:
			value = f.Value
		}

		if h.UseColor {
			fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, f.Name, value)
		} else {
			fmt.Fprintf(h.Writer, " %s=%v", f.Name, value)
		}

	}

	fmt.Fprintln(h.Writer)

	return nil
}

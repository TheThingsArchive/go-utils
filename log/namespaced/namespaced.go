package namespaced

import (
	"sync"

	"github.com/TheThingsNetwork/go-utils/log"
)

var NamespaceKey = "namespace"

// Namespaced is a logger that only logs an entry when the namespace of that
// entry is enabled
type Namespaced struct {
	sync.RWMutex
	log.Interface
	namespaces *ns
	namespace  string
}

// WithNamespace adds a namespace to the logging context
func WithNamespace(namespace string, ctx log.Interface) log.Interface {
	return ctx.WithField(NamespaceKey, namespace)
}

// Wrap wraps the logger in a Namespaced logger and enables the specified
// namespaces. See SetNamespaces for information on how to set the namspaces
func Wrap(ctx log.Interface, namespaces ...string) *Namespaced {
	return &Namespaced{
		Interface: ctx,
		namespaces: &ns{
			namespaces: namespaces,
		},
	}
}

// SetNamespaces replaces the set of enabled namespaces
// The namespaces follow this format:
// "*"   enables all namespaces
// "a"   enables namespace a
// "-a"  disables namespace a (even if "a" or "*" is set)
// For example:
//  - SetNamespaces("*", "-foo") enables every namespace but foo
//  - SetNamespaces("foo") enables only foo
//  - SetNamespaces() disables all namespaced entries
//
// Note that entries without a namespace will always be logged.
func (n *Namespaced) SetNamespaces(namespaces ...string) {
	n.namespaces.Set(namespaces)
}

// WithField adds a field to the logger
func (n *Namespaced) WithField(k string, v interface{}) log.Interface {
	if k == NamespaceKey {
		if str, ok := v.(string); ok {
			return &Namespaced{
				Interface:  n.Interface,
				namespaces: n.namespaces,
				namespace:  str,
			}
		}
	}

	return &Namespaced{
		Interface:  n.Interface.WithField(k, v),
		namespaces: n.namespaces,
		namespace:  n.namespace,
	}
}

// WithFields adds multiple fields to the logger
func (n *Namespaced) WithFields(fields log.Fields) log.Interface {
	return &Namespaced{
		Interface:  n.Interface.WithFields(fields),
		namespaces: n.namespaces,
		namespace:  n.namespace,
	}
}

func (n *Namespaced) WithError(err error) log.Interface {
	return &Namespaced{
		Interface:  n.Interface.WithError(err),
		namespaces: n.namespaces,
		namespace:  n.namespace,
	}
}

// isEnabdled returns whether or not this Namespaced logger should log,
// based on the enabled namespaces and the namespace
func (n *Namespaced) isEnabled() bool {
	return n.namespaces.IsEnabled(n.namespace)
}

func (n *Namespaced) Debug(msg string) {
	if n.isEnabled() {
		n.Interface.Debug(msg)
	}
}

func (n *Namespaced) Debugf(msg string, v ...interface{}) {
	if n.isEnabled() {
		n.Interface.Debugf(msg, v...)
	}
}

func (n *Namespaced) Info(msg string) {
	if n.isEnabled() {
		n.Interface.Info(msg)
	}
}

func (n *Namespaced) Infof(msg string, v ...interface{}) {
	if n.isEnabled() {
		n.Interface.Infof(msg, v...)
	}
}

func (n *Namespaced) Warn(msg string) {
	if n.isEnabled() {
		n.Interface.Warn(msg)
	}
}

func (n *Namespaced) Warnf(msg string, v ...interface{}) {
	if n.isEnabled() {
		n.Interface.Warnf(msg, v...)
	}
}

func (n *Namespaced) Error(msg string) {
	if n.isEnabled() {
		n.Interface.Error(msg)
	}
}

func (n *Namespaced) Errorf(msg string, v ...interface{}) {
	if n.isEnabled() {
		n.Interface.Errorf(msg, v...)
	}
}

func (n *Namespaced) Fatal(msg string) {
	if n.isEnabled() {
		n.Interface.Fatal(msg)
	}
}

func (n *Namespaced) Fatalf(msg string, v ...interface{}) {
	if n.isEnabled() {
		n.Interface.Fatalf(msg, v...)
	}
}

package namespaced

import "sync"

type ns struct {
	sync.RWMutex
	namespaces []string
}

// negate negates the namspace
func negate(namespace string) string {
	return "-" + namespace
}

// isEnabled checks wether or not the namespace is enabled
func (n *ns) IsEnabled(namespace string) bool {
	n.RLock()
	defer n.RUnlock()

	if namespace == "" {
		return true
	}

	hasStar := false
	included := false

	for _, ns := range n.namespaces {
		// if the namspace is negated, it can never be enabled
		if ns == negate(namespace) {
			return false
		}

		// if the namespace is explicitly enabled, mark it as included
		if ns == namespace {
			included = true
		}

		// mark that we have a *
		if ns == "*" {
			hasStar = true
		}
	}

	// non-mentioned namespaces are only enabled if we got the catch-all *
	return hasStar || included
}

// Set updates the namespaces
func (n *ns) Set(namespaces []string) {
	n.Lock()
	defer n.Unlock()
	n.namespaces = namespaces
}

// Package null is a no-op Notifier. It implements the Notifier interface,
// but does nothing on Changed() or Deleted() calls. Useful for
// testing.
package null

import (
	"github.com/bpineau/kube-deployments-notifier/config"
)

// Notifier implements notifiers.Notifier.
type Notifier struct {
}

// Changed do nothing.
func (l *Notifier) Changed(c *config.KdnConfig, msg string) error {
	return nil
}

// Deleted do nothing.
func (l *Notifier) Deleted(c *config.KdnConfig, msg string) error {
	return nil
}

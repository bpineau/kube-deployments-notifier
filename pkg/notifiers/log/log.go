package log

import (
	"github.com/bpineau/kube-deployments-notifier/config"
)

// Notifier implements notifiers.Notifier
type Notifier struct {
}

// Changed sends notification to the configured logrus logger
func (l *Notifier) Changed(c *config.KdnConfig, msg string) error {
	c.Logger.Infof("Changed: %s", msg)
	return nil
}

// Deleted sends notification to the configured logrus logger
func (l *Notifier) Deleted(c *config.KdnConfig, msg string) error {
	c.Logger.Infof("Deleted: %s", msg)
	return nil
}

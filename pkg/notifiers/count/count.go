package count

import (
	"sync"

	"github.com/bpineau/kube-deployments-notifier/config"
)

// Notifier implements notifiers.Notifier
type Notifier struct {
	sync.RWMutex
	counter int
}

// Changed increment the event counter
func (l *Notifier) Changed(c *config.KdnConfig, msg string) error {
	l.Lock()
	l.counter++
	l.Unlock()
	return nil
}

// Deleted increment the event counter
func (l *Notifier) Deleted(c *config.KdnConfig, msg string) error {
	l.Lock()
	l.counter++
	l.Unlock()
	return nil
}

// Count return the number of notification handled
func (l *Notifier) Count() int {
	l.RLock()
	defer l.RUnlock()
	return l.counter
}

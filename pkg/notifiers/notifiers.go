package notifiers

import (
	"fmt"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/http"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/log"
)

// Notifier sends message to the API endpoint
type Notifier interface {
	Changed(c *config.KdnConfig, msg string) error
	Deleted(c *config.KdnConfig, msg string) error
}

// Notifiers maps all known notifiers
var Notifiers = []Notifier{
	&log.Notifier{},
	&http.Notifier{},
}

// Changed send creation/change events to the notifiers
func Changed(c *config.KdnConfig, msg string) {
	if c.DryRun {
		fmt.Printf("Changed: %s\n", msg)
		return
	}

	for _, notifier := range Notifiers {
		err := notifier.Changed(c, msg)
		if err != nil {
			c.Logger.Warningf("Failed to notify: %s", err)
		}
	}
}

// Deleted send deletion events to the notifiers
func Deleted(c *config.KdnConfig, msg string) {
	if c.DryRun {
		fmt.Printf("Deleted: %s\n", msg)
		return
	}

	for _, notifier := range Notifiers {
		err := notifier.Deleted(c, msg)
		if err != nil {
			c.Logger.Warningf("Failed to notify: %s", err)
		}
	}
}

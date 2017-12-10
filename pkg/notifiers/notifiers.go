package notifiers

import (
	"fmt"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/http"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/log"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/null"
)

// Notifier sends message to the API endpoint (or other backend)
type Notifier interface {
	Changed(c *config.KdnConfig, msg string) error
	Deleted(c *config.KdnConfig, msg string) error
}

// Fakes map all test, fake notifiers
var Fakes = []Notifier{
	&count.Notifier{},
	&null.Notifier{},
}

// Backends maps all real, effective notifiers
var Backends = []Notifier{
	&log.Notifier{},
	&http.Notifier{},
}

// Composite combine and chain several notifiers
type Composite struct {
	Notifiers []Notifier
}

// Init initialize a Composite structure
func Init(notifiers []Notifier) *Composite {
	return &Composite{Notifiers: notifiers}
}

// Changed send creation/change events to the notifiers
func (n *Composite) Changed(c *config.KdnConfig, msg string) error {
	return n.invoke("Changed", c, msg)
}

// Deleted send deletion events to the notifiers
func (n *Composite) Deleted(c *config.KdnConfig, msg string) error {
	return n.invoke("Deleted", c, msg)
}

func (n *Composite) invoke(method string, c *config.KdnConfig, msg string) error {
	if c.DryRun {
		fmt.Printf("%s: %s\n", method, msg)
		return nil
	}

	var res error = nil
	for _, notifier := range n.Notifiers {
		var err error = nil
		if method == "Changed" {
			err = notifier.Changed(c, msg)
		} else {
			err = notifier.Deleted(c, msg)
		}

		if err != nil {
			c.Logger.Warningf("Failed to notify: %s", err)
			res = err
		}
	}
	return res
}

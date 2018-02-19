package notifiers

import (
	"fmt"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/http"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/log"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/null"
)

// Notifier convey Kubernetes events (creation/changes, deletion)
// as messages to dedicated backends (like an API endpoint).
type Notifier interface {
	Changed(c *config.KdnConfig, msg string) error
	Deleted(c *config.KdnConfig, msg string) error
}

// Fakes map all test, fake Notifiers.
var Fakes = []Notifier{
	&count.Notifier{},
	&null.Notifier{},
}

// Backends maps all real, effective Notifiers.
var Backends = []Notifier{
	&log.Notifier{},
	&http.Notifier{},
}

// Composite combine and chain several Notifiers, while implementing the
// Notifier interface itself.
type Composite struct {
	Notifiers []Notifier
}

// Init initialize a Composite structure
func Init(notifiers []Notifier) *Composite {
	return &Composite{Notifiers: notifiers}
}

// Changed send creation/change events notifications to the Notifier.
func (n *Composite) Changed(c *config.KdnConfig, msg string) error {
	return n.invoke("Changed", c, msg)
}

// Deleted send deletion events notifications to the Notifier.
func (n *Composite) Deleted(c *config.KdnConfig, msg string) error {
	return n.invoke("Deleted", c, msg)
}

func (n *Composite) invoke(method string, c *config.KdnConfig, msg string) error {
	if c.DryRun {
		fmt.Printf("%s: %s\n", method, msg)
		return nil
	}

	var res error
	for _, notifier := range n.Notifiers {
		var err error
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

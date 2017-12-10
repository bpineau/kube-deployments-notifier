package null

import (
	"testing"

	"github.com/bpineau/kube-deployments-notifier/config"
)

func TestNullNotifier(t *testing.T) {
	conf := new(config.KdnConfig)
	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a change")
	}

	err = notifier.Deleted(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a deletion")
	}
}

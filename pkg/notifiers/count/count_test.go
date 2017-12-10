package count

import (
	"testing"

	"github.com/bpineau/kube-deployments-notifier/config"
)

func TestCountNotifier(t *testing.T) {
	conf := new(config.KdnConfig)
	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err != nil || notifier.Count() != 1 {
		t.Errorf("Failed to notify a change")
	}

	err = notifier.Changed(conf, "bar")
	if err != nil || notifier.Count() != 2 {
		t.Errorf("Failed to notify a change")
	}

	err = notifier.Deleted(conf, "foo")
	if err != nil || notifier.Count() != 3 {
		t.Errorf("Failed to notify a deletion")
	}

}

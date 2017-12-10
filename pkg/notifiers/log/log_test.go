package log

import (
	"testing"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestLogNotifier(t *testing.T) {
	conf := new(config.KdnConfig)

	logger, hook := test.NewNullLogger()
	conf.Logger = logger

	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err != nil || len(hook.Entries) != 1 {
		t.Errorf("failed to notify a change")
	}

	err = notifier.Changed(conf, "bar")
	if err != nil || hook.LastEntry().Message != "Changed: bar" {
		t.Errorf("failed to notify a change %s", hook.LastEntry().Message)
	}

	err = notifier.Deleted(conf, "foo")
	if err != nil || hook.LastEntry().Message != "Deleted: foo" {
		t.Errorf("failed to notify a deletion %s", hook.LastEntry().Message)
	}
}

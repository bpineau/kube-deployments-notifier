package notifiers

import (
	"testing"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
	"github.com/sirupsen/logrus/hooks/test"
)

func initComposite() (*Composite, *config.KdnConfig) {
	conf := new(config.KdnConfig)
	backends := Init(Fakes)
	return backends, conf
}

func TestNotifyChangedComposite(t *testing.T) {
	backends, conf := initComposite()
	err := backends.Changed(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a change")
	}
}

func TestNotifyDeletedComposite(t *testing.T) {
	backends, conf := initComposite()
	err := backends.Deleted(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a deletion")
	}
}

func TestNotifyCounted(t *testing.T) {
	conf := new(config.KdnConfig)
	countNotifier := &count.Notifier{}
	backends := &Composite{Notifiers: []Notifier{countNotifier}}

	err := backends.Changed(conf, "foo")
	if err != nil || countNotifier.Count() != 1 {
		t.Errorf("Wrong number of notification recieved after a change")
	}

	err = backends.Deleted(conf, "foo")
	if err != nil || countNotifier.Count() != 2 {
		t.Errorf("Wrong number of notification recieved after a deletion")
	}

	conf.DryRun = true
	err = backends.Changed(conf, "foo")
	if err != nil || countNotifier.Count() != 2 {
		t.Errorf("Dry run was not honored")
	}
}

func TestNotifyForwardError(t *testing.T) {
	logger, _ := test.NewNullLogger()
	conf := new(config.KdnConfig)
	conf.Endpoint = "i'm a broken url"
	conf.Logger = logger
	backends := Init(Backends)

	err := backends.Changed(conf, "foo")
	if err == nil {
		t.Errorf("Failed to forward up errors from backends")
	}
}

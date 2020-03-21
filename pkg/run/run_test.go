package run

import (
	"syscall"
	"testing"
	"time"

	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers/deployment"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
)

var (
	obj1 = &apps_v1.Deployment{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test",
		Labels:    config.Labels,
		Namespace: v1.NamespaceDefault},
	}
)

func TestRun(t *testing.T) {
	conf := config.FakeConfig(obj1)
	runFromConf(t, conf)
}

func TestRunFailingHealthcheck(t *testing.T) {
	conf := config.FakeConfig(obj1)
	conf.Endpoint = ""
	conf.HealthPort = -1
	runFromConf(t, conf)
}

func runFromConf(t *testing.T, conf *config.KdnConfig) {
	conts = []controllers.Controller{
		&deployment.Controller{},
	}
	notif := new(count.Notifier)
	go Run(conf, notif)

	ch := make(chan int, 1)
	defer close(ch)

	go func() {
		for notif.Count() == 0 {
			time.Sleep(time.Second)
		}
		ch <- notif.Count()
	}()

	select {
	case res := <-ch:
		if res != 1 {
			t.Errorf("Failed to convert an event to a notification")
		}
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for a event to pop up as a notification")
	}

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //nolint:errcheck
	time.Sleep(300 * time.Millisecond)              // controllers wait for 200ms before stopping
}

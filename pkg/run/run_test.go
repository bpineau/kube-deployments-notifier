package run

import (
	"syscall"
	"testing"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
)

var (
	obj1 = &v1beta1.Deployment{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test",
		Labels:    config.Labels,
		Namespace: v1.NamespaceDefault},
	}
)

func TestRun(t *testing.T) {
	conf := config.FakeConfig(obj1)
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

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}

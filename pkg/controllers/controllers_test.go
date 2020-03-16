package controllers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers/count"
)

var (
	obj1 = &v1.Pod{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test1",
		Labels:    config.Labels,
		Namespace: v1.NamespaceDefault},
	}
	obj2 = &v1.Pod{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test2",
		Labels:    config.Labels,
		Namespace: v1.NamespaceDefault},
	}
	obj3 = &v1.Pod{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test3",
		Labels:    map[string]string{"wrong": "label"},
		Namespace: v1.NamespaceDefault},
	}
	obj4 = &v1.Pod{ObjectMeta: meta_v1.ObjectMeta{
		Name:      "test4",
		Labels:    config.Labels,
		Namespace: v1.NamespaceDefault},
	}
)

type testController struct {
	CommonController
}

// Init initialize pod controller
func (c *testController) Init(conf *config.KdnConfig, n notifiers.Notifier) Controller {
	c.CommonController = CommonController{
		Conf:      conf,
		Name:      "pod",
		Notifiers: n,
	}

	client := c.Conf.ClientSet
	c.ObjType = &v1.Pod{}
	selector := meta_v1.ListOptions{LabelSelector: conf.Filter}

	ls := new(cache.ListWatch)
	ls.ListFunc = func(options meta_v1.ListOptions) (runtime.Object, error) {
		return client.CoreV1().Pods(meta_v1.NamespaceAll).List(selector)
	}
	ls.WatchFunc = func(options meta_v1.ListOptions) (watch.Interface, error) {
		return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(selector)
	}
	c.ListWatch = ls

	return c
}

type failingNotifier struct {
	sync.Mutex
	changeCalls int
	deleteCalls int
}

// Changed is a failing and counting change notifier
func (l *failingNotifier) Changed(c *config.KdnConfig, msg string) error {
	l.Lock()
	l.changeCalls++
	l.Unlock()
	return fmt.Errorf("Normal and expected error on change event")
}

// Deleted is a failing and counting delete notifier
func (l *failingNotifier) Deleted(c *config.KdnConfig, msg string) error {
	l.Lock()
	l.deleteCalls++
	l.Unlock()
	return fmt.Errorf("Normal and expected error on delete event")
}

func (l *failingNotifier) countChange() int {
	l.Lock()
	defer l.Unlock()
	return l.changeCalls
}

func (l *failingNotifier) countDelete() int {
	l.Lock()
	defer l.Unlock()
	return l.deleteCalls
}

func TestController(t *testing.T) {
	c := config.FakeConfig(obj1, obj2, obj3)
	n := new(count.Notifier)
	cont := &testController{}
	cont.Init(c, n)

	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go cont.Start(&wg)
	defer func(cont Controller) {
		go cont.Stop()
	}(cont)

	// test initial event list and filtering
	time.Sleep(config.FakeResyncInterval + config.FakeResyncInterval/2)
	if n.Count() != 4 { // 2 initial + 2 after resync
		t.Errorf("Expected 4 notification, got %d", n.Count())
	}

	// test deletion
	store := cont.Informer.GetStore()
	err := store.Delete(obj2)
	if err != nil {
		t.Fatalf("Unexcepted error %v", err)
	}
	time.Sleep(config.FakeResyncInterval + config.FakeResyncInterval/2)

	if n.Count() < 5 {
		t.Errorf("Expected 5 notification, got %d", n.Count())
	}

	// test deletion notification path
	cont.Queue.Add("fake/fake")

	// test retries on notifiers failure
	fnotif := new(failingNotifier)
	cont.Notifiers = fnotif
	err = store.Add(obj4)
	if err != nil {
		t.Fatalf("Unexcepted error %v", err)
	}

	time.Sleep(config.FakeResyncInterval + config.FakeResyncInterval/2)
	if fnotif.countChange() < maxProcessRetry {
		t.Errorf("Should retry %d times on failing notifiers, got %d retries",
			maxProcessRetry,
			fnotif.countChange())
	}
	if fnotif.countDelete() < maxProcessRetry {
		t.Errorf("Notified should get a deletion event")
	}
}

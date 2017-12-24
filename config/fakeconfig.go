package config

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/bpineau/kube-deployments-notifier/pkg/log"
)

var (
	// FakeResyncInterval is the interval between resyncs during tests
	FakeResyncInterval = time.Duration(time.Second)

	// Labels use to filter objets in the tests runs
	Labels = map[string]string{"foo": "bar", "spam": "egg"}
)

// FakeConfig returns a config objet using a fake clientset (for tests)
func FakeConfig(objects ...runtime.Object) *KdnConfig {
	c := &KdnConfig{
		DryRun:     true,
		Logger:     log.New("", "", "test"),
		ClientSet:  fake.NewSimpleClientset(objects...),
		Endpoint:   "http://example.com",
		TokenHdr:   "",
		TokenVal:   "",
		Filter:     "foo=bar,spam=egg",
		ResyncIntv: FakeResyncInterval,
	}

	return c
}

// FakeClientSet provides a fake.NewSimpleClientset, useful for testing
func FakeClientSet() *fake.Clientset {
	return fake.NewSimpleClientset()
}

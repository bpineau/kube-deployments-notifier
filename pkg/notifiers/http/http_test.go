package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bpineau/kube-deployments-notifier/config"
)

func TestHttpNotifier(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Foo") != "Bar" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	conf := new(config.KdnConfig)
	conf.Endpoint = ts.URL
	conf.TokenHdr = "X-Foo"
	conf.TokenVal = "Bar"

	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a change")
	}

	err = notifier.Deleted(conf, "foo")
	if err != nil {
		t.Errorf("Failed to notify a deletion")
	}

	conf.Endpoint = ""
	err = notifier.Changed(conf, "foo")
	if err != nil {
		t.Errorf("HTTP notifier should ignore null endpoints")
	}
}

func TestHttpNotifierNoEndpoint(t *testing.T) {
	conf := new(config.KdnConfig)
	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err != nil {
		t.Errorf("Http shouldn't fail on nil endpoint")
	}
}

func TestHttpNotifierFailures(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	conf := new(config.KdnConfig)
	conf.Endpoint = ts.URL

	notifier := new(Notifier)

	err := notifier.Changed(conf, "foo")
	if err == nil {
		t.Errorf("Failed to notice request failure")
	}

	conf.Endpoint = "i'm a broken url"

	err = notifier.Changed(conf, "bar")
	if err == nil {
		t.Errorf("Failed to notice request failure")
	}

	conf.Endpoint = ":"

	err = notifier.Changed(conf, "baz")
	if err == nil {
		t.Errorf("Failed to notice a unparsable url")
	}

}

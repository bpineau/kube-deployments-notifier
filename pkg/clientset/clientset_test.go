package clientset

import (
	"testing"
)

const nonExistentPath = "\\/hopefully/non/existent/path"

func TestClientSet(t *testing.T) {
	_, err := NewClientSet("http://127.0.0.1", "", nonExistentPath)
	if err == nil {
		t.Fatal("NewClientSet() should fail on non existent kubeconfig path")
	}
}

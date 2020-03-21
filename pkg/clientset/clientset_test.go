package clientset

import (
	"os"
	"testing"
)

const nonExistentPath = "\\/hopefully/non/existent/path"

func TestClientSet(t *testing.T) {
	_, err := NewClientSet("http://127.0.0.1", nonExistentPath)
	if err == nil {
		t.Fatal("NewClientSet() should fail on non existent kubeconfig path")
	}

	_ = os.Unsetenv("KUBERNETES_SERVICE_HOST")
	_ = os.Setenv("HOME", nonExistentPath)
	_, err = NewClientSet("", "")
	if err == nil {
		t.Fatal("NewClientSet() should fail to load InClusterConfig without kube address env")
	}
}

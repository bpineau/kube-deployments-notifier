package clientset

import (
	"os"
	"fmt"
	"testing"
)

const nonExistentPath = "\\/hopefuly/non/existent/path\t"

func TestClientSet(t *testing.T) {
	cs, err := NewClientSet("", "")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%T", cs) != "*kubernetes.Clientset" {
		t.Errorf("NewClientSet() didn't return a *kubernetes.Clientset: %T", cs)
	}

	cs, err = NewClientSet("http://127.0.0.1", "/dev/null")
	if fmt.Sprintf("%T", cs) != "*kubernetes.Clientset" {
		t.Errorf("NewClientSet(server) didn't return a *kubernetes.Clientset: %T", cs)
	}

	cs, err = NewClientSet("http://127.0.0.1", nonExistentPath)
	if err == nil {
		t.Fatal("NewClientSet() should fail on non existent kubeconfig path")
	}

	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Setenv("HOME", nonExistentPath)
	cs, err = NewClientSet("", "")
	if err == nil {
		t.Fatal("NewClientSet() should fail to load InClusterConfig without kube address env")
	}
}
